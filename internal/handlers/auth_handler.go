package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/codedbyshoe/goxclone/internal/models"
	"github.com/codedbyshoe/goxclone/internal/services/forms"
	"github.com/codedbyshoe/goxclone/internal/services/hash"
	"github.com/codedbyshoe/goxclone/internal/services/hash/passwordhash"
	"github.com/codedbyshoe/goxclone/internal/views"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo     models.UserRepo
	sessionRepo  models.SessionRepo
	passwordhash hash.PasswordHash
	*chi.Mux
	db                *gorm.DB
	createUserForm    *forms.CreateUserForm
	authUserForm      *forms.AuthenticateUserForm
	sessionCookieName string
}

func NewAuthHandler(db *gorm.DB, ur models.UserRepo, sr models.SessionRepo, scn string) *AuthHandler {
	h := AuthHandler{
		Mux:               chi.NewMux(),
		db:                db,
		userRepo:          ur,
		sessionRepo:       sr,
		passwordhash:      passwordhash.NewHPasswordHash(),
		createUserForm:    forms.NewCreateUserForm(),
		authUserForm:      forms.NewAuthenticateUserForm(),
		sessionCookieName: scn,
	}

	h.Route("/", func(r chi.Router) {
		r.Get("/", h.Index)
		r.Post("/create", h.CreateUser)
		r.Post("/login", h.PostLogin)
	})

	return &h
}

// GET auth endpoint: /auth
func (h *AuthHandler) Index(w http.ResponseWriter, r *http.Request) {
	views.AuthLayout("Welcome", h.authUserForm, h.createUserForm).Render(r.Context(), w)
}

// POST auth create user endpoint: /auth/create
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	redirectWithQuery := "/auth?showsignupform"
	h.createUserForm.Name = r.FormValue("name")
	h.createUserForm.Username = r.FormValue("username")
	h.createUserForm.Email = r.FormValue("email")
	h.createUserForm.Password = r.FormValue("password")
	h.createUserForm.ConfirmPassword = r.FormValue("confirm_password")

	if !h.createUserForm.Validate() {
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	hash, err := h.passwordhash.GenerateFromPassword(h.createUserForm.Password)
	if err != nil {
		h.createUserForm.FormErrors.Global = "Errors processing request something went wrong"
		http.Redirect(w, r, redirectWithQuery, http.StatusPermanentRedirect)
		return
	}

	user := models.User{
		Name:     h.createUserForm.Name,
		Username: h.createUserForm.Username,
		Email:    h.createUserForm.Email,
		Password: string(hash),
	}

	existingUser, _ := h.userRepo.GetUser(user.Email)
	if existingUser != nil {
		h.createUserForm.FormErrors.Global = "User with this email already exists."
		http.Redirect(w, r, redirectWithQuery, http.StatusPermanentRedirect)
		return
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		h.createUserForm.FormErrors.Global = "Errors processing request something went wrong"
		http.Redirect(w, r, redirectWithQuery, http.StatusPermanentRedirect)
		return
	}

	session, err := h.sessionRepo.CreateSession(&models.Session{
		UserID: user.ID,
	})
	if err != nil {
		h.authUserForm.FormErrors.Global = "Internal Server Error. Something went wrong"
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	cookie := createCookie(user.ID, session.SessionID, h.sessionCookieName)

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	redirectWithQuery := "/auth?showloginform"
	h.authUserForm.Email = r.FormValue("email")
	h.authUserForm.Password = r.FormValue("password")

	if !h.authUserForm.Validate() {
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	user, err := h.userRepo.GetUser(h.authUserForm.Email)
	if err != nil {
		h.authUserForm.FormErrors.Add("email", "Invalid login credentials")
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	// check password hash
	ok, _ := h.passwordhash.ComparePasswordAndHash(h.authUserForm.Password, user.Password)
	if !ok {
		// NOTE: testing purposes only
		h.authUserForm.FormErrors.Add("password", "Invalid login credentials")
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	session, err := h.sessionRepo.CreateSession(&models.Session{
		UserID: user.ID,
	})
	if err != nil {
		h.authUserForm.FormErrors.Global = "Internal Server Error. Something went wrong"
		http.Redirect(w, r, redirectWithQuery, http.StatusSeeOther)
		return
	}

	cookie := createCookie(user.ID, session.SessionID, h.sessionCookieName)

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createCookie(userid uint, sessionid string, cookiename string) *http.Cookie {
	cookieValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d", sessionid, userid)))
	expiration := time.Now().Add(365 * 24 * time.Hour)
	return &http.Cookie{
		Name:     cookiename,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/codedbyshoe/goxclone/internal/models"
)

type AuthMiddleware struct {
	sessionRepo       models.SessionRepo
	sessionCookieName string
}

func NewAuthMiddleware(r models.SessionRepo, c string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo:       r,
		sessionCookieName: c,
	}
}

type UserContextKey string

var UserKey UserContextKey = "user"

func (m *AuthMiddleware) RequireGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := m.hasValidCookie(r)
		if ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := m.hasValidCookie(r)
		if !ok {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := m.hasValidCookie(r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) hasValidCookie(r *http.Request) (*models.User, bool) {
	sessionCookie, err := r.Cookie(m.sessionCookieName)
	if err != nil {
		return nil, false
	}

	decodedValue, err := base64.StdEncoding.DecodeString(sessionCookie.Value)
	if err != nil {
		return nil, false
	}

	splitValue := strings.Split(string(decodedValue), ":")

	if len(splitValue) != 2 {
		return nil, false
	}

	sessionID := splitValue[0]
	userID := splitValue[1]

	user, err := m.sessionRepo.GetUserFromSession(sessionID, userID)
	if err != nil {
		return nil, false
	}
	return user, true
}

func GetUser(ctx context.Context) *models.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}

	return user.(*models.User)
}

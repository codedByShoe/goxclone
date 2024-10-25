package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
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
		if hasCookie(r) {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasCookie(r) {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(m.sessionCookieName)
		if err != nil {
			fmt.Println("error getting session cookie", err)
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := base64.StdEncoding.DecodeString(sessionCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := splitValue[0]
		userID := splitValue[1]

		fmt.Println("sessionID", sessionID)
		fmt.Println("userID", userID)

		user, err := m.sessionRepo.GetUserFromSession(sessionID, userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hasCookie(r *http.Request) bool {
	_, err := r.Cookie("access_token")
	return err == nil
}

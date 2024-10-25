package main

import (
	"log"
	"net/http"
	"time"

	"github.com/codedbyshoe/goxclone/internal/handlers"
	"github.com/codedbyshoe/goxclone/internal/middleware"
	"github.com/codedbyshoe/goxclone/internal/models"
	"github.com/codedbyshoe/goxclone/internal/models/repo"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database Error: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Notification{},
	)

	ur := repo.NewUserRepo(db)
	sr := repo.NewSessionRepo(db)
	ph := handlers.NewPostHandler(db)
	ah := handlers.NewAuthHandler(db, ur, sr)
	auth := middleware.NewAuthMiddleware(sr, "access_token")

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(auth.AddUserToContext)
		router.Mount("/", auth.RequireAuth(ph))
	})

	router.Group(func(r chi.Router) {
		r.Use()
		router.Mount("/auth", auth.RequireGuest(ah))
	})

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

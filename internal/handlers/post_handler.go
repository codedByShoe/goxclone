package handlers

import (
	"net/http"

	"github.com/codedbyshoe/goxclone/internal/models"
	"github.com/codedbyshoe/goxclone/internal/views"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type PostHandler struct {
	*chi.Mux
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	h := PostHandler{
		Mux: chi.NewMux(),
		db:  db,
	}

	h.Route("/", func(r chi.Router) {
		r.Get("/", h.Home)
	})

	return &h
}

// App home view
func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	post1 := models.Post{
		Model:   gorm.Model{ID: 1},
		Content: "This is the body of a post",
		User: models.User{
			Model:    gorm.Model{ID: 1},
			Name:     "Andrew Shoemaker",
			Username: "codedbyshoe",
			Email:    "andrew.shoemaker9@gmail.com",
			Password: "somethingsecret",
		},
	}

	post2 := models.Post{
		Model:   gorm.Model{ID: 2},
		Content: "This is the body of a second post",
		User: models.User{
			Model:    gorm.Model{ID: 2},
			Name:     "Someone Something",
			Username: "somethingRandom",
			Email:    "timmyTester@gmail.com",
			Password: "somethingsecret",
		},
	}

	posts = append(posts, post1, post2)
	views.Layout(views.IndexPage(posts), "Home").Render(r.Context(), w)
}

// GET  list all post resources path: /posts
func (h *PostHandler) IndexPost(w http.ResponseWriter, r *http.Request) {
}

// POST create a post resource path: /posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
}

// GET show post resource path: /posts/{id}
func (h *PostHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
}

// GET edit post resource path: /posts/{id}/edit
func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
}

// PUT edit post resource path: /posts/{id}
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
}

// DELETE remove post resource path: /posts/{id}
func (h *PostHandler) DestroyPost(w http.ResponseWriter, r *http.Request) {
}

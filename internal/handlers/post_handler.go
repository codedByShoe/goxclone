package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/codedbyshoe/goxclone/internal/middleware"
	"github.com/codedbyshoe/goxclone/internal/models"
	"github.com/codedbyshoe/goxclone/internal/services/forms"
	"github.com/codedbyshoe/goxclone/internal/views"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	*chi.Mux
	postRepo       models.PostRepo
	createPostForm *forms.CreatePostForm
}

func NewPostHandler(pr models.PostRepo) *PostHandler {
	h := PostHandler{
		Mux:            chi.NewMux(),
		postRepo:       pr,
		createPostForm: forms.NewCreatePostForm(),
	}

	h.Route("/", func(r chi.Router) {
		r.Get("/", h.Home)
		r.Get("/posts/{id}", h.EditPost)
		r.Post("/posts/{id}", h.UpdatePost)
		r.Post("/posts", h.CreatePost)
		r.Post("/posts/delete", h.DestroyPost)
	})

	return &h
}

// App home view
func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	posts, err := h.postRepo.GetUsersPosts(user.ID)
	if err != nil {
		h.createPostForm.FormErrors.Global = err.Error()
		return
	}

	views.Layout(views.IndexPage(posts), "Home", h.createPostForm).Render(r.Context(), w)
}

// GET  list all post resources path: /posts
func (h *PostHandler) IndexPost(w http.ResponseWriter, r *http.Request) {
}

// POST create a post resource path: /posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	h.createPostForm.Content = r.FormValue("content")
	h.createPostForm.ConvertUserId(r.FormValue("user_id"))
	post := &models.Post{
		Content: h.createPostForm.Content,
		UserId:  h.createPostForm.UserId,
	}

	if err := h.postRepo.CreatePost(post); err != nil {
		h.createPostForm.FormErrors.Global = "Internal Server Error. Please try again"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET show post resource path: /posts/{id}
func (h *PostHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
}

// GET edit post resource path: /posts/{id}/edit
func (h *PostHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	post, err := h.postRepo.GetPost(uint(id))
	if err != nil {
		h.createPostForm.FormErrors.Global = "Could not find post"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	views.Layout(views.EditPostPage(post), "Edit", h.createPostForm).Render(r.Context(), w)
}

// PUT edit post resource path: /posts/{id}
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	postId, _ := strconv.ParseUint(idStr, 10, 64)
	content := r.FormValue("content")
	userIdStr := r.FormValue("user_id")
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	if strings.TrimSpace(content) == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err := h.postRepo.UpdatePost(uint(postId), uint(userId), content)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DELETE remove post resource path: /posts/{id}
// NOTE: This is currently a post request for now until I add _method middleware
func (h *PostHandler) DestroyPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.FormValue("post_id")
	userIdStr := r.FormValue("user_id")
	postId, _ := strconv.ParseUint(postIdStr, 10, 64)
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	if err := h.postRepo.DeletePost(uint(postId), uint(userId)); err != nil {
		h.createPostForm.FormErrors.Global = "Internal Server Error"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package forms

import (
	"strconv"
	"strings"
)

type CreatePostForm struct {
	FormErrors *FormErrors
	Content    string
	UserId     uint
}

func NewCreatePostForm() *CreatePostForm {
	return &CreatePostForm{
		FormErrors: &FormErrors{},
	}
}

func (f *CreatePostForm) Validate() {
	f.FormErrors = &FormErrors{}

	// Username validation
	if strings.TrimSpace(f.Content) == "" {
		f.FormErrors.Add("Content", "Post Content is required")
	}
}

func (f *CreatePostForm) ConvertUserId(stringId string) {
	id, _ := strconv.ParseUint(stringId, 10, 64)
	f.UserId = uint(id)
}

package forms

type FormErrors struct {
	Errors map[string]string
	Global string
}

func (f *FormErrors) Add(field, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[field] = message
}

func (f *FormErrors) HasErrors() bool {
	return len(f.Errors) > 0 || f.Global != ""
}

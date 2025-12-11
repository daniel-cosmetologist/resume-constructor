package resume

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldError описывает ошибку конкретного поля.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError агрегирует ошибки валидации.
type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation error"
}

func (e *ValidationError) Add(field, msg string) {
	e.Errors = append(e.Errors, FieldError{
		Field:   field,
		Message: msg,
	})
}

func (e *ValidationError) Empty() bool {
	return len(e.Errors) == 0
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// ValidateResume выполняет базовую валидацию резюме.
func ValidateResume(r Resume) error {
	var ve ValidationError

	if strings.TrimSpace(r.FullName) == "" {
		ve.Add("fullName", "Full name is required")
	}
	if len(r.FullName) > 100 {
		ve.Add("fullName", "Full name is too long (max 100 characters)")
	}

	if strings.TrimSpace(r.Position) == "" {
		ve.Add("position", "Position is required")
	}
	if len(r.Position) > 100 {
		ve.Add("position", "Position is too long (max 100 characters)")
	}

	if len(strings.TrimSpace(r.Summary)) == 0 {
		ve.Add("summary", "Summary is required")
	}
	if len(r.Summary) > 1500 {
		ve.Add("summary", "Summary is too long (max 1500 characters)")
	}

	validateContacts(r.Contacts, &ve)

	if len(r.Skills) > 50 {
		ve.Add("skills", "Too many skills (max 50)")
	}
	for i, skill := range r.Skills {
		if len(skill) > 50 {
			ve.Add(fmt.Sprintf("skills[%d]", i), "Skill is too long (max 50 characters)")
		}
	}

	if len(r.Experience) > 10 {
		ve.Add("experience", "Too many experience entries (max 10)")
	}
	if len(r.Education) > 10 {
		ve.Add("education", "Too many education entries (max 10)")
	}
	if len(r.CustomSections) > 10 {
		ve.Add("customSections", "Too many custom sections (max 10)")
	}

	if !ve.Empty() {
		return &ve
	}
	return nil
}

func validateContacts(c Contacts, ve *ValidationError) {
	email := strings.TrimSpace(c.Email)
	if email != "" && !emailRe.MatchString(email) {
		ve.Add("contacts.email", "Invalid email format")
	}

	for i, l := range c.Links {
		url := strings.TrimSpace(l.URL)
		if url == "" {
			continue
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			ve.Add(fmt.Sprintf("contacts.links[%d].url", i), "URL must start with http:// or https://")
		}
	}
}

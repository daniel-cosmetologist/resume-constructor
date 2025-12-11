package resume

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// ValidationError описывает ошибку валидации с набором полей.
type ValidationError struct {
	fields map[string]string
}

// Error реализует интерфейс error.
func (e *ValidationError) Error() string {
	return "validation error"
}

// Fields возвращает карту ошибок по полям.
func (e *ValidationError) Fields() map[string]string {
	if e == nil {
		return nil
	}
	return e.fields
}

func (e *ValidationError) add(field, msg string) {
	if e.fields == nil {
		e.fields = make(map[string]string)
	}
	if _, exists := e.fields[field]; !exists {
		e.fields[field] = msg
	}
}

// Validate выполняет базовую валидацию резюме по требованиям MVP.
func Validate(resume Resume) error {
	var verr ValidationError

	// fullName
	if strings.TrimSpace(resume.FullName) == "" {
		verr.add("fullName", "Full name is required")
	} else if utf8.RuneCountInString(resume.FullName) > 100 {
		verr.add("fullName", "Full name must be at most 100 characters")
	}

	// position
	if strings.TrimSpace(resume.Position) == "" {
		verr.add("position", "Position is required")
	} else if utf8.RuneCountInString(resume.Position) > 100 {
		verr.add("position", "Position must be at most 100 characters")
	}

	// summary
	if utf8.RuneCountInString(resume.Summary) > 1500 {
		verr.add("summary", "Summary must be at most 1500 characters")
	}

	// contacts
	if strings.TrimSpace(resume.Contacts.Email) == "" {
		verr.add("contacts.email", "Email is required")
	} else if !isPlausibleEmail(resume.Contacts.Email) {
		verr.add("contacts.email", "Email format is invalid")
	}
	if utf8.RuneCountInString(resume.Contacts.Location) > 200 {
		verr.add("contacts.location", "Location is too long")
	}

	// skills
	if len(resume.Skills) > 50 {
		verr.add("skills", "Too many skills (maximum is 50)")
	}
	for i, skill := range resume.Skills {
		if strings.TrimSpace(skill) == "" {
			verr.add(fieldIndex("skills", i), "Skill cannot be empty")
		} else if utf8.RuneCountInString(skill) > 50 {
			verr.add(fieldIndex("skills", i), "Skill must be at most 50 characters")
		}
	}

	// experience
	if len(resume.Experience) > 10 {
		verr.add("experience", "Too many experience entries (maximum is 10)")
	}
	for i, exp := range resume.Experience {
		prefix := fieldIndex("experience", i)

		if strings.TrimSpace(exp.Company) == "" {
			verr.add(prefix+".company", "Company is required")
		}
		if strings.TrimSpace(exp.Position) == "" {
			verr.add(prefix+".position", "Position is required")
		}
		if strings.TrimSpace(exp.StartDate) == "" {
			verr.add(prefix+".startDate", "Start date is required")
		}
		if !isYearMonth(exp.StartDate) {
			verr.add(prefix+".startDate", "Start date must be in YYYY-MM format")
		}
		if strings.TrimSpace(exp.EndDate) != "" && !isYearMonth(exp.EndDate) {
			verr.add(prefix+".endDate", "End date must be in YYYY-MM format or empty")
		}

		if utf8.RuneCountInString(exp.Description) > 500 {
			verr.add(prefix+".description", "Description is too long")
		}

		if len(exp.Bullets) > 20 {
			verr.add(prefix+".bullets", "Too many bullet points (maximum is 20)")
		}
		for j, b := range exp.Bullets {
			if utf8.RuneCountInString(b) > 300 {
				verr.add(fieldIndex(prefix+".bullets", j), "Bullet must be at most 300 characters")
			}
		}
	}

	// education
	if len(resume.Education) > 10 {
		verr.add("education", "Too many education entries (maximum is 10)")
	}
	for i, ed := range resume.Education {
		prefix := fieldIndex("education", i)

		if strings.TrimSpace(ed.Institution) == "" {
			verr.add(prefix+".institution", "Institution is required")
		}
		if strings.TrimSpace(ed.StartDate) == "" {
			verr.add(prefix+".startDate", "Start date is required")
		}
		if !isYearMonth(ed.StartDate) {
			verr.add(prefix+".startDate", "Start date must be in YYYY-MM format")
		}
		if strings.TrimSpace(ed.EndDate) != "" && !isYearMonth(ed.EndDate) {
			verr.add(prefix+".endDate", "End date must be in YYYY-MM format or empty")
		}

		if utf8.RuneCountInString(ed.Details) > 500 {
			verr.add(prefix+".details", "Details are too long")
		}
	}

	// custom sections
	if len(resume.CustomSections) > 10 {
		verr.add("customSections", "Too many custom sections (maximum is 10)")
	}
	for i, cs := range resume.CustomSections {
		prefix := fieldIndex("customSections", i)

		if strings.TrimSpace(cs.Title) == "" {
			verr.add(prefix+".title", "Title is required")
		}
		if utf8.RuneCountInString(cs.Title) > 100 {
			verr.add(prefix+".title", "Title must be at most 100 characters")
		}

		if len(cs.Items) > 10 {
			verr.add(prefix+".items", "Too many items (maximum is 10)")
		}
		for j, item := range cs.Items {
			if strings.TrimSpace(item) == "" {
				verr.add(fieldIndex(prefix+".items", j), "Item cannot be empty")
			} else if utf8.RuneCountInString(item) > 300 {
				verr.add(fieldIndex(prefix+".items", j), "Item must be at most 300 characters")
			}
		}
	}

	// photo
	if resume.Photo != nil {
		mt := strings.ToLower(strings.TrimSpace(resume.Photo.MimeType))
		if mt != "image/jpeg" && mt != "image/png" {
			verr.add("photo.mimeType", "Photo mimeType must be image/jpeg or image/png")
		}
		if strings.TrimSpace(resume.Photo.Data) == "" {
			verr.add("photo.data", "Photo data is empty")
		}
	}

	if len(verr.fields) == 0 {
		return nil
	}
	return &verr
}

// isPlausibleEmail выполняет простую проверку email-адреса.
func isPlausibleEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

// isYearMonth проверяет формат YYYY-MM.
func isYearMonth(value string) bool {
	if len(value) != 7 {
		return false
	}
	for i, r := range value {
		if i == 4 {
			if r != '-' {
				return false
			}
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// fieldIndex формирует имя поля с индексом, например "skills[0]".
func fieldIndex(prefix string, index int) string {
	return prefix + "[" + strconv.Itoa(index) + "]"
}

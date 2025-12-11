package model

// Resume описывает структуру резюме, используемую latex-service.
type Resume struct {
	FullName       string            `json:"fullName"`
	Position       string            `json:"position"`
	Summary        string            `json:"summary"`
	Contacts       Contacts          `json:"contacts"`
	Skills         []string          `json:"skills"`
	Experience     []ExperienceEntry `json:"experience"`
	Education      []EducationEntry  `json:"education"`
	CustomSections []CustomSection   `json:"customSections"`
	Photo          *Photo            `json:"photo,omitempty"`
}

// Contacts описывает контактные данные пользователя.
type Contacts struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Links    []Link `json:"links"`
}

// Link описывает внешнюю ссылку (GitHub, LinkedIn и т.п.).
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ExperienceEntry описывает запись об опыте работы.
type ExperienceEntry struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Bullets     []string `json:"bullets"`
}

// EducationEntry описывает запись об образовании.
type EducationEntry struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Details     string `json:"details"`
}

// CustomSection описывает кастомную секцию с буллетами.
type CustomSection struct {
	Title        string   `json:"title"`
	BulletSymbol string   `json:"bulletSymbol"`
	Items        []string `json:"items"`
}

// Photo описывает загруженное фото в base64.
type Photo struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

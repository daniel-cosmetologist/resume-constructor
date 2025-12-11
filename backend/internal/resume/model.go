package resume

// Link описывает одну ссылку в контактах (GitHub, LinkedIn и т.п.).
type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Contacts содержит контактную информацию.
type Contacts struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Links    []Link `json:"links"`
}

// ExperienceItem описывает один блок опыта работы.
type ExperienceItem struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"` // формат YYYY-MM
	EndDate     string   `json:"endDate"`   // формат YYYY-MM или пусто, если по настоящее время
	Description string   `json:"description"`
	Bullets     []string `json:"bullets"`
}

// EducationItem описывает одну запись об обучении.
type EducationItem struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"` // формат YYYY-MM
	EndDate     string `json:"endDate"`   // формат YYYY-MM
	Details     string `json:"details"`
}

// CustomSection — кастомный раздел (например, Homelab).
type CustomSection struct {
	Title        string   `json:"title"`
	BulletSymbol string   `json:"bulletSymbol"`
	Items        []string `json:"items"`
}

// Photo — опциональное фото в base64.
type Photo struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// Resume — основная доменная модель резюме.
type Resume struct {
	FullName       string           `json:"fullName"`
	Position       string           `json:"position"`
	Summary        string           `json:"summary"`
	Contacts       Contacts         `json:"contacts"`
	Skills         []string         `json:"skills"`
	Experience     []ExperienceItem `json:"experience"`
	Education      []EducationItem  `json:"education"`
	CustomSections []CustomSection  `json:"customSections"`
	Photo          *Photo           `json:"photo"` // опционально
}

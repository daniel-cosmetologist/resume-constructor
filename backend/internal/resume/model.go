package resume

// Resume описывает агрегированное резюме в доменном слое.
type Resume struct {
	FullName       string
	Position       string
	Summary        string
	Contacts       Contacts
	Skills         []string
	Experience     []ExperienceEntry
	Education      []EducationEntry
	CustomSections []CustomSection
	Photo          *Photo
}

// Contacts описывает контактные данные пользователя.
type Contacts struct {
	Email    string
	Phone    string
	Location string
	Links    []Link
}

// Link описывает внешнюю ссылку (GitHub, LinkedIn и т.п.).
type Link struct {
	Label string
	URL   string
}

// ExperienceEntry описывает запись об опыте работы.
type ExperienceEntry struct {
	Company     string
	Position    string
	Location    string
	StartDate   string
	EndDate     string
	Description string
	Bullets     []string
}

// EducationEntry описывает запись об образовании.
type EducationEntry struct {
	Institution string
	Degree      string
	Location    string
	StartDate   string
	EndDate     string
	Details     string
}

// CustomSection описывает кастомную секцию с буллетами.
type CustomSection struct {
	Title        string
	BulletSymbol string
	Items        []string
}

// Photo описывает загруженное фото в base64.
type Photo struct {
	MimeType string
	Data     string
}

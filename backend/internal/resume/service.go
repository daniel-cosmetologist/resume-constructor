package resume

import (
	"context"
	"errors"
	"log"

	httptransport "resume_backend/internal/http"
	"resume_backend/internal/latexclient"
)

// Service реализует бизнес-логику генерации резюме через latex-service.
type Service struct {
	client *latexclient.Client
	logger *log.Logger
}

// NewService создаёт новый экземпляр сервиса резюме.
func NewService(client *latexclient.Client, logger *log.Logger) *Service {
	if logger == nil {
		logger = log.Default()
	}
	return &Service{
		client: client,
		logger: logger,
	}
}

// GeneratePDF реализует интерфейс ResumeService из пакета internal/http.
// 1) Проецирует DTO уровня HTTP в доменную модель,
// 2) Валидирует,
// 3) Делегирует генерацию PDF latex-сервису.
func (s *Service) GeneratePDF(ctx context.Context, req httptransport.ResumeRequest) ([]byte, error) {
	if s.client == nil {
		return nil, errors.New("latex client is not configured")
	}

	// Проекция HTTP-модели в доменную.
	dom := Resume{
		FullName: req.FullName,
		Position: req.Position,
		Summary:  req.Summary,
		Contacts: Contacts{
			Email:    req.Contacts.Email,
			Phone:    req.Contacts.Phone,
			Location: req.Contacts.Location,
			Links: func() []Link {
				links := make([]Link, len(req.Contacts.Links))
				for i, l := range req.Contacts.Links {
					links[i] = Link{
						Label: l.Label,
						URL:   l.URL,
					}
				}
				return links
			}(),
		},
		Skills: func() []string {
			skills := make([]string, len(req.Skills))
			copy(skills, req.Skills)
			return skills
		}(),
		Experience: func() []ExperienceEntry {
			items := make([]ExperienceEntry, len(req.Experience))
			for i, e := range req.Experience {
				bullets := make([]string, len(e.Bullets))
				copy(bullets, e.Bullets)
				items[i] = ExperienceEntry{
					Company:     e.Company,
					Position:    e.Position,
					Location:    e.Location,
					StartDate:   e.StartDate,
					EndDate:     e.EndDate,
					Description: e.Description,
					Bullets:     bullets,
				}
			}
			return items
		}(),
		Education: func() []EducationEntry {
			items := make([]EducationEntry, len(req.Education))
			for i, e := range req.Education {
				items[i] = EducationEntry{
					Institution: e.Institution,
					Degree:      e.Degree,
					Location:    e.Location,
					StartDate:   e.StartDate,
					EndDate:     e.EndDate,
					Details:     e.Details,
				}
			}
			return items
		}(),
		CustomSections: func() []CustomSection {
			items := make([]CustomSection, len(req.CustomSections))
			for i, cs := range req.CustomSections {
				sectionItems := make([]string, len(cs.Items))
				copy(sectionItems, cs.Items)
				items[i] = CustomSection{
					Title:        cs.Title,
					BulletSymbol: cs.BulletSymbol,
					Items:        sectionItems,
				}
			}
			return items
		}(),
	}

	if req.Photo != nil {
		dom.Photo = &Photo{
			MimeType: req.Photo.MimeType,
			Data:     req.Photo.Data,
		}
	}

	// Валидация доменной модели.
	if err := Validate(dom); err != nil {
		// Ошибка реализует метод Fields() и будет корректно интерпретирована слоем HTTP.
		return nil, err
	}

	// Делегируем фактическую генерацию PDF latex-сервису.
	pdf, err := s.client.RenderPDF(ctx, req)
	if err != nil {
		s.logger.Printf("latex-service RenderPDF failed: %v", err)
		return nil, err
	}

	return pdf, nil
}

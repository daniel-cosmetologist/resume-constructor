package latex

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"latex_service/internal/model"
)

// Renderer отвечает за генерацию LaTeX и PDF.
type Renderer struct {
	templatePath string
	logger       *log.Logger
}

// NewRenderer создаёт новый Renderer.
func NewRenderer(templatePath string, logger *log.Logger) *Renderer {
	if logger == nil {
		logger = log.Default()
	}
	return &Renderer{
		templatePath: templatePath,
		logger:       logger,
	}
}

// Render генерирует PDF по данным резюме.
func (r *Renderer) Render(ctx context.Context, resume model.Resume) ([]byte, error) {
	templateBytes, err := os.ReadFile(r.templatePath)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	latexSource := string(templateBytes)

	placeholders := map[string]string{
		"FullName":       escapeLatex(resume.FullName),
		"Position":       escapeLatex(resume.Position),
		"Summary":        buildSummary(resume.Summary),
		"Contacts":       buildContacts(resume.Contacts),
		"Skills":         buildSkills(resume.Skills),
		"Experience":     buildExperience(resume.Experience),
		"Education":      buildEducation(resume.Education),
		"CustomSections": buildCustomSections(resume.CustomSections),
	}

	var photoBytes []byte
	var photoExt string

	if resume.Photo != nil && strings.TrimSpace(resume.Photo.Data) != "" {
		pb, ext, err := processPhoto(resume.Photo.Data, resume.Photo.MimeType)
		if err != nil {
			r.logger.Printf("photo processing failed: %v", err)
			placeholders["Photo"] = ""
		} else {
			photoBytes = pb
			photoExt = ext
			placeholders["Photo"] = `\includegraphics[width=3cm,height=4cm,keepaspectratio]{photo.` + photoExt + `}`
		}
	} else {
		placeholders["Photo"] = ""
	}

	for key, val := range placeholders {
		latexSource = strings.ReplaceAll(latexSource, "{{"+key+"}}", val)
	}

	workDir, err := os.MkdirTemp("", "resume-latex-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	texPath := filepath.Join(workDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(latexSource), 0o644); err != nil {
		return nil, fmt.Errorf("write tex file: %w", err)
	}

	if len(photoBytes) > 0 && photoExt != "" {
		photoPath := filepath.Join(workDir, "photo."+photoExt)
		if err := os.WriteFile(photoPath, photoBytes, 0o644); err != nil {
			r.logger.Printf("failed to write photo file: %v", err)
		}
	}

	cmd := exec.CommandContext(ctx, "latexmk", "-pdf", "-interaction=nonstopmode", "resume.tex")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("latexmk failed: %w", err)
	}

	pdfPath := filepath.Join(workDir, "resume.pdf")
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}

	if len(pdfBytes) == 0 {
		return nil, fmt.Errorf("generated PDF is empty")
	}

	return pdfBytes, nil
}

func buildSummary(summary string) string {
	if strings.TrimSpace(summary) == "" {
		return ""
	}
	return escapeLatex(summary)
}

func buildContacts(c model.Contacts) string {
	var parts []string

	if strings.TrimSpace(c.Email) != "" {
		parts = append(parts, `\href{mailto:`+c.Email+`}{`+escapeLatex(c.Email)+`}`)
	}
	if strings.TrimSpace(c.Phone) != "" {
		parts = append(parts, escapeLatex(c.Phone))
	}
	if strings.TrimSpace(c.Location) != "" {
		parts = append(parts, escapeLatex(c.Location))
	}

	for _, l := range c.Links {
		if strings.TrimSpace(l.URL) == "" {
			continue
		}
		label := l.Label
		if strings.TrimSpace(label) == "" {
			label = l.URL
		}
		parts = append(parts, `\href{`+l.URL+`}{`+escapeLatex(label)+`}`)
	}

	return strings.Join(parts, ` \textbullet{} `)
}

func buildSkills(skills []string) string {
	if len(skills) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\\begin{itemize}[leftmargin=*]\n")
	for _, s := range skills {
		if strings.TrimSpace(s) == "" {
			continue
		}
		b.WriteString("\\item " + escapeLatex(s) + "\n")
	}
	b.WriteString("\\end{itemize}\n")
	return b.String()
}

func buildExperience(exps []model.ExperienceEntry) string {
	if len(exps) == 0 {
		return ""
	}

	var b strings.Builder
	for _, e := range exps {
		if strings.TrimSpace(e.Company) == "" && strings.TrimSpace(e.Position) == "" {
			continue
		}

		title := strings.TrimSpace(e.Position)
		company := strings.TrimSpace(e.Company)
		location := strings.TrimSpace(e.Location)
		dates := strings.TrimSpace(e.StartDate)
		if strings.TrimSpace(e.EndDate) != "" {
			dates = dates + " -- " + strings.TrimSpace(e.EndDate)
		}

		b.WriteString("\\textbf{" + escapeLatex(title) + "}")
		if company != "" {
			b.WriteString(" at " + escapeLatex(company))
		}
		if dates != "" {
			b.WriteString("\\hfill " + escapeLatex(dates))
		}
		b.WriteString("\\\\\n")

		if location != "" {
			b.WriteString(escapeLatex(location) + "\\\\\n")
		}

		if strings.TrimSpace(e.Description) != "" {
			b.WriteString(escapeLatex(e.Description) + "\\\\[0.2cm]\n")
		}

		if len(e.Bullets) > 0 {
			b.WriteString("\\begin{itemize}[leftmargin=*]\n")
			for _, item := range e.Bullets {
				if strings.TrimSpace(item) == "" {
					continue
				}
				b.WriteString("\\item " + escapeLatex(item) + "\n")
			}
			b.WriteString("\\end{itemize}\n")
		}

		b.WriteString("\\vspace{0.4cm}\n")
	}

	return b.String()
}

func buildEducation(eds []model.EducationEntry) string {
	if len(eds) == 0 {
		return ""
	}

	var b strings.Builder
	for _, e := range eds {
		if strings.TrimSpace(e.Institution) == "" && strings.TrimSpace(e.Degree) == "" {
			continue
		}

		institution := strings.TrimSpace(e.Institution)
		degree := strings.TrimSpace(e.Degree)
		location := strings.TrimSpace(e.Location)
		dates := strings.TrimSpace(e.StartDate)
		if strings.TrimSpace(e.EndDate) != "" {
			dates = dates + " -- " + strings.TrimSpace(e.EndDate)
		}

		if institution != "" {
			b.WriteString("\\textbf{" + escapeLatex(institution) + "}")
		}
		if degree != "" {
			b.WriteString(" -- " + escapeLatex(degree))
		}
		if dates != "" {
			b.WriteString("\\hfill " + escapeLatex(dates))
		}
		b.WriteString("\\\\\n")

		if location != "" {
			b.WriteString(escapeLatex(location) + "\\\\\n")
		}

		if strings.TrimSpace(e.Details) != "" {
			b.WriteString(escapeLatex(e.Details) + "\\\\[0.2cm]\n")
		}

		b.WriteString("\\vspace{0.4cm}\n")
	}

	return b.String()
}

func buildCustomSections(sections []model.CustomSection) string {
	if len(sections) == 0 {
		return ""
	}

	var b strings.Builder
	for _, cs := range sections {
		if strings.TrimSpace(cs.Title) == "" {
			continue
		}
		b.WriteString("\\section*{" + escapeLatex(cs.Title) + "}\n")

		if len(cs.Items) == 0 {
			continue
		}

		bullet := strings.TrimSpace(cs.BulletSymbol)
		if bullet == "" {
			bullet = "•"
		}

		b.WriteString("\\begin{itemize}[leftmargin=*]\n")
		for _, item := range cs.Items {
			if strings.TrimSpace(item) == "" {
				continue
			}
			// В LaTeX маркер bullet задаётся окружением, символ просто часть текста.
			b.WriteString("\\item " + escapeLatex(item) + "\n")
		}
		b.WriteString("\\end{itemize}\n")
	}

	return b.String()
}

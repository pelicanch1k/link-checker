package pdf

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	
	"github.com/pelicanch1k/link-checker/internal/domain"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) GenerateReport(tasks []*domain.Task) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Заголовок
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(12)
	
	// Информация о задачах
	pdf.SetFont("Arial", "", 12)
	
	for _, task := range tasks {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Task ID: %d", task.ID))
		pdf.Ln(8)
		
		pdf.SetFont("Arial", "", 11)
		for _, link := range task.Links {
			status := "available"
			if link.Status == domain.StatusNotAvailable {
				status = "not available"
			}
			
			text := fmt.Sprintf("  %s - %s", link.URL, status)
			pdf.Cell(40, 7, text)
			pdf.Ln(6)
		}
		
		pdf.Ln(4)
	}
	
	// Сохраняем в buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}
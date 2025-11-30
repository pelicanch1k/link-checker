package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/pelicanch1k/link-checker/internal/checker"
)

type HTTPController struct {
	useCase *checker.LinkCheckerUseCase
}

func NewHTTPController(uc *checker.LinkCheckerUseCase) *HTTPController {
	return &HTTPController{useCase: uc}
}

func (ctrl *HTTPController) CheckLinks(c *fiber.Ctx) error {
	var req CheckLinksRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.Links) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Links array is empty",
		})
	}

	input := checker.CheckLinksInput{
		URLs: req.Links,
	}

	output, err := ctrl.useCase.CheckLinks(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := CheckLinksResponse{
		Links:    output.Links,
		LinksNum: output.TaskID,
	}
	return c.JSON(response)
}

func (ctrl *HTTPController) CheckLinksByIDs(c *fiber.Ctx) error {
	var req CheckLinksByIDsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.LinksList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Links list is empty",
		})
	}

	input := checker.CheckLinksByIDsInput{
		LinksList: req.LinksList,
	}

	output, err := ctrl.useCase.CheckLinksByIDs(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=link-report.pdf")
	return c.Send(output.PDFData)
}

func SetupRoutes(app *fiber.App, controller *HTTPController) {
	// Initialize default config
	app.Use(logger.New())

	app.Post("/check", controller.CheckLinks)
	app.Post("/info", controller.CheckLinksByIDs)
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
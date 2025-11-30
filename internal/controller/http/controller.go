package http

import (
	"github.com/pelicanch1k/link-checker/internal/checker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		LinksNum: output.LinksNum,
	}
	return c.JSON(response)
}

func SetupRoutes(app *fiber.App, controller *HTTPController) {
	// Initialize default config
	app.Use(logger.New())

	app.Post("/check", controller.CheckLinks)
	
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
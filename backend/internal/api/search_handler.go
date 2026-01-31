package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	domain "github.com/new-carmen/backend/internal/domain"
	"github.com/new-carmen/backend/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		searchService: services.NewSearchService(),
	}
}

func (h *SearchHandler) AnalyzeClarity(c *fiber.Ctx) error {
	var req domain.ClarificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.searchService.AnalyzeQuestionClarity(req.Question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	results, err := h.searchService.Search(query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}

func (h *SearchHandler) SearchWithContext(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	if limit > 10 {
		limit = 10
	}

	answer, err := h.searchService.SearchWithContext(query, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"answer": answer,
		"query":  query,
	})
}

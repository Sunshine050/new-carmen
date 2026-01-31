package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/new-carmen/backend/internal/api"
	"github.com/new-carmen/backend/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Auth routes (public)
	authHandler := api.NewAuthHandler()
	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	apiGroup := app.Group("/api", middleware.AuthMiddleware())

	// Document routes
	docHandler := api.NewDocumentHandler()
	documents := apiGroup.Group("/documents")
	documents.Post("/", docHandler.CreateDocument)
	documents.Get("/:id", docHandler.GetDocument)
	documents.Put("/:id", docHandler.UpdateDocument)
	documents.Delete("/:id", docHandler.DeleteDocument)
	documents.Post("/:id/versions", docHandler.AddVersion)

	// Search routes
	searchHandler := api.NewSearchHandler()
	search := apiGroup.Group("/search")
	search.Post("/analyze", searchHandler.AnalyzeClarity)
	search.Get("/", searchHandler.Search)
	search.Get("/context", searchHandler.SearchWithContext)

	// Indexing routes (admin only)
	indexingHandler := api.NewIndexingHandler()
	indexing := apiGroup.Group("/indexing", middleware.RoleMiddleware("admin", "superadmin"))
	indexing.Post("/load-github", indexingHandler.LoadFromGitHub)
	indexing.Post("/index-all", indexingHandler.IndexAll)

	// Admin routes
	_ = apiGroup.Group("/admin", middleware.RoleMiddleware("admin", "superadmin"))
	// Add admin routes here
}

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type controller interface {
	authController
	documentsController
}

type authController interface {
	Register(ctx *fiber.Ctx) error
	Auth(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type documentsController interface {
	UploadDocument(ctx *fiber.Ctx) error
	GetDocument(ctx *fiber.Ctx) error
	GetDocumentHead(ctx *fiber.Ctx) error
	GetDocuments(ctx *fiber.Ctx) error
	GetDocumentsHead(ctx *fiber.Ctx) error
	DeleteDocument(ctx *fiber.Ctx) error
}

func SetRoutes(app *fiber.App, controller controller) {
	app.Use(cors.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	api := app.Group("api")

	api.Post("register", controller.Register)
	api.Post("auth", controller.Auth)
	api.Delete("auth/:token", controller.Logout)

	docs := api.Group("docs")
	docs.Post("", controller.UploadDocument)
	docs.Get("/:id", controller.GetDocument)
	docs.Head("/:id", controller.GetDocumentHead)
	docs.Get("", controller.GetDocument)
	docs.Head("", controller.GetDocumentsHead)
	docs.Delete("/:id", controller.DeleteDocument)
}

package api

import (
	"github.com/gofiber/fiber/v2"
)

type controller interface {
	authController
	fileController
}

type authController interface {
	Register(ctx *fiber.Ctx) error
	Auth(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type fileController interface {
	UploadFile(ctx *fiber.Ctx) error
	GetFile(ctx *fiber.Ctx) error
	GetFileHead(ctx *fiber.Ctx) error
	GetFiles(ctx *fiber.Ctx) error
	GetFilesHead(ctx *fiber.Ctx) error
	DeleteFile(ctx *fiber.Ctx) error
}

func SetRoutes(app *fiber.App, controller controller) {
	api := app.Group("api")

	api.Post("register", controller.Register)
	api.Post("auth", controller.Auth)
	api.Delete("auth/:token", controller.Logout)

	docs := api.Group("docs")
	docs.Post("", controller.UploadFile)
	docs.Get("/:id", controller.GetFile)
	docs.Head("/:id", controller.GetFileHead)
	docs.Get("", controller.GetFiles)
	docs.Head("", controller.GetFilesHead)
	docs.Delete("/:id", controller.DeleteFile)
}

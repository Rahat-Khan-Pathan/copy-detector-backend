package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func ExamsRoute(route fiber.Router) {
	route.Post("/new", Controllers.ExamsNew)
	route.Post("/get_all", Controllers.ExamsGetAll)
	route.Post("/get_all_populated", Controllers.ExamsGetAllPopulated)
	route.Post("/modify/:id", Controllers.ExamsModify)
	route.Post("/delete/:id", Controllers.ExamsDelete)
}

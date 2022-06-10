package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func CoursesRoute(route fiber.Router) {
	route.Post("/new", Controllers.CoursesNew)
	route.Post("/get_all", Controllers.CoursesGetAll)
	route.Post("/get_all_populated", Controllers.CoursesGetAllPopulated)
	route.Post("/modify/:id", Controllers.CoursesModify)
	route.Post("/delete/:id", Controllers.CoursesDelete)
}

package Routes

import (
	"example.com/example/Controllers"
	"github.com/gofiber/fiber/v2"
)

func CodesRoute(route fiber.Router) {
	route.Post("/submit", Controllers.CodesSubmit)
	route.Post("/get_results", Controllers.CodesGetResults)

}

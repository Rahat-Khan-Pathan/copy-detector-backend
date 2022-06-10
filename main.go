package main

import (
	"fmt"
	"os"

	"example.com/example/DBManager"
	"example.com/example/Routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.UsersRoute(app.Group("/users"))
	Routes.CoursesRoute(app.Group("/courses"))
	Routes.ExamsRoute(app.Group("/exams"))
	Routes.CodesRoute(app.Group("/codes"))
}

func main() {
	fmt.Println("Hello Copy Detector")

	fmt.Print("Initializing DataBase Connections ... ")
	initState := DBManager.InitCollections()
	if initState {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")
	app := fiber.New()
	app.Use(cors.New())
	// app.Use(Middlewares.Auth)
	app.Use(pprof.New())

	SetupRoutes(app)
	app.Static("/Resources", "./Resources")
	fmt.Println("[OK]")
	fmt.Println(os.Getenv("PORT"))
	app.Listen(":" + os.Getenv("PORT"))
}

package main

import (
	"backend/database"
	"backend/routes"
)

func main() {
	database.Connect()
	r := routes.SetupRoutes()
	r.Run(":8080")
}

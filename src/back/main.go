package main

import (
	"backend/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
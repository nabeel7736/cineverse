package main

import (
	"cineverse/config"
	"cineverse/routes"
	"fmt"
	"log"
	"os"
)

func main() {
	config.Loadenv()
	config.ConnectDatabase()

	config.MigrateAll()

	r := routes.SetupRouter(config.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	add := fmt.Sprintf(":%s", port)
	log.Printf("server running at http://localhost%s", add)

	r.Run(add)

}

package main

import (
	"os"

	"adminPanel.com/db"
	"adminPanel.com/routes"
	"adminPanel.com/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.InitDB()

	server := gin.Default()
	routes.ChechForRoutes(server)
	server.Run(os.Getenv("PORT"))

	utils.StartMailGarbageCollector()
	utils.StartUnverifiedUsersGarbagerCollector()

}

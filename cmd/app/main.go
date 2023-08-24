package main

import (
	"fmt"
	"goproject/internal/infrastructure/database"
	httproute "goproject/internal/infrastructure/http"
	"goproject/internal/utils"
	"os"

	"goproject/docs"
)

//	@title			Go-Blog
//	@version		1.0
//	@description	A simple medium-like blog API, writen in Go

// @securityDefinitions.apikey	BearerToken
// @in							header
// @name						Authorization
// @scheme						bearer
// @bearerFormat				JWT
// @description				JWT Bearer Token. Need to Login to get the token. Usage: "Bearer <your-token-here>"
func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err)
	// }

	host := os.Getenv("HOST")
	port := os.Getenv("HTTP_PORT")

	db, err := database.NewPostgresDB()
	if err != nil {
		panic(err)
	}

	logger := utils.NewLogger()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := httproute.NewRoute(db.DB, logger)
	r.Run(fmt.Sprintf(":%s", port))
}

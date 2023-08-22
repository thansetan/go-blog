package main

import (
	"fmt"
	"goproject/internal/infrastructure/database"
	httproute "goproject/internal/infrastructure/http"
	"os"

	"goproject/docs"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("HTTP_PORT")

	db, err := database.NewPostgresDB()
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.Title = "Sanbercode Final Project"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := httproute.NewRoute(db.DB)
	r.Run(fmt.Sprintf(":%s", port))

}

package main

import (
	_ "arvore-genealogica-golang/cmd/docs"
	"arvore-genealogica-golang/config"
	"arvore-genealogica-golang/internal/adapters"
	"arvore-genealogica-golang/internal/app/repository"
	"arvore-genealogica-golang/internal/app/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Arvore Genealógica NEO4J
// @version 1.0
// @description Aplicação desenvolvida em GoLang para busca de arvores utilizando Neo4j.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	driver := config.InitializeDatabase()
	repo := repository.NewPersonRepository(driver)
	srv := service.NewPersonService(repo)
	adapter := adapters.NewFamilyTreeAdapter(srv)
	router := gin.Default()

	userApi := router.Group("/api/v1")
	{
		userApi.GET("/person/:id", adapter.GetPersonById)
		userApi.POST("/person", adapter.CreatePerson)
		userApi.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	}

	router.Run(":8080")
}

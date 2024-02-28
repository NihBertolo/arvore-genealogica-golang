package adapters

import (
	"fmt"

	_ "arvore-genealogica-golang/cmd/docs"
	"arvore-genealogica-golang/internal/adapters/response"
	models "arvore-genealogica-golang/internal/app/domain/models"
	"arvore-genealogica-golang/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FamilyTreeAdapter struct {
	service *service.PersonService
}

func NewFamilyTreeAdapter(service *service.PersonService) *FamilyTreeAdapter {
	return &FamilyTreeAdapter{service}
}

// GetPersonByID godoc
// @Summary Get person By ID
// @Description Get user details by providing user ID
// @Tags person
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} models.Person
// @Router /person/{id} [get]
func (a *FamilyTreeAdapter) GetPersonById(c *gin.Context) {
	id := c.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("Erro ao fazer o parse do UUID:", err)
		response.FailResponse(c, 400, err.Error())
		return
	}
	person, err := a.service.FindByID(c, parsedUUID)
	if err != nil {
		fmt.Printf("Erro ao procurar o usuário com UUID %s, error: %s", parsedUUID, err)
		response.FailResponse(c, 400, err.Error())
		return
	}

	response.SuccessResponse(c, person)
}

// Criar pessoa
// @Summary                    Create Person
// @Description                Create a new Person
// @Tags                       person
// @Accept                     json
// @Produce                    json
// @Param                      person body     models.Person true "Person Data"
// @Success                    200  {object} models.Person
// @Router                     /person [post]
func (a *FamilyTreeAdapter) CreatePerson(c *gin.Context) {
	var newPerson models.Person

	if err := c.BindJSON(&newPerson); err != nil {
		response.FailResponse(c, 400, err.Error())
		return
	}

	a.service.Create(c, &newPerson)

	response.SuccessResponse(c, newPerson)
}

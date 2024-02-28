package models

import "github.com/google/uuid"

type Person struct {
	ID       uuid.UUID `json:"id" bson:"id" swaggerignore:"true"`
	Name     string    `json:"name" bson:"name" binding:"required"`
	Parents  []Person  `json:"parents" bson:"parents"`
	Children []Person  `json:"children" bson:"children"`
}

func NewPerson(name string, parents []Person, children []Person) *Person {
	return &Person{
		ID:       uuid.New(),
		Name:     name,
		Parents:  parents,
		Children: children,
	}
}

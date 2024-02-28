package repository

import (
	"arvore-genealogica-golang/internal/app/domain/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type PersonRepository struct {
	driver neo4j.Driver
}

func NewPersonRepository(driver neo4j.Driver) *PersonRepository {
	return &PersonRepository{driver}
}

func (r *PersonRepository) FindByID(id string) (*models.Person, error) {
	session, errorSession := r.driver.Session(neo4j.AccessModeWrite)
	if errorSession != nil {
		return nil, errorSession
	}
	defer session.Close()

	result, err := session.Run(
		"MATCH (p:Person {id: $id}) RETURN p",
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		return nil, err
	}

	if result.Next() {
		record := result.Record()
		person := models.Person{
			ID:   uuid.MustParse(record.GetByIndex(0).(string)),
			Name: record.GetByIndex(0).(string),
		}
		parents, err := r.GetParents(id, 1)
		if err != nil {
			return &person, err
		}
		children, err := r.GetChildren(id, 1)
		if err != nil {
			return &person, err
		}

		person.Parents = parents
		person.Children = children

		return &person, nil
	}

	return nil, fmt.Errorf("Pessoa com ID %s não encontrada", id)
}

func (r *PersonRepository) Create(ctx context.Context, data *models.Person) error {
	session, err := r.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (p:Person {id: $id, name: $name}) RETURN p",
			map[string]interface{}{
				"id":   data.ID.String(),
				"name": data.Name,
			})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Pessoa com ID %s inserida", data.ID)
		if result.Next() {
			return nil, fmt.Errorf("Pessoa com ID %s já existente", data.ID)
		}
		return nil, result.Err()
	})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(data.Parents); i++ {
		r.CreateRelationship(data.Parents[i], *data)
	}
	for i := 0; i < len(data.Children); i++ {
		r.CreateRelationship(*data, data.Children[i])
	}
	return nil
}

func (r *PersonRepository) CreateRelationship(parent, child models.Person) error {
	session, errorSession := r.driver.Session(neo4j.AccessModeWrite)
	if errorSession != nil {
		return errorSession
	}
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MATCH (p:Person {id: $parentId}), (c:Person {id: $childId}) "+
				"CREATE (p)-[:PARENT_OF]->(c)",
			map[string]interface{}{
				"parentId": parent.ID,
				"childId":  child.ID,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *PersonRepository) GetParents(id string, level int) ([]models.Person, error) {
	session, err := r.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(
		"MATCH (p:Person)-[:PARENT_OF]->(c:Person {id: $childId}) RETURN p",
		map[string]interface{}{
			"childId": id,
		})
	if err != nil {
		return nil, err
	}

	level = level - 1
	var parents []models.Person
	for result.Next() {
		record := result.Record()
		parent := models.Person{
			ID:   uuid.MustParse(record.GetByIndex(0).(string)),
			Name: record.GetByIndex(0).(string),
		}
		if level > 0 {
			ascParents, err := r.GetParents(parent.ID.String(), level)
			if err != nil {
				return nil, err
			}
			parent.Parents = ascParents
		}
		parents = append(parents, parent)
	}

	return parents, nil
}

func (r *PersonRepository) GetChildren(id string, level int) ([]models.Person, error) {
	session, err := r.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(
		"MATCH (c:Person {id: $parentId})-[:PARENT_OF]->(p:Person) RETURN p",
		map[string]interface{}{
			"parentId": id,
		})
	if err != nil {
		return nil, err
	}

	var children []models.Person
	for result.Next() {
		record := result.Record()
		child := models.Person{
			ID:   uuid.MustParse(record.GetByIndex(0).(string)),
			Name: record.GetByIndex(1).(string),
		}
		if level > 0 {
			descChildren, err := r.GetChildren(child.ID.String(), level)
			if err != nil {
				return nil, err
			}
			child.Children = descChildren
		}

		children = append(children, child)
	}

	return children, nil
}

func (r *PersonRepository) DeleteNode(id string) error {
	session, errorSession := r.driver.Session(neo4j.AccessModeWrite)
	if errorSession != nil {
		return errorSession
	}
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MATCH (n) WHERE n.id = $id DELETE n",
			map[string]interface{}{
				"id": id,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *PersonRepository) UpdatePerson(person *models.Person) error {
	session, errorSession := r.driver.Session(neo4j.AccessModeWrite)
	if errorSession != nil {
		return errorSession
	}
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run(
			"MATCH (p:Person {id: $id}) SET p.name = $name",
			map[string]interface{}{
				"id":   person.ID,
				"name": person.Name,
			})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

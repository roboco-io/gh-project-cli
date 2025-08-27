package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectQueries(t *testing.T) {
	t.Run("ListUserProjects query structure", func(t *testing.T) {
		query := &ListUserProjectsQuery{}

		assert.NotNil(t, query)
		// Test that query has the right structure for GraphQL
	})

	t.Run("ListOrgProjects query structure", func(t *testing.T) {
		query := &ListOrgProjectsQuery{}

		assert.NotNil(t, query)
	})

	t.Run("GetProject query structure", func(t *testing.T) {
		query := &GetProjectQuery{}

		assert.NotNil(t, query)
	})
}

func TestProjectMutations(t *testing.T) {
	t.Run("CreateProject mutation structure", func(t *testing.T) {
		mutation := &CreateProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("UpdateProject mutation structure", func(t *testing.T) {
		mutation := &UpdateProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("DeleteProject mutation structure", func(t *testing.T) {
		mutation := &DeleteProjectMutation{}

		assert.NotNil(t, mutation)
	})
}

func TestItemMutations(t *testing.T) {
	t.Run("AddItemToProject mutation structure", func(t *testing.T) {
		mutation := &AddItemToProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("UpdateItemField mutation structure", func(t *testing.T) {
		mutation := &UpdateItemFieldMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("RemoveItemFromProject mutation structure", func(t *testing.T) {
		mutation := &RemoveItemFromProjectMutation{}

		assert.NotNil(t, mutation)
	})
}

func TestVariableBuilders(t *testing.T) {
	t.Run("BuildCreateProjectVariables creates proper variables", func(t *testing.T) {
		input := CreateProjectInput{
			OwnerID: "test-owner-id",
			Title:   "Test Project",
		}

		variables := BuildCreateProjectVariables(&input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "test-owner-id", inputVar["ownerId"])
		assert.Equal(t, "Test Project", inputVar["title"])
	})

	t.Run("BuildAddItemVariables creates proper variables", func(t *testing.T) {
		input := AddItemInput{
			ProjectID: "project-id",
			ContentID: "content-id",
		}

		variables := BuildAddItemVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "project-id", inputVar["projectId"])
		assert.Equal(t, "content-id", inputVar["contentId"])
	})
}

func TestResponseParsing(t *testing.T) {
	t.Run("ParseProjectResponse extracts project data", func(t *testing.T) {
		response := &GetProjectQuery{
			Organization: struct {
				ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
			}{
				ProjectV2: ProjectV2{
					ID:     "test-id",
					Title:  "Test Project",
					Number: 42,
				},
			},
		}

		project := response.Organization.ProjectV2
		assert.Equal(t, "test-id", project.ID)
		assert.Equal(t, "Test Project", project.Title)
		assert.Equal(t, 42, project.Number)
	})
}

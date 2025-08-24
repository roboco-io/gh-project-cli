package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectV2ViewTypes(t *testing.T) {
	t.Run("ProjectV2ViewLayout constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2ViewLayout("TABLE_VIEW"), ProjectV2ViewLayoutTable)
		assert.Equal(t, ProjectV2ViewLayout("BOARD_VIEW"), ProjectV2ViewLayoutBoard)
		assert.Equal(t, ProjectV2ViewLayout("ROADMAP_VIEW"), ProjectV2ViewLayoutRoadmap)
	})

	t.Run("ProjectV2ViewSortDirection constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2ViewSortDirection("ASC"), ProjectV2ViewSortDirectionASC)
		assert.Equal(t, ProjectV2ViewSortDirection("DESC"), ProjectV2ViewSortDirectionDESC)
	})
}

func TestViewVariableBuilders(t *testing.T) {
	t.Run("BuildCreateViewVariables creates proper variables", func(t *testing.T) {
		input := CreateViewInput{
			ProjectID: "test-project-id",
			Name:      "Test View",
			Layout:    ProjectV2ViewLayoutTable,
		}

		variables := BuildCreateViewVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"projectId": "test-project-id",
				"name":      "Test View",
				"layout":    ProjectV2ViewLayoutTable,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewVariables creates proper variables", func(t *testing.T) {
		name := "Updated View"
		filter := "status:todo"
		input := UpdateViewInput{
			ViewID: "test-view-id",
			Name:   &name,
			Filter: &filter,
		}

		variables := BuildUpdateViewVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId": "test-view-id",
				"name":   "Updated View",
				"filter": "status:todo",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewVariables with minimal input", func(t *testing.T) {
		input := UpdateViewInput{
			ViewID: "test-view-id",
		}

		variables := BuildUpdateViewVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId": "test-view-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildDeleteViewVariables creates proper variables", func(t *testing.T) {
		input := DeleteViewInput{
			ViewID: "test-view-id",
		}

		variables := BuildDeleteViewVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId": "test-view-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildCopyViewVariables creates proper variables", func(t *testing.T) {
		input := CopyViewInput{
			ProjectID: "test-project-id",
			ViewID:    "test-view-id",
			Name:      "Copied View",
		}

		variables := BuildCopyViewVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"projectId": "test-project-id",
				"viewId":    "test-view-id",
				"name":      "Copied View",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewSortByVariables creates proper variables", func(t *testing.T) {
		sortByID := "test-field-id"
		input := UpdateViewSortByInput{
			ViewID:    "test-view-id",
			SortByID:  &sortByID,
			Direction: ProjectV2ViewSortDirectionASC,
		}

		variables := BuildUpdateViewSortByVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId":    "test-view-id",
				"sortById":  "test-field-id",
				"direction": ProjectV2ViewSortDirectionASC,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewSortByVariables without sortById", func(t *testing.T) {
		input := UpdateViewSortByInput{
			ViewID:    "test-view-id",
			Direction: ProjectV2ViewSortDirectionDESC,
		}

		variables := BuildUpdateViewSortByVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId":    "test-view-id",
				"direction": ProjectV2ViewSortDirectionDESC,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewGroupByVariables creates proper variables", func(t *testing.T) {
		groupByID := "test-field-id"
		input := UpdateViewGroupByInput{
			ViewID:    "test-view-id",
			GroupByID: &groupByID,
			Direction: ProjectV2ViewSortDirectionASC,
		}

		variables := BuildUpdateViewGroupByVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId":    "test-view-id",
				"groupById": "test-field-id",
				"direction": ProjectV2ViewSortDirectionASC,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateViewGroupByVariables without groupById", func(t *testing.T) {
		input := UpdateViewGroupByInput{
			ViewID:    "test-view-id",
			Direction: ProjectV2ViewSortDirectionDESC,
		}

		variables := BuildUpdateViewGroupByVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"viewId":    "test-view-id",
				"direction": ProjectV2ViewSortDirectionDESC,
			},
		}

		assert.Equal(t, expected, variables)
	})
}

func TestViewHelperFunctions(t *testing.T) {
	t.Run("ValidViewLayouts returns all valid layouts", func(t *testing.T) {
		layouts := ValidViewLayouts()
		expected := []string{
			string(ProjectV2ViewLayoutTable),
			string(ProjectV2ViewLayoutBoard),
			string(ProjectV2ViewLayoutRoadmap),
		}

		assert.Equal(t, expected, layouts)
		assert.Len(t, layouts, 3)
	})

	t.Run("ValidSortDirections returns all valid directions", func(t *testing.T) {
		directions := ValidSortDirections()
		expected := []string{
			string(ProjectV2ViewSortDirectionASC),
			string(ProjectV2ViewSortDirectionDESC),
		}

		assert.Equal(t, expected, directions)
		assert.Len(t, directions, 2)
	})

	t.Run("FormatViewLayout formats correctly", func(t *testing.T) {
		assert.Equal(t, "Table", FormatViewLayout(ProjectV2ViewLayoutTable))
		assert.Equal(t, "Board", FormatViewLayout(ProjectV2ViewLayoutBoard))
		assert.Equal(t, "Roadmap", FormatViewLayout(ProjectV2ViewLayoutRoadmap))
		assert.Equal(t, "UNKNOWN_VIEW", FormatViewLayout(ProjectV2ViewLayout("UNKNOWN_VIEW")))
	})

	t.Run("FormatSortDirection formats correctly", func(t *testing.T) {
		assert.Equal(t, "Ascending", FormatSortDirection(ProjectV2ViewSortDirectionASC))
		assert.Equal(t, "Descending", FormatSortDirection(ProjectV2ViewSortDirectionDESC))
		assert.Equal(t, "UNKNOWN", FormatSortDirection(ProjectV2ViewSortDirection("UNKNOWN")))
	})
}

func TestViewStructures(t *testing.T) {
	t.Run("ProjectV2View structure validation", func(t *testing.T) {
		view := ProjectV2View{
			ID:     "view-id",
			Name:   "Test View",
			Layout: ProjectV2ViewLayoutTable,
			Number: 1,
		}

		assert.Equal(t, "view-id", view.ID)
		assert.Equal(t, "Test View", view.Name)
		assert.Equal(t, ProjectV2ViewLayoutTable, view.Layout)
		assert.Equal(t, 1, view.Number)
	})

	t.Run("ProjectV2ViewGroupBy structure validation", func(t *testing.T) {
		groupBy := ProjectV2ViewGroupBy{
			Direction: ProjectV2ViewSortDirectionASC,
		}
		groupBy.Field.ID = "field-id"
		groupBy.Field.Name = "Status"

		assert.Equal(t, "field-id", groupBy.Field.ID)
		assert.Equal(t, "Status", groupBy.Field.Name)
		assert.Equal(t, ProjectV2ViewSortDirectionASC, groupBy.Direction)
	})

	t.Run("ProjectV2ViewSortBy structure validation", func(t *testing.T) {
		sortBy := ProjectV2ViewSortBy{
			Direction: ProjectV2ViewSortDirectionDESC,
		}
		sortBy.Field.ID = "field-id"
		sortBy.Field.Name = "Priority"

		assert.Equal(t, "field-id", sortBy.Field.ID)
		assert.Equal(t, "Priority", sortBy.Field.Name)
		assert.Equal(t, ProjectV2ViewSortDirectionDESC, sortBy.Direction)
	})

	t.Run("ProjectV2ViewColumn structure validation", func(t *testing.T) {
		column := ProjectV2ViewColumn{
			ID:       "column-id",
			Name:     "Column Name",
			Width:    200,
			IsHidden: false,
		}
		column.Field.ID = "field-id"
		column.Field.Name = "Field Name"

		assert.Equal(t, "column-id", column.ID)
		assert.Equal(t, "Column Name", column.Name)
		assert.Equal(t, 200, column.Width)
		assert.False(t, column.IsHidden)
		assert.Equal(t, "field-id", column.Field.ID)
		assert.Equal(t, "Field Name", column.Field.Name)
	})
}
package graphql

import "time"

// ProjectV2View represents a view in a GitHub Project v2
type ProjectV2View struct {
	ID           string                 `graphql:"id"`
	Name         string                 `graphql:"name"`
	Layout       ProjectV2ViewLayout    `graphql:"layout"`
	Number       int                    `graphql:"number"`
	CreatedAt    time.Time             `graphql:"createdAt"`
	UpdatedAt    time.Time             `graphql:"updatedAt"`
	DatabaseID   int                   `graphql:"databaseId"`
	Filter       *string               `graphql:"filter"`
	GroupBy      []ProjectV2ViewGroupBy `graphql:"groupBy"`
	SortBy       []ProjectV2ViewSortBy  `graphql:"sortBy"`
	VerticalGroupBy []ProjectV2ViewGroupBy `graphql:"verticalGroupBy"`
	
	// Group by and sort by fields for different view layouts
	Fields struct {
		Nodes []ProjectV2Field `graphql:"nodes"`
	} `graphql:"groupByFields(first: 20)"`
}

// ProjectV2ViewLayout represents the layout type of a view
type ProjectV2ViewLayout string

const (
	ProjectV2ViewLayoutTable    ProjectV2ViewLayout = "TABLE_VIEW"
	ProjectV2ViewLayoutBoard    ProjectV2ViewLayout = "BOARD_VIEW"
	ProjectV2ViewLayoutRoadmap  ProjectV2ViewLayout = "ROADMAP_VIEW"
)

// ProjectV2ViewGroupBy represents a group by configuration
type ProjectV2ViewGroupBy struct {
	Field struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"field"`
	Direction ProjectV2ViewSortDirection `graphql:"direction"`
}

// ProjectV2ViewSortBy represents a sort by configuration
type ProjectV2ViewSortBy struct {
	Field struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"field"`
	Direction ProjectV2ViewSortDirection `graphql:"direction"`
}

// ProjectV2ViewSortDirection represents sort direction
type ProjectV2ViewSortDirection string

const (
	ProjectV2ViewSortDirectionASC  ProjectV2ViewSortDirection = "ASC"
	ProjectV2ViewSortDirectionDESC ProjectV2ViewSortDirection = "DESC"
)

// ProjectV2ViewColumn represents a column in a view
type ProjectV2ViewColumn struct {
	ID       string `graphql:"id"`
	Name     string `graphql:"name"`
	Width    int    `graphql:"width"`
	IsHidden bool   `graphql:"isHidden"`
	
	Field struct {
		ID   string `graphql:"id"`
		Name string `graphql:"name"`
	} `graphql:"field"`
}

// Queries

// GetProjectViewsQuery gets all views for a project
type GetProjectViewsQuery struct {
	Node struct {
		ProjectV2 struct {
			Views struct {
				Nodes []ProjectV2View `graphql:"nodes"`
			} `graphql:"views(first: 20)"`
		} `graphql:"... on ProjectV2"`
	} `graphql:"node(id: $projectId)"`
}

// GetProjectViewQuery gets a specific view by ID
type GetProjectViewQuery struct {
	Node struct {
		ProjectV2View ProjectV2View `graphql:"... on ProjectV2View"`
	} `graphql:"node(id: $viewId)"`
}

// Mutations

// CreateProjectViewMutation creates a new view
type CreateProjectViewMutation struct {
	CreateProjectV2View struct {
		ProjectV2View ProjectV2View `graphql:"projectV2View"`
	} `graphql:"createProjectV2View(input: $input)"`
}

// UpdateProjectViewMutation updates an existing view
type UpdateProjectViewMutation struct {
	UpdateProjectV2View struct {
		ProjectV2View ProjectV2View `graphql:"projectV2View"`
	} `graphql:"updateProjectV2View(input: $input)"`
}

// DeleteProjectViewMutation deletes a view
type DeleteProjectViewMutation struct {
	DeleteProjectV2View struct {
		ProjectV2View ProjectV2View `graphql:"projectV2View"`
	} `graphql:"deleteProjectV2View(input: $input)"`
}

// CopyProjectViewMutation creates a copy of an existing view
type CopyProjectViewMutation struct {
	CopyProjectV2View struct {
		ProjectV2View ProjectV2View `graphql:"projectV2View"`
	} `graphql:"copyProjectV2View(input: $input)"`
}

// Input Types

// CreateViewInput represents input for creating a view
type CreateViewInput struct {
	ProjectID string              `json:"projectId"`
	Name      string              `json:"name"`
	Layout    ProjectV2ViewLayout `json:"layout"`
}

// UpdateViewInput represents input for updating a view
type UpdateViewInput struct {
	ViewID   string  `json:"viewId"`
	Name     *string `json:"name,omitempty"`
	Filter   *string `json:"filter,omitempty"`
}

// DeleteViewInput represents input for deleting a view
type DeleteViewInput struct {
	ViewID string `json:"viewId"`
}

// CopyViewInput represents input for copying a view
type CopyViewInput struct {
	ProjectID string `json:"projectId"`
	ViewID    string `json:"viewId"`
	Name      string `json:"name"`
}

// UpdateViewSortByInput represents input for updating view sort configuration
type UpdateViewSortByInput struct {
	ViewID    string                     `json:"viewId"`
	SortByID  *string                    `json:"sortById,omitempty"`
	Direction ProjectV2ViewSortDirection `json:"direction"`
}

// UpdateViewGroupByInput represents input for updating view group configuration
type UpdateViewGroupByInput struct {
	ViewID      string                     `json:"viewId"`
	GroupByID   *string                    `json:"groupById,omitempty"`
	Direction   ProjectV2ViewSortDirection `json:"direction"`
}

// Variable Builders

// BuildCreateViewVariables builds variables for view creation
func BuildCreateViewVariables(input CreateViewInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"name":      input.Name,
			"layout":    input.Layout,
		},
	}
}

// BuildUpdateViewVariables builds variables for view update
func BuildUpdateViewVariables(input UpdateViewInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"viewId": input.ViewID,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.Name != nil {
		inputMap["name"] = *input.Name
	}
	if input.Filter != nil {
		inputMap["filter"] = *input.Filter
	}

	return vars
}

// BuildDeleteViewVariables builds variables for view deletion
func BuildDeleteViewVariables(input DeleteViewInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"viewId": input.ViewID,
		},
	}
}

// BuildCopyViewVariables builds variables for view copying
func BuildCopyViewVariables(input CopyViewInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"viewId":    input.ViewID,
			"name":      input.Name,
		},
	}
}

// BuildUpdateViewSortByVariables builds variables for updating view sort configuration
func BuildUpdateViewSortByVariables(input UpdateViewSortByInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"viewId":    input.ViewID,
			"direction": input.Direction,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.SortByID != nil {
		inputMap["sortById"] = *input.SortByID
	}

	return vars
}

// BuildUpdateViewGroupByVariables builds variables for updating view group configuration
func BuildUpdateViewGroupByVariables(input UpdateViewGroupByInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"viewId":    input.ViewID,
			"direction": input.Direction,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.GroupByID != nil {
		inputMap["groupById"] = *input.GroupByID
	}

	return vars
}

// Helper Functions

// ValidViewLayouts returns all valid view layout types
func ValidViewLayouts() []string {
	return []string{
		string(ProjectV2ViewLayoutTable),
		string(ProjectV2ViewLayoutBoard),
		string(ProjectV2ViewLayoutRoadmap),
	}
}

// ValidSortDirections returns all valid sort directions
func ValidSortDirections() []string {
	return []string{
		string(ProjectV2ViewSortDirectionASC),
		string(ProjectV2ViewSortDirectionDESC),
	}
}

// FormatViewLayout formats view layout for display
func FormatViewLayout(layout ProjectV2ViewLayout) string {
	switch layout {
	case ProjectV2ViewLayoutTable:
		return "Table"
	case ProjectV2ViewLayoutBoard:
		return "Board"
	case ProjectV2ViewLayoutRoadmap:
		return "Roadmap"
	default:
		return string(layout)
	}
}

// FormatSortDirection formats sort direction for display
func FormatSortDirection(direction ProjectV2ViewSortDirection) string {
	switch direction {
	case ProjectV2ViewSortDirectionASC:
		return "Ascending"
	case ProjectV2ViewSortDirectionDESC:
		return "Descending"
	default:
		return string(direction)
	}
}
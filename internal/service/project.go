package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
)

// ProjectService handles project-related operations
type ProjectService struct {
	client *api.Client
}

// NewProjectService creates a new project service
func NewProjectService(client *api.Client) *ProjectService {
	return &ProjectService{
		client: client,
	}
}

// ProjectInfo represents simplified project information for display
type ProjectInfo struct {
	Description *string
	ID          string
	Title       string
	URL         string
	Owner       string
	Number      int
	ItemCount   int
	FieldCount  int
	Closed      bool
}

// ListUserProjectsOptions represents options for listing user projects
type ListUserProjectsOptions struct {
	After *string
	Login string
	First int
}

// ListOrgProjectsOptions represents options for listing organization projects
type ListOrgProjectsOptions struct {
	After *string
	Login string
	First int
}

// convertProjectNodes converts GraphQL project nodes to ProjectInfo slice
func convertProjectNodes(nodes []graphql.ProjectV2) []ProjectInfo {
	projects := make([]ProjectInfo, len(nodes))
	for i := range nodes {
		project := &nodes[i]
		projects[i] = ProjectInfo{
			ID:          project.ID,
			Number:      project.Number,
			Title:       project.Title,
			Description: project.Description,
			URL:         project.URL,
			Closed:      project.Closed,
			Owner:       project.Owner.Login,
			ItemCount:   len(project.Items.Nodes),
			FieldCount:  len(project.Fields.Nodes),
		}
	}
	return projects
}

// buildProjectVariables builds common GraphQL variables for project listing
func buildProjectVariables(login string, first int, after *string) map[string]interface{} {
	variables := map[string]interface{}{
		"login": login,
		"first": first,
	}
	if after != nil {
		variables["after"] = *after
	}
	return variables
}

// ListUserProjects lists projects for a user
func (s *ProjectService) ListUserProjects(ctx context.Context, opts ListUserProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := buildProjectVariables(opts.Login, opts.First, opts.After)

	var query graphql.ListUserProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list user projects: %w", err)
	}

	return convertProjectNodes(query.User.ProjectsV2.Nodes), nil
}

// ListOrgProjects lists projects for an organization
func (s *ProjectService) ListOrgProjects(ctx context.Context, opts ListOrgProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := buildProjectVariables(opts.Login, opts.First, opts.After)

	var query graphql.ListOrgProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization projects: %w", err)
	}

	return convertProjectNodes(query.Organization.ProjectsV2.Nodes), nil
}

// GetProject gets a specific project by number
func (s *ProjectService) GetProject(ctx context.Context, owner string, number int, isOrg bool) (*graphql.ProjectV2, error) {
	if isOrg {
		variables := map[string]interface{}{
			"orgLogin": owner,
			"number":   number,
		}

		var query graphql.GetProjectQuery
		err := s.client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization project: %w", err)
		}

		return &query.Organization.ProjectV2, nil
	}

	variables := map[string]interface{}{
		"userLogin": owner,
		"number":    number,
	}

	var query graphql.GetUserProjectQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get user project: %w", err)
	}

	return &query.User.ProjectV2, nil
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	OwnerID string
	Title   string
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, input CreateProjectInput) (*graphql.ProjectV2, error) {
	variables := graphql.BuildCreateProjectVariables(graphql.CreateProjectInput{
		OwnerID: input.OwnerID,
		Title:   input.Title,
	})

	var mutation graphql.CreateProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &mutation.CreateProjectV2.ProjectV2, nil
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Title     *string
	Closed    *bool
	ProjectID string
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, input UpdateProjectInput) (*graphql.ProjectV2, error) {
	variables := graphql.BuildUpdateProjectVariables(graphql.UpdateProjectInput{
		ProjectID: input.ProjectID,
		Title:     input.Title,
		Closed:    input.Closed,
	})

	var mutation graphql.UpdateProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &mutation.UpdateProjectV2.ProjectV2, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, projectID string) error {
	variables := graphql.BuildDeleteProjectVariables(graphql.DeleteProjectInput{
		ProjectID: projectID,
	})

	var mutation graphql.DeleteProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// AddItemInput represents input for adding an item to a project
type AddItemInput struct {
	ProjectID string
	ContentID string
}

// AddItem adds an item to a project
func (s *ProjectService) AddItem(ctx context.Context, input AddItemInput) (*graphql.ProjectV2Item, error) {
	variables := graphql.BuildAddItemVariables(graphql.AddItemInput{
		ProjectID: input.ProjectID,
		ContentID: input.ContentID,
	})

	var mutation graphql.AddItemToProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to project: %w", err)
	}

	return &mutation.AddProjectV2ItemByID.Item, nil
}

// UpdateItemFieldInput represents input for updating an item field
type UpdateItemFieldInput struct {
	Value     interface{}
	ProjectID string
	ItemID    string
	FieldID   string
}

// UpdateItemField updates a field value for an item
func (s *ProjectService) UpdateItemField(ctx context.Context, input UpdateItemFieldInput) (*graphql.ProjectV2Item, error) {
	variables := graphql.BuildUpdateItemFieldVariables(graphql.UpdateItemFieldInput{
		ProjectID: input.ProjectID,
		ItemID:    input.ItemID,
		FieldID:   input.FieldID,
		Value:     input.Value,
	})

	var mutation graphql.UpdateItemFieldMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update item field: %w", err)
	}

	return &mutation.UpdateProjectV2ItemFieldValue.ProjectV2Item, nil
}

// RemoveItemInput represents input for removing an item from a project
type RemoveItemInput struct {
	ProjectID string
	ItemID    string
}

// RemoveItem removes an item from a project
func (s *ProjectService) RemoveItem(ctx context.Context, input RemoveItemInput) error {
	variables := graphql.BuildRemoveItemVariables(graphql.RemoveItemInput{
		ProjectID: input.ProjectID,
		ItemID:    input.ItemID,
	})

	var mutation graphql.RemoveItemFromProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to remove item from project: %w", err)
	}

	return nil
}

// ParseProjectReference parses a project reference in the format "owner/number"
func ParseProjectReference(ref string) (owner string, number int, err error) {
	// Simple parsing - in practice, this might need more sophisticated handling
	var numStr string
	for i := len(ref) - 1; i >= 0; i-- {
		if ref[i] == '/' {
			owner = ref[:i]
			numStr = ref[i+1:]
			break
		}
	}

	if owner == "" || numStr == "" {
		return "", 0, fmt.Errorf("invalid project reference format: %s (expected owner/number)", ref)
	}

	number, err = strconv.Atoi(numStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid project number in reference: %s", numStr)
	}

	return owner, number, nil
}

// FormatProjectReference formats owner and number into a project reference
func FormatProjectReference(owner string, number int) string {
	return fmt.Sprintf("%s/%d", owner, number)
}

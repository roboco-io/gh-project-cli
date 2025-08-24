package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
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
	ID          string
	Number      int
	Title       string
	Description *string
	URL         string
	Closed      bool
	Owner       string
	ItemCount   int
	FieldCount  int
}

// ListUserProjectsOptions represents options for listing user projects
type ListUserProjectsOptions struct {
	Login string
	First int
	After *string
}

// ListOrgProjectsOptions represents options for listing organization projects
type ListOrgProjectsOptions struct {
	Login string
	First int
	After *string
}

// ListUserProjects lists projects for a user
func (s *ProjectService) ListUserProjects(ctx context.Context, opts ListUserProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"login": opts.Login,
		"first": opts.First,
	}
	if opts.After != nil {
		variables["after"] = *opts.After
	}

	var query graphql.ListUserProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list user projects: %w", err)
	}

	projects := make([]ProjectInfo, len(query.User.ProjectsV2.Nodes))
	for i, project := range query.User.ProjectsV2.Nodes {
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

	return projects, nil
}

// ListOrgProjects lists projects for an organization
func (s *ProjectService) ListOrgProjects(ctx context.Context, opts ListOrgProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := map[string]interface{}{
		"login": opts.Login,
		"first": opts.First,
	}
	if opts.After != nil {
		variables["after"] = *opts.After
	}

	var query graphql.ListOrgProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization projects: %w", err)
	}

	projects := make([]ProjectInfo, len(query.Organization.ProjectsV2.Nodes))
	for i, project := range query.Organization.ProjectsV2.Nodes {
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

	return projects, nil
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
	} else {
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
	ProjectID string
	Title     *string
	Closed    *bool
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
	ProjectID string
	ItemID    string
	FieldID   string
	Value     interface{}
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

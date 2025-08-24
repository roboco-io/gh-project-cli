package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
)

// ViewService handles view-related operations
type ViewService struct {
	client *api.Client
}

// NewViewService creates a new view service
func NewViewService(client *api.Client) *ViewService {
	return &ViewService{
		client: client,
	}
}

// ViewInfo represents simplified view information for display
type ViewInfo struct {
	ID          string
	Name        string
	Layout      graphql.ProjectV2ViewLayout
	Number      int
	Filter      *string
	ProjectID   string
	ProjectName string
	GroupBy     []ViewGroupByInfo
	SortBy      []ViewSortByInfo
}

// ViewGroupByInfo represents group by configuration information
type ViewGroupByInfo struct {
	FieldID   string
	FieldName string
	Direction graphql.ProjectV2ViewSortDirection
}

// ViewSortByInfo represents sort by configuration information
type ViewSortByInfo struct {
	FieldID   string
	FieldName string
	Direction graphql.ProjectV2ViewSortDirection
}

// CreateViewInput represents input for creating a view
type CreateViewInput struct {
	ProjectID string
	Name      string
	Layout    graphql.ProjectV2ViewLayout
}

// UpdateViewInput represents input for updating a view
type UpdateViewInput struct {
	ViewID string
	Name   *string
	Filter *string
}

// DeleteViewInput represents input for deleting a view
type DeleteViewInput struct {
	ViewID string
}

// CopyViewInput represents input for copying a view
type CopyViewInput struct {
	ProjectID string
	ViewID    string
	Name      string
}

// UpdateViewSortInput represents input for updating view sort configuration
type UpdateViewSortInput struct {
	ViewID    string
	SortByID  *string
	Direction graphql.ProjectV2ViewSortDirection
}

// UpdateViewGroupInput represents input for updating view group configuration
type UpdateViewGroupInput struct {
	ViewID      string
	GroupByID   *string
	Direction   graphql.ProjectV2ViewSortDirection
}

// CreateView creates a new project view
func (s *ViewService) CreateView(ctx context.Context, input CreateViewInput) (*graphql.ProjectV2View, error) {
	variables := graphql.BuildCreateViewVariables(graphql.CreateViewInput{
		ProjectID: input.ProjectID,
		Name:      input.Name,
		Layout:    input.Layout,
	})

	var mutation graphql.CreateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	return &mutation.CreateProjectV2View.ProjectV2View, nil
}

// UpdateView updates an existing project view
func (s *ViewService) UpdateView(ctx context.Context, input UpdateViewInput) (*graphql.ProjectV2View, error) {
	variables := graphql.BuildUpdateViewVariables(graphql.UpdateViewInput{
		ViewID: input.ViewID,
		Name:   input.Name,
		Filter: input.Filter,
	})

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update view: %w", err)
	}

	return &mutation.UpdateProjectV2View.ProjectV2View, nil
}

// DeleteView deletes a project view
func (s *ViewService) DeleteView(ctx context.Context, input DeleteViewInput) error {
	variables := graphql.BuildDeleteViewVariables(graphql.DeleteViewInput{
		ViewID: input.ViewID,
	})

	var mutation graphql.DeleteProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete view: %w", err)
	}

	return nil
}

// CopyView creates a copy of an existing view
func (s *ViewService) CopyView(ctx context.Context, input CopyViewInput) (*graphql.ProjectV2View, error) {
	variables := graphql.BuildCopyViewVariables(graphql.CopyViewInput{
		ProjectID: input.ProjectID,
		ViewID:    input.ViewID,
		Name:      input.Name,
	})

	var mutation graphql.CopyProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to copy view: %w", err)
	}

	return &mutation.CopyProjectV2View.ProjectV2View, nil
}

// UpdateViewSort updates the sort configuration for a view
func (s *ViewService) UpdateViewSort(ctx context.Context, input UpdateViewSortInput) error {
	variables := graphql.BuildUpdateViewSortByVariables(graphql.UpdateViewSortByInput{
		ViewID:    input.ViewID,
		SortByID:  input.SortByID,
		Direction: input.Direction,
	})

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update view sort: %w", err)
	}

	return nil
}

// UpdateViewGroup updates the group configuration for a view
func (s *ViewService) UpdateViewGroup(ctx context.Context, input UpdateViewGroupInput) error {
	variables := graphql.BuildUpdateViewGroupByVariables(graphql.UpdateViewGroupByInput{
		ViewID:    input.ViewID,
		GroupByID: input.GroupByID,
		Direction: input.Direction,
	})

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update view group: %w", err)
	}

	return nil
}

// GetProjectViews gets all views for a project
func (s *ViewService) GetProjectViews(ctx context.Context, projectID string) ([]ViewInfo, error) {
	variables := map[string]interface{}{
		"projectId": projectID,
	}

	var query graphql.GetProjectViewsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get project views: %w", err)
	}

	views := make([]ViewInfo, len(query.Node.ProjectV2.Views.Nodes))
	for i, view := range query.Node.ProjectV2.Views.Nodes {
		groupBy := make([]ViewGroupByInfo, len(view.GroupBy))
		for j, gb := range view.GroupBy {
			groupBy[j] = ViewGroupByInfo{
				FieldID:   gb.Field.ID,
				FieldName: gb.Field.Name,
				Direction: gb.Direction,
			}
		}

		sortBy := make([]ViewSortByInfo, len(view.SortBy))
		for j, sb := range view.SortBy {
			sortBy[j] = ViewSortByInfo{
				FieldID:   sb.Field.ID,
				FieldName: sb.Field.Name,
				Direction: sb.Direction,
			}
		}

		views[i] = ViewInfo{
			ID:        view.ID,
			Name:      view.Name,
			Layout:    view.Layout,
			Number:    view.Number,
			Filter:    view.Filter,
			ProjectID: projectID,
			GroupBy:   groupBy,
			SortBy:    sortBy,
		}
	}

	return views, nil
}

// GetView gets a specific view by ID
func (s *ViewService) GetView(ctx context.Context, viewID string) (*ViewInfo, error) {
	variables := map[string]interface{}{
		"viewId": viewID,
	}

	var query graphql.GetProjectViewQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get view: %w", err)
	}

	view := query.Node.ProjectV2View

	groupBy := make([]ViewGroupByInfo, len(view.GroupBy))
	for i, gb := range view.GroupBy {
		groupBy[i] = ViewGroupByInfo{
			FieldID:   gb.Field.ID,
			FieldName: gb.Field.Name,
			Direction: gb.Direction,
		}
	}

	sortBy := make([]ViewSortByInfo, len(view.SortBy))
	for i, sb := range view.SortBy {
		sortBy[i] = ViewSortByInfo{
			FieldID:   sb.Field.ID,
			FieldName: sb.Field.Name,
			Direction: sb.Direction,
		}
	}

	viewInfo := &ViewInfo{
		ID:      view.ID,
		Name:    view.Name,
		Layout:  view.Layout,
		Number:  view.Number,
		Filter:  view.Filter,
		GroupBy: groupBy,
		SortBy:  sortBy,
	}

	return viewInfo, nil
}

// ValidateViewName validates a view name
func ValidateViewName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("view name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("view name cannot exceed 100 characters")
	}
	return nil
}

// ValidateViewLayout validates a view layout
func ValidateViewLayout(layout string) (graphql.ProjectV2ViewLayout, error) {
	switch strings.ToUpper(layout) {
	case "TABLE", "TABLE_VIEW":
		return graphql.ProjectV2ViewLayoutTable, nil
	case "BOARD", "BOARD_VIEW":
		return graphql.ProjectV2ViewLayoutBoard, nil
	case "ROADMAP", "ROADMAP_VIEW":
		return graphql.ProjectV2ViewLayoutRoadmap, nil
	default:
		validLayouts := graphql.ValidViewLayouts()
		return "", fmt.Errorf("invalid view layout: %s (valid layouts: %s)", layout, strings.ToLower(strings.Join(validLayouts, ", ")))
	}
}

// ValidateSortDirection validates a sort direction
func ValidateSortDirection(direction string) (graphql.ProjectV2ViewSortDirection, error) {
	switch strings.ToUpper(direction) {
	case "ASC", "ASCENDING":
		return graphql.ProjectV2ViewSortDirectionASC, nil
	case "DESC", "DESCENDING":
		return graphql.ProjectV2ViewSortDirectionDESC, nil
	default:
		validDirections := graphql.ValidSortDirections()
		return "", fmt.Errorf("invalid sort direction: %s (valid directions: %s)", direction, strings.ToLower(strings.Join(validDirections, ", ")))
	}
}

// NormalizeSortDirection normalizes a sort direction string to the proper format
func NormalizeSortDirection(direction string) string {
	switch strings.ToUpper(direction) {
	case "ASC", "ASCENDING":
		return string(graphql.ProjectV2ViewSortDirectionASC)
	case "DESC", "DESCENDING":
		return string(graphql.ProjectV2ViewSortDirectionDESC)
	default:
		return strings.ToUpper(direction)
	}
}

// FormatViewLayout formats view layout for display
func FormatViewLayout(layout graphql.ProjectV2ViewLayout) string {
	return graphql.FormatViewLayout(layout)
}

// FormatSortDirection formats sort direction for display
func FormatSortDirection(direction graphql.ProjectV2ViewSortDirection) string {
	return graphql.FormatSortDirection(direction)
}
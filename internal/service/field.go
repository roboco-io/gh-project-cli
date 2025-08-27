package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
)

// FieldService handles field-related operations
type FieldService struct {
	client *api.Client
}

// NewFieldService creates a new field service
func NewFieldService(client *api.Client) *FieldService {
	return &FieldService{
		client: client,
	}
}

// FieldInfo represents simplified field information for display
type FieldInfo struct {
	ID          string
	Name        string
	DataType    graphql.ProjectV2FieldDataType
	ProjectID   string
	ProjectName string
	Options     []FieldOptionInfo
}

// FieldOptionInfo represents field option information
type FieldOptionInfo struct {
	Description *string
	ID          string
	Name        string
	Color       string
}

// CreateFieldInput represents input for creating a field
type CreateFieldInput struct {
	ProjectID           string
	Name                string
	DataType            graphql.ProjectV2FieldDataType
	SingleSelectOptions []string
	Duration            string
}

// UpdateFieldInput represents input for updating a field
type UpdateFieldInput struct {
	Name    *string
	FieldID string
}

// DeleteFieldInput represents input for deleting a field
type DeleteFieldInput struct {
	FieldID string
}

// CreateFieldOptionInput represents input for creating a field option
type CreateFieldOptionInput struct {
	Description *string
	FieldID     string
	Name        string
	Color       string
}

// UpdateFieldOptionInput represents input for updating a field option
type UpdateFieldOptionInput struct {
	Name        *string
	Color       *string
	Description *string
	OptionID    string
}

// DeleteFieldOptionInput represents input for deleting a field option
type DeleteFieldOptionInput struct {
	OptionID string
}

// CreateField creates a new project field
func (s *FieldService) CreateField(ctx context.Context, input CreateFieldInput) (*graphql.ProjectV2Field, error) {
	variables := graphql.BuildCreateFieldVariables(graphql.CreateFieldInput{
		ProjectID:           input.ProjectID,
		Name:                input.Name,
		DataType:            input.DataType,
		SingleSelectOptions: input.SingleSelectOptions,
		Duration:            input.Duration,
	})

	var mutation graphql.CreateFieldMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create field: %w", err)
	}

	return &mutation.CreateProjectV2Field.ProjectV2Field, nil
}

// UpdateField updates an existing project field
func (s *FieldService) UpdateField(ctx context.Context, input UpdateFieldInput) (*graphql.ProjectV2Field, error) {
	variables := graphql.BuildUpdateFieldVariables(graphql.UpdateFieldInput{
		FieldID: input.FieldID,
		Name:    input.Name,
	})

	var mutation graphql.UpdateFieldMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update field: %w", err)
	}

	return &mutation.UpdateProjectV2Field.ProjectV2Field, nil
}

// DeleteField deletes a project field
func (s *FieldService) DeleteField(ctx context.Context, input DeleteFieldInput) error {
	variables := graphql.BuildDeleteFieldVariables(graphql.DeleteFieldInput{
		FieldID: input.FieldID,
	})

	var mutation graphql.DeleteFieldMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete field: %w", err)
	}

	return nil
}

// CreateFieldOption creates a new single select field option
func (s *FieldService) CreateFieldOption(
	ctx context.Context,
	input CreateFieldOptionInput,
) (*graphql.ProjectV2SingleSelectFieldOption, error) {
	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	variables := graphql.BuildCreateSingleSelectFieldOptionVariables(graphql.CreateSingleSelectFieldOptionInput{
		FieldID:     input.FieldID,
		Name:        input.Name,
		Color:       input.Color,
		Description: description,
	})

	var mutation graphql.CreateSingleSelectFieldOptionMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create field option: %w", err)
	}

	return &mutation.CreateProjectV2SingleSelectFieldOption.ProjectV2SingleSelectFieldOption, nil
}

// UpdateFieldOption updates a single select field option
func (s *FieldService) UpdateFieldOption(
	ctx context.Context,
	input UpdateFieldOptionInput,
) (*graphql.ProjectV2SingleSelectFieldOption, error) {
	variables := graphql.BuildUpdateSingleSelectFieldOptionVariables(graphql.UpdateSingleSelectFieldOptionInput{
		OptionID:    input.OptionID,
		Name:        input.Name,
		Color:       input.Color,
		Description: input.Description,
	})

	var mutation graphql.UpdateSingleSelectFieldOptionMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update field option: %w", err)
	}

	return &mutation.UpdateProjectV2SingleSelectFieldOption.ProjectV2SingleSelectFieldOption, nil
}

// DeleteFieldOption deletes a single select field option
func (s *FieldService) DeleteFieldOption(ctx context.Context, input DeleteFieldOptionInput) error {
	variables := graphql.BuildDeleteSingleSelectFieldOptionVariables(graphql.DeleteSingleSelectFieldOptionInput{
		OptionID: input.OptionID,
	})

	var mutation graphql.DeleteSingleSelectFieldOptionMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete field option: %w", err)
	}

	return nil
}

// GetProjectFields gets all fields for a project
func (s *FieldService) GetProjectFields(ctx context.Context, owner string, number int, isOrg bool) ([]FieldInfo, error) {
	// Get project first to get fields
	projectService := NewProjectService(s.client)
	project, err := projectService.GetProject(ctx, owner, number, isOrg)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	fields := make([]FieldInfo, len(project.Fields.Nodes))
	for i, field := range project.Fields.Nodes {
		options := make([]FieldOptionInfo, len(field.Options.Nodes))
		for j, option := range field.Options.Nodes {
			options[j] = FieldOptionInfo{
				ID:          option.ID,
				Name:        option.Name,
				Color:       option.Color,
				Description: option.Description,
			}
		}

		fields[i] = FieldInfo{
			ID:          field.ID,
			Name:        field.Name,
			DataType:    field.DataType,
			Options:     options,
			ProjectID:   project.ID,
			ProjectName: project.Title,
		}
	}

	return fields, nil
}

// ValidateFieldName validates a field name
func ValidateFieldName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("field name cannot be empty")
	}
	if len(name) > maxFieldNameLength {
		return fmt.Errorf("field name cannot exceed %d characters", maxFieldNameLength)
	}
	return nil
}

// ValidateFieldType validates a field data type
func ValidateFieldType(dataType string) (graphql.ProjectV2FieldDataType, error) {
	switch strings.ToUpper(dataType) {
	case "TEXT":
		return graphql.ProjectV2FieldDataTypeText, nil
	case "NUMBER":
		return graphql.ProjectV2FieldDataTypeNumber, nil
	case "DATE":
		return graphql.ProjectV2FieldDataTypeDate, nil
	case "SINGLE_SELECT":
		return graphql.ProjectV2FieldDataTypeSingleSelect, nil
	case "ITERATION":
		return graphql.ProjectV2FieldDataTypeIteration, nil
	default:
		return "", fmt.Errorf("invalid field type: %s (valid types: text, number, date, single_select, iteration)", dataType)
	}
}

// ValidateColor validates a single select field option color
func ValidateColor(color string) error {
	validColors := graphql.ValidSingleSelectColors()
	colorUpper := strings.ToUpper(color)

	for _, validColor := range validColors {
		if colorUpper == validColor {
			return nil
		}
	}

	return fmt.Errorf("invalid color: %s (valid colors: %s)", color, strings.ToLower(strings.Join(validColors, ", ")))
}

// NormalizeColor normalizes a color string to the proper format
func NormalizeColor(color string) string {
	return strings.ToUpper(color)
}

// FormatFieldDataType formats field data type for display
func FormatFieldDataType(dataType graphql.ProjectV2FieldDataType) string {
	switch dataType {
	case graphql.ProjectV2FieldDataTypeText:
		return "Text"
	case graphql.ProjectV2FieldDataTypeNumber:
		return "Number"
	case graphql.ProjectV2FieldDataTypeDate:
		return "Date"
	case graphql.ProjectV2FieldDataTypeSingleSelect:
		return "Single Select"
	case graphql.ProjectV2FieldDataTypeIteration:
		return "Iteration"
	default:
		return string(dataType)
	}
}

// FormatColor formats color for display
func FormatColor(color string) string {
	// Convert to title case for display
	lower := strings.ToLower(color)
	if lower != "" {
		return strings.ToUpper(lower[:1]) + lower[1:]
	}
	return color
}

package graphql

import "time"

// Field creation mutations and queries

// CreateFieldMutation represents the createProjectV2Field mutation
type CreateFieldMutation struct {
	CreateProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"createProjectV2Field(input: $input)"`
}

// UpdateFieldMutation represents the updateProjectV2Field mutation
type UpdateFieldMutation struct {
	UpdateProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"updateProjectV2Field(input: $input)"`
}

// DeleteFieldMutation represents the deleteProjectV2Field mutation
type DeleteFieldMutation struct {
	DeleteProjectV2Field struct {
		ProjectV2Field ProjectV2Field `graphql:"projectV2Field"`
	} `graphql:"deleteProjectV2Field(input: $input)"`
}

// CreateSingleSelectFieldOptionMutation represents the createProjectV2SingleSelectFieldOption mutation
type CreateSingleSelectFieldOptionMutation struct {
	CreateProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"createProjectV2SingleSelectFieldOption(input: $input)"`
}

// UpdateSingleSelectFieldOptionMutation represents the updateProjectV2SingleSelectFieldOption mutation
type UpdateSingleSelectFieldOptionMutation struct {
	UpdateProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"updateProjectV2SingleSelectFieldOption(input: $input)"`
}

// DeleteSingleSelectFieldOptionMutation represents the deleteProjectV2SingleSelectFieldOption mutation
type DeleteSingleSelectFieldOptionMutation struct {
	DeleteProjectV2SingleSelectFieldOption struct {
		ProjectV2SingleSelectFieldOption ProjectV2SingleSelectFieldOption `graphql:"projectV2SingleSelectFieldOption"`
	} `graphql:"deleteProjectV2SingleSelectFieldOption(input: $input)"`
}

// Field input types
type CreateFieldInput struct {
	ProjectID  string                 `json:"projectId"`
	Name       string                 `json:"name"`
	DataType   ProjectV2FieldDataType `json:"dataType"`
	SingleSelectOptions []string      `json:"singleSelectOptions,omitempty"`
}

type UpdateFieldInput struct {
	FieldID string  `json:"fieldId"`
	Name    *string `json:"name,omitempty"`
}

type DeleteFieldInput struct {
	FieldID string `json:"fieldId"`
}

type CreateSingleSelectFieldOptionInput struct {
	FieldID     string `json:"fieldId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

type UpdateSingleSelectFieldOptionInput struct {
	OptionID    string  `json:"singleSelectOptionId"`
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
}

type DeleteSingleSelectFieldOptionInput struct {
	OptionID string `json:"singleSelectOptionId"`
}

// Variable builders
func BuildCreateFieldVariables(input CreateFieldInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"projectId": input.ProjectID,
		"name":      input.Name,
		"dataType":  input.DataType,
	}

	if input.DataType == ProjectV2FieldDataTypeSingleSelect && len(input.SingleSelectOptions) > 0 {
		options := make([]map[string]interface{}, len(input.SingleSelectOptions))
		for i, option := range input.SingleSelectOptions {
			options[i] = map[string]interface{}{
				"name":  option,
				"color": "GRAY", // Default color
			}
		}
		inputMap["singleSelectOptions"] = options
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

func BuildUpdateFieldVariables(input UpdateFieldInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"fieldId": input.FieldID,
	}

	if input.Name != nil {
		inputMap["name"] = *input.Name
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

func BuildDeleteFieldVariables(input DeleteFieldInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"fieldId": input.FieldID,
		},
	}
}

func BuildCreateSingleSelectFieldOptionVariables(input CreateSingleSelectFieldOptionInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"fieldId": input.FieldID,
		"name":    input.Name,
		"color":   input.Color,
	}

	if input.Description != "" {
		inputMap["description"] = input.Description
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

func BuildUpdateSingleSelectFieldOptionVariables(input UpdateSingleSelectFieldOptionInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"singleSelectOptionId": input.OptionID,
	}

	if input.Name != nil {
		inputMap["name"] = *input.Name
	}
	if input.Color != nil {
		inputMap["color"] = *input.Color
	}
	if input.Description != nil {
		inputMap["description"] = *input.Description
	}

	return map[string]interface{}{
		"input": inputMap,
	}
}

func BuildDeleteSingleSelectFieldOptionVariables(input DeleteSingleSelectFieldOptionInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"singleSelectOptionId": input.OptionID,
		},
	}
}

// Extended field info for display
type FieldInfo struct {
	ID          string
	Name        string
	DataType    ProjectV2FieldDataType
	Options     []FieldOptionInfo
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FieldOptionInfo struct {
	ID          string
	Name        string
	Color       string
	Description *string
}

// Color constants for single select options
const (
	SingleSelectColorGray   = "GRAY"
	SingleSelectColorRed    = "RED"
	SingleSelectColorOrange = "ORANGE"
	SingleSelectColorYellow = "YELLOW"
	SingleSelectColorGreen  = "GREEN"
	SingleSelectColorBlue   = "BLUE"
	SingleSelectColorPurple = "PURPLE"
	SingleSelectColorPink   = "PINK"
)

// ValidSingleSelectColors returns all valid colors for single select options
func ValidSingleSelectColors() []string {
	return []string{
		SingleSelectColorGray,
		SingleSelectColorRed,
		SingleSelectColorOrange,
		SingleSelectColorYellow,
		SingleSelectColorGreen,
		SingleSelectColorBlue,
		SingleSelectColorPurple,
		SingleSelectColorPink,
	}
}
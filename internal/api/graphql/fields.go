package graphql

import (
	"strconv"
	"strings"
	"time"
)

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
	ProjectID           string                 `json:"projectId"`
	Name                string                 `json:"name"`
	DataType            ProjectV2FieldDataType `json:"dataType"`
	SingleSelectOptions []string               `json:"singleSelectOptions,omitempty"`
	Duration            string                 `json:"duration,omitempty"`
}

type UpdateFieldInput struct {
	Name    *string `json:"name,omitempty"`
	FieldID string  `json:"fieldId"`
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
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
	OptionID    string  `json:"singleSelectOptionId"`
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

	// Add iteration field configuration
	if input.DataType == ProjectV2FieldDataTypeIteration && input.Duration != "" {
		inputMap["iterationSetting"] = map[string]interface{}{
			"duration": parseDuration(input.Duration),
		}
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
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string
	Name      string
	DataType  ProjectV2FieldDataType
	Options   []FieldOptionInfo
}

type FieldOptionInfo struct {
	Description *string
	ID          string
	Name        string
	Color       string
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

// parseDuration parses duration string like "2w", "1m" into days
func parseDuration(duration string) int {
	if duration == "" {
		return 14 // Default 2 weeks
	}

	duration = strings.ToLower(strings.TrimSpace(duration))

	// Handle numeric part and unit
	var numStr string
	var unit string

	for i, char := range duration {
		if char >= '0' && char <= '9' {
			numStr += string(char)
		} else {
			unit = duration[i:]
			break
		}
	}

	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		return 14 // Default fallback
	}

	switch unit {
	case "d", "day", "days":
		return num
	case "w", "week", "weeks":
		return num * 7
	case "m", "month", "months":
		return num * 30
	default:
		return 14 // Default fallback
	}
}

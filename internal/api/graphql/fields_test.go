package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldMutations(t *testing.T) {
	t.Run("CreateField mutation structure", func(t *testing.T) {
		mutation := &CreateFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("UpdateField mutation structure", func(t *testing.T) {
		mutation := &UpdateFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("DeleteField mutation structure", func(t *testing.T) {
		mutation := &DeleteFieldMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("CreateSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &CreateSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("UpdateSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &UpdateSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})

	t.Run("DeleteSingleSelectFieldOption mutation structure", func(t *testing.T) {
		mutation := &DeleteSingleSelectFieldOptionMutation{}
		assert.NotNil(t, mutation)
	})
}

func TestFieldVariableBuilders(t *testing.T) {
	t.Run("BuildCreateFieldVariables creates proper variables", func(t *testing.T) {
		input := CreateFieldInput{
			ProjectID: "project-id",
			Name:      "Priority",
			DataType:  ProjectV2FieldDataTypeText,
		}

		variables := BuildCreateFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "project-id", inputVar["projectId"])
		assert.Equal(t, "Priority", inputVar["name"])
		assert.Equal(t, ProjectV2FieldDataTypeText, inputVar["dataType"])
	})

	t.Run("BuildCreateFieldVariables with single select options", func(t *testing.T) {
		input := CreateFieldInput{
			ProjectID:           "project-id",
			Name:                "Priority",
			DataType:            ProjectV2FieldDataTypeSingleSelect,
			SingleSelectOptions: []string{"High", "Medium", "Low"},
		}

		variables := BuildCreateFieldVariables(input)

		assert.NotNil(t, variables)
		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, ProjectV2FieldDataTypeSingleSelect, inputVar["dataType"])
		
		options := inputVar["singleSelectOptions"].([]map[string]interface{})
		assert.Len(t, options, 3)
		assert.Equal(t, "High", options[0]["name"])
		assert.Equal(t, "GRAY", options[0]["color"])
	})

	t.Run("BuildUpdateFieldVariables creates proper variables", func(t *testing.T) {
		newName := "Updated Priority"
		input := UpdateFieldInput{
			FieldID: "field-id",
			Name:    &newName,
		}

		variables := BuildUpdateFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "field-id", inputVar["fieldId"])
		assert.Equal(t, "Updated Priority", inputVar["name"])
	})

	t.Run("BuildDeleteFieldVariables creates proper variables", func(t *testing.T) {
		input := DeleteFieldInput{
			FieldID: "field-id",
		}

		variables := BuildDeleteFieldVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "field-id", inputVar["fieldId"])
	})

	t.Run("BuildCreateSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		input := CreateSingleSelectFieldOptionInput{
			FieldID:     "field-id",
			Name:        "Critical",
			Color:       SingleSelectColorRed,
			Description: "Critical priority items",
		}

		variables := BuildCreateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "field-id", inputVar["fieldId"])
		assert.Equal(t, "Critical", inputVar["name"])
		assert.Equal(t, SingleSelectColorRed, inputVar["color"])
		assert.Equal(t, "Critical priority items", inputVar["description"])
	})

	t.Run("BuildCreateSingleSelectFieldOptionVariables without description", func(t *testing.T) {
		input := CreateSingleSelectFieldOptionInput{
			FieldID: "field-id",
			Name:    "Critical",
			Color:   SingleSelectColorRed,
		}

		variables := BuildCreateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		inputVar := variables["input"].(map[string]interface{})
		assert.NotContains(t, inputVar, "description")
	})

	t.Run("BuildUpdateSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		newName := "Very High"
		newColor := SingleSelectColorOrange
		newDescription := "Very high priority"
		
		input := UpdateSingleSelectFieldOptionInput{
			OptionID:    "option-id",
			Name:        &newName,
			Color:       &newColor,
			Description: &newDescription,
		}

		variables := BuildUpdateSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "option-id", inputVar["singleSelectOptionId"])
		assert.Equal(t, "Very High", inputVar["name"])
		assert.Equal(t, SingleSelectColorOrange, inputVar["color"])
		assert.Equal(t, "Very high priority", inputVar["description"])
	})

	t.Run("BuildDeleteSingleSelectFieldOptionVariables creates proper variables", func(t *testing.T) {
		input := DeleteSingleSelectFieldOptionInput{
			OptionID: "option-id",
		}

		variables := BuildDeleteSingleSelectFieldOptionVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")

		inputVar := variables["input"].(map[string]interface{})
		assert.Equal(t, "option-id", inputVar["singleSelectOptionId"])
	})
}

func TestFieldDataTypes(t *testing.T) {
	t.Run("All field data types defined", func(t *testing.T) {
		assert.Equal(t, "TEXT", string(ProjectV2FieldDataTypeText))
		assert.Equal(t, "NUMBER", string(ProjectV2FieldDataTypeNumber))
		assert.Equal(t, "DATE", string(ProjectV2FieldDataTypeDate))
		assert.Equal(t, "SINGLE_SELECT", string(ProjectV2FieldDataTypeSingleSelect))
		assert.Equal(t, "ITERATION", string(ProjectV2FieldDataTypeIteration))
	})
}

func TestValidSingleSelectColors(t *testing.T) {
	t.Run("All valid colors returned", func(t *testing.T) {
		colors := ValidSingleSelectColors()
		
		assert.Len(t, colors, 8)
		assert.Contains(t, colors, SingleSelectColorGray)
		assert.Contains(t, colors, SingleSelectColorRed)
		assert.Contains(t, colors, SingleSelectColorOrange)
		assert.Contains(t, colors, SingleSelectColorYellow)
		assert.Contains(t, colors, SingleSelectColorGreen)
		assert.Contains(t, colors, SingleSelectColorBlue)
		assert.Contains(t, colors, SingleSelectColorPurple)
		assert.Contains(t, colors, SingleSelectColorPink)
	})
}

func TestFieldInfo(t *testing.T) {
	t.Run("FieldInfo structure", func(t *testing.T) {
		info := FieldInfo{
			ID:       "field-id",
			Name:     "Priority",
			DataType: ProjectV2FieldDataTypeSingleSelect,
		}

		assert.Equal(t, "field-id", info.ID)
		assert.Equal(t, "Priority", info.Name)
		assert.Equal(t, ProjectV2FieldDataTypeSingleSelect, info.DataType)
	})

	t.Run("FieldOptionInfo structure", func(t *testing.T) {
		description := "High priority option"
		option := FieldOptionInfo{
			ID:          "option-id",
			Name:        "High",
			Color:       SingleSelectColorRed,
			Description: &description,
		}

		assert.Equal(t, "option-id", option.ID)
		assert.Equal(t, "High", option.Name)
		assert.Equal(t, SingleSelectColorRed, option.Color)
		assert.Equal(t, "High priority option", *option.Description)
	})
}
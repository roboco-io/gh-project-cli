package graphql

import (
	"reflect"
	"testing"
)

func TestProjectV2ExportFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		format   ProjectV2ExportFormat
		expected string
	}{
		{"JSON format", ProjectV2ExportFormatJSON, "JSON"},
		{"CSV format", ProjectV2ExportFormatCSV, "CSV"},
		{"XML format", ProjectV2ExportFormatXML, "XML"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.format) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.format))
			}
		})
	}
}

func TestBulkOperationTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		opType   BulkOperationType
		expected string
	}{
		{"Update type", BulkOperationTypeUpdate, "UPDATE"},
		{"Delete type", BulkOperationTypeDelete, "DELETE"},
		{"Import type", BulkOperationTypeImport, "IMPORT"},
		{"Export type", BulkOperationTypeExport, "EXPORT"},
		{"Archive type", BulkOperationTypeArchive, "ARCHIVE"},
		{"Move type", BulkOperationTypeMove, "MOVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.opType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.opType))
			}
		})
	}
}

func TestBulkOperationStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   BulkOperationStatus
		expected string
	}{
		{"Pending status", BulkOperationStatusPending, "PENDING"},
		{"Running status", BulkOperationStatusRunning, "RUNNING"},
		{"Completed status", BulkOperationStatusCompleted, "COMPLETED"},
		{"Failed status", BulkOperationStatusFailed, "FAILED"},
		{"Cancelled status", BulkOperationStatusCancelled, "CANCELLED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.status))
			}
		})
	}
}

func TestBuildExportProjectVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    ExportProjectInput
		expected map[string]interface{}
	}{
		{
			name: "Complete input",
			input: ExportProjectInput{
				ProjectID:        "project-123",
				Format:           ProjectV2ExportFormatJSON,
				IncludeItems:     true,
				IncludeFields:    true,
				IncludeViews:     true,
				IncludeWorkflows: true,
				Filter:           stringPtr("status:open"),
			},
			expected: map[string]interface{}{
				"input": map[string]interface{}{
					"projectId":        "project-123",
					"format":           ProjectV2ExportFormatJSON,
					"includeItems":     true,
					"includeFields":    true,
					"includeViews":     true,
					"includeWorkflows": true,
					"filter":           "status:open",
				},
			},
		},
		{
			name: "Minimal input",
			input: ExportProjectInput{
				ProjectID:        "project-456",
				Format:           ProjectV2ExportFormatCSV,
				IncludeItems:     false,
				IncludeFields:    false,
				IncludeViews:     false,
				IncludeWorkflows: false,
			},
			expected: map[string]interface{}{
				"input": map[string]interface{}{
					"projectId":        "project-456",
					"format":           ProjectV2ExportFormatCSV,
					"includeItems":     false,
					"includeFields":    false,
					"includeViews":     false,
					"includeWorkflows": false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildExportProjectVariables(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestBuildImportProjectVariables(t *testing.T) {
	input := ImportProjectInput{
		ProjectID:     "project-123",
		Format:        ProjectV2ExportFormatJSON,
		Data:          `{"items": []}`,
		MergeStrategy: "merge",
	}

	expected := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId":     "project-123",
			"format":        ProjectV2ExportFormatJSON,
			"data":          `{"items": []}`,
			"mergeStrategy": "merge",
		},
	}

	result := BuildImportProjectVariables(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestBuildBulkUpdateItemsVariables(t *testing.T) {
	input := BulkUpdateItemsInput{
		ProjectID: "project-123",
		ItemIDs:   []string{"item1", "item2", "item3"},
		Updates: map[string]interface{}{
			"status":   "Done",
			"priority": "High",
		},
	}

	expected := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": "project-123",
			"itemIds":   []string{"item1", "item2", "item3"},
			"updates": map[string]interface{}{
				"status":   "Done",
				"priority": "High",
			},
		},
	}

	result := BuildBulkUpdateItemsVariables(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestBuildBulkDeleteItemsVariables(t *testing.T) {
	input := BulkDeleteItemsInput{
		ProjectID: "project-123",
		ItemIDs:   []string{"item1", "item2"},
	}

	expected := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": "project-123",
			"itemIds":   []string{"item1", "item2"},
		},
	}

	result := BuildBulkDeleteItemsVariables(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestBuildBulkArchiveItemsVariables(t *testing.T) {
	input := BulkArchiveItemsInput{
		ProjectID: "project-123",
		ItemIDs:   []string{"item1", "item2", "item3"},
	}

	expected := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": "project-123",
			"itemIds":   []string{"item1", "item2", "item3"},
		},
	}

	result := BuildBulkArchiveItemsVariables(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestValidExportFormats(t *testing.T) {
	expected := []string{"JSON", "CSV", "XML"}
	result := ValidExportFormats()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 export formats, got %d", len(result))
	}
}

func TestValidBulkOperationTypes(t *testing.T) {
	expected := []string{"UPDATE", "DELETE", "IMPORT", "EXPORT", "ARCHIVE", "MOVE"}
	result := ValidBulkOperationTypes()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	if len(result) != 6 {
		t.Errorf("Expected 6 bulk operation types, got %d", len(result))
	}
}

func TestValidBulkOperationStatuses(t *testing.T) {
	expected := []string{"PENDING", "RUNNING", "COMPLETED", "FAILED", "CANCELLED"}
	result := ValidBulkOperationStatuses()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 bulk operation statuses, got %d", len(result))
	}
}

func TestFormatExportFormat(t *testing.T) {
	tests := []struct {
		format   ProjectV2ExportFormat
		expected string
	}{
		{ProjectV2ExportFormatJSON, "JSON"},
		{ProjectV2ExportFormatCSV, "CSV"},
		{ProjectV2ExportFormatXML, "XML"},
		{ProjectV2ExportFormat("UNKNOWN"), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			result := FormatExportFormat(tt.format)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatBulkOperationType(t *testing.T) {
	tests := []struct {
		opType   BulkOperationType
		expected string
	}{
		{BulkOperationTypeUpdate, "Update"},
		{BulkOperationTypeDelete, "Delete"},
		{BulkOperationTypeImport, "Import"},
		{BulkOperationTypeExport, "Export"},
		{BulkOperationTypeArchive, "Archive"},
		{BulkOperationTypeMove, "Move"},
		{BulkOperationType("UNKNOWN"), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.opType), func(t *testing.T) {
			result := FormatBulkOperationType(tt.opType)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatBulkOperationStatus(t *testing.T) {
	tests := []struct {
		status   BulkOperationStatus
		expected string
	}{
		{BulkOperationStatusPending, "Pending"},
		{BulkOperationStatusRunning, "Running"},
		{BulkOperationStatusCompleted, "Completed"},
		{BulkOperationStatusFailed, "Failed"},
		{BulkOperationStatusCancelled, "Cancelled"},
		{BulkOperationStatus("UNKNOWN"), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := FormatBulkOperationStatus(tt.status)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}

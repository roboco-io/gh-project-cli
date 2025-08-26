package graphql

import "time"

// ProjectV2Analytics represents analytics data for a GitHub Project v2
type ProjectV2Analytics struct {
	Timeline         ProjectV2Timeline    `graphql:"timeline"`
	ProjectID        string               `graphql:"id"`
	Title            string               `graphql:"title"`
	ItemsByStatus    []ItemStatusCount    `graphql:"itemsByStatus"`
	ItemsByAssignee  []ItemAssigneeCount  `graphql:"itemsByAssignee"`
	ItemsByLabel     []ItemLabelCount     `graphql:"itemsByLabel"`
	ItemsByMilestone []ItemMilestoneCount `graphql:"itemsByMilestone"`
	Velocity         ProjectV2Velocity    `graphql:"velocity"`
	ItemCount        int                  `graphql:"totalItemCount"`
	FieldCount       int                  `graphql:"totalFieldCount"`
	ViewCount        int                  `graphql:"totalViewCount"`
}

// ItemStatusCount represents count of items by status
type ItemStatusCount struct {
	Status string `graphql:"status"`
	Count  int    `graphql:"count"`
}

// ItemAssigneeCount represents count of items by assignee
type ItemAssigneeCount struct {
	Assignee string `graphql:"assignee"`
	Count    int    `graphql:"count"`
}

// ItemLabelCount represents count of items by label
type ItemLabelCount struct {
	Label string `graphql:"label"`
	Count int    `graphql:"count"`
}

// ItemMilestoneCount represents count of items by milestone
type ItemMilestoneCount struct {
	Milestone string `graphql:"milestone"`
	Count     int    `graphql:"count"`
}

// ProjectV2Timeline represents project timeline data
type ProjectV2Timeline struct {
	StartDate  *time.Time          `graphql:"startDate"`
	EndDate    *time.Time          `graphql:"endDate"`
	Milestones []TimelineMilestone `graphql:"milestones"`
	Activities []TimelineActivity  `graphql:"activities"`
	Duration   int                 `graphql:"durationDays"`
}

// TimelineMilestone represents a milestone in the project timeline
type TimelineMilestone struct {
	ID          string     `graphql:"id"`
	Title       string     `graphql:"title"`
	DueDate     *time.Time `graphql:"dueDate"`
	State       string     `graphql:"state"`
	Progress    float64    `graphql:"progressPercentage"`
	ItemCount   int        `graphql:"itemCount"`
	ClosedCount int        `graphql:"closedItemCount"`
}

// TimelineActivity represents activity in the project timeline
type TimelineActivity struct {
	Date        time.Time `graphql:"date"`
	Type        string    `graphql:"type"`
	Description string    `graphql:"description"`
	Count       int       `graphql:"count"`
}

// ProjectV2Velocity represents project velocity metrics
type ProjectV2Velocity struct {
	Period          string            `graphql:"period"`
	WeeklyVelocity  []WeeklyVelocity  `graphql:"weeklyVelocity"`
	MonthlyVelocity []MonthlyVelocity `graphql:"monthlyVelocity"`
	LeadTime        VelocityMetric    `graphql:"leadTime"`
	CycleTime       VelocityMetric    `graphql:"cycleTime"`
	ClosureRate     float64           `graphql:"closureRate"`
	CompletedItems  int               `graphql:"completedItems"`
	AddedItems      int               `graphql:"addedItems"`
}

// WeeklyVelocity represents velocity metrics for a week
type WeeklyVelocity struct {
	Week      string  `graphql:"week"`
	Completed int     `graphql:"completed"`
	Added     int     `graphql:"added"`
	Velocity  float64 `graphql:"velocity"`
}

// MonthlyVelocity represents velocity metrics for a month
type MonthlyVelocity struct {
	Month     string  `graphql:"month"`
	Completed int     `graphql:"completed"`
	Added     int     `graphql:"added"`
	Velocity  float64 `graphql:"velocity"`
}

// VelocityMetric represents time-based velocity metrics
type VelocityMetric struct {
	Unit    string  `graphql:"unit"`
	Average float64 `graphql:"average"`
	Median  float64 `graphql:"median"`
	P95     float64 `graphql:"p95"`
}

// ProjectV2Export represents export data for a project
type ProjectV2Export struct {
	ProjectID   string                `graphql:"id"`
	Title       string                `graphql:"title"`
	Description *string               `graphql:"description"`
	ExportDate  time.Time             `graphql:"exportDate"`
	Format      ProjectV2ExportFormat `graphql:"format"`

	Items     []ExportItem     `graphql:"items"`
	Fields    []ExportField    `graphql:"fields"`
	Views     []ExportView     `graphql:"views"`
	Workflows []ExportWorkflow `graphql:"workflows"`
}

// ProjectV2ExportFormat represents export format types
type ProjectV2ExportFormat string

const (
	ProjectV2ExportFormatJSON ProjectV2ExportFormat = "JSON"
	ProjectV2ExportFormatCSV  ProjectV2ExportFormat = "CSV"
	ProjectV2ExportFormatXML  ProjectV2ExportFormat = "XML"
)

// ExportItem represents an item in the export
type ExportItem struct {
	ID          string                 `graphql:"id"`
	Title       string                 `graphql:"title"`
	Type        string                 `graphql:"type"`
	State       string                 `graphql:"state"`
	CreatedAt   time.Time              `graphql:"createdAt"`
	UpdatedAt   time.Time              `graphql:"updatedAt"`
	FieldValues []ExportItemFieldValue `graphql:"fieldValues"`
}

// ExportItemFieldValue represents a field value in the export
type ExportItemFieldValue struct {
	Value     interface{} `graphql:"value"`
	FieldID   string      `graphql:"fieldId"`
	FieldName string      `graphql:"fieldName"`
	FieldType string      `graphql:"fieldType"`
}

// ExportField represents a field in the export
type ExportField struct {
	ID       string `graphql:"id"`
	Name     string `graphql:"name"`
	DataType string `graphql:"dataType"`
	Options  []struct {
		ID    string `graphql:"id"`
		Name  string `graphql:"name"`
		Color string `graphql:"color"`
	} `graphql:"options"`
}

// ExportView represents a view in the export
type ExportView struct {
	ID     string `graphql:"id"`
	Name   string `graphql:"name"`
	Layout string `graphql:"layout"`
	Filter string `graphql:"filter"`
	SortBy []struct {
		FieldID   string `graphql:"fieldId"`
		Direction string `graphql:"direction"`
	} `graphql:"sortBy"`
	GroupBy []struct {
		FieldID   string `graphql:"fieldId"`
		Direction string `graphql:"direction"`
	} `graphql:"groupBy"`
}

// ExportWorkflow represents a workflow in the export
type ExportWorkflow struct {
	ID       string `graphql:"id"`
	Name     string `graphql:"name"`
	Triggers []struct {
		Event   *string `graphql:"event"`
		FieldID *string `graphql:"fieldId"`
		Value   *string `graphql:"value"`
		Type    string  `graphql:"type"`
	} `graphql:"triggers"`
	Actions []struct {
		FieldID *string `graphql:"fieldId"`
		Value   *string `graphql:"value"`
		ViewID  *string `graphql:"viewId"`
		Column  *string `graphql:"column"`
		Message *string `graphql:"message"`
		Type    string  `graphql:"type"`
	} `graphql:"actions"`
	Enabled bool `graphql:"enabled"`
}

// BulkOperation represents a bulk operation request
type BulkOperation struct {
	CreatedAt      time.Time             `graphql:"createdAt"`
	CompletedAt    *time.Time            `graphql:"completedAt"`
	ErrorMessage   *string               `graphql:"errorMessage"`
	ID             string                `graphql:"id"`
	Type           BulkOperationType     `graphql:"type"`
	Status         BulkOperationStatus   `graphql:"status"`
	Results        []BulkOperationResult `graphql:"results"`
	Progress       float64               `graphql:"progress"`
	TotalItems     int                   `graphql:"totalItems"`
	ProcessedItems int                   `graphql:"processedItems"`
	FailedItems    int                   `graphql:"failedItems"`
}

// BulkOperationType represents the type of bulk operation
type BulkOperationType string

const (
	BulkOperationTypeUpdate  BulkOperationType = "UPDATE"
	BulkOperationTypeDelete  BulkOperationType = "DELETE"
	BulkOperationTypeImport  BulkOperationType = "IMPORT"
	BulkOperationTypeExport  BulkOperationType = "EXPORT"
	BulkOperationTypeArchive BulkOperationType = "ARCHIVE"
	BulkOperationTypeMove    BulkOperationType = "MOVE"
)

// BulkOperationStatus represents the status of bulk operation
type BulkOperationStatus string

const (
	BulkOperationStatusPending   BulkOperationStatus = "PENDING"
	BulkOperationStatusRunning   BulkOperationStatus = "RUNNING"
	BulkOperationStatusCompleted BulkOperationStatus = "COMPLETED"
	BulkOperationStatusFailed    BulkOperationStatus = "FAILED"
	BulkOperationStatusCancelled BulkOperationStatus = "CANCELED"
)

// BulkOperationResult represents the result of a single item in bulk operation
type BulkOperationResult struct {
	ErrorMessage *string `graphql:"errorMessage"`
	ItemID       string  `graphql:"itemId"`
	Success      bool    `graphql:"success"`
}

// Queries

// GetProjectAnalyticsQuery gets analytics data for a project
type GetProjectAnalyticsQuery struct {
	Node struct {
		ProjectV2 ProjectV2Analytics `graphql:"... on ProjectV2"`
	} `graphql:"node(id: $projectId)"`
}

// GetBulkOperationQuery gets bulk operation status
type GetBulkOperationQuery struct {
	Node struct {
		BulkOperation BulkOperation `graphql:"... on BulkOperation"`
	} `graphql:"node(id: $operationId)"`
}

// Mutations

// ExportProjectMutation exports a project
type ExportProjectMutation struct {
	ExportProjectV2 struct {
		Export ProjectV2Export `graphql:"export"`
	} `graphql:"exportProjectV2(input: $input)"`
}

// ImportProjectMutation imports project data
type ImportProjectMutation struct {
	ImportProjectV2 struct {
		BulkOperation BulkOperation `graphql:"bulkOperation"`
	} `graphql:"importProjectV2(input: $input)"`
}

// BulkUpdateItemsMutation performs bulk update on items
type BulkUpdateItemsMutation struct {
	BulkUpdateProjectV2Items struct {
		BulkOperation BulkOperation `graphql:"bulkOperation"`
	} `graphql:"bulkUpdateProjectV2Items(input: $input)"`
}

// BulkDeleteItemsMutation performs bulk delete on items
type BulkDeleteItemsMutation struct {
	BulkDeleteProjectV2Items struct {
		BulkOperation BulkOperation `graphql:"bulkOperation"`
	} `graphql:"bulkDeleteProjectV2Items(input: $input)"`
}

// BulkArchiveItemsMutation performs bulk archive on items
type BulkArchiveItemsMutation struct {
	BulkArchiveProjectV2Items struct {
		BulkOperation BulkOperation `graphql:"bulkOperation"`
	} `graphql:"bulkArchiveProjectV2Items(input: $input)"`
}

// Input Types

// ExportProjectInput represents input for project export
type ExportProjectInput struct {
	Filter           *string               `json:"filter,omitempty"`
	ProjectID        string                `json:"projectId"`
	Format           ProjectV2ExportFormat `json:"format"`
	IncludeItems     bool                  `json:"includeItems"`
	IncludeFields    bool                  `json:"includeFields"`
	IncludeViews     bool                  `json:"includeViews"`
	IncludeWorkflows bool                  `json:"includeWorkflows"`
}

// ImportProjectInput represents input for project import
type ImportProjectInput struct {
	ProjectID     string                `json:"projectId"`
	Format        ProjectV2ExportFormat `json:"format"`
	Data          string                `json:"data"`
	MergeStrategy string                `json:"mergeStrategy"`
}

// BulkUpdateItemsInput represents input for bulk item update
type BulkUpdateItemsInput struct {
	Updates   map[string]interface{} `json:"updates"`
	ProjectID string                 `json:"projectId"`
	ItemIDs   []string               `json:"itemIds"`
}

// BulkDeleteItemsInput represents input for bulk item delete
type BulkDeleteItemsInput struct {
	ProjectID string   `json:"projectId"`
	ItemIDs   []string `json:"itemIds"`
}

// BulkArchiveItemsInput represents input for bulk item archive
type BulkArchiveItemsInput struct {
	ProjectID string   `json:"projectId"`
	ItemIDs   []string `json:"itemIds"`
}

// Variable Builders

// BuildExportProjectVariables builds variables for project export
func BuildExportProjectVariables(input ExportProjectInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId":        input.ProjectID,
			"format":           input.Format,
			"includeItems":     input.IncludeItems,
			"includeFields":    input.IncludeFields,
			"includeViews":     input.IncludeViews,
			"includeWorkflows": input.IncludeWorkflows,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.Filter != nil {
		inputMap["filter"] = *input.Filter
	}

	return vars
}

// BuildImportProjectVariables builds variables for project import
func BuildImportProjectVariables(input ImportProjectInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId":     input.ProjectID,
			"format":        input.Format,
			"data":          input.Data,
			"mergeStrategy": input.MergeStrategy,
		},
	}
}

// BuildBulkUpdateItemsVariables builds variables for bulk item update
func BuildBulkUpdateItemsVariables(input BulkUpdateItemsInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"itemIds":   input.ItemIDs,
			"updates":   input.Updates,
		},
	}
}

// BuildBulkDeleteItemsVariables builds variables for bulk item delete
func BuildBulkDeleteItemsVariables(input BulkDeleteItemsInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"itemIds":   input.ItemIDs,
		},
	}
}

// BuildBulkArchiveItemsVariables builds variables for bulk item archive
func BuildBulkArchiveItemsVariables(input BulkArchiveItemsInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"itemIds":   input.ItemIDs,
		},
	}
}

// Helper Functions

// ValidExportFormats returns all valid export formats
func ValidExportFormats() []string {
	return []string{
		string(ProjectV2ExportFormatJSON),
		string(ProjectV2ExportFormatCSV),
		string(ProjectV2ExportFormatXML),
	}
}

// ValidBulkOperationTypes returns all valid bulk operation types
func ValidBulkOperationTypes() []string {
	return []string{
		string(BulkOperationTypeUpdate),
		string(BulkOperationTypeDelete),
		string(BulkOperationTypeImport),
		string(BulkOperationTypeExport),
		string(BulkOperationTypeArchive),
		string(BulkOperationTypeMove),
	}
}

// ValidBulkOperationStatuses returns all valid bulk operation statuses
func ValidBulkOperationStatuses() []string {
	return []string{
		string(BulkOperationStatusPending),
		string(BulkOperationStatusRunning),
		string(BulkOperationStatusCompleted),
		string(BulkOperationStatusFailed),
		string(BulkOperationStatusCancelled),
	}
}

// FormatExportFormat formats export format for display
func FormatExportFormat(format ProjectV2ExportFormat) string {
	switch format {
	case ProjectV2ExportFormatJSON:
		return "JSON"
	case ProjectV2ExportFormatCSV:
		return "CSV"
	case ProjectV2ExportFormatXML:
		return "XML"
	default:
		return string(format)
	}
}

// FormatBulkOperationType formats bulk operation type for display
func FormatBulkOperationType(opType BulkOperationType) string {
	switch opType {
	case BulkOperationTypeUpdate:
		return "Update"
	case BulkOperationTypeDelete:
		return "Delete"
	case BulkOperationTypeImport:
		return "Import"
	case BulkOperationTypeExport:
		return "Export"
	case BulkOperationTypeArchive:
		return "Archive"
	case BulkOperationTypeMove:
		return "Move"
	default:
		return string(opType)
	}
}

// FormatBulkOperationStatus formats bulk operation status for display
func FormatBulkOperationStatus(status BulkOperationStatus) string {
	switch status {
	case BulkOperationStatusPending:
		return "Pending"
	case BulkOperationStatusRunning:
		return "Running"
	case BulkOperationStatusCompleted:
		return "Completed"
	case BulkOperationStatusFailed:
		return "Failed"
	case BulkOperationStatusCancelled:
		return "Canceled"
	default:
		return string(status)
	}
}

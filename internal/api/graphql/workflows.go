package graphql

import "time"

// ProjectV2Workflow represents a workflow automation in a GitHub Project v2
type ProjectV2Workflow struct {
	CreatedAt time.Time                  `graphql:"createdAt"`
	UpdatedAt time.Time                  `graphql:"updatedAt"`
	ID        string                     `graphql:"id"`
	Name      string                     `graphql:"name"`
	Triggers  []ProjectV2WorkflowTrigger `graphql:"triggers"`
	Actions   []ProjectV2WorkflowAction  `graphql:"actions"`
	Enabled   bool                       `graphql:"enabled"`
}

// ProjectV2WorkflowTrigger represents a trigger that starts a workflow
type ProjectV2WorkflowTrigger struct {
	Field    *ProjectV2Field              `graphql:"field"`
	Value    *string                      `graphql:"value"`
	Schedule *ProjectV2WorkflowSchedule   `graphql:"schedule"`
	ID       string                       `graphql:"id"`
	Type     ProjectV2WorkflowTriggerType `graphql:"type"`
	Event    ProjectV2WorkflowEvent       `graphql:"event"`
}

// ProjectV2WorkflowAction represents an action performed by a workflow
type ProjectV2WorkflowAction struct {
	ID   string                      `graphql:"id"`
	Type ProjectV2WorkflowActionType `graphql:"type"`

	// Field update actions
	Field *ProjectV2Field `graphql:"field"`
	Value *string         `graphql:"value"`

	// Item movement actions
	View   *ProjectV2View `graphql:"view"`
	Column *string        `graphql:"column"`

	// Notification actions
	Message    *string  `graphql:"message"`
	Recipients []string `graphql:"recipients"`
}

// ProjectV2WorkflowSchedule represents a time-based trigger schedule
type ProjectV2WorkflowSchedule struct {
	DayOfWeek *int                  `graphql:"dayOfWeek"`
	Hour      *int                  `graphql:"hour"`
	Minute    *int                  `graphql:"minute"`
	Type      ProjectV2ScheduleType `graphql:"type"`
	Interval  int                   `graphql:"interval"`
}

// ProjectV2WorkflowTriggerType represents the type of workflow trigger
type ProjectV2WorkflowTriggerType string

const (
	ProjectV2WorkflowTriggerTypeItemAdded       ProjectV2WorkflowTriggerType = "ITEM_ADDED"
	ProjectV2WorkflowTriggerTypeItemUpdated     ProjectV2WorkflowTriggerType = "ITEM_UPDATED"
	ProjectV2WorkflowTriggerTypeItemArchived    ProjectV2WorkflowTriggerType = "ITEM_ARCHIVED"
	ProjectV2WorkflowTriggerTypeFieldChanged    ProjectV2WorkflowTriggerType = "FIELD_CHANGED"
	ProjectV2WorkflowTriggerTypeStatusChanged   ProjectV2WorkflowTriggerType = "STATUS_CHANGED"
	ProjectV2WorkflowTriggerTypeAssigneeChanged ProjectV2WorkflowTriggerType = "ASSIGNEE_CHANGED"
	ProjectV2WorkflowTriggerTypeScheduled       ProjectV2WorkflowTriggerType = "SCHEDULED"
)

// ProjectV2WorkflowEvent represents specific events within trigger types
type ProjectV2WorkflowEvent string

const (
	ProjectV2WorkflowEventIssueOpened   ProjectV2WorkflowEvent = "ISSUE_OPENED"
	ProjectV2WorkflowEventIssueClosed   ProjectV2WorkflowEvent = "ISSUE_CLOSED"
	ProjectV2WorkflowEventIssueReopened ProjectV2WorkflowEvent = "ISSUE_REOPENED"
	ProjectV2WorkflowEventPROpened      ProjectV2WorkflowEvent = "PR_OPENED"
	ProjectV2WorkflowEventPRClosed      ProjectV2WorkflowEvent = "PR_CLOSED"
	ProjectV2WorkflowEventPRMerged      ProjectV2WorkflowEvent = "PR_MERGED"
	ProjectV2WorkflowEventPRDraft       ProjectV2WorkflowEvent = "PR_DRAFT"
	ProjectV2WorkflowEventPRReady       ProjectV2WorkflowEvent = "PR_READY"
)

// ProjectV2WorkflowActionType represents the type of workflow action
type ProjectV2WorkflowActionType string

const (
	ProjectV2WorkflowActionTypeSetField     ProjectV2WorkflowActionType = "SET_FIELD"
	ProjectV2WorkflowActionTypeClearField   ProjectV2WorkflowActionType = "CLEAR_FIELD"
	ProjectV2WorkflowActionTypeMoveToColumn ProjectV2WorkflowActionType = "MOVE_TO_COLUMN"
	ProjectV2WorkflowActionTypeArchiveItem  ProjectV2WorkflowActionType = "ARCHIVE_ITEM"
	ProjectV2WorkflowActionTypeAddToProject ProjectV2WorkflowActionType = "ADD_TO_PROJECT"
	ProjectV2WorkflowActionTypeNotify       ProjectV2WorkflowActionType = "NOTIFY"
	ProjectV2WorkflowActionTypeAssign       ProjectV2WorkflowActionType = "ASSIGN"
	ProjectV2WorkflowActionTypeAddComment   ProjectV2WorkflowActionType = "ADD_COMMENT"
)

// ProjectV2ScheduleType represents the type of schedule
type ProjectV2ScheduleType string

const (
	ProjectV2ScheduleTypeDaily   ProjectV2ScheduleType = "DAILY"
	ProjectV2ScheduleTypeWeekly  ProjectV2ScheduleType = "WEEKLY"
	ProjectV2ScheduleTypeMonthly ProjectV2ScheduleType = "MONTHLY"
	ProjectV2ScheduleTypeCustom  ProjectV2ScheduleType = "CUSTOM"
)

// Queries

// GetProjectWorkflowsQuery gets all workflows for a project
type GetProjectWorkflowsQuery struct {
	Node struct {
		ProjectV2 struct {
			Workflows struct {
				Nodes []ProjectV2Workflow `graphql:"nodes"`
			} `graphql:"workflows(first: 20)"`
		} `graphql:"... on ProjectV2"`
	} `graphql:"node(id: $projectId)"`
}

// GetWorkflowQuery gets a specific workflow by ID
type GetWorkflowQuery struct {
	Node struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"... on ProjectV2Workflow"`
	} `graphql:"node(id: $workflowId)"`
}

// Mutations

// CreateProjectWorkflowMutation creates a new workflow
type CreateProjectWorkflowMutation struct {
	CreateProjectV2Workflow struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"projectV2Workflow"`
	} `graphql:"createProjectV2Workflow(input: $input)"`
}

// UpdateProjectWorkflowMutation updates an existing workflow
type UpdateProjectWorkflowMutation struct {
	UpdateProjectV2Workflow struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"projectV2Workflow"`
	} `graphql:"updateProjectV2Workflow(input: $input)"`
}

// DeleteProjectWorkflowMutation deletes a workflow
type DeleteProjectWorkflowMutation struct {
	DeleteProjectV2Workflow struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"projectV2Workflow"`
	} `graphql:"deleteProjectV2Workflow(input: $input)"`
}

// EnableProjectWorkflowMutation enables a workflow
type EnableProjectWorkflowMutation struct {
	EnableProjectV2Workflow struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"projectV2Workflow"`
	} `graphql:"enableProjectV2Workflow(input: $input)"`
}

// DisableProjectWorkflowMutation disables a workflow
type DisableProjectWorkflowMutation struct {
	DisableProjectV2Workflow struct {
		ProjectV2Workflow ProjectV2Workflow `graphql:"projectV2Workflow"`
	} `graphql:"disableProjectV2Workflow(input: $input)"`
}

// Input Types

// CreateWorkflowInput represents input for creating a workflow
type CreateWorkflowInput struct {
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
}

// UpdateWorkflowInput represents input for updating a workflow
type UpdateWorkflowInput struct {
	Name       *string `json:"name,omitempty"`
	Enabled    *bool   `json:"enabled,omitempty"`
	WorkflowID string  `json:"workflowId"`
}

// DeleteWorkflowInput represents input for deleting a workflow
type DeleteWorkflowInput struct {
	WorkflowID string `json:"workflowId"`
}

// EnableWorkflowInput represents input for enabling a workflow
type EnableWorkflowInput struct {
	WorkflowID string `json:"workflowId"`
}

// DisableWorkflowInput represents input for disabling a workflow
type DisableWorkflowInput struct {
	WorkflowID string `json:"workflowId"`
}

// CreateTriggerInput represents input for creating a workflow trigger
type CreateTriggerInput struct {
	FieldID    *string                      `json:"fieldId,omitempty"`
	Value      *string                      `json:"value,omitempty"`
	WorkflowID string                       `json:"workflowId"`
	Type       ProjectV2WorkflowTriggerType `json:"type"`
	Event      ProjectV2WorkflowEvent       `json:"event,omitempty"`
}

// CreateActionInput represents input for creating a workflow action
type CreateActionInput struct {
	FieldID    *string                     `json:"fieldId,omitempty"`
	Value      *string                     `json:"value,omitempty"`
	ViewID     *string                     `json:"viewId,omitempty"`
	Column     *string                     `json:"column,omitempty"`
	Message    *string                     `json:"message,omitempty"`
	WorkflowID string                      `json:"workflowId"`
	Type       ProjectV2WorkflowActionType `json:"type"`
}

// Variable Builders

// BuildCreateWorkflowVariables builds variables for workflow creation
func BuildCreateWorkflowVariables(input CreateWorkflowInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": input.ProjectID,
			"name":      input.Name,
			"enabled":   input.Enabled,
		},
	}
}

// BuildUpdateWorkflowVariables builds variables for workflow update
func BuildUpdateWorkflowVariables(input UpdateWorkflowInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.Name != nil {
		inputMap["name"] = *input.Name
	}
	if input.Enabled != nil {
		inputMap["enabled"] = *input.Enabled
	}

	return vars
}

// BuildDeleteWorkflowVariables builds variables for workflow deletion
func BuildDeleteWorkflowVariables(input DeleteWorkflowInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
		},
	}
}

// BuildEnableWorkflowVariables builds variables for workflow enable
func BuildEnableWorkflowVariables(input EnableWorkflowInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
		},
	}
}

// BuildDisableWorkflowVariables builds variables for workflow disable
func BuildDisableWorkflowVariables(input DisableWorkflowInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
		},
	}
}

// BuildCreateTriggerVariables builds variables for trigger creation
func BuildCreateTriggerVariables(input CreateTriggerInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
			"type":       input.Type,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.Event != "" {
		inputMap["event"] = input.Event
	}
	if input.FieldID != nil {
		inputMap["fieldId"] = *input.FieldID
	}
	if input.Value != nil {
		inputMap["value"] = *input.Value
	}

	return vars
}

// BuildCreateActionVariables builds variables for action creation
func BuildCreateActionVariables(input CreateActionInput) map[string]interface{} {
	vars := map[string]interface{}{
		"input": map[string]interface{}{
			"workflowId": input.WorkflowID,
			"type":       input.Type,
		},
	}

	inputMap := vars["input"].(map[string]interface{})
	if input.FieldID != nil {
		inputMap["fieldId"] = *input.FieldID
	}
	if input.Value != nil {
		inputMap["value"] = *input.Value
	}
	if input.ViewID != nil {
		inputMap["viewId"] = *input.ViewID
	}
	if input.Column != nil {
		inputMap["column"] = *input.Column
	}
	if input.Message != nil {
		inputMap["message"] = *input.Message
	}

	return vars
}

// Helper Functions

// ValidTriggerTypes returns all valid workflow trigger types
func ValidTriggerTypes() []string {
	return []string{
		string(ProjectV2WorkflowTriggerTypeItemAdded),
		string(ProjectV2WorkflowTriggerTypeItemUpdated),
		string(ProjectV2WorkflowTriggerTypeItemArchived),
		string(ProjectV2WorkflowTriggerTypeFieldChanged),
		string(ProjectV2WorkflowTriggerTypeStatusChanged),
		string(ProjectV2WorkflowTriggerTypeAssigneeChanged),
		string(ProjectV2WorkflowTriggerTypeScheduled),
	}
}

// ValidActionTypes returns all valid workflow action types
func ValidActionTypes() []string {
	return []string{
		string(ProjectV2WorkflowActionTypeSetField),
		string(ProjectV2WorkflowActionTypeClearField),
		string(ProjectV2WorkflowActionTypeMoveToColumn),
		string(ProjectV2WorkflowActionTypeArchiveItem),
		string(ProjectV2WorkflowActionTypeAddToProject),
		string(ProjectV2WorkflowActionTypeNotify),
		string(ProjectV2WorkflowActionTypeAssign),
		string(ProjectV2WorkflowActionTypeAddComment),
	}
}

// ValidEventTypes returns all valid workflow event types
func ValidEventTypes() []string {
	return []string{
		string(ProjectV2WorkflowEventIssueOpened),
		string(ProjectV2WorkflowEventIssueClosed),
		string(ProjectV2WorkflowEventIssueReopened),
		string(ProjectV2WorkflowEventPROpened),
		string(ProjectV2WorkflowEventPRClosed),
		string(ProjectV2WorkflowEventPRMerged),
		string(ProjectV2WorkflowEventPRDraft),
		string(ProjectV2WorkflowEventPRReady),
	}
}

// ValidScheduleTypes returns all valid schedule types
func ValidScheduleTypes() []string {
	return []string{
		string(ProjectV2ScheduleTypeDaily),
		string(ProjectV2ScheduleTypeWeekly),
		string(ProjectV2ScheduleTypeMonthly),
		string(ProjectV2ScheduleTypeCustom),
	}
}

// FormatTriggerType formats trigger type for display
func FormatTriggerType(triggerType ProjectV2WorkflowTriggerType) string {
	switch triggerType {
	case ProjectV2WorkflowTriggerTypeItemAdded:
		return "Item Added"
	case ProjectV2WorkflowTriggerTypeItemUpdated:
		return "Item Updated"
	case ProjectV2WorkflowTriggerTypeItemArchived:
		return "Item Archived"
	case ProjectV2WorkflowTriggerTypeFieldChanged:
		return "Field Changed"
	case ProjectV2WorkflowTriggerTypeStatusChanged:
		return "Status Changed"
	case ProjectV2WorkflowTriggerTypeAssigneeChanged:
		return "Assignee Changed"
	case ProjectV2WorkflowTriggerTypeScheduled:
		return "Scheduled"
	default:
		return string(triggerType)
	}
}

// FormatActionType formats action type for display
func FormatActionType(actionType ProjectV2WorkflowActionType) string {
	switch actionType {
	case ProjectV2WorkflowActionTypeSetField:
		return "Set Field"
	case ProjectV2WorkflowActionTypeClearField:
		return "Clear Field"
	case ProjectV2WorkflowActionTypeMoveToColumn:
		return "Move to Column"
	case ProjectV2WorkflowActionTypeArchiveItem:
		return "Archive Item"
	case ProjectV2WorkflowActionTypeAddToProject:
		return "Add to Project"
	case ProjectV2WorkflowActionTypeNotify:
		return "Send Notification"
	case ProjectV2WorkflowActionTypeAssign:
		return "Assign User"
	case ProjectV2WorkflowActionTypeAddComment:
		return "Add Comment"
	default:
		return string(actionType)
	}
}

// FormatEvent formats event type for display
func FormatEvent(event ProjectV2WorkflowEvent) string {
	switch event {
	case ProjectV2WorkflowEventIssueOpened:
		return "Issue Opened"
	case ProjectV2WorkflowEventIssueClosed:
		return "Issue Closed"
	case ProjectV2WorkflowEventIssueReopened:
		return "Issue Reopened"
	case ProjectV2WorkflowEventPROpened:
		return "PR Opened"
	case ProjectV2WorkflowEventPRClosed:
		return "PR Closed"
	case ProjectV2WorkflowEventPRMerged:
		return "PR Merged"
	case ProjectV2WorkflowEventPRDraft:
		return "PR Draft"
	case ProjectV2WorkflowEventPRReady:
		return "PR Ready"
	default:
		return string(event)
	}
}

// FormatScheduleType formats schedule type for display
func FormatScheduleType(scheduleType ProjectV2ScheduleType) string {
	switch scheduleType {
	case ProjectV2ScheduleTypeDaily:
		return "Daily"
	case ProjectV2ScheduleTypeWeekly:
		return "Weekly"
	case ProjectV2ScheduleTypeMonthly:
		return "Monthly"
	case ProjectV2ScheduleTypeCustom:
		return "Custom"
	default:
		return string(scheduleType)
	}
}

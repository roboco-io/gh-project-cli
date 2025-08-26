package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectV2WorkflowTypes(t *testing.T) {
	t.Run("ProjectV2WorkflowTriggerType constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2WorkflowTriggerType("ITEM_ADDED"), ProjectV2WorkflowTriggerTypeItemAdded)
		assert.Equal(t, ProjectV2WorkflowTriggerType("ITEM_UPDATED"), ProjectV2WorkflowTriggerTypeItemUpdated)
		assert.Equal(t, ProjectV2WorkflowTriggerType("ITEM_ARCHIVED"), ProjectV2WorkflowTriggerTypeItemArchived)
		assert.Equal(t, ProjectV2WorkflowTriggerType("FIELD_CHANGED"), ProjectV2WorkflowTriggerTypeFieldChanged)
		assert.Equal(t, ProjectV2WorkflowTriggerType("STATUS_CHANGED"), ProjectV2WorkflowTriggerTypeStatusChanged)
		assert.Equal(t, ProjectV2WorkflowTriggerType("ASSIGNEE_CHANGED"), ProjectV2WorkflowTriggerTypeAssigneeChanged)
		assert.Equal(t, ProjectV2WorkflowTriggerType("SCHEDULED"), ProjectV2WorkflowTriggerTypeScheduled)
	})

	t.Run("ProjectV2WorkflowActionType constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2WorkflowActionType("SET_FIELD"), ProjectV2WorkflowActionTypeSetField)
		assert.Equal(t, ProjectV2WorkflowActionType("CLEAR_FIELD"), ProjectV2WorkflowActionTypeClearField)
		assert.Equal(t, ProjectV2WorkflowActionType("MOVE_TO_COLUMN"), ProjectV2WorkflowActionTypeMoveToColumn)
		assert.Equal(t, ProjectV2WorkflowActionType("ARCHIVE_ITEM"), ProjectV2WorkflowActionTypeArchiveItem)
		assert.Equal(t, ProjectV2WorkflowActionType("ADD_TO_PROJECT"), ProjectV2WorkflowActionTypeAddToProject)
		assert.Equal(t, ProjectV2WorkflowActionType("NOTIFY"), ProjectV2WorkflowActionTypeNotify)
		assert.Equal(t, ProjectV2WorkflowActionType("ASSIGN"), ProjectV2WorkflowActionTypeAssign)
		assert.Equal(t, ProjectV2WorkflowActionType("ADD_COMMENT"), ProjectV2WorkflowActionTypeAddComment)
	})

	t.Run("ProjectV2WorkflowEvent constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2WorkflowEvent("ISSUE_OPENED"), ProjectV2WorkflowEventIssueOpened)
		assert.Equal(t, ProjectV2WorkflowEvent("ISSUE_CLOSED"), ProjectV2WorkflowEventIssueClosed)
		assert.Equal(t, ProjectV2WorkflowEvent("ISSUE_REOPENED"), ProjectV2WorkflowEventIssueReopened)
		assert.Equal(t, ProjectV2WorkflowEvent("PR_OPENED"), ProjectV2WorkflowEventPROpened)
		assert.Equal(t, ProjectV2WorkflowEvent("PR_CLOSED"), ProjectV2WorkflowEventPRClosed)
		assert.Equal(t, ProjectV2WorkflowEvent("PR_MERGED"), ProjectV2WorkflowEventPRMerged)
		assert.Equal(t, ProjectV2WorkflowEvent("PR_DRAFT"), ProjectV2WorkflowEventPRDraft)
		assert.Equal(t, ProjectV2WorkflowEvent("PR_READY"), ProjectV2WorkflowEventPRReady)
	})

	t.Run("ProjectV2ScheduleType constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2ScheduleType("DAILY"), ProjectV2ScheduleTypeDaily)
		assert.Equal(t, ProjectV2ScheduleType("WEEKLY"), ProjectV2ScheduleTypeWeekly)
		assert.Equal(t, ProjectV2ScheduleType("MONTHLY"), ProjectV2ScheduleTypeMonthly)
		assert.Equal(t, ProjectV2ScheduleType("CUSTOM"), ProjectV2ScheduleTypeCustom)
	})
}

func TestWorkflowVariableBuilders(t *testing.T) {
	t.Run("BuildCreateWorkflowVariables creates proper variables", func(t *testing.T) {
		input := CreateWorkflowInput{
			ProjectID: "test-project-id",
			Name:      "Test Workflow",
			Enabled:   true,
		}

		variables := BuildCreateWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"projectId": "test-project-id",
				"name":      "Test Workflow",
				"enabled":   true,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateWorkflowVariables creates proper variables", func(t *testing.T) {
		name := "Updated Workflow"
		enabled := false
		input := UpdateWorkflowInput{
			WorkflowID: "test-workflow-id",
			Name:       &name,
			Enabled:    &enabled,
		}

		variables := BuildUpdateWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
				"name":       "Updated Workflow",
				"enabled":    false,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildUpdateWorkflowVariables with minimal input", func(t *testing.T) {
		input := UpdateWorkflowInput{
			WorkflowID: "test-workflow-id",
		}

		variables := BuildUpdateWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildDeleteWorkflowVariables creates proper variables", func(t *testing.T) {
		input := DeleteWorkflowInput{
			WorkflowID: "test-workflow-id",
		}

		variables := BuildDeleteWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildEnableWorkflowVariables creates proper variables", func(t *testing.T) {
		input := EnableWorkflowInput{
			WorkflowID: "test-workflow-id",
		}

		variables := BuildEnableWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildDisableWorkflowVariables creates proper variables", func(t *testing.T) {
		input := DisableWorkflowInput{
			WorkflowID: "test-workflow-id",
		}

		variables := BuildDisableWorkflowVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildCreateTriggerVariables creates proper variables", func(t *testing.T) {
		fieldID := "test-field-id"
		value := "test-value"
		input := CreateTriggerInput{
			WorkflowID: "test-workflow-id",
			Type:       ProjectV2WorkflowTriggerTypeFieldChanged,
			Event:      ProjectV2WorkflowEventIssueOpened,
			FieldID:    &fieldID,
			Value:      &value,
		}

		variables := BuildCreateTriggerVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
				"type":       ProjectV2WorkflowTriggerTypeFieldChanged,
				"event":      ProjectV2WorkflowEventIssueOpened,
				"fieldId":    "test-field-id",
				"value":      "test-value",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildCreateTriggerVariables with minimal input", func(t *testing.T) {
		input := CreateTriggerInput{
			WorkflowID: "test-workflow-id",
			Type:       ProjectV2WorkflowTriggerTypeItemAdded,
		}

		variables := BuildCreateTriggerVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
				"type":       ProjectV2WorkflowTriggerTypeItemAdded,
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildCreateActionVariables creates proper variables", func(t *testing.T) {
		fieldID := "test-field-id"
		value := "test-value"
		viewID := "test-view-id"
		column := "test-column"
		message := "test-message"
		input := CreateActionInput{
			WorkflowID: "test-workflow-id",
			Type:       ProjectV2WorkflowActionTypeSetField,
			FieldID:    &fieldID,
			Value:      &value,
			ViewID:     &viewID,
			Column:     &column,
			Message:    &message,
		}

		variables := BuildCreateActionVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
				"type":       ProjectV2WorkflowActionTypeSetField,
				"fieldId":    "test-field-id",
				"value":      "test-value",
				"viewId":     "test-view-id",
				"column":     "test-column",
				"message":    "test-message",
			},
		}

		assert.Equal(t, expected, variables)
	})

	t.Run("BuildCreateActionVariables with minimal input", func(t *testing.T) {
		input := CreateActionInput{
			WorkflowID: "test-workflow-id",
			Type:       ProjectV2WorkflowActionTypeArchiveItem,
		}

		variables := BuildCreateActionVariables(input)

		expected := map[string]interface{}{
			"input": map[string]interface{}{
				"workflowId": "test-workflow-id",
				"type":       ProjectV2WorkflowActionTypeArchiveItem,
			},
		}

		assert.Equal(t, expected, variables)
	})
}

func TestWorkflowHelperFunctions(t *testing.T) {
	t.Run("ValidTriggerTypes returns all valid trigger types", func(t *testing.T) {
		triggerTypes := ValidTriggerTypes()
		expected := []string{
			string(ProjectV2WorkflowTriggerTypeItemAdded),
			string(ProjectV2WorkflowTriggerTypeItemUpdated),
			string(ProjectV2WorkflowTriggerTypeItemArchived),
			string(ProjectV2WorkflowTriggerTypeFieldChanged),
			string(ProjectV2WorkflowTriggerTypeStatusChanged),
			string(ProjectV2WorkflowTriggerTypeAssigneeChanged),
			string(ProjectV2WorkflowTriggerTypeScheduled),
		}

		assert.Equal(t, expected, triggerTypes)
		assert.Len(t, triggerTypes, 7)
	})

	t.Run("ValidActionTypes returns all valid action types", func(t *testing.T) {
		actionTypes := ValidActionTypes()
		expected := []string{
			string(ProjectV2WorkflowActionTypeSetField),
			string(ProjectV2WorkflowActionTypeClearField),
			string(ProjectV2WorkflowActionTypeMoveToColumn),
			string(ProjectV2WorkflowActionTypeArchiveItem),
			string(ProjectV2WorkflowActionTypeAddToProject),
			string(ProjectV2WorkflowActionTypeNotify),
			string(ProjectV2WorkflowActionTypeAssign),
			string(ProjectV2WorkflowActionTypeAddComment),
		}

		assert.Equal(t, expected, actionTypes)
		assert.Len(t, actionTypes, 8)
	})

	t.Run("ValidEventTypes returns all valid event types", func(t *testing.T) {
		eventTypes := ValidEventTypes()
		expected := []string{
			string(ProjectV2WorkflowEventIssueOpened),
			string(ProjectV2WorkflowEventIssueClosed),
			string(ProjectV2WorkflowEventIssueReopened),
			string(ProjectV2WorkflowEventPROpened),
			string(ProjectV2WorkflowEventPRClosed),
			string(ProjectV2WorkflowEventPRMerged),
			string(ProjectV2WorkflowEventPRDraft),
			string(ProjectV2WorkflowEventPRReady),
		}

		assert.Equal(t, expected, eventTypes)
		assert.Len(t, eventTypes, 8)
	})

	t.Run("ValidScheduleTypes returns all valid schedule types", func(t *testing.T) {
		scheduleTypes := ValidScheduleTypes()
		expected := []string{
			string(ProjectV2ScheduleTypeDaily),
			string(ProjectV2ScheduleTypeWeekly),
			string(ProjectV2ScheduleTypeMonthly),
			string(ProjectV2ScheduleTypeCustom),
		}

		assert.Equal(t, expected, scheduleTypes)
		assert.Len(t, scheduleTypes, 4)
	})

	t.Run("FormatTriggerType formats correctly", func(t *testing.T) {
		assert.Equal(t, "Item Added", FormatTriggerType(ProjectV2WorkflowTriggerTypeItemAdded))
		assert.Equal(t, "Item Updated", FormatTriggerType(ProjectV2WorkflowTriggerTypeItemUpdated))
		assert.Equal(t, "Item Archived", FormatTriggerType(ProjectV2WorkflowTriggerTypeItemArchived))
		assert.Equal(t, "Field Changed", FormatTriggerType(ProjectV2WorkflowTriggerTypeFieldChanged))
		assert.Equal(t, "Status Changed", FormatTriggerType(ProjectV2WorkflowTriggerTypeStatusChanged))
		assert.Equal(t, "Assignee Changed", FormatTriggerType(ProjectV2WorkflowTriggerTypeAssigneeChanged))
		assert.Equal(t, "Scheduled", FormatTriggerType(ProjectV2WorkflowTriggerTypeScheduled))
		assert.Equal(t, "UNKNOWN_TRIGGER", FormatTriggerType(ProjectV2WorkflowTriggerType("UNKNOWN_TRIGGER")))
	})

	t.Run("FormatActionType formats correctly", func(t *testing.T) {
		assert.Equal(t, "Set Field", FormatActionType(ProjectV2WorkflowActionTypeSetField))
		assert.Equal(t, "Clear Field", FormatActionType(ProjectV2WorkflowActionTypeClearField))
		assert.Equal(t, "Move to Column", FormatActionType(ProjectV2WorkflowActionTypeMoveToColumn))
		assert.Equal(t, "Archive Item", FormatActionType(ProjectV2WorkflowActionTypeArchiveItem))
		assert.Equal(t, "Add to Project", FormatActionType(ProjectV2WorkflowActionTypeAddToProject))
		assert.Equal(t, "Send Notification", FormatActionType(ProjectV2WorkflowActionTypeNotify))
		assert.Equal(t, "Assign User", FormatActionType(ProjectV2WorkflowActionTypeAssign))
		assert.Equal(t, "Add Comment", FormatActionType(ProjectV2WorkflowActionTypeAddComment))
		assert.Equal(t, "UNKNOWN_ACTION", FormatActionType(ProjectV2WorkflowActionType("UNKNOWN_ACTION")))
	})

	t.Run("FormatEvent formats correctly", func(t *testing.T) {
		assert.Equal(t, "Issue Opened", FormatEvent(ProjectV2WorkflowEventIssueOpened))
		assert.Equal(t, "Issue Closed", FormatEvent(ProjectV2WorkflowEventIssueClosed))
		assert.Equal(t, "Issue Reopened", FormatEvent(ProjectV2WorkflowEventIssueReopened))
		assert.Equal(t, "PR Opened", FormatEvent(ProjectV2WorkflowEventPROpened))
		assert.Equal(t, "PR Closed", FormatEvent(ProjectV2WorkflowEventPRClosed))
		assert.Equal(t, "PR Merged", FormatEvent(ProjectV2WorkflowEventPRMerged))
		assert.Equal(t, "PR Draft", FormatEvent(ProjectV2WorkflowEventPRDraft))
		assert.Equal(t, "PR Ready", FormatEvent(ProjectV2WorkflowEventPRReady))
		assert.Equal(t, "UNKNOWN_EVENT", FormatEvent(ProjectV2WorkflowEvent("UNKNOWN_EVENT")))
	})

	t.Run("FormatScheduleType formats correctly", func(t *testing.T) {
		assert.Equal(t, "Daily", FormatScheduleType(ProjectV2ScheduleTypeDaily))
		assert.Equal(t, "Weekly", FormatScheduleType(ProjectV2ScheduleTypeWeekly))
		assert.Equal(t, "Monthly", FormatScheduleType(ProjectV2ScheduleTypeMonthly))
		assert.Equal(t, "Custom", FormatScheduleType(ProjectV2ScheduleTypeCustom))
		assert.Equal(t, "UNKNOWN_SCHEDULE", FormatScheduleType(ProjectV2ScheduleType("UNKNOWN_SCHEDULE")))
	})
}

func TestWorkflowStructures(t *testing.T) {
	t.Run("ProjectV2Workflow structure validation", func(t *testing.T) {
		workflow := ProjectV2Workflow{
			ID:      "workflow-id",
			Name:    "Test Workflow",
			Enabled: true,
		}

		assert.Equal(t, "workflow-id", workflow.ID)
		assert.Equal(t, "Test Workflow", workflow.Name)
		assert.True(t, workflow.Enabled)
	})

	t.Run("ProjectV2WorkflowTrigger structure validation", func(t *testing.T) {
		trigger := ProjectV2WorkflowTrigger{
			ID:    "trigger-id",
			Type:  ProjectV2WorkflowTriggerTypeFieldChanged,
			Event: ProjectV2WorkflowEventIssueOpened,
		}

		assert.Equal(t, "trigger-id", trigger.ID)
		assert.Equal(t, ProjectV2WorkflowTriggerTypeFieldChanged, trigger.Type)
		assert.Equal(t, ProjectV2WorkflowEventIssueOpened, trigger.Event)
	})

	t.Run("ProjectV2WorkflowAction structure validation", func(t *testing.T) {
		action := ProjectV2WorkflowAction{
			ID:   "action-id",
			Type: ProjectV2WorkflowActionTypeSetField,
		}

		assert.Equal(t, "action-id", action.ID)
		assert.Equal(t, ProjectV2WorkflowActionTypeSetField, action.Type)
	})

	t.Run("ProjectV2WorkflowSchedule structure validation", func(t *testing.T) {
		hour := 9
		minute := 30
		schedule := ProjectV2WorkflowSchedule{
			Type:     ProjectV2ScheduleTypeDaily,
			Interval: 1,
			Hour:     &hour,
			Minute:   &minute,
		}

		assert.Equal(t, ProjectV2ScheduleTypeDaily, schedule.Type)
		assert.Equal(t, 1, schedule.Interval)
		assert.Equal(t, 9, *schedule.Hour)
		assert.Equal(t, 30, *schedule.Minute)
	})
}

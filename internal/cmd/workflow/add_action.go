package workflow

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// AddActionOptions holds options for the add-action command
type AddActionOptions struct {
	WorkflowID string
	Type       string
	FieldID    string
	Value      string
	ViewID     string
	Column     string
	Message    string
	Format     string
}

// NewAddActionCmd creates the add-action command
func NewAddActionCmd() *cobra.Command {
	opts := &AddActionOptions{}

	cmd := &cobra.Command{
		Use:   "add-action <workflow-id> <action-type>",
		Short: "Add action to workflow",
		Long: `Add an action to an existing workflow.

Actions define what should happen when a workflow is triggered. Different
action types require different parameters.

Action Types:
  set-field         - Set field to specific value (requires --field and --value)
  clear-field       - Clear field value (requires --field)
  move-to-column    - Move item to different column (requires --view and --column)
  archive-item      - Archive the item (no additional parameters)
  add-to-project    - Add item to another project (future implementation)
  notify            - Send notification to users (requires --message)
  assign            - Assign user to item (requires --value with username)
  add-comment       - Add comment to issue/PR (requires --message)

Examples:
  ghp workflow add-action workflow-id set-field --field priority-id --value "High"
  ghp workflow add-action workflow-id clear-field --field status-id
  ghp workflow add-action workflow-id move-to-column --view board-id --column "In Progress"
  ghp workflow add-action workflow-id archive-item
  ghp workflow add-action workflow-id notify --message "Item needs attention"
  ghp workflow add-action workflow-id assign --value "octocat"
  ghp workflow add-action workflow-id add-comment --message "Automatically triaged"`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Type = args[1]
			opts.Format = cmd.Flag("format").Value.String()
			return runAddAction(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.FieldID, "field", "", "Field ID for field-based actions")
	cmd.Flags().StringVar(&opts.Value, "value", "", "Value for set-field and assign actions")
	cmd.Flags().StringVar(&opts.ViewID, "view", "", "View ID for move-to-column actions")
	cmd.Flags().StringVar(&opts.Column, "column", "", "Column name for move-to-column actions")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Message for notify and add-comment actions")

	return cmd
}

func runAddAction(ctx context.Context, opts *AddActionOptions) error {
	// Validate action type
	actionType, err := service.ValidateActionType(opts.Type)
	if err != nil {
		return err
	}

	// Validate required parameters based on action type
	switch actionType {
	case graphql.ProjectV2WorkflowActionTypeSetField:
		if opts.FieldID == "" || opts.Value == "" {
			return fmt.Errorf("set-field action requires both --field and --value parameters")
		}
	case graphql.ProjectV2WorkflowActionTypeClearField:
		if opts.FieldID == "" {
			return fmt.Errorf("clear-field action requires --field parameter")
		}
	case graphql.ProjectV2WorkflowActionTypeMoveToColumn:
		if opts.ViewID == "" || opts.Column == "" {
			return fmt.Errorf("move-to-column action requires both --view and --column parameters")
		}
	case graphql.ProjectV2WorkflowActionTypeNotify:
		if opts.Message == "" {
			return fmt.Errorf("notify action requires --message parameter")
		}
	case graphql.ProjectV2WorkflowActionTypeAssign:
		if opts.Value == "" {
			return fmt.Errorf("assign action requires --value parameter with username")
		}
	case graphql.ProjectV2WorkflowActionTypeAddComment:
		if opts.Message == "" {
			return fmt.Errorf("add-comment action requires --message parameter")
		}
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	// Prepare input
	input := service.CreateActionInput{
		WorkflowID: opts.WorkflowID,
		Type:       actionType,
	}

	if opts.FieldID != "" {
		input.FieldID = &opts.FieldID
	}
	if opts.Value != "" {
		input.Value = &opts.Value
	}
	if opts.ViewID != "" {
		input.ViewID = &opts.ViewID
	}
	if opts.Column != "" {
		input.Column = &opts.Column
	}
	if opts.Message != "" {
		input.Message = &opts.Message
	}

	// Create action
	err = workflowService.CreateAction(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to add action: %w", err)
	}

	// Output result
	return outputAddedAction(&input, opts.Format)
}

func outputAddedAction(input *service.CreateActionInput, format string) error {
	switch format {
	case "json":
		return outputAddedActionJSON(input)
	case "table":
		return outputAddedActionTable(input)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputAddedActionTable(input *service.CreateActionInput) error {
	fmt.Printf("✅ Action added to workflow successfully\n\n")

	fmt.Printf("Action Details:\n")
	fmt.Printf("  Workflow ID: %s\n", input.WorkflowID)
	fmt.Printf("  Type: %s\n", service.FormatActionType(input.Type))

	if input.FieldID != nil {
		fmt.Printf("  Field ID: %s\n", *input.FieldID)
	}

	if input.Value != nil {
		fmt.Printf("  Value: %s\n", *input.Value)
	}

	if input.ViewID != nil {
		fmt.Printf("  View ID: %s\n", *input.ViewID)
	}

	if input.Column != nil {
		fmt.Printf("  Column: %s\n", *input.Column)
	}

	if input.Message != nil {
		fmt.Printf("  Message: %s\n", *input.Message)
	}

	return nil
}

func outputAddedActionJSON(input *service.CreateActionInput) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"workflowId\": \"%s\",\n", input.WorkflowID)
	fmt.Printf("  \"actionType\": \"%s\"", input.Type)

	if input.FieldID != nil {
		fmt.Printf(",\n  \"fieldId\": \"%s\"", *input.FieldID)
	}

	if input.Value != nil {
		fmt.Printf(",\n  \"value\": \"%s\"", *input.Value)
	}

	if input.ViewID != nil {
		fmt.Printf(",\n  \"viewId\": \"%s\"", *input.ViewID)
	}

	if input.Column != nil {
		fmt.Printf(",\n  \"column\": \"%s\"", *input.Column)
	}

	if input.Message != nil {
		fmt.Printf(",\n  \"message\": \"%s\"", *input.Message)
	}

	fmt.Printf("\n}\n")

	return nil
}
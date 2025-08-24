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

// AddTriggerOptions holds options for the add-trigger command
type AddTriggerOptions struct {
	WorkflowID string
	Type       string
	Event      string
	FieldID    string
	Value      string
	Format     string
}

// NewAddTriggerCmd creates the add-trigger command
func NewAddTriggerCmd() *cobra.Command {
	opts := &AddTriggerOptions{}

	cmd := &cobra.Command{
		Use:   "add-trigger <workflow-id> <trigger-type>",
		Short: "Add trigger to workflow",
		Long: `Add a trigger to an existing workflow.

Triggers define when a workflow should execute. Different trigger types
support different options and events.

Trigger Types:
  item-added        - When items are added to the project
  item-updated      - When items are modified  
  item-archived     - When items are archived
  field-changed     - When specific field values change
  status-changed    - When issue/PR status changes
  assignee-changed  - When assignee is modified
  scheduled         - Time-based triggers (future implementation)

Event Types (for item triggers):
  issue-opened      - When issues are opened
  issue-closed      - When issues are closed
  issue-reopened    - When issues are reopened
  pr-opened         - When pull requests are opened
  pr-closed         - When pull requests are closed
  pr-merged         - When pull requests are merged
  pr-draft          - When pull requests are converted to draft
  pr-ready          - When pull requests are marked ready for review

Examples:
  ghp workflow add-trigger workflow-id item-added --event issue-opened
  ghp workflow add-trigger workflow-id field-changed --field priority-id --value "High"
  ghp workflow add-trigger workflow-id status-changed --event pr-merged`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Type = args[1]
			opts.Format = cmd.Flag("format").Value.String()
			return runAddTrigger(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Event, "event", "", "Event type for the trigger")
	cmd.Flags().StringVar(&opts.FieldID, "field", "", "Field ID for field-based triggers")
	cmd.Flags().StringVar(&opts.Value, "value", "", "Value for field-based triggers")

	return cmd
}

func runAddTrigger(ctx context.Context, opts *AddTriggerOptions) error {
	// Validate trigger type
	triggerType, err := service.ValidateTriggerType(opts.Type)
	if err != nil {
		return err
	}

	// Validate event type if provided
	var event graphql.ProjectV2WorkflowEvent
	if opts.Event != "" {
		event, err = service.ValidateEventType(opts.Event)
		if err != nil {
			return err
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
	input := service.CreateTriggerInput{
		WorkflowID: opts.WorkflowID,
		Type:       triggerType,
		Event:      event,
	}

	if opts.FieldID != "" {
		input.FieldID = &opts.FieldID
	}
	if opts.Value != "" {
		input.Value = &opts.Value
	}

	// Create trigger
	err = workflowService.CreateTrigger(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to add trigger: %w", err)
	}

	// Output result
	return outputAddedTrigger(&input, opts.Format)
}

func outputAddedTrigger(input *service.CreateTriggerInput, format string) error {
	switch format {
	case "json":
		return outputAddedTriggerJSON(input)
	case "table":
		return outputAddedTriggerTable(input)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputAddedTriggerTable(input *service.CreateTriggerInput) error {
	fmt.Printf("✅ Trigger added to workflow successfully\n\n")

	fmt.Printf("Trigger Details:\n")
	fmt.Printf("  Workflow ID: %s\n", input.WorkflowID)
	fmt.Printf("  Type: %s\n", service.FormatTriggerType(input.Type))

	if input.Event != "" {
		fmt.Printf("  Event: %s\n", service.FormatEvent(input.Event))
	}

	if input.FieldID != nil {
		fmt.Printf("  Field ID: %s\n", *input.FieldID)
	}

	if input.Value != nil {
		fmt.Printf("  Value: %s\n", *input.Value)
	}

	return nil
}

func outputAddedTriggerJSON(input *service.CreateTriggerInput) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"workflowId\": \"%s\",\n", input.WorkflowID)
	fmt.Printf("  \"triggerType\": \"%s\"", input.Type)

	if input.Event != "" {
		fmt.Printf(",\n  \"event\": \"%s\"", input.Event)
	}

	if input.FieldID != nil {
		fmt.Printf(",\n  \"fieldId\": \"%s\"", *input.FieldID)
	}

	if input.Value != nil {
		fmt.Printf(",\n  \"value\": \"%s\"", *input.Value)
	}

	fmt.Printf("\n}\n")

	return nil
}
package workflow

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// ToggleOptions holds options for workflow enable/disable commands
type ToggleOptions struct {
	WorkflowID string
	Format     string
}

// ToggleConfig holds configuration for workflow toggle commands
type ToggleConfig struct {
	Use     string
	Short   string
	Long    string
	Action  func(context.Context, *service.WorkflowService, string) (*graphql.ProjectV2Workflow, error)
	Success string
}

// createWorkflowToggleCmd creates a workflow enable/disable command with shared logic
func createWorkflowToggleCmd(config ToggleConfig) *cobra.Command {
	opts := &ToggleOptions{}

	cmd := &cobra.Command{
		Use:   config.Use,
		Short: config.Short,
		Long:  config.Long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runWorkflowToggle(cmd.Context(), opts, config)
		},
	}

	return cmd
}

func runWorkflowToggle(ctx context.Context, opts *ToggleOptions, config ToggleConfig) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	// Execute the action (enable/disable)
	workflow, err := config.Action(ctx, workflowService, opts.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed to %s workflow: %w", config.Success, err)
	}

	// Output result
	return outputWorkflowToggleResult(workflow, config.Success, opts.Format)
}

func outputWorkflowToggleResult(workflow *graphql.ProjectV2Workflow, action, format string) error {
	switch format {
	case formatJSON:
		return outputWorkflowToggleResultJSON(workflow, action)
	case formatTable:
		return outputWorkflowToggleResultTable(workflow, action)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputWorkflowToggleResultTable(workflow *graphql.ProjectV2Workflow, action string) error {
	status := statusDisabled
	if workflow.Enabled {
		status = statusEnabled
	}

	fmt.Printf("✅ Workflow '%s' %sd successfully\n\n", workflow.Name, action)
	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Status: %s\n", status)
	fmt.Printf("  Triggers: %d configured\n", len(workflow.Triggers))
	fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))

	return nil
}

func outputWorkflowToggleResultJSON(workflow *graphql.ProjectV2Workflow, action string) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", workflow.ID)
	fmt.Printf("  \"name\": \"%s\",\n", workflow.Name)
	fmt.Printf("  \"enabled\": %t,\n", workflow.Enabled)
	fmt.Printf("  \"action\": \"%s\",\n", action)
	fmt.Printf("  \"triggerCount\": %d,\n", len(workflow.Triggers))
	fmt.Printf("  \"actionCount\": %d\n", len(workflow.Actions))
	fmt.Printf("}\n")

	return nil
}

// outputWorkflowDetails outputs common workflow details in table format
func outputWorkflowDetails(workflow *graphql.ProjectV2Workflow) {
	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)

	status := statusDisabled
	if workflow.Enabled {
		status = statusEnabled
	}
	fmt.Printf("  Status: %s\n", status)
	fmt.Printf("  Triggers: %d configured\n", len(workflow.Triggers))
	fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))
}

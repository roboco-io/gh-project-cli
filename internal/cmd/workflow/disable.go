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

// DisableOptions holds options for the disable command
type DisableOptions struct {
	WorkflowID string
	Format     string
}

// NewDisableCmd creates the disable command
func NewDisableCmd() *cobra.Command {
	opts := &DisableOptions{}

	cmd := &cobra.Command{
		Use:   "disable <workflow-id>",
		Short: "Disable a project workflow",
		Long: `Disable a project workflow to stop automation.

Disabled workflows will not respond to triggers or execute any actions.
The workflow configuration is preserved and can be re-enabled later.

Examples:
  ghp workflow disable workflow-id
  ghp workflow disable workflow-id --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runDisable(cmd.Context(), opts)
		},
	}

	return cmd
}

func runDisable(ctx context.Context, opts *DisableOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	// Disable workflow
	workflow, err := workflowService.DisableWorkflow(ctx, opts.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed to disable workflow: %w", err)
	}

	// Output result
	return outputDisabledWorkflow(workflow, opts.Format)
}

func outputDisabledWorkflow(workflow *graphql.ProjectV2Workflow, format string) error {
	switch format {
	case "json":
		return outputDisabledWorkflowJSON(workflow)
	case "table":
		return outputDisabledWorkflowTable(workflow)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputDisabledWorkflowTable(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("✅ Workflow '%s' disabled successfully\n\n", workflow.Name)

	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Status: Disabled\n")
	fmt.Printf("  Triggers: %d configured\n", len(workflow.Triggers))
	fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))

	return nil
}

func outputDisabledWorkflowJSON(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", workflow.ID)
	fmt.Printf("  \"name\": \"%s\",\n", workflow.Name)
	fmt.Printf("  \"enabled\": false,\n")
	fmt.Printf("  \"triggerCount\": %d,\n", len(workflow.Triggers))
	fmt.Printf("  \"actionCount\": %d\n", len(workflow.Actions))
	fmt.Printf("}\n")

	return nil
}
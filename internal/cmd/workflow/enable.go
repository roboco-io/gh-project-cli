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

// EnableOptions holds options for the enable command
type EnableOptions struct {
	WorkflowID string
	Format     string
}

// NewEnableCmd creates the enable command
func NewEnableCmd() *cobra.Command {
	opts := &EnableOptions{}

	cmd := &cobra.Command{
		Use:   "enable <workflow-id>",
		Short: "Enable a project workflow",
		Long: `Enable a project workflow to start automation.

Enabled workflows will respond to triggers and execute their configured
actions automatically.

Examples:
  ghp workflow enable workflow-id
  ghp workflow enable workflow-id --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runEnable(cmd.Context(), opts)
		},
	}

	return cmd
}

func runEnable(ctx context.Context, opts *EnableOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	// Enable workflow
	workflow, err := workflowService.EnableWorkflow(ctx, opts.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed to enable workflow: %w", err)
	}

	// Output result
	return outputEnabledWorkflow(workflow, opts.Format)
}

func outputEnabledWorkflow(workflow *graphql.ProjectV2Workflow, format string) error {
	switch format {
	case "json":
		return outputEnabledWorkflowJSON(workflow)
	case "table":
		return outputEnabledWorkflowTable(workflow)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputEnabledWorkflowTable(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("✅ Workflow '%s' enabled successfully\n\n", workflow.Name)

	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Status: Enabled\n")
	fmt.Printf("  Triggers: %d configured\n", len(workflow.Triggers))
	fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))

	return nil
}

func outputEnabledWorkflowJSON(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", workflow.ID)
	fmt.Printf("  \"name\": \"%s\",\n", workflow.Name)
	fmt.Printf("  \"enabled\": true,\n")
	fmt.Printf("  \"triggerCount\": %d,\n", len(workflow.Triggers))
	fmt.Printf("  \"actionCount\": %d\n", len(workflow.Actions))
	fmt.Printf("}\n")

	return nil
}
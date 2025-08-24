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

// UpdateOptions holds options for the update command
type UpdateOptions struct {
	WorkflowID string
	Name       string
	Enabled    *bool
	Format     string
}

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	opts := &UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <workflow-id>",
		Short: "Update a project workflow",
		Long: `Update properties of an existing project workflow.

You can update the workflow name and enable/disable status. At least one
property must be specified.

Examples:
  ghp workflow update workflow-id --name "Updated Auto-assign"
  ghp workflow update workflow-id --enable
  ghp workflow update workflow-id --disable
  ghp workflow update workflow-id --name "Status Sync" --enable --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			
			// Handle enable/disable flags
			if enable, _ := cmd.Flags().GetBool("enable"); enable {
				enabled := true
				opts.Enabled = &enabled
			}
			if disable, _ := cmd.Flags().GetBool("disable"); disable {
				enabled := false
				opts.Enabled = &enabled
			}
			
			return runUpdate(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the workflow")
	cmd.Flags().Bool("enable", false, "Enable the workflow")
	cmd.Flags().Bool("disable", false, "Disable the workflow")

	return cmd
}

func runUpdate(ctx context.Context, opts *UpdateOptions) error {
	// Validate at least one field is provided
	if opts.Name == "" && opts.Enabled == nil {
		return fmt.Errorf("at least one of --name, --enable, or --disable must be provided")
	}

	// Validate name if provided
	if opts.Name != "" {
		if err := service.ValidateWorkflowName(opts.Name); err != nil {
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
	input := service.UpdateWorkflowInput{
		WorkflowID: opts.WorkflowID,
	}

	if opts.Name != "" {
		input.Name = &opts.Name
	}
	if opts.Enabled != nil {
		input.Enabled = opts.Enabled
	}

	// Update workflow
	workflow, err := workflowService.UpdateWorkflow(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	// Output updated workflow
	return outputUpdatedWorkflow(workflow, opts.Format)
}

func outputUpdatedWorkflow(workflow *graphql.ProjectV2Workflow, format string) error {
	switch format {
	case "json":
		return outputUpdatedWorkflowJSON(workflow)
	case "table":
		return outputUpdatedWorkflowTable(workflow)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputUpdatedWorkflowTable(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("✅ Workflow '%s' updated successfully\n\n", workflow.Name)

	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	
	status := "Enabled"
	if !workflow.Enabled {
		status = "Disabled"
	}
	fmt.Printf("  Status: %s\n", status)

	if len(workflow.Triggers) > 0 {
		fmt.Printf("  Triggers: %d configured\n", len(workflow.Triggers))
	} else {
		fmt.Printf("  Triggers: None\n")
	}

	if len(workflow.Actions) > 0 {
		fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))
	} else {
		fmt.Printf("  Actions: None\n")
	}

	return nil
}

func outputUpdatedWorkflowJSON(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", workflow.ID)
	fmt.Printf("  \"name\": \"%s\",\n", workflow.Name)
	fmt.Printf("  \"enabled\": %t,\n", workflow.Enabled)
	fmt.Printf("  \"triggerCount\": %d,\n", len(workflow.Triggers))
	fmt.Printf("  \"actionCount\": %d\n", len(workflow.Actions))
	fmt.Printf("}\n")

	return nil
}
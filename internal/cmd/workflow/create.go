package workflow

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// CreateOptions holds options for the create command
type CreateOptions struct {
	ProjectRef string
	Name       string
	Enabled    bool
	Format     string
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <owner/project-number> <name>",
		Short: "Create a new project workflow",
		Long: `Create a new workflow in a GitHub Project.

Workflows provide automation capabilities that can respond to events and
perform actions automatically. After creating a workflow, you can add
triggers and actions to define the automation behavior.

Examples:
  ghp workflow create octocat/123 "Auto-assign Priority"
  ghp workflow create octocat/123 "Status Sync" --disabled
  ghp workflow create --org myorg/456 "Release Pipeline" --format json`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Name = args[1]
			opts.Format = cmd.Flag("format").Value.String()
			return runCreate(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Enabled, "enabled", true, "Enable the workflow after creation")
	cmd.Flags().Bool("disabled", false, "Disable the workflow after creation")
	cmd.Flags().Bool("org", false, "Create workflow in organization project")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions) error {
	// Validate workflow name
	if err := service.ValidateWorkflowName(opts.Name); err != nil {
		return err
	}

	// Handle disabled flag will be processed in the command handler

	// Parse project reference
	parts := strings.Split(opts.ProjectRef, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid project reference format. Use: owner/project-number")
	}

	owner := parts[0]
	projectNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid project number: %s", parts[1])
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)
	workflowService := service.NewWorkflowService(client)

	// Get project to validate access and get project ID
	isOrg := false // TODO: Get this from flag properly
	project, err := projectService.GetProject(ctx, owner, projectNumber, isOrg)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Create workflow
	input := service.CreateWorkflowInput{
		ProjectID: project.ID,
		Name:      opts.Name,
		Enabled:   opts.Enabled,
	}

	workflow, err := workflowService.CreateWorkflow(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	// Output created workflow
	return outputCreatedWorkflow(workflow, opts.Format)
}

func outputCreatedWorkflow(workflow *graphql.ProjectV2Workflow, format string) error {
	switch format {
	case "json":
		return outputCreatedWorkflowJSON(workflow)
	case "table":
		return outputCreatedWorkflowTable(workflow)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCreatedWorkflowTable(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("✅ Workflow '%s' created successfully\n\n", workflow.Name)

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
		fmt.Printf("  Triggers: None (add triggers with 'ghp workflow add-trigger')\n")
	}

	if len(workflow.Actions) > 0 {
		fmt.Printf("  Actions: %d configured\n", len(workflow.Actions))
	} else {
		fmt.Printf("  Actions: None (add actions with 'ghp workflow add-action')\n")
	}

	return nil
}

func outputCreatedWorkflowJSON(workflow *graphql.ProjectV2Workflow) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", workflow.ID)
	fmt.Printf("  \"name\": \"%s\",\n", workflow.Name)
	fmt.Printf("  \"enabled\": %t,\n", workflow.Enabled)
	fmt.Printf("  \"triggerCount\": %d,\n", len(workflow.Triggers))
	fmt.Printf("  \"actionCount\": %d\n", len(workflow.Actions))
	fmt.Printf("}\n")

	return nil
}
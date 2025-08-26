package workflow

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// DeleteOptions holds options for the delete command
type DeleteOptions struct {
	WorkflowID string
	Format     string
	Force      bool
}

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <workflow-id>",
		Short: "Delete a project workflow",
		Long: `Delete an existing project workflow.

This operation cannot be undone. By default, you will be prompted for
confirmation unless you use the --force flag.

WARNING: Deleting a workflow will remove all its triggers and actions.
The automation will stop working immediately.

Examples:
  ghp workflow delete workflow-id
  ghp workflow delete workflow-id --force
  ghp workflow delete workflow-id --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runDelete(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *DeleteOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	// Get workflow details for confirmation
	workflowInfo, err := workflowService.GetWorkflow(ctx, opts.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow details: %w", err)
	}

	// Confirm deletion unless force flag is used
	if !opts.Force {
		status := "enabled"
		if !workflowInfo.Enabled {
			status = "disabled"
		}

		fmt.Printf("Are you sure you want to delete workflow '%s' (%s)?\n", workflowInfo.Name, status)
		fmt.Printf("This will remove %d trigger(s) and %d action(s).\n", len(workflowInfo.Triggers), len(workflowInfo.Actions))
		fmt.Printf("This action cannot be undone. [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("failed to read confirmation: %w", readErr)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Deletion canceled.")
			return nil
		}
	}

	// Delete workflow
	input := service.DeleteWorkflowInput{
		WorkflowID: opts.WorkflowID,
	}

	err = workflowService.DeleteWorkflow(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	// Output confirmation
	return outputDeleteConfirmation(workflowInfo, opts.Format)
}

func outputDeleteConfirmation(workflowInfo *service.WorkflowInfo, format string) error {
	switch format {
	case formatJSON:
		return outputDeleteConfirmationJSON(workflowInfo)
	case formatTable:
		return outputDeleteConfirmationTable(workflowInfo)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputDeleteConfirmationTable(workflowInfo *service.WorkflowInfo) error {
	fmt.Printf("✅ Workflow '%s' deleted successfully\n", workflowInfo.Name)
	return nil
}

func outputDeleteConfirmationJSON(workflowInfo *service.WorkflowInfo) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"deleted\": true,\n")
	fmt.Printf("  \"workflowId\": \"%s\",\n", workflowInfo.ID)
	fmt.Printf("  \"workflowName\": \"%s\"\n", workflowInfo.Name)
	fmt.Printf("}\n")
	return nil
}

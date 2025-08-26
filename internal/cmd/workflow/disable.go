package workflow

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// NewDisableCmd creates the disable command
func NewDisableCmd() *cobra.Command {
	config := ToggleConfig{
		Use:   "disable <workflow-id>",
		Short: "Disable a project workflow",
		Long: `Disable a project workflow to stop automation.

Disabled workflows will not respond to triggers or execute any actions.
The workflow configuration is preserved and can be re-enabled later.

Examples:
  ghp workflow disable workflow-id
  ghp workflow disable workflow-id --format json`,

		Action: func(ctx context.Context, service *service.WorkflowService, workflowID string) (*graphql.ProjectV2Workflow, error) {
			return service.DisableWorkflow(ctx, workflowID)
		},
		Success: "disable",
	}

	return createWorkflowToggleCmd(config)
}

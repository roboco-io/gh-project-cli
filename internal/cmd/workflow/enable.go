package workflow

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// NewEnableCmd creates the enable command
func NewEnableCmd() *cobra.Command {
	config := ToggleConfig{
		Use:   "enable <workflow-id>",
		Short: "Enable a project workflow",
		Long: `Enable a project workflow to start automation.

Enabled workflows will respond to triggers and execute their configured
actions automatically.

Examples:
  ghp workflow enable workflow-id
  ghp workflow enable workflow-id --format json`,

		Action: func(ctx context.Context, service *service.WorkflowService, workflowID string) (*graphql.ProjectV2Workflow, error) {
			return service.EnableWorkflow(ctx, workflowID)
		},
		Success: "enable",
	}

	return createWorkflowToggleCmd(config)
}

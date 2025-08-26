package workflow

import (
	"github.com/spf13/cobra"
)

// NewWorkflowCmd creates the workflow command
func NewWorkflowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage project workflows and automation",
		Long: `Manage workflows and automation in GitHub Projects.

Workflows provide powerful automation capabilities that can respond to events
and perform actions automatically. You can create workflows that:

• React to item changes (added, updated, archived)
• Respond to field value changes 
• Execute on schedule (daily, weekly, monthly)
• Perform automatic field updates
• Move items between views and columns
• Send notifications and assign users
• Add comments to issues and pull requests

This command group provides comprehensive workflow management capabilities:

• List existing workflows in projects
• Create new workflows with triggers and actions
• Update workflow names and enable/disable status
• Delete workflows when no longer needed
• Enable and disable workflows independently
• Add triggers to respond to specific events
• Add actions to perform automated tasks

Trigger Types:
  item-added        - When items are added to the project
  item-updated      - When items are modified
  item-archived     - When items are archived
  field-changed     - When specific field values change
  status-changed    - When issue/PR status changes
  assignee-changed  - When assignee is modified
  scheduled         - Time-based triggers (daily, weekly, monthly)

Action Types:
  set-field         - Set field to specific value
  clear-field       - Clear field value
  move-to-column    - Move item to different column/view
  archive-item      - Archive the item
  add-to-project    - Add item to another project
  notify            - Send notification to users
  assign            - Assign user to item
  add-comment       - Add comment to issue/PR

Workflow Operations:
  list              - List all workflows in a project
  create            - Create a new workflow
  update            - Update workflow name or status
  delete            - Delete a workflow
  enable            - Enable a workflow
  disable           - Disable a workflow
  add-trigger       - Add trigger to workflow
  add-action        - Add action to workflow`,

		Example: `  # List all workflows in a project
  ghp workflow list octocat/123

  # Create a new workflow
  ghp workflow create octocat/123 "Auto-assign Priority"

  # Enable/disable workflows
  ghp workflow enable workflow-id
  ghp workflow disable workflow-id

  # Add trigger to workflow
  ghp workflow add-trigger workflow-id item-added --event issue-opened

  # Add action to workflow
  ghp workflow add-action workflow-id set-field --field priority-id --value "High"`,
	}

	// Add format flag to all subcommands
	cmd.PersistentFlags().String("format", "table", "Output format (table, json)")

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewEnableCmd())
	cmd.AddCommand(NewDisableCmd())
	cmd.AddCommand(NewAddTriggerCmd())
	cmd.AddCommand(NewAddActionCmd())

	return cmd
}

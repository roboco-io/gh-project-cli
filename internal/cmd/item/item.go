package item

import (
	"github.com/spf13/cobra"
)

// NewItemCmd creates the item command group
func NewItemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item <command>",
		Short: "Manage project items",
		Long: `Manage items in GitHub Projects.

Items are the core content of GitHub Projects - they can be existing issues,
pull requests, or draft issues created directly in the project.

This command group provides comprehensive item management capabilities:

• Add existing issues and pull requests to projects
• Create draft issues directly in projects
• List and search items across repositories
• View detailed item information
• Remove items from projects
• Update item field values

For more information about GitHub Projects, visit:
https://docs.github.com/en/issues/planning-and-tracking-with-projects`,
		Example: `  ghp item list octocat/Hello-World               # List items from repository
  ghp item add octocat/1 octocat/Hello-World#123  # Add issue to project
  ghp item view octocat/Hello-World#456           # View item details
  ghp item remove myorg/2 item-id --force         # Remove item from project
  ghp item add octocat/1 --draft --title "Task"   # Create draft issue`,
	}

	// Add subcommands
	cmd.AddCommand(NewAddCmd())
	cmd.AddCommand(NewEditCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewRemoveCmd())
	cmd.AddCommand(NewViewCmd())

	return cmd
}

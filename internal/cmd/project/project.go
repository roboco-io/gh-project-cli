package project

import (
	"github.com/spf13/cobra"
)

// NewProjectCmd creates the project command group
func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project <command>",
		Short: "Manage GitHub Projects",
		Long: `Work with GitHub Projects v2.

GitHub Projects are flexible project management tools integrated with your GitHub repositories.
This command group provides comprehensive project management capabilities including:

• List, view, create, edit, and delete projects
• Manage project items (issues, pull requests, draft issues)
• Configure custom fields and views
• Bulk operations and automation

For more information about GitHub Projects, visit:
https://docs.github.com/en/issues/planning-and-tracking-with-projects`,
		Example: `  ghp project list                    # List projects for authenticated user
  ghp project list octocat            # List projects for user octocat
  ghp project view octocat/123        # View project details
  ghp project create "My Project"     # Create a new project
  ghp project edit 123 --title "New"  # Edit project title
  ghp project delete 123 --force      # Delete a project`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewViewCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewEditCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewLinkCmd())
	cmd.AddCommand(NewExportCmd())
	cmd.AddCommand(NewImportCmd())
	cmd.AddCommand(NewWorkflowCmd())
	cmd.AddCommand(NewTemplateCmd())

	return cmd
}

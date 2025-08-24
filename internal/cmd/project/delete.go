package project

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// DeleteOptions holds options for the delete command
type DeleteOptions struct {
	Owner  string
	Number int
	Org    bool
	Force  bool
}

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete {<owner>/<number> | <number>}",
		Short: "Delete a project",
		Long: `Delete an existing project.

⚠️  WARNING: This action cannot be undone. All project data will be permanently deleted.

Examples:
  ghp project delete 123 --force           # Delete project 123 (with confirmation)
  ghp project delete octocat/123 --force   # Delete project owned by octocat
  ghp project delete myorg/456 --org --force  # Delete org project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *DeleteOptions, args []string) error {
	// Parse project reference
	projectRef := args[0]
	var err error

	if strings.Contains(projectRef, "/") {
		opts.Owner, opts.Number, err = service.ParseProjectReference(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project reference: %w", err)
		}
	} else {
		// Just a number, need to determine owner from context
		opts.Number, err = strconv.Atoi(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project number: %s", projectRef)
		}

		// For now, require owner to be specified
		return fmt.Errorf("owner must be specified in format owner/number")
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)

	// First, get the current project to obtain its ID and show details
	currentProject, err := projectService.GetProject(ctx, opts.Owner, opts.Number, opts.Org)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Show what will be deleted
	fmt.Printf("⚠️  You are about to delete the following project:\n\n")
	fmt.Printf("Project #%d: %s\n", currentProject.Number, currentProject.Title)
	fmt.Printf("Owner: %s (%s)\n", currentProject.Owner.Login, currentProject.Owner.Type)
	fmt.Printf("URL: %s\n", currentProject.URL)
	fmt.Printf("Items: %d\n", len(currentProject.Items.Nodes))
	fmt.Printf("Fields: %d\n", len(currentProject.Fields.Nodes))

	// Confirm deletion unless --force is used
	if !opts.Force {
		fmt.Printf("\n⚠️  This action cannot be undone. All project data will be permanently deleted.\n")
		fmt.Printf("Type 'DELETE' to confirm: ")

		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "DELETE" {
			fmt.Println("❌ Deletion canceled.")
			return nil
		}
	}

	// Delete project
	err = projectService.DeleteProject(ctx, currentProject.ID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	fmt.Printf("✅ Project #%d '%s' deleted successfully.\n", currentProject.Number, currentProject.Title)
	return nil
}

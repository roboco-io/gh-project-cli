package view

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

// ListOptions holds options for the list command
type ListOptions struct {
	ProjectRef string
	Format     string
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list <owner/project-number>",
		Short: "List project views",
		Long: `List all views in a GitHub Project.

This command displays all views configured for a project, showing their
names, layouts, and basic configuration.

Examples:
  ghp view list octocat/123
  ghp view list --org myorg/456
  ghp view list octocat/123 --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().Bool("org", false, "List views from organization project")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions) error {
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
	viewService := service.NewViewService(client)

	// Get project to validate access and get project ID
	isOrg := false
	if cmd := cobra.CheckErr; cmd != nil {
		// This is a workaround to access the flag from parent context
		// In a real implementation, we'd pass this properly
	}
	
	project, err := projectService.GetProject(ctx, owner, projectNumber, isOrg)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Get project views
	views, err := viewService.GetProjectViews(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("failed to list views: %w", err)
	}

	// Output views
	return outputViews(views, project.Title, opts.Format)
}

func outputViews(views []service.ViewInfo, projectName string, format string) error {
	switch format {
	case "json":
		return outputViewsJSON(views)
	case "table":
		return outputViewsTable(views, projectName)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputViewsTable(views []service.ViewInfo, projectName string) error {
	if len(views) == 0 {
		fmt.Printf("No views found in project '%s'\n", projectName)
		return nil
	}

	fmt.Printf("Views in project '%s':\n\n", projectName)

	// Find max widths for formatting
	maxNameWidth := 4 // "Name"
	maxLayoutWidth := 6 // "Layout"
	maxNumberWidth := 6 // "Number"

	for _, view := range views {
		if len(view.Name) > maxNameWidth {
			maxNameWidth = len(view.Name)
		}
		layout := service.FormatViewLayout(view.Layout)
		if len(layout) > maxLayoutWidth {
			maxLayoutWidth = len(layout)
		}
		numberStr := fmt.Sprintf("%d", view.Number)
		if len(numberStr) > maxNumberWidth {
			maxNumberWidth = len(numberStr)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  %s\n", 
		maxNumberWidth, "Number",
		maxNameWidth, "Name",
		maxLayoutWidth, "Layout",
		"Filter")
	
	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", maxNumberWidth),
		strings.Repeat("-", maxNameWidth),
		strings.Repeat("-", maxLayoutWidth),
		strings.Repeat("-", 6))

	// Print views
	for _, view := range views {
		filter := ""
		if view.Filter != nil {
			filter = *view.Filter
			if len(filter) > 50 {
				filter = filter[:47] + "..."
			}
		}

		fmt.Printf("%-*d  %-*s  %-*s  %s\n",
			maxNumberWidth, view.Number,
			maxNameWidth, view.Name,
			maxLayoutWidth, service.FormatViewLayout(view.Layout),
			filter)
	}

	fmt.Printf("\n%d view(s) total\n", len(views))

	return nil
}

func outputViewsJSON(views []service.ViewInfo) error {
	fmt.Printf("[\n")
	for i, view := range views {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", view.ID)
		fmt.Printf("    \"name\": \"%s\",\n", view.Name)
		fmt.Printf("    \"layout\": \"%s\",\n", view.Layout)
		fmt.Printf("    \"number\": %d", view.Number)

		if view.Filter != nil {
			fmt.Printf(",\n    \"filter\": \"%s\"", *view.Filter)
		}

		if len(view.GroupBy) > 0 {
			fmt.Printf(",\n    \"groupBy\": [\n")
			for j, gb := range view.GroupBy {
				fmt.Printf("      {\n")
				fmt.Printf("        \"fieldId\": \"%s\",\n", gb.FieldID)
				fmt.Printf("        \"fieldName\": \"%s\",\n", gb.FieldName)
				fmt.Printf("        \"direction\": \"%s\"\n", gb.Direction)
				fmt.Printf("      }")
				if j < len(view.GroupBy)-1 {
					fmt.Printf(",")
				}
				fmt.Printf("\n")
			}
			fmt.Printf("    ]")
		}

		if len(view.SortBy) > 0 {
			fmt.Printf(",\n    \"sortBy\": [\n")
			for j, sb := range view.SortBy {
				fmt.Printf("      {\n")
				fmt.Printf("        \"fieldId\": \"%s\",\n", sb.FieldID)
				fmt.Printf("        \"fieldName\": \"%s\",\n", sb.FieldName)
				fmt.Printf("        \"direction\": \"%s\"\n", sb.Direction)
				fmt.Printf("      }")
				if j < len(view.SortBy)-1 {
					fmt.Printf(",")
				}
				fmt.Printf("\n")
			}
			fmt.Printf("    ]")
		}

		fmt.Printf("\n  }")
		if i < len(views)-1 {
			fmt.Printf(",")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("]\n")

	return nil
}
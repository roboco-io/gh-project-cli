package view

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// GroupOptions holds options for the group command
type GroupOptions struct {
	ViewID    string
	FieldID   string
	Direction string
	Clear     bool
	Format    string
}

// NewGroupCmd creates the group command
func NewGroupCmd() *cobra.Command {
	opts := &GroupOptions{}

	cmd := &cobra.Command{
		Use:   "group <view-id>",
		Short: "Configure view grouping",
		Long: `Configure grouping for a project view.

You can set the field to group by and the group direction. Use --clear to
remove grouping from the view. Grouping is particularly useful for board
and roadmap views.

Group Directions:
  asc, ascending    - Group in ascending order (A-Z, 1-9, oldest first)
  desc, descending  - Group in descending order (Z-A, 9-1, newest first)

Examples:
  ghp view group view-id --field status-field-id --direction asc
  ghp view group view-id --field assignee-field-id --direction desc
  ghp view group view-id --clear
  ghp view group view-id --field priority-field-id --direction desc --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runGroup(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.FieldID, "field", "", "Field ID to group by")
	cmd.Flags().StringVar(&opts.Direction, "direction", "asc", "Group direction (asc, desc)")
	cmd.Flags().BoolVar(&opts.Clear, "clear", false, "Clear grouping from the view")

	return cmd
}

func runGroup(ctx context.Context, opts *GroupOptions) error {
	// Validate input
	if opts.Clear && opts.FieldID != "" {
		return fmt.Errorf("cannot use --clear with --field")
	}

	if !opts.Clear && opts.FieldID == "" {
		return fmt.Errorf("must specify --field or --clear")
	}

	// Validate direction if not clearing
	var direction graphql.ProjectV2ViewSortDirection
	if !opts.Clear {
		var err error
		direction, err = service.ValidateSortDirection(opts.Direction)
		if err != nil {
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
	viewService := service.NewViewService(client)

	// Update view group
	input := service.UpdateViewGroupInput{
		ViewID:    opts.ViewID,
		Direction: direction,
	}

	if !opts.Clear {
		input.GroupByID = &opts.FieldID
	}

	err = viewService.UpdateViewGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update view group: %w", err)
	}

	// Get updated view for output
	viewInfo, err := viewService.GetView(ctx, opts.ViewID)
	if err != nil {
		return fmt.Errorf("failed to get updated view: %w", err)
	}

	// Output result
	return outputGroupResult(viewInfo, opts.Clear, opts.Format)
}

func outputGroupResult(viewInfo *service.ViewInfo, cleared bool, format string) error {
	switch format {
	case "json":
		return outputGroupResultJSON(viewInfo, cleared)
	case "table":
		return outputGroupResultTable(viewInfo, cleared)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputGroupResultTable(viewInfo *service.ViewInfo, cleared bool) error {
	if cleared {
		fmt.Printf("✅ Grouping cleared from view '%s'\n", viewInfo.Name)
	} else {
		fmt.Printf("✅ View '%s' group configuration updated\n", viewInfo.Name)
	}

	fmt.Printf("\nView Details:\n")
	fmt.Printf("  ID: %s\n", viewInfo.ID)
	fmt.Printf("  Name: %s\n", viewInfo.Name)
	fmt.Printf("  Layout: %s\n", service.FormatViewLayout(viewInfo.Layout))

	if len(viewInfo.GroupBy) > 0 {
		fmt.Printf("  Group By:\n")
		for _, gb := range viewInfo.GroupBy {
			fmt.Printf("    - %s (%s)\n", gb.FieldName, service.FormatSortDirection(gb.Direction))
		}
	} else {
		fmt.Printf("  Group By: None\n")
	}

	return nil
}

func outputGroupResultJSON(viewInfo *service.ViewInfo, cleared bool) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"cleared\": %t,\n", cleared)
	fmt.Printf("  \"viewId\": \"%s\",\n", viewInfo.ID)
	fmt.Printf("  \"viewName\": \"%s\",\n", viewInfo.Name)
	fmt.Printf("  \"layout\": \"%s\"", viewInfo.Layout)

	if len(viewInfo.GroupBy) > 0 {
		fmt.Printf(",\n  \"groupBy\": [\n")
		for i, gb := range viewInfo.GroupBy {
			fmt.Printf("    {\n")
			fmt.Printf("      \"fieldId\": \"%s\",\n", gb.FieldID)
			fmt.Printf("      \"fieldName\": \"%s\",\n", gb.FieldName)
			fmt.Printf("      \"direction\": \"%s\"\n", gb.Direction)
			fmt.Printf("    }")
			if i < len(viewInfo.GroupBy)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]")
	}

	fmt.Printf("\n}\n")

	return nil
}
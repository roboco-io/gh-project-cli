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

// SortOptions holds options for the sort command
type SortOptions struct {
	ViewID    string
	FieldID   string
	Direction string
	Clear     bool
	Format    string
}

// NewSortCmd creates the sort command
func NewSortCmd() *cobra.Command {
	opts := &SortOptions{}

	cmd := &cobra.Command{
		Use:   "sort <view-id>",
		Short: "Configure view sorting",
		Long: `Configure sorting for a project view.

You can set the field to sort by and the sort direction. Use --clear to
remove sorting from the view.

Sort Directions:
  asc, ascending    - Sort in ascending order (A-Z, 1-9, oldest first)
  desc, descending  - Sort in descending order (Z-A, 9-1, newest first)

Examples:
  ghp view sort view-id --field priority-field-id --direction desc
  ghp view sort view-id --field status-field-id --direction asc
  ghp view sort view-id --clear
  ghp view sort view-id --field due-date-field-id --direction asc --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runSort(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.FieldID, "field", "", "Field ID to sort by")
	cmd.Flags().StringVar(&opts.Direction, "direction", "asc", "Sort direction (asc, desc)")
	cmd.Flags().BoolVar(&opts.Clear, "clear", false, "Clear sorting from the view")

	return cmd
}

func runSort(ctx context.Context, opts *SortOptions) error {
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

	// Update view sort
	input := service.UpdateViewSortInput{
		ViewID:    opts.ViewID,
		Direction: direction,
	}

	if !opts.Clear {
		input.SortByID = &opts.FieldID
	}

	err = viewService.UpdateViewSort(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update view sort: %w", err)
	}

	// Get updated view for output
	viewInfo, err := viewService.GetView(ctx, opts.ViewID)
	if err != nil {
		return fmt.Errorf("failed to get updated view: %w", err)
	}

	// Output result
	return outputSortResult(viewInfo, opts.Clear, opts.Format)
}

func outputSortResult(viewInfo *service.ViewInfo, cleared bool, format string) error {
	switch format {
	case "json":
		return outputSortResultJSON(viewInfo, cleared)
	case "table":
		return outputSortResultTable(viewInfo, cleared)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputSortResultTable(viewInfo *service.ViewInfo, cleared bool) error {
	if cleared {
		fmt.Printf("✅ Sorting cleared from view '%s'\n", viewInfo.Name)
	} else {
		fmt.Printf("✅ View '%s' sort configuration updated\n", viewInfo.Name)
	}

	fmt.Printf("\nView Details:\n")
	fmt.Printf("  ID: %s\n", viewInfo.ID)
	fmt.Printf("  Name: %s\n", viewInfo.Name)
	fmt.Printf("  Layout: %s\n", service.FormatViewLayout(viewInfo.Layout))

	if len(viewInfo.SortBy) > 0 {
		fmt.Printf("  Sort By:\n")
		for _, sb := range viewInfo.SortBy {
			fmt.Printf("    - %s (%s)\n", sb.FieldName, service.FormatSortDirection(sb.Direction))
		}
	} else {
		fmt.Printf("  Sort By: None\n")
	}

	return nil
}

func outputSortResultJSON(viewInfo *service.ViewInfo, cleared bool) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"cleared\": %t,\n", cleared)
	fmt.Printf("  \"viewId\": \"%s\",\n", viewInfo.ID)
	fmt.Printf("  \"viewName\": \"%s\",\n", viewInfo.Name)
	fmt.Printf("  \"layout\": \"%s\"", viewInfo.Layout)

	if len(viewInfo.SortBy) > 0 {
		fmt.Printf(",\n  \"sortBy\": [\n")
		for i, sb := range viewInfo.SortBy {
			fmt.Printf("    {\n")
			fmt.Printf("      \"fieldId\": \"%s\",\n", sb.FieldID)
			fmt.Printf("      \"fieldName\": \"%s\",\n", sb.FieldName)
			fmt.Printf("      \"direction\": \"%s\"\n", sb.Direction)
			fmt.Printf("    }")
			if i < len(viewInfo.SortBy)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]")
	}

	fmt.Printf("\n}\n")

	return nil
}
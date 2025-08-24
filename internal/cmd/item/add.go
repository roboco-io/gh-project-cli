package item

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// AddOptions holds options for the add command
type AddOptions struct {
	ProjectRef string
	ItemRef    string
	Draft      bool
	Title      string
	Body       string
	Format     string
}

// NewAddCmd creates the add command
func NewAddCmd() *cobra.Command {
	opts := &AddOptions{}

	cmd := &cobra.Command{
		Use:   "add <project> <item>",
		Short: "Add an item to a project",
		Long: `Add an existing issue, pull request, or create a draft issue in a project.

Item references can be in the following formats:
• owner/repo#123 (issue or PR reference)
• https://github.com/owner/repo/issues/123 (GitHub issue URL)
• https://github.com/owner/repo/pull/456 (GitHub PR URL)

Project references should be in owner/number format (e.g., octocat/1).

Examples:
  ghp item add octocat/1 octocat/Hello-World#123     # Add issue to project
  ghp item add myorg/2 myorg/repo#456 --format json  # Add PR with JSON output
  ghp item add octocat/1 --draft --title "New task"  # Create draft issue`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			if len(args) > 1 {
				opts.ItemRef = args[1]
			}
			return runAdd(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Draft, "draft", false, "Create a draft issue instead of adding existing item")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for draft issue (required when --draft is used)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for draft issue")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func runAdd(ctx context.Context, opts *AddOptions) error {
	// Validate options
	if opts.Draft {
		if opts.Title == "" {
			return fmt.Errorf("title is required when creating draft issue (use --title)")
		}
		if opts.ItemRef != "" {
			return fmt.Errorf("cannot specify both --draft and item reference")
		}
	} else {
		if opts.ItemRef == "" {
			return fmt.Errorf("item reference is required (or use --draft to create draft issue)")
		}
	}

	// Parse project reference
	projectOwner, projectNumber, err := service.ParseProjectReference(opts.ProjectRef)
	if err != nil {
		return fmt.Errorf("invalid project reference: %w", err)
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	itemService := service.NewItemService(client)
	projectService := service.NewProjectService(client)

	// Get project details to obtain project ID
	project, err := projectService.GetProject(ctx, projectOwner, projectNumber, false) // TODO: detect org vs user
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	if opts.Draft {
		// Create draft issue
		var body *string
		if opts.Body != "" {
			body = &opts.Body
		}

		item, err := itemService.CreateDraftIssue(ctx, project.ID, opts.Title, body)
		if err != nil {
			return fmt.Errorf("failed to create draft issue: %w", err)
		}

		fmt.Printf("✅ Draft issue created and added to project!\n\n")
		return outputAddedItem(item, opts.Format, "DraftIssue", opts.Title)
	} else {
		// Parse item reference and get item details
		itemOwner, itemRepo, itemNumber, err := service.ParseItemReference(opts.ItemRef)
		if err != nil {
			return fmt.Errorf("invalid item reference: %w", err)
		}

		// Try to get as issue first, then as PR
		var contentID string
		var itemType string
		var itemTitle string

		issue, err := itemService.GetIssue(ctx, itemOwner, itemRepo, itemNumber)
		if err == nil {
			contentID = issue.ID
			itemType = "Issue"
			itemTitle = issue.Title
		} else {
			// Try as pull request
			pr, err := itemService.GetPullRequest(ctx, itemOwner, itemRepo, itemNumber)
			if err != nil {
				return fmt.Errorf("failed to find issue or pull request: %w", err)
			}
			contentID = pr.ID
			itemType = "PullRequest"
			itemTitle = pr.Title
		}

		// Add item to project
		item, err := itemService.AddItemToProject(ctx, project.ID, contentID)
		if err != nil {
			return fmt.Errorf("failed to add item to project: %w", err)
		}

		fmt.Printf("✅ %s added to project!\n\n", itemType)
		return outputAddedItem(item, opts.Format, itemType, itemTitle)
	}
}

func outputAddedItem(item interface{}, format, itemType, title string) error {
	switch format {
	case "json":
		return outputAddedItemJSON(item)
	case "table":
		return outputAddedItemTable(itemType, title)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputAddedItemTable(itemType, title string) error {
	fmt.Printf("Type: %s\n", itemType)

	if title != "" {
		displayTitle := title
		if len(displayTitle) > 80 {
			displayTitle = displayTitle[:77] + "..."
		}
		fmt.Printf("Title: %s\n", displayTitle)
	}

	return nil
}

func outputAddedItemJSON(item interface{}) error {
	// In a real implementation, we'd properly serialize the item
	fmt.Printf("{\n")
	fmt.Printf("  \"status\": \"added\",\n")
	fmt.Printf("  \"item\": {\n")
	fmt.Printf("    \"id\": \"<item-id>\",\n")
	fmt.Printf("    \"type\": \"<item-type>\"\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
	return nil
}

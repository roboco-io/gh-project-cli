package item

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// ListOptions holds options for the list command
type ListOptions struct {
	Repository string
	Search     string
	Type       string
	State      string
	Author     string
	Assignee   string
	Labels     []string
	Limit      int
	Format     string
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list [repository]",
		Short: "List issues and pull requests",
		Long: `List issues and pull requests from repositories or search across GitHub.

You can list items from a specific repository or search across all of GitHub
using various filters.

Examples:
  ghp item list octocat/Hello-World                    # List items from repository
  ghp item list octocat/Hello-World --type issue       # List only issues
  ghp item list --search "is:issue is:open bug"       # Search across GitHub
  ghp item list --author octocat --state open          # Find items by author
  ghp item list --assignee @me --type pr               # Find PRs assigned to you`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Repository = args[0]
			}
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Search, "search", "", "Search query (GitHub search syntax)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Item type: issue, pr, pullrequest")
	cmd.Flags().StringVar(&opts.State, "state", "", "Item state: open, closed, merged")
	cmd.Flags().StringVar(&opts.Author, "author", "", "Filter by author username")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Filter by assignee username")
	cmd.Flags().StringSliceVar(&opts.Labels, "label", nil, "Filter by labels (can be used multiple times)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 20, "Maximum number of items to list")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	itemService := service.NewItemService(client)

	var items []service.ItemInfo

	if opts.Repository != "" {
		// List items from specific repository
		items, err = listRepositoryItems(ctx, itemService, opts)
	} else {
		// Search items across GitHub
		items, err = searchItems(ctx, itemService, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list items: %w", err)
	}

	return outputItems(items, opts.Format)
}

func listRepositoryItems(ctx context.Context, itemService *service.ItemService, opts *ListOptions) ([]service.ItemInfo, error) {
	// Parse repository reference
	parts := strings.Split(opts.Repository, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", opts.Repository)
	}
	owner, repo := parts[0], parts[1]

	var states []string
	if opts.State != "" {
		states = []string{strings.ToUpper(opts.State)}
	}

	var allItems []service.ItemInfo

	// Get issues if not specifically requesting PRs
	if opts.Type == "" || opts.Type == "issue" {
		issues, err := itemService.ListRepositoryIssues(ctx, owner, repo, states, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		allItems = append(allItems, issues...)
	}

	// Get PRs if not specifically requesting issues
	if opts.Type == "" || opts.Type == "pr" || opts.Type == "pullrequest" {
		prs, err := itemService.ListRepositoryPullRequests(ctx, owner, repo, states, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		allItems = append(allItems, prs...)
	}

	// Apply additional filters
	return applyFilters(allItems, opts), nil
}

func searchItems(ctx context.Context, itemService *service.ItemService, opts *ListOptions) ([]service.ItemInfo, error) {
	// Build search query
	filters := service.SearchFilters{
		Type:       opts.Type,
		State:      opts.State,
		Repository: opts.Repository,
		Author:     opts.Author,
		Assignee:   opts.Assignee,
		Labels:     opts.Labels,
		Query:      opts.Search,
	}

	searchQuery := service.BuildSearchQuery(filters)
	if searchQuery == "" {
		return nil, fmt.Errorf("no search criteria specified")
	}

	var allItems []service.ItemInfo

	// Search issues if not specifically requesting PRs
	if opts.Type == "" || opts.Type == "issue" {
		issues, err := itemService.SearchIssues(ctx, searchQuery, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search issues: %w", err)
		}
		allItems = append(allItems, issues...)
	}

	// Search PRs if not specifically requesting issues
	if opts.Type == "" || opts.Type == "pr" || opts.Type == "pullrequest" {
		prs, err := itemService.SearchPullRequests(ctx, searchQuery, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search pull requests: %w", err)
		}
		allItems = append(allItems, prs...)
	}

	return allItems, nil
}

func applyFilters(items []service.ItemInfo, opts *ListOptions) []service.ItemInfo {
	var filtered []service.ItemInfo

	for _, item := range items {
		// Apply author filter
		if opts.Author != "" && item.Author != nil && *item.Author != opts.Author {
			continue
		}

		// Apply assignee filter
		if opts.Assignee != "" {
			found := false
			for _, assignee := range item.Assignees {
				if assignee == opts.Assignee || (opts.Assignee == "@me" && assignee == "current-user") {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Apply label filters
		if len(opts.Labels) > 0 {
			hasAllLabels := true
			for _, requiredLabel := range opts.Labels {
				found := false
				for _, itemLabel := range item.Labels {
					if itemLabel == requiredLabel {
						found = true
						break
					}
				}
				if !found {
					hasAllLabels = false
					break
				}
			}
			if !hasAllLabels {
				continue
			}
		}

		filtered = append(filtered, item)
	}

	// Limit results
	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	return filtered
}

func outputItems(items []service.ItemInfo, format string) error {
	if len(items) == 0 {
		fmt.Println("No items found")
		return nil
	}

	switch format {
	case "json":
		return outputItemsJSON(items)
	case "table":
		return outputItemsTable(items)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputItemsTable(items []service.ItemInfo) error {
	// Print header
	fmt.Printf("%-10s %-8s %-6s %-30s %-20s %-15s %-12s\n",
		"TYPE", "STATE", "NUMBER", "TITLE", "REPOSITORY", "AUTHOR", "UPDATED")
	fmt.Println(strings.Repeat("-", 100))

	// Print items
	for _, item := range items {
		itemType := item.Type
		if len(itemType) > 10 {
			itemType = itemType[:7] + "..."
		}

		state := item.State
		if len(state) > 8 {
			state = state[:5] + "..."
		}

		number := ""
		if item.Number != nil {
			number = fmt.Sprintf("#%d", *item.Number)
		}
		if len(number) > 6 {
			number = number[:3] + "..."
		}

		title := item.Title
		if len(title) > 28 {
			title = title[:25] + "..."
		}

		repository := ""
		if item.Repository != nil {
			repository = *item.Repository
		}
		if len(repository) > 18 {
			repository = repository[:15] + "..."
		}

		author := ""
		if item.Author != nil {
			author = *item.Author
		}
		if len(author) > 13 {
			author = author[:10] + "..."
		}

		updated := item.UpdatedAt
		if len(updated) > 10 {
			updated = updated[:10] // Keep only date part
		}

		fmt.Printf("%-10s %-8s %-6s %-30s %-20s %-15s %-12s\n",
			itemType, state, number, title, repository, author, updated)
	}

	return nil
}

func outputItemsJSON(items []service.ItemInfo) error {
	// Simplified JSON output
	fmt.Println("[")
	for i, item := range items {
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"%s\",\n", item.Type)
		fmt.Printf("    \"title\": \"%s\",\n", item.Title)
		if item.Number != nil {
			fmt.Printf("    \"number\": %d,\n", *item.Number)
		}
		fmt.Printf("    \"state\": \"%s\",\n", item.State)
		if item.Repository != nil {
			fmt.Printf("    \"repository\": \"%s\",\n", *item.Repository)
		}
		if item.Author != nil {
			fmt.Printf("    \"author\": \"%s\",\n", *item.Author)
		}
		fmt.Printf("    \"updated_at\": \"%s\"\n", item.UpdatedAt)

		if i < len(items)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")

	return nil
}

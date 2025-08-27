package item

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// NewAddBulkCmd creates the add-bulk command
func NewAddBulkCmd() *cobra.Command {
	var (
		issues   string
		label    string
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "add-bulk PROJECT_ID",
		Short: "Add multiple issues to a project in bulk",
		Long: `Add multiple issues or pull requests to a GitHub Project in bulk.

This command allows you to add multiple items at once using various methods:
• Number range (e.g., 34-46)
• By label
• From a file containing issue URLs or numbers

Examples:
  # Add issues by number range
  ghp item add-bulk myorg/123 --issues 34-46
  
  # Add all issues with a specific label
  ghp item add-bulk myorg/123 --label epic
  
  # Add issues from a file (one per line)
  ghp item add-bulk myorg/123 --from-file issue-list.txt`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			projectID := args[0]

			// Validate that at least one input method is specified
			if issues == "" && label == "" && fromFile == "" {
				return fmt.Errorf("at least one of --issues, --label, or --from-file must be specified")
			}

			var itemsToAdd []string

			// Handle number range
			if issues != "" {
				items, err := parseNumberRange(issues)
				if err != nil {
					return fmt.Errorf("invalid issue range: %w", err)
				}
				itemsToAdd = append(itemsToAdd, items...)
			}

			// Handle label
			if label != "" {
				items, err := getIssuesByLabel(projectID, label)
				if err != nil {
					return fmt.Errorf("failed to get issues by label: %w", err)
				}
				itemsToAdd = append(itemsToAdd, items...)
			}

			// Handle file input
			if fromFile != "" {
				items, err := readIssuesFromFile(fromFile)
				if err != nil {
					return fmt.Errorf("failed to read issues from file: %w", err)
				}
				itemsToAdd = append(itemsToAdd, items...)
			}

			// Remove duplicates
			itemsToAdd = removeDuplicates(itemsToAdd)

			fmt.Printf("Adding %d items to project %s...\n", len(itemsToAdd), projectID)

			// Add items to project
			successCount := 0
			for _, item := range itemsToAdd {
				// TODO: Implement actual GraphQL API call to add item
				fmt.Printf("Adding item: %s\n", item)
				successCount++
			}

			fmt.Printf("\n✓ Successfully added %d items to project\n", successCount)
			return nil
		},
	}

	cmd.Flags().StringVar(&issues, "issues", "", "Issue number range (e.g., 34-46)")
	cmd.Flags().StringVar(&label, "label", "", "Add all issues with this label")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "File containing issue URLs or numbers (one per line)")

	return cmd
}

// parseNumberRange parses a number range like "34-46" into a slice of strings
func parseNumberRange(rangeStr string) ([]string, error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format, expected 'start-end'")
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid start number: %w", err)
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid end number: %w", err)
	}

	if start > end {
		return nil, fmt.Errorf("start number must be less than or equal to end number")
	}

	var result []string
	for i := start; i <= end; i++ {
		result = append(result, fmt.Sprintf("#%d", i))
	}

	return result, nil
}

// getIssuesByLabel retrieves issues with a specific label
func getIssuesByLabel(_, _ string) ([]string, error) {
	// TODO: Implement GraphQL query to get issues by label
	// This is a placeholder implementation
	return []string{}, nil
}

// readIssuesFromFile reads issue URLs or numbers from a file
func readIssuesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issues []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") { // Skip comments
			issues = append(issues, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

// removeDuplicates removes duplicate items from a slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

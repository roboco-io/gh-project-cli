package field

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
	ProjectRef string
	Owner      string
	Number     int
	Org        bool
	Format     string
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list <owner>/<number>",
		Short: "List project fields",
		Long: `List all custom fields in a GitHub Project.

This command shows all fields defined in the project along with their
data types and configuration. For single select fields, it also shows
the available options.

Examples:
  ghp field list octocat/123        # List fields in project 123
  ghp field list --org myorg/456    # List fields in org project 456
  ghp field list octocat/123 --format json  # JSON output`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions) error {
	// Parse project reference
	var err error
	if strings.Contains(opts.ProjectRef, "/") {
		opts.Owner, opts.Number, err = service.ParseProjectReference(opts.ProjectRef)
		if err != nil {
			return fmt.Errorf("invalid project reference: %w", err)
		}
	} else {
		return fmt.Errorf("project reference must be in format owner/number")
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)

	// Get project fields
	fields, err := fieldService.GetProjectFields(ctx, opts.Owner, opts.Number, opts.Org)
	if err != nil {
		return fmt.Errorf("failed to get project fields: %w", err)
	}

	// Output fields
	return outputFields(fields, opts.Format)
}

func outputFields(fields []service.FieldInfo, format string) error {
	switch format {
	case "json":
		return outputFieldsJSON(fields)
	case "table":
		return outputFieldsTable(fields)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputFieldsTable(fields []service.FieldInfo) error {
	if len(fields) == 0 {
		fmt.Println("No custom fields found")
		return nil
	}

	fmt.Printf("Fields in project '%s':\n\n", fields[0].ProjectName)
	fmt.Printf("%-20s %-15s %-10s %s\n", "NAME", "TYPE", "OPTIONS", "ID")
	fmt.Println(strings.Repeat("-", 70))

	for _, field := range fields {
		optionCount := len(field.Options)
		optionsStr := ""
		if optionCount > 0 {
			optionsStr = fmt.Sprintf("%d options", optionCount)
		}

		fmt.Printf("%-20s %-15s %-10s %s\n",
			truncate(field.Name, 18),
			service.FormatFieldDataType(field.DataType),
			optionsStr,
			field.ID)

		// Show options for single select fields
		if len(field.Options) > 0 && len(field.Options) <= 5 {
			fmt.Printf("%-20s   Options: ", "")
			optionNames := make([]string, len(field.Options))
			for i, option := range field.Options {
				optionNames[i] = fmt.Sprintf("%s (%s)", option.Name, service.FormatColor(option.Color))
			}
			fmt.Printf("%s\n", strings.Join(optionNames, ", "))
		}
	}

	return nil
}

func outputFieldsJSON(fields []service.FieldInfo) error {
	if len(fields) == 0 {
		fmt.Println("[]")
		return nil
	}

	fmt.Println("[")
	for i, field := range fields {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", field.ID)
		fmt.Printf("    \"name\": \"%s\",\n", field.Name)
		fmt.Printf("    \"dataType\": \"%s\",\n", field.DataType)
		fmt.Printf("    \"projectId\": \"%s\",\n", field.ProjectID)
		fmt.Printf("    \"projectName\": \"%s\",\n", field.ProjectName)
		
		if len(field.Options) > 0 {
			fmt.Printf("    \"options\": [\n")
			for j, option := range field.Options {
				fmt.Printf("      {\n")
				fmt.Printf("        \"id\": \"%s\",\n", option.ID)
				fmt.Printf("        \"name\": \"%s\",\n", option.Name)
				fmt.Printf("        \"color\": \"%s\"", option.Color)
				if option.Description != nil {
					fmt.Printf(",\n        \"description\": \"%s\"\n", *option.Description)
				} else {
					fmt.Printf("\n")
				}
				if j < len(field.Options)-1 {
					fmt.Printf("      },\n")
				} else {
					fmt.Printf("      }\n")
				}
			}
			fmt.Printf("    ]\n")
		} else {
			fmt.Printf("    \"options\": []\n")
		}

		if i < len(fields)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")

	return nil
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
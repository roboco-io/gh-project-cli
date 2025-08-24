package field

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/api/graphql"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// CreateOptions holds options for the create command
type CreateOptions struct {
	ProjectRef string
	Owner      string
	Number     int
	Org        bool
	Name       string
	FieldType  string
	Options    []string
	Format     string
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <owner>/<number> <name> <type>",
		Short: "Create a new project field",
		Long: `Create a new custom field in a GitHub Project.

Fields allow you to track additional metadata for your project items.
Different field types support different kinds of data:

Field Types:
  text         - Text field for arbitrary text input
  number       - Numeric field for numbers and calculations  
  date         - Date field for deadlines and milestones
  single_select - Single select field with predefined options
  iteration    - Iteration field for sprint/cycle planning

For single select fields, you can provide initial options using --options.

Examples:
  ghp field create octocat/123 "Priority" text
  ghp field create octocat/123 "Story Points" number
  ghp field create octocat/123 "Due Date" date
  ghp field create octocat/123 "Status" single_select --options "Todo,In Progress,Done"
  ghp field create --org myorg/456 "Sprint" iteration`,

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Name = args[1]
			opts.FieldType = args[2]
			opts.Format = cmd.Flag("format").Value.String()
			return runCreate(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")
	cmd.Flags().StringSliceVar(&opts.Options, "options", []string{}, "Options for single select field (comma-separated)")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions) error {
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

	// Validate field name
	if err := service.ValidateFieldName(opts.Name); err != nil {
		return err
	}

	// Validate field type
	dataType, err := service.ValidateFieldType(opts.FieldType)
	if err != nil {
		return err
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)
	projectService := service.NewProjectService(client)

	// Get project to verify it exists and get project ID
	project, err := projectService.GetProject(ctx, opts.Owner, opts.Number, opts.Org)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Create field
	input := service.CreateFieldInput{
		ProjectID:           project.ID,
		Name:                opts.Name,
		DataType:            dataType,
		SingleSelectOptions: opts.Options,
	}

	field, err := fieldService.CreateField(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create field: %w", err)
	}

	// Output created field
	return outputCreatedField(field, project.Title, opts.Format)
}

func outputCreatedField(field *graphql.ProjectV2Field, projectName, format string) error {
	switch format {
	case "json":
		return outputCreatedFieldJSON(field)
	case "table":
		return outputCreatedFieldTable(field, projectName)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCreatedFieldTable(field *graphql.ProjectV2Field, projectName string) error {
	fmt.Printf("✅ Field '%s' created successfully in project '%s'\n\n", field.Name, projectName)
	
	fmt.Printf("Field Details:\n")
	fmt.Printf("  ID: %s\n", field.ID)
	fmt.Printf("  Name: %s\n", field.Name)
	fmt.Printf("  Type: %s\n", service.FormatFieldDataType(field.DataType))

	// Show options if single select field
	if len(field.Options.Nodes) > 0 {
		fmt.Printf("  Options:\n")
		for _, option := range field.Options.Nodes {
			fmt.Printf("    • %s (%s)", option.Name, service.FormatColor(option.Color))
			if option.Description != nil && *option.Description != "" {
				fmt.Printf(" - %s", *option.Description)
			}
			fmt.Printf("\n")
		}
	}

	return nil
}

func outputCreatedFieldJSON(field *graphql.ProjectV2Field) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", field.ID)
	fmt.Printf("  \"name\": \"%s\",\n", field.Name)
	fmt.Printf("  \"dataType\": \"%s\"", field.DataType)

	if len(field.Options.Nodes) > 0 {
		fmt.Printf(",\n  \"options\": [\n")
		for i, option := range field.Options.Nodes {
			fmt.Printf("    {\n")
			fmt.Printf("      \"id\": \"%s\",\n", option.ID)
			fmt.Printf("      \"name\": \"%s\",\n", option.Name)
			fmt.Printf("      \"color\": \"%s\"", option.Color)
			if option.Description != nil {
				fmt.Printf(",\n      \"description\": \"%s\"", *option.Description)
			}
			fmt.Printf("\n    }")
			if i < len(field.Options.Nodes)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("}\n")

	return nil
}
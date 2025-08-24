package workflow

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
		Short: "List project workflows",
		Long: `List all workflows in a GitHub Project.

This command displays all workflows configured for a project, showing their
names, enabled status, triggers, and actions.

Examples:
  ghp workflow list octocat/123
  ghp workflow list --org myorg/456
  ghp workflow list octocat/123 --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().Bool("org", false, "List workflows from organization project")

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
	workflowService := service.NewWorkflowService(client)

	// Get project to validate access and get project ID
	isOrg := false
	// TODO: Get this from flag properly
	
	project, err := projectService.GetProject(ctx, owner, projectNumber, isOrg)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Get project workflows
	workflows, err := workflowService.GetProjectWorkflows(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	// Output workflows
	return outputWorkflows(workflows, project.Title, opts.Format)
}

func outputWorkflows(workflows []service.WorkflowInfo, projectName string, format string) error {
	switch format {
	case "json":
		return outputWorkflowsJSON(workflows)
	case "table":
		return outputWorkflowsTable(workflows, projectName)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputWorkflowsTable(workflows []service.WorkflowInfo, projectName string) error {
	if len(workflows) == 0 {
		fmt.Printf("No workflows found in project '%s'\n", projectName)
		return nil
	}

	fmt.Printf("Workflows in project '%s':\n\n", projectName)

	// Find max widths for formatting
	maxNameWidth := 4 // "Name"
	maxTriggersWidth := 8 // "Triggers"
	maxActionsWidth := 7 // "Actions"

	for _, workflow := range workflows {
		if len(workflow.Name) > maxNameWidth {
			maxNameWidth = len(workflow.Name)
		}
		
		triggersStr := formatTriggersSummary(workflow.Triggers)
		if len(triggersStr) > maxTriggersWidth {
			maxTriggersWidth = len(triggersStr)
		}
		
		actionsStr := formatActionsSummary(workflow.Actions)
		if len(actionsStr) > maxActionsWidth {
			maxActionsWidth = len(actionsStr)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-7s  %-*s  %-*s\n", 
		maxNameWidth, "Name",
		"Status",
		maxTriggersWidth, "Triggers",
		maxActionsWidth, "Actions")
	
	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", maxNameWidth),
		strings.Repeat("-", 7),
		strings.Repeat("-", maxTriggersWidth),
		strings.Repeat("-", maxActionsWidth))

	// Print workflows
	for _, workflow := range workflows {
		status := "Enabled"
		if !workflow.Enabled {
			status = "Disabled"
		}

		triggersStr := formatTriggersSummary(workflow.Triggers)
		if len(triggersStr) > 50 {
			triggersStr = triggersStr[:47] + "..."
		}

		actionsStr := formatActionsSummary(workflow.Actions)
		if len(actionsStr) > 50 {
			actionsStr = actionsStr[:47] + "..."
		}

		fmt.Printf("%-*s  %-7s  %-*s  %-*s\n",
			maxNameWidth, workflow.Name,
			status,
			maxTriggersWidth, triggersStr,
			maxActionsWidth, actionsStr)
	}

	fmt.Printf("\n%d workflow(s) total\n", len(workflows))

	return nil
}

func outputWorkflowsJSON(workflows []service.WorkflowInfo) error {
	fmt.Printf("[\n")
	for i, workflow := range workflows {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", workflow.ID)
		fmt.Printf("    \"name\": \"%s\",\n", workflow.Name)
		fmt.Printf("    \"enabled\": %t", workflow.Enabled)

		if len(workflow.Triggers) > 0 {
			fmt.Printf(",\n    \"triggers\": [\n")
			for j, trigger := range workflow.Triggers {
				fmt.Printf("      {\n")
				fmt.Printf("        \"id\": \"%s\",\n", trigger.ID)
				fmt.Printf("        \"type\": \"%s\"", trigger.Type)
				
				if trigger.Event != "" {
					fmt.Printf(",\n        \"event\": \"%s\"", trigger.Event)
				}
				if trigger.FieldName != nil {
					fmt.Printf(",\n        \"fieldName\": \"%s\"", *trigger.FieldName)
				}
				if trigger.Value != nil {
					fmt.Printf(",\n        \"value\": \"%s\"", *trigger.Value)
				}
				
				fmt.Printf("\n      }")
				if j < len(workflow.Triggers)-1 {
					fmt.Printf(",")
				}
				fmt.Printf("\n")
			}
			fmt.Printf("    ]")
		}

		if len(workflow.Actions) > 0 {
			fmt.Printf(",\n    \"actions\": [\n")
			for j, action := range workflow.Actions {
				fmt.Printf("      {\n")
				fmt.Printf("        \"id\": \"%s\",\n", action.ID)
				fmt.Printf("        \"type\": \"%s\"", action.Type)
				
				if action.FieldName != nil {
					fmt.Printf(",\n        \"fieldName\": \"%s\"", *action.FieldName)
				}
				if action.Value != nil {
					fmt.Printf(",\n        \"value\": \"%s\"", *action.Value)
				}
				if action.ViewName != nil {
					fmt.Printf(",\n        \"viewName\": \"%s\"", *action.ViewName)
				}
				if action.Column != nil {
					fmt.Printf(",\n        \"column\": \"%s\"", *action.Column)
				}
				if action.Message != nil {
					fmt.Printf(",\n        \"message\": \"%s\"", *action.Message)
				}
				
				fmt.Printf("\n      }")
				if j < len(workflow.Actions)-1 {
					fmt.Printf(",")
				}
				fmt.Printf("\n")
			}
			fmt.Printf("    ]")
		}

		fmt.Printf("\n  }")
		if i < len(workflows)-1 {
			fmt.Printf(",")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("]\n")

	return nil
}

func formatTriggersSummary(triggers []service.TriggerInfo) string {
	if len(triggers) == 0 {
		return "None"
	}

	var parts []string
	for _, trigger := range triggers {
		parts = append(parts, service.FormatTriggerType(trigger.Type))
	}

	if len(parts) > 2 {
		return fmt.Sprintf("%s, %s, +%d more", parts[0], parts[1], len(parts)-2)
	}

	return strings.Join(parts, ", ")
}

func formatActionsSummary(actions []service.ActionInfo) string {
	if len(actions) == 0 {
		return "None"
	}

	var parts []string
	for _, action := range actions {
		parts = append(parts, service.FormatActionType(action.Type))
	}

	if len(parts) > 2 {
		return fmt.Sprintf("%s, %s, +%d more", parts[0], parts[1], len(parts)-2)
	}

	return strings.Join(parts, ", ")
}
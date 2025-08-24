package analytics

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

// OverviewOptions holds options for the overview command
type OverviewOptions struct {
	ProjectRef string
	Format     string
}

// NewOverviewCmd creates the overview command
func NewOverviewCmd() *cobra.Command {
	opts := &OverviewOptions{}

	cmd := &cobra.Command{
		Use:   "overview <owner/project-number>",
		Short: "Generate project overview analytics",
		Long: `Generate comprehensive overview analytics for a GitHub Project.

The overview report includes:
• Project basic information (title, item count, field count, view count)
• Item distribution by status with counts and percentages
• Item distribution by assignee with workload analysis
• Item distribution by labels and milestones
• Basic velocity metrics and timeline information

This report provides a high-level view of project health and progress,
making it easy to understand the current state and identify potential
bottlenecks or areas that need attention.

Examples:
  ghp analytics overview octocat/123
  ghp analytics overview octocat/123 --format json
  ghp analytics overview --org myorg/456 --format table`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runOverview(cmd.Context(), opts)
		},
	}

	cmd.Flags().Bool("org", false, "Target organization project")

	return cmd
}

func runOverview(ctx context.Context, opts *OverviewOptions) error {
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
	analyticsService := service.NewAnalyticsService(client)

	// Get project to validate access and get project ID
	isOrg := false // TODO: Get this from flag properly
	project, err := projectService.GetProject(ctx, owner, projectNumber, isOrg)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Get analytics data
	analytics, err := analyticsService.GetProjectAnalytics(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("failed to get project analytics: %w", err)
	}

	// Format and output analytics
	analyticsInfo := service.FormatAnalytics(analytics)
	return outputOverview(analyticsInfo, opts.Format)
}

func outputOverview(analytics *service.AnalyticsInfo, format string) error {
	switch format {
	case "json":
		return outputOverviewJSON(analytics)
	case "table":
		return outputOverviewTable(analytics)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputOverviewTable(analytics *service.AnalyticsInfo) error {
	fmt.Printf("📊 Project Overview: %s\n\n", analytics.Title)

	// Basic Statistics
	fmt.Printf("Basic Statistics:\n")
	fmt.Printf("  Project ID: %s\n", analytics.ProjectID)
	fmt.Printf("  Total Items: %d\n", analytics.ItemCount)
	fmt.Printf("  Total Fields: %d\n", analytics.FieldCount)
	fmt.Printf("  Total Views: %d\n", analytics.ViewCount)

	// Status Distribution
	if len(analytics.StatusStats) > 0 {
		fmt.Printf("\n📈 Item Distribution by Status:\n")
		for _, stat := range analytics.StatusStats {
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * 100
			fmt.Printf("  %-20s %3d items (%.1f%%)\n", stat.Status, stat.Count, percentage)
		}
	}

	// Assignee Distribution
	if len(analytics.AssigneeStats) > 0 {
		fmt.Printf("\n👥 Item Distribution by Assignee:\n")
		for _, stat := range analytics.AssigneeStats {
			assignee := stat.Assignee
			if assignee == "" {
				assignee = "Unassigned"
			}
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * 100
			fmt.Printf("  %-20s %3d items (%.1f%%)\n", assignee, stat.Count, percentage)
		}
	}

	// Velocity Information
	if analytics.VelocityData != nil {
		fmt.Printf("\n⚡ Velocity Metrics (%s):\n", analytics.VelocityData.Period)
		fmt.Printf("  Completed Items: %d\n", analytics.VelocityData.CompletedItems)
		fmt.Printf("  Added Items: %d\n", analytics.VelocityData.AddedItems)
		fmt.Printf("  Closure Rate: %.1f%%\n", analytics.VelocityData.ClosureRate*100)
		fmt.Printf("  Average Lead Time: %s\n", analytics.VelocityData.LeadTime)
		fmt.Printf("  Average Cycle Time: %s\n", analytics.VelocityData.CycleTime)
	}

	// Timeline Information
	if analytics.TimelineData != nil {
		fmt.Printf("\n📅 Timeline Information:\n")
		if analytics.TimelineData.StartDate != nil {
			fmt.Printf("  Start Date: %s\n", *analytics.TimelineData.StartDate)
		}
		if analytics.TimelineData.EndDate != nil {
			fmt.Printf("  End Date: %s\n", *analytics.TimelineData.EndDate)
		}
		if analytics.TimelineData.Duration > 0 {
			fmt.Printf("  Duration: %d days\n", analytics.TimelineData.Duration)
		}
		fmt.Printf("  Milestones: %d\n", analytics.TimelineData.MilestoneCount)
		fmt.Printf("  Activities: %d\n", analytics.TimelineData.ActivityCount)
	}

	return nil
}

func outputOverviewJSON(analytics *service.AnalyticsInfo) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"projectId\": \"%s\",\n", analytics.ProjectID)
	fmt.Printf("  \"title\": \"%s\",\n", analytics.Title)
	fmt.Printf("  \"itemCount\": %d,\n", analytics.ItemCount)
	fmt.Printf("  \"fieldCount\": %d,\n", analytics.FieldCount)
	fmt.Printf("  \"viewCount\": %d,\n", analytics.ViewCount)

	// Status statistics
	if len(analytics.StatusStats) > 0 {
		fmt.Printf("  \"statusDistribution\": [\n")
		for i, stat := range analytics.StatusStats {
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * 100
			fmt.Printf("    {\n")
			fmt.Printf("      \"status\": \"%s\",\n", stat.Status)
			fmt.Printf("      \"count\": %d,\n", stat.Count)
			fmt.Printf("      \"percentage\": %.1f\n", percentage)
			if i < len(analytics.StatusStats)-1 {
				fmt.Printf("    },\n")
			} else {
				fmt.Printf("    }\n")
			}
		}
		fmt.Printf("  ],\n")
	}

	// Assignee statistics
	if len(analytics.AssigneeStats) > 0 {
		fmt.Printf("  \"assigneeDistribution\": [\n")
		for i, stat := range analytics.AssigneeStats {
			assignee := stat.Assignee
			if assignee == "" {
				assignee = "Unassigned"
			}
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * 100
			fmt.Printf("    {\n")
			fmt.Printf("      \"assignee\": \"%s\",\n", assignee)
			fmt.Printf("      \"count\": %d,\n", stat.Count)
			fmt.Printf("      \"percentage\": %.1f\n", percentage)
			if i < len(analytics.AssigneeStats)-1 {
				fmt.Printf("    },\n")
			} else {
				fmt.Printf("    }\n")
			}
		}
		fmt.Printf("  ],\n")
	}

	// Velocity data
	if analytics.VelocityData != nil {
		fmt.Printf("  \"velocity\": {\n")
		fmt.Printf("    \"period\": \"%s\",\n", analytics.VelocityData.Period)
		fmt.Printf("    \"completedItems\": %d,\n", analytics.VelocityData.CompletedItems)
		fmt.Printf("    \"addedItems\": %d,\n", analytics.VelocityData.AddedItems)
		fmt.Printf("    \"closureRate\": %.3f,\n", analytics.VelocityData.ClosureRate)
		fmt.Printf("    \"leadTime\": \"%s\",\n", analytics.VelocityData.LeadTime)
		fmt.Printf("    \"cycleTime\": \"%s\"\n", analytics.VelocityData.CycleTime)
		fmt.Printf("  },\n")
	}

	// Timeline data
	if analytics.TimelineData != nil {
		fmt.Printf("  \"timeline\": {\n")
		if analytics.TimelineData.StartDate != nil {
			fmt.Printf("    \"startDate\": \"%s\",\n", *analytics.TimelineData.StartDate)
		}
		if analytics.TimelineData.EndDate != nil {
			fmt.Printf("    \"endDate\": \"%s\",\n", *analytics.TimelineData.EndDate)
		}
		if analytics.TimelineData.Duration > 0 {
			fmt.Printf("    \"duration\": %d,\n", analytics.TimelineData.Duration)
		}
		fmt.Printf("    \"milestoneCount\": %d,\n", analytics.TimelineData.MilestoneCount)
		fmt.Printf("    \"activityCount\": %d\n", analytics.TimelineData.ActivityCount)
		fmt.Printf("  }\n")
	} else {
		// Remove trailing comma
		fmt.Printf("  \"timeline\": null\n")
	}

	fmt.Printf("}\n")
	return nil
}

package field

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// DeleteOptions holds options for the delete command
type DeleteOptions struct {
	FieldID string
	Force   bool
}

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <field-id>",
		Short: "Delete a project field",
		Long: `Delete a project field and all its data.

⚠️  WARNING: This action is irreversible. All field data for project items
will be permanently lost. Use with caution.

By default, this command will prompt for confirmation. Use --force to skip
the confirmation prompt.

Examples:
  ghp field delete field-id
  ghp field delete field-id --force`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FieldID = args[0]
			return runDelete(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *DeleteOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)

	// Show field information and confirm deletion unless --force is used
	if !opts.Force {
		fmt.Printf("⚠️  You are about to delete field: %s\n", opts.FieldID)
		fmt.Printf("\nThis action cannot be undone. All field data will be permanently lost.\n")
		fmt.Printf("Type 'DELETE' to confirm: ")

		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "DELETE" {
			fmt.Println("❌ Deletion canceled.")
			return nil
		}
	}

	// Delete field
	input := service.DeleteFieldInput{
		FieldID: opts.FieldID,
	}

	err = fieldService.DeleteField(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete field: %w", err)
	}

	fmt.Printf("✅ Field %s deleted successfully.\n", opts.FieldID)
	return nil
}
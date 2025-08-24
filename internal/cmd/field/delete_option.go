package field

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghp-cli/internal/api"
	"github.com/roboco-io/ghp-cli/internal/auth"
	"github.com/roboco-io/ghp-cli/internal/service"
)

// DeleteOptionOptions holds options for the delete-option command
type DeleteOptionOptions struct {
	OptionID string
	Force    bool
}

// NewDeleteOptionCmd creates the delete-option command
func NewDeleteOptionCmd() *cobra.Command {
	opts := &DeleteOptionOptions{}

	cmd := &cobra.Command{
		Use:   "delete-option <option-id>",
		Short: "Delete single select field option",
		Long: `Delete an option from a single select field.

⚠️  WARNING: This action is irreversible. Items that currently have this
option selected will lose their field value. Use with caution.

By default, this command will prompt for confirmation. Use --force to skip
the confirmation prompt.

Examples:
  ghp field delete-option option-id
  ghp field delete-option option-id --force`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.OptionID = args[0]
			return runDeleteOption(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDeleteOption(ctx context.Context, opts *DeleteOptionOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)

	// Show option information and confirm deletion unless --force is used
	if !opts.Force {
		fmt.Printf("⚠️  You are about to delete field option: %s\n", opts.OptionID)
		fmt.Printf("\nThis action cannot be undone. Items using this option will lose their field value.\n")
		fmt.Printf("Type 'DELETE' to confirm: ")

		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "DELETE" {
			fmt.Println("❌ Deletion canceled.")
			return nil
		}
	}

	// Delete field option
	input := service.DeleteFieldOptionInput{
		OptionID: opts.OptionID,
	}

	err = fieldService.DeleteFieldOption(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete field option: %w", err)
	}

	fmt.Printf("✅ Field option %s deleted successfully.\n", opts.OptionID)
	return nil
}
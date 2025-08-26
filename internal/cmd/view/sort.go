package view

import "github.com/spf13/cobra"

// NewSortCmd creates the sort command
func NewSortCmd() *cobra.Command {
	return createSortCmd()
}

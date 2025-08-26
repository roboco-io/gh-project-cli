package view

import "github.com/spf13/cobra"

// NewGroupCmd creates the group command
func NewGroupCmd() *cobra.Command {
	return createGroupCmd()
}

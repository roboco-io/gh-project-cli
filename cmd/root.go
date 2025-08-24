package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/roboco-io/ghp-cli/internal/cmd/auth"
	"github.com/roboco-io/ghp-cli/internal/cmd/field"
	"github.com/roboco-io/ghp-cli/internal/cmd/item"
	"github.com/roboco-io/ghp-cli/internal/cmd/project"
)

var (
	cfgFile string
	rootCmd *cobra.Command

	// Version information
	version   string
	commit    string
	buildTime string
)

func init() {
	rootCmd = NewRootCmd()
}

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghp",
		Short: "GitHub Projects CLI - Powerful command-line tool for GitHub Projects",
		Long: `ghp-cli is a powerful command-line interface for managing GitHub Projects.
		
It provides complete control over GitHub Projects features including:
- Project management (create, list, edit, delete)
- Item management (add, update, archive)
- Field management (create custom fields)
- View management (table, board, roadmap)
- Automation workflows
- Reporting and analytics
- Bulk operations

Example:
  ghp project list --org myorg
  ghp project view owner/123
  ghp project create "My Project"`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime),
	}

	// Add persistent flags
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ghp.yaml)")
	cmd.PersistentFlags().String("token", "", "GitHub Personal Access Token")
	cmd.PersistentFlags().String("org", "", "GitHub organization")
	cmd.PersistentFlags().String("user", "", "GitHub user")
	cmd.PersistentFlags().String("format", "table", "Output format (table, json, yaml)")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug output")
	cmd.PersistentFlags().Bool("no-cache", false, "Disable caching")

	// Bind flags to viper
	viper.BindPFlag("token", cmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("org", cmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("user", cmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("format", cmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("no-cache", cmd.PersistentFlags().Lookup("no-cache"))

	// Add subcommands
	cmd.AddCommand(auth.NewAuthCmd())
	cmd.AddCommand(field.NewFieldCmd())
	cmd.AddCommand(item.NewItemCmd())
	cmd.AddCommand(project.NewProjectCmd())

	// Initialize config
	cobra.OnInitialize(initConfig)

	return cmd
}

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(v, c, b string) {
	version = v
	commit = c
	buildTime = b
	if rootCmd != nil {
		rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime)
	}
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ghp" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ghp")
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("GHP")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Check for GitHub token in environment if not set
	if viper.GetString("token") == "" {
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			viper.Set("token", token)
		} else if token := os.Getenv("GH_TOKEN"); token != "" {
			viper.Set("token", token)
		}
	}
}

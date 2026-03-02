package cli

import (
	"github.com/nickpricks/ft/internal/config"
	"github.com/nickpricks/ft/internal/constants"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     constants.RootUse,
	Short:   constants.RootShort,
	Long:    constants.RootLong,
	Example: constants.RootExample,
	Version: constants.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.LoadOrInit()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags can be added here
}

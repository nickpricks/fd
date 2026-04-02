package cli

import (
	"github.com/nickpricks/ft/internal/config"
	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     constants.RootUse,
		Short:   constants.RootShort,
		Long:    constants.RootLong,
		Example: constants.RootExample,
		Version: constants.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "help" || cmd.CalledAs() == "help" {
				return nil
			}
			dir, err := config.LoadOrInit()
			if err != nil {
				return err
			}
			s, err := core.NewNoteStore(dir)
			if err != nil {
				return err
			}
			store = s
			return nil
		},
	}

	store *core.NoteStore // set by PersistentPreRunE, used by all subcommands
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags can be added here
}

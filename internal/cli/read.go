package cli

import (
	"fmt"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:     constants.ReadUse,
	Short:   constants.ReadShort,
	Long:    constants.ReadLong,
	Example: constants.ReadExample,
	Args:    cobra.ExactArgs(1),
	RunE:    runRead,
}

func runRead(cmd *cobra.Command, args []string) error {
	noteID, err := core.ParseNoteID(args[0])
	if err != nil {
		return err
	}
	content, err := store.Read(noteID.String())
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}

func init() {
	rootCmd.AddCommand(readCmd)
}

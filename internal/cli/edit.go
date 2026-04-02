package cli

import (
	"fmt"
	"strings"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     constants.EditUse,
	Short:   constants.EditShort,
	Long:    constants.EditLong,
	Example: constants.EditExample,
	Args:    cobra.MinimumNArgs(2),
	RunE:    runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
	noteID, err := core.ParseNoteID(args[0])
	if err != nil {
		return err
	}
	text := strings.Join(args[1:], " ")
	path, err := store.Edit(noteID.String(), text)
	if err != nil {
		return err
	}
	fmt.Printf(constants.LogNoteUpdated, path)
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}

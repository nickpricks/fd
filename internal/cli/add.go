package cli

import (
	"fmt"
	"strings"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/core"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     constants.AddUse,
	Short:   constants.AddShort,
	Long:    constants.AddLong,
	Example: constants.AddExample,
	Args:    cobra.MinimumNArgs(1),
	RunE:    runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	text := strings.Join(args, " ")
	path, err := core.Add(text)
	if err != nil {
		return err
	}
	fmt.Printf(constants.LogNoteCreated, path)
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}

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
	id := args[0]
	content, err := core.Read(id)
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}

func init() {
	rootCmd.AddCommand(readCmd)
}

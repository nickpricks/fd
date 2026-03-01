package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nickpricks/ft/internal/constants"
	"github.com/nickpricks/ft/internal/notes"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     constants.ListUse,
	Short:   constants.ListShort,
	Long:    constants.ListLong,
	Example: constants.ListExample,
	RunE:    runList,
}

func runList(cmd *cobra.Command, args []string) error {
	items, err := notes.List()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Println(constants.LogNoNotes)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tID\tSLUG\tPATH")
	for _, note := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", note.Date, note.ID, note.Slug, note.Path)
	}
	w.Flush()
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}

package facade

import (
	"os"

	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/gpt-cli/gpt/chat"
	"github.com/spf13/cobra"
)

// newVersionCmd returns cobra.Command instance for show sub-command
func newHistoryCmd(ui *rwi.RWI) *cobra.Command {
	historyCmd := &cobra.Command{
		Use:     "history",
		Aliases: []string{"hist", "h"},
		Short:   "Print chat history",
		Long:    "Print chat history.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// local options
			histPath, err := cmd.Flags().GetString("history-file")
			if err != nil {
				return debugPrint(ui, err)
			}
			userName, err := cmd.Flags().GetString("user-name")
			if err != nil {
				return debugPrint(ui, err)
			}
			assistantName, err := cmd.Flags().GetString("assistant-name")
			if err != nil {
				return debugPrint(ui, err)
			}

			file, err := os.Open(histPath)
			if err != nil {
				return debugPrint(ui, errs.Wrap(err, errs.WithContext("histPath", histPath)))
			}
			defer file.Close()

			// Output history
			if err := chat.OutputHistory(file, ui.Writer(), userName, assistantName); err != nil {
				return debugPrint(ui, err)
			}
			return nil
		},
	}
	historyCmd.Flags().StringP("history-file", "f", "", "Path of history file (JSON format)")
	historyCmd.Flags().StringP("user-name", "u", "", "User name (display name)")
	historyCmd.Flags().StringP("assistant-name", "a", "", "Assistant name (display name)")

	return historyCmd
}

package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/gpt-cli/gpt/chat"
	"github.com/spf13/cobra"
)

// newVersionCmd returns cobra.Command instance for show sub-command
func newInteractiveCmd(ui *rwi.RWI) *cobra.Command {
	interactiveCmd := &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Short:   "Interactive mode",
		Long:    "Interactive mode in chat.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			opts, err := getOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			prepPath, err := cmd.Flags().GetString("prepare-file")
			if err != nil {
				return debugPrint(ui, err)
			}
			savePath, err := cmd.Flags().GetString("output-file")
			if err != nil {
				return debugPrint(ui, err)
			}
			multiLine, err := cmd.Flags().GetBool("multi-line")
			if err != nil {
				return debugPrint(ui, err)
			}

			// create Chat context
			chatCtx, err := chat.New(opts.APIKey, opts.CacheDir, opts.Logger, prepPath, savePath)
			if err != nil {
				opts.Logger.Error().Interface("error", errs.Wrap(err)).Send()
				return debugPrint(ui, err)
			}

			// kicking interactive mode
			if multiLine {
				if err := chatCtx.InteractiveMulti(cmd.Context(), ui.Writer()); err != nil {
					return debugPrint(ui, err)
				}
			} else {
				if err := chatCtx.Interactive(cmd.Context(), ui.Writer()); err != nil {
					return debugPrint(ui, err)
				}
			}
			if len(chatCtx.SavePath()) > 0 {
				return ui.Outputln("\nsave to", chatCtx.SavePath())
			}
			return nil
		},
	}
	interactiveCmd.Flags().StringP("prepare-file", "p", "", "Path of prepare file (JSON format)")
	interactiveCmd.Flags().StringP("output-file", "o", "", "Path of save file (JSON format)")
	interactiveCmd.Flags().BoolP("multi-line", "m", false, "Editing with multi-line mode")

	return interactiveCmd
}

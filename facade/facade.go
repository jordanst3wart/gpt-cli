package facade

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/goark/errs"

	"github.com/goark/gocli/config"
	"github.com/goark/gocli/exitcode"
	"github.com/goark/gocli/rwi"
	"github.com/goark/gpt-cli/ecode"
	"github.com/goark/gpt-cli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//Name is application name
	Name = "gpt-cli"
	//Version is version for application
	Version = "dev-version"
)

var (
	debugFlag         bool
	cfgFile           string
	configFile        = "config"
	defaultConfigPath = config.Path(Name, configFile+".yaml")
)

// newRootCmd returns cobra.Command instance for root command
func newRootCmd(ui *rwi.RWI, args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   Name,
		Short: "CLI tool for GPT with OpenAI API",
		Long:  "CLI tool for GPT with OpenAI API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	// global options (binding)
	rootCmd.PersistentFlags().StringP("api-key", "", "", "OpenAI API key")
	rootCmd.PersistentFlags().StringP("log-dir", "", logger.DefaultLogDir(Name), "Directory for log files")
	rootCmd.PersistentFlags().StringP("log-level", "", "nop", fmt.Sprintf("Log level [%s]", strings.Join(logger.LevelList(), "|")))

	//Bind config file
	_ = viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	_ = viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	cobra.OnInitialize(initConfig)

	// global options (other)
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "", false, "for debug")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Config file (default %v)", defaultConfigPath))

	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetArgs(args)
	rootCmd.SetIn(ui.Reader())       //Stdin
	rootCmd.SetOut(ui.ErrorWriter()) //Stdout -> Stderr
	rootCmd.SetErr(ui.ErrorWriter()) //Stderr
	rootCmd.AddCommand(
		newVersionCmd(ui),
		newChatCmd(ui),
	)

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory.
		confDir := config.Dir(Name)
		if len(confDir) == 0 {
			confDir = "."
		}
		// Search config in home directory with name "config.yaml" (without extension).
		viper.AddConfigPath(confDir)
		viper.SetConfigName(configFile)
	}
	viper.AutomaticEnv()     // read in environment variables that match
	_ = viper.ReadInConfig() // If a config file is found, read it in.
}

// Execute is called from main function
func Execute(ui *rwi.RWI, args []string) (exit exitcode.ExitCode) {
	defer func() {
		//panic hundling
		if r := recover(); r != nil {
			_ = ui.OutputErrln("Panic:", r)
			for depth := 0; ; depth++ {
				pc, src, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				_ = ui.OutputErrln(" ->", depth, ":", runtime.FuncForPC(pc).Name(), ":", src, ":", line)
			}
			exit = exitcode.Abnormal
		}
	}()

	//execution
	exit = exitcode.Normal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := newRootCmd(ui, args).ExecuteContext(ctx); err != nil {
		exit = exitcode.Abnormal
	}
	return
}

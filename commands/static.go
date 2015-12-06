package commands

import (
	"fmt"
	"os"

	"github.com/michaeldwan/static/printer"
	"github.com/michaeldwan/static/staticlib"
	"github.com/spf13/cobra"
)

func Execute() {
	// StaticCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", context.ConfigFileName, "path to a configuration file")
	// StaticCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	// StaticCmd.PersistentFlags().StringVarP(&flags.AWSAccessKeyId, "aws-access-key-id", "", "", "AWS Access Key ID")
	// StaticCmd.PersistentFlags().StringVarP(&flags.AWSSecretAccessKey, "aws-secret-access-key", "", "", "AWS Secret Access Key")
	// StaticCmd.PersistentFlags().StringVarP(&flags.AWSSessionToken, "aws-session-token", "", "", "AWS Session Token")

	if err := staticCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var (
	configFilePath string
	verboseOutput  bool
)

var staticCmd = &cobra.Command{
	Use:   "Static",
	Short: "Static is a simple tool for managing static content on AWS",
	Long:  `Static is a simple tool for managing static content on AWS.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verboseOutput {
			printer.SetLevel(printer.LevelDebug)
		}
		return
	},
}

func init() {
	staticCmd.PersistentFlags().StringVarP(&configFilePath, "config", "", staticlib.ConfigFileName, "path to a configuration file")
	staticCmd.PersistentFlags().BoolVarP(&verboseOutput, "verbose", "v", false, "verbose output")
}

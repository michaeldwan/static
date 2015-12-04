package main

import (
	"github.com/michaeldwan/static/context"
	"github.com/michaeldwan/static/printer"
	"github.com/michaeldwan/static/push"
	"github.com/spf13/cobra"
)

var configFilePath string
var ctx *context.Context

// var Verbose bool
var flags = context.Flags{}

var StaticCmd = &cobra.Command{
	Use:   "Static",
	Short: "Static is a simple tool for managing static content on AWS",
	Long:  `Static is a simple tool for managing static content on AWS.`,
	Run: func(cmd *cobra.Command, args []string) {
		printer.Infoln("Hello from static!")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = context.New(configFilePath)
		ctx.Flags = flags
		if ctx.Flags.Verbose {
			printer.SetLevel(printer.LevelDebug)
		}
		return
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		ctx.Clean()
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to AWS",
	Long:  `Static will upload changed files to S3, and more..`,
	Run: func(cmd *cobra.Command, args []string) {
		push.Perform(ctx)
	},
}

// var diffCmd = &cobra.Command{
// 	Use: "diff",
// 	Short: "Print diff",
// 	Long: "Print diff long...",
// 	Run: func(cmd *cobra.Command, args []string) {
//
// 	}
// }

func main() {
	StaticCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", context.ConfigFileName, "path to a configuration file")
	StaticCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	pushCmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "n", false, "simulate operations")
	pushCmd.Flags().BoolVarP(&flags.Force, "force", "f", false, "push everything")
	pushCmd.Flags().IntVarP(&flags.Concurrency, "parallel-uploads", "p", 2, "parallel uploads")
	StaticCmd.PersistentFlags().StringVarP(&flags.AWSAccessKeyId, "aws-access-key-id", "", "", "AWS Access Key ID")
	StaticCmd.PersistentFlags().StringVarP(&flags.AWSSecretAccessKey, "aws-secret-access-key", "", "", "AWS Secret Access Key")
	StaticCmd.PersistentFlags().StringVarP(&flags.AWSSessionToken, "aws-session-token", "", "", "AWS Session Token")
	StaticCmd.AddCommand(pushCmd)
	StaticCmd.Execute()
}

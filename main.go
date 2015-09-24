package main

import (
	"fmt"

	"github.com/michaeldwan/webmaster/context"
	"github.com/michaeldwan/webmaster/push"
	"github.com/spf13/cobra"
)

var configFilePath string
var ctx *context.Context
var logger Logger

// var Verbose bool
var flags = context.Flags{}

var WebmasterCmd = &cobra.Command{
	Use:   "webmaster",
	Short: "Webmaster is a simple tool for managing static sites on AWS",
	Long:  `Webmaster is a simple tool for managing static sites on AWS.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from webmaster!")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = context.New(configFilePath)
		ctx.Flags = flags
		logger = newLogger(flags.Verbose)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		ctx.Clean()
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to AWS",
	Long:  `Webmaster will upload changed files to S3, and more..`,
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
	WebmasterCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", context.ConfigFileName, "path to a configuration file")
	WebmasterCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	pushCmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "n", false, "simulate operations")
	pushCmd.Flags().BoolVarP(&flags.Force, "force", "f", false, "push everything")
	pushCmd.Flags().IntVarP(&flags.Concurrency, "parallel-uploads", "p", 2, "parallel uploads")
	WebmasterCmd.AddCommand(pushCmd)
	WebmasterCmd.Execute()
}

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"webmaster/context"
	"webmaster/push"
)

var configFilePath string
var ctx *context.Context
// var Verbose bool
var dryRun bool

var WebmasterCmd = &cobra.Command{
	Use:   "webmaster",
	Short: "Webmaster is a simple tool for managing static sites on AWS",
	Long:  `Webmaster is a simple tool for managing static sites on AWS.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from webmaster!")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = context.New(configFilePath)
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

func main() {
	WebmasterCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", context.ConfigFileName, "path to a configuration file")
	// WebmasterCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	pushCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "simulate operations")
	pushCmd.Flags().BoolVarP(&dryRun, "force", "f", false, "push everything")
	WebmasterCmd.AddCommand(pushCmd)
	WebmasterCmd.Execute()
}

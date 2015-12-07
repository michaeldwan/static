package commands

import (
	"fmt"

	"github.com/michaeldwan/static/staticlib"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a static.yml config file",
	Long:  `Creates a new static.yml file with placeholders and documentation`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := staticlib.WriteConfigFile(configFilePath); err != nil {
			panic(err)
		}
		fmt.Println("Initialized a static deployment configuration at", configFilePath)
	},
}

func init() {
	staticCmd.AddCommand(initCmd)
}

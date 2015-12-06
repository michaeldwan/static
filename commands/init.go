package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new static deployment",
	Long:  `Creates a new static.yml file`,
	Run:   push,
}

func Init(cmd *cobra.Command, args []string) {
	fmt.Println("Initialize a new static deployment")
}

func init() {
	staticCmd.AddCommand(initCmd)
}

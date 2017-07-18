package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.0.1")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

}

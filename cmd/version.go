package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of aas-ctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aas-ctl " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

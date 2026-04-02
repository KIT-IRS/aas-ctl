package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// getCmd represents the "get" command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Execute a GET request on the given endpoint.",
	Long: `Execute a GET request on the given endpoint.
Returns the raw json response.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.RequireSingleArg(args)
		if err != nil {
			log.Fatal(err)
		}
		response, err := utils.GetRequest(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(response))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

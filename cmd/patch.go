package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// patchCmd represents the "patch" command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Execute a PATCH request on the given endpoint with the provided data as request body.",
	Long: `patch <endpoint> <data>
Execute a PATCH request on the given endpoint with the provided data as request body.
<endpoint> is the target URL where the PUT request will be sent.
<data> is the data to be sent in the request body in JSON format.
Allows editing of submodel element values.
Returns the raw json response.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.RequireMinArgs(args, 2)
		if err != nil {
			log.Fatal(err)
		}
		endpoint := args[0]
		response, err := utils.PatchRequest(endpoint, []byte(strconv.Quote(args[1])))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(response))
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)
}

package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// putCmd represents the "put" command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Execute a PUT request on the given endpoint with the provided data as request body.",
	Long: `put <endpoint> [<data>]
Execute a PUT request on the given endpoint with the provided data as request body.
<endpoint> is the target URL where the PUT request will be sent.
<data> is the data to be sent in the request body in JSON format.
If <data> is not provided, the command reads from standard input, which allows for piping data into the command.
Allows editing and setting of shell, submodel and submodel element attributes.
Returns the raw json response.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.RequireMinArgs(args, 1)
		if err != nil {
			log.Fatal(err)
		}
		endpoint := args[0]
		var body []byte
		if len(args) < 2 {
			body, _ = io.ReadAll(os.Stdin)
		} else {
			body = []byte(args[1])
		}
		response, err := utils.PutRequest(endpoint, body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(response))
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}

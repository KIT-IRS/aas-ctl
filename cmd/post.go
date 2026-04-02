package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// postCmd represents the "post" command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Execute a POST request on the given endpoint.",
	Long: `post <endpoint> [<data>]
Execute a POST request on the given endpoint.
<endpoint> is the target URL where the POST request will be sent.
<data> is the data to be sent in the request body in JSON format.
If <data> is not provided, the command reads from standard input, which allows for piping data into the command.
Allows creating submodels and submodel elements.
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
		response, err := utils.PostRequest(endpoint, body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(response))
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
}

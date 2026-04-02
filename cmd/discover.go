package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover the contents of the repositories",
	Long: `discover <identifiable> [args]
Discover the entries of the repositories (repository and sm-repository).
First argument (<identifiable>) needs to be an ID or an IDShort of a Shell or SM.
Navigate through the repository like a file explorer using the following args.`,
	Run: discover,
}

func init() {
	rootCmd.AddCommand(discoverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	discoverCmd.Flags().Bool("url", false, "Only return the url of the requested element")
	discoverCmd.Flags().Bool("json", false, "Return the whole element json")
	discoverCmd.MarkFlagsMutuallyExclusive("url", "json")
}

// discover funcion executes the discover command
func discover(cmd *cobra.Command, args []string) {
	err := utils.RequireMinArgs(args, 1)
	if err != nil {
		log.Fatal(err)
	}
	flags, err := utils.NewFlagsDiscoverFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	endpoint, err := utils.ResolveDiscovery(args)
	if err != nil {
		log.Fatal(err)
	}
	if flags.OnlyURL {
		fmt.Println(endpoint)
		return
	}
	if flags.OnlyJSON {
		elementRaw, err := utils.GetRequest(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(elementRaw))
	} else {
		discoveryTree, err := utils.BuildDiscoveryTree(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		for _, branch := range discoveryTree {
			fmt.Println(branch)
		}
	}
}

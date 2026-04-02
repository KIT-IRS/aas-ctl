package cmd

import (
	"aas-ctl/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Shells in the repository that match a given pattern",
	Long: `Returns a list of shells that match the pattern provided by the flags.

Shells can be filtered for submodels with a given IDShort (--sm)
[AND the submodel contains a submodel element with a given IDShort (--elementId or --elementIdx)
[AND the value of the submodel element equals a specified value (--value)]]
`,
	Run: search,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	searchCmd.Flags().String("sm", "", "IDShort of the Submodel the Shells must contain")
	searchCmd.Flags().String("elementId", "", "IDShort of the SubmodelElement the Submodel must contain")
	searchCmd.Flags().Int("elementIdx", -1, "Index of the SubmodelElement the Submodel must contain")
	searchCmd.MarkFlagsMutuallyExclusive("elementId", "elementIdx")
	searchCmd.Flags().String("value", "", "Value the value of the SubmodelElement must match")

	// Output format flags, see aas command
	searchCmd.Flags().Bool("id", false, "Only return the <id>(s)")
	searchCmd.Flags().Bool("url", false, "Only return the URL(s) of the AAS(s)")
	searchCmd.Flags().Bool("json", false, "Return the whole AAS JSON")
	searchCmd.MarkFlagsMutuallyExclusive("id", "url", "json")
}

// search function executes the search command
// It extracts the flags from the command and creates a Filter from it.
// The Filter is then applied on the list of all shells to filter for the ones matching the criteria.
func search(cmd *cobra.Command, args []string) {
	flags, err := utils.NewFlagsSearchFromCMD(cmd)
	if err != nil {
		log.Fatal(err)
	}
	filter := utils.SearchFilterFromFlags(flags)
	shells, err := utils.GetAllShells()
	if err != nil {
		log.Fatal(err)
	}
	shells, err = filter.Apply(shells)
	if err != nil {
		log.Fatal(err)
	}
	for _, shell := range shells {
		utils.PrintIdentifiable(shell, flags.Flags, false)
		fmt.Println()
	}
}

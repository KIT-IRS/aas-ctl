package cmd

import (
	"github.com/spf13/cobra"
)

// smCmd represents the sm command
var smCmd = &cobra.Command{
	Use:   "sm",
	Short: "Explore Submodels in the selected repository",
	Long:  "Explore Submodels in the selected repository",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("sm called")
	// },
}

func init() {
	rootCmd.AddCommand(smCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// smCmd.PersistentFlags().String("foo", "", "A help for foo")
	smCmd.PersistentFlags().String("aas", "", `Shows only Submodel(s) of the Shell with the provided <Identifier>; useful if there are multiple Submodels with the same <idShort>
If the flag is provided the submodel <Identifier> must be an <idShort>`)
	smCmd.PersistentFlags().Bool("id", false, "Only return the <id>(s)")
	smCmd.PersistentFlags().Bool("url", false, "Only return the URL(s) of the Submodel(s)")
	smCmd.PersistentFlags().Bool("json", false, "Return the whole element JSON")
	smCmd.MarkFlagsMutuallyExclusive("id", "json")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// smCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

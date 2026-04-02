package cmd

import (
	"github.com/spf13/cobra"
)

// aasCmd represents the aas command
var aasCmd = &cobra.Command{
	Use:   "aas",
	Short: "Explore AssetAdministrationShells in the selected repository",
	Long:  "Explore AssetAdministrationShells in the selected repository",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("aas called")
	// },
}

func init() {
	rootCmd.AddCommand(aasCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aasCmd.PersistentFlags().String("foo", "", "A help for foo")
	aasCmd.PersistentFlags().Bool("id", false, "Only return the <id>(s)")
	aasCmd.PersistentFlags().Bool("url", false, "Only return the URL(s) of the AAS(s)")
	aasCmd.PersistentFlags().Bool("json", false, "Return the whole AAS JSON")
	aasCmd.MarkFlagsMutuallyExclusive("id", "url", "json")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

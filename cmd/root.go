/*
Package cmd implements all executable commands of the module.
The commands are

	config [aas-ctl/cmd/config] to show and select the repository
	aas [aas-ctl/cmd/aas] to access the Shells in the selected repository
	sm [aas-ctl/cmd/sm] to access the Submodels in the selected repository
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aas-ctl",
	Short: "CLI-Tool to explore AAS-Repositories",
	Long: `aas-ctl offers tools to explore AAS repositories
	
config: Show and select different repositories
aas: 	Explore AssetAdministrationShells in the selected repository
sm: 	Explore Submodels in the selected repository
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aas/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

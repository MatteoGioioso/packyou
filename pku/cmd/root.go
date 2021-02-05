package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "packyou",
	Short: "Package javascript files for serverless application without bundling",
	Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		initializeCommand(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	var Entry string
	var ProjectRoot string
	var Output string
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringVarP(&Entry, "entry", "e", "", "file entry point (required)")
	rootCmd.MarkFlagRequired("entry")
	rootCmd.Flags().StringVarP(&ProjectRoot, "project-root", "r", "", "project root folder")
	rootCmd.Flags().StringVarP(&Output, "output", "o", "", "output path")
	rootCmd.Flags().BoolP("compile-commonjs", "c", false, "compile ESM to commonjs modules")
	rootCmd.Flags().BoolP("add-extension", "x", false, "add .js extension to the import path in case you use ESM")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pku.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".packyou" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".packyou")
	}

	viper.Set("Hello", "world")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

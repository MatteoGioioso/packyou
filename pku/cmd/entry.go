package cmd

import (
	"github.com/spf13/cobra"
	"packyou/pku/fileCollector"
)

func initializeCommand(cmd *cobra.Command, getConfig func(key string) interface{}) {
	entry := cmd.Flag("entry").Value.String()
	projectRoot := cmd.Flag("projectRoot").Value.String()
	output := cmd.Flag("output").Value.String()

	fileCollector.New(entry, projectRoot, output).Collect()
}

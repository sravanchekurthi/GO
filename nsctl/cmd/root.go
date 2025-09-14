package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nsctl",
	Short: "nsctl is a CLI to manage Kubernetes namespaces",
	Long:  `nsctl is a command-line tool to manage Kubernetes namespaces using client-go.`,
}

// Execute starts the root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
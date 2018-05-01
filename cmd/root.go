package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "kubetransform",
	Short: "kubetrainsform transforms various kubernetes manifests",
	Long:  "kubetrainsform transforms various kubernetes manifests",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

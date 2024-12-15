package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cli/tcp"
)

var rootCmd = &cobra.Command{
	Short: "Socks is a very fast tcp server",
	Run: func(cmd *cobra.Command, args []string) {
		tcp.StartServer()
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start tcp client",
	Run: func(cmd *cobra.Command, args []string) {
		tcp.StartClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

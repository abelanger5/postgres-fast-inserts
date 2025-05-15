package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inserts",
	Short: "inserts demonstrates simple commands for testing different insert methods in Postgres.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var maxPayloadSize int

func init() {
	rootCmd.PersistentFlags().IntVar(
		&maxPayloadSize,
		"max-payload-size",
		1000,
		"maximum size of the payload in kilobytes",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

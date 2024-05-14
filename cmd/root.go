package cmd

import (
	"os"

	"github.com/marciomarquesdesouza/go-stresstest/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-stresstest",
	Short: "Sistema teste de estresse",
	Long:  "Sistema para fazer teste de estresse em requisições",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.PersistentFlags().GetString("url")
		requests, _ := cmd.PersistentFlags().GetInt64("requests")
		concurrency, _ := cmd.PersistentFlags().GetInt64("concurrency")
		internal.RunStressTester(url, requests, concurrency)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("url", "u", "", "URL a ser testada")
	rootCmd.PersistentFlags().Int64P("requests", "r", 0, "Total de requisições que serão feitas")
	rootCmd.PersistentFlags().Int64P("concurrency", "c", 0, "Numéro máximo de requisições simultâneas")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("requests")
	rootCmd.MarkPersistentFlagRequired("concurrency")
}

package cmd

import (
	"github.com/mzelenkin/go-calendar/internal/restapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := restapi.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Здесь можно объявиить флаги и настройки
	viper.SetDefault("http.listen", "localhost:7879")
	viper.SetDefault("http.enable_cors", true)
	viper.SetDefault("log.level", "debug")
}

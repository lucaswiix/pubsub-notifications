package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lucaswiix/meli/notifications/api"
	"github.com/lucaswiix/meli/notifications/repository"
	"github.com/lucaswiix/meli/notifications/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "MeliNotification",
	Short: "meli notification, processor and indexing",
	Long:  `meli notification,, processor and indexing`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		utils.InitLogger()

		if err = repository.InitDatabase(); err != nil {
			utils.Log.Error("error to open database connection", zap.Error(err))
			os.Exit(1)
		}

		if err = repository.NewRabbitMQ(); err != nil {
			utils.Log.Error("error to open rabbitMQ connection", zap.Error(err))
			os.Exit(1)
		}

		server := api.InitWebServer()

		if err = server.Run(":8080"); err != nil {
			utils.Log.Error("error to initialize web server", zap.Error(err))
			os.Exit(1)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(configInitialize)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.api.yaml)")

	viper.BindPFlags(rootCmd.PersistentFlags())

	rootCmd.Flags().String("input.api.port", "", "http server api port")
	rootCmd.Flags().String("input.redis.addr", "", "input redis address")
	rootCmd.Flags().String("input.redis.password", "", "input redis password")
	rootCmd.Flags().String("input.rabbitmq.subscription_url", "", "input redis subscribe url")
	rootCmd.Flags().String("logger.level", "", "logger level")

	rootCmd.Flags().String("prom.addr", ":9090", "prometheus address")

	viper.BindPFlags(rootCmd.Flags())

}

func configInitialize() {

	// enable ability to specify config file via flag
	if cfgFile != "" {
		log.Println(cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		// directory path of config file
		viper.AddConfigPath("/etc/api")
		// directory path of config file
		viper.AddConfigPath(".")
		// name of config file (without extension)
		viper.SetConfigName("api")
	}

	// read in environment variables that match
	viper.AutomaticEnv()
	// all environment variables that contains _ will be replaced by - (example - env: HOST_NAME -> access: host-name)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}

}

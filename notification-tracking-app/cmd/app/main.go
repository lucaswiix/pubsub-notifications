package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lucaswiix/notifications-tracking-app/notification/client"
	_notificationClient "github.com/lucaswiix/notifications-tracking-app/notification/client"
	_notificationWebDelivery "github.com/lucaswiix/notifications-tracking-app/notification/delivery/web"
	_notificationUcase "github.com/lucaswiix/notifications-tracking-app/notification/usecase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "MeliNotification",
	Short: "BFF for Meli Notifications",
	Long:  `BFF for Meli Notifications`,
	Run: func(cmd *cobra.Command, args []string) {
		config := client.DefaultRabbitConfig

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Error on Unmarshal configs: %s", err)
		}

		pc, err := _notificationClient.NewRabbitMQClient(config.Address, config.QueueName)
		if err != nil {
			log.Panicln(err)
		}
		defer pc.Close()

		pu := _notificationUcase.NewNotificationUsecase(pc)

		//web
		_notificationWebDelivery.NewNotificationHandler(pu)
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	viper.BindPFlags(rootCmd.PersistentFlags())

	rootCmd.Flags().String("rabbitmq_addr", "", "rabbit address")
	rootCmd.Flags().String("queue_name", "", "redis queue name")
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
		viper.AddConfigPath("/etc/config")
		// directory path of config file
		viper.AddConfigPath(".")
		// name of config file (without extension)
		viper.SetConfigName("config")
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

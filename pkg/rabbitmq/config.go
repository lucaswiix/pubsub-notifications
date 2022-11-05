package rabbitmq

//Config holds bucket config params
type RabbitMQConfig struct {
	SubscriptionURL string `mapstructure:"subscription_url"`
}

package PubSub

type PubSubInterface interface {
	PublishToService(serviceName string, message PublishMessageType) error
}

package PubSub

type PubSubInterface interface {
	PublishToService(serviceName string, message PubSubMessageType) error
}

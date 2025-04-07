package PubSub

type PubSubMessageType struct {
	Data map[string]interface{} `json:"data"`
	Type string                 `json:"type"`
}

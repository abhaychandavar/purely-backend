package PubSub

type PublishMessageType struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

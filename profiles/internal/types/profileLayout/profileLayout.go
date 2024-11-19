package profileLayoutConstant

// leaf element
type Element struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Label    string `json:"label"`
}

// Supported elements
type InputElement struct {
	Element
	Placeholder string `json:"placeholder"`
	InputType   string `json:"inputType"`
}

type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Id    string `json:"id"`
}

type SelectElement struct {
	Element
	Options     []SelectOption `json:"options"`
	Placeholder string         `json:"placeholder"`
}

type PromptSelect struct {
	Options     []SelectOption `json:"options"`
	Placeholder string         `json:"placeholder"`
}

type PromptInput struct {
	Placeholder string `json:"placeholder"`
	InputType   string `json:"inputType"`
}

type Prompt struct {
	Element
	SelectElement       PromptSelect `json:"selectElement"`
	InputElement        PromptInput  `json:"inputElement"`
	PromptFetchEndpoint string       `json:"promptFetchEndpoint"`
	Count               int          `json:"count"`
	UniquePrompts       bool         `json:"uniquePrompts"`
	PromptOptions       []string     `json:"staticPrompts"`
}

type LayoutElement interface{}

// element group (optional)

type ElementGroup struct {
	Id             string          `json:"id"`
	Label          string          `json:"label"`
	IsLogicalGroup bool            `json:"isLogicalGroup"`
	Elements       []LayoutElement `json:"elements"`
}

package profileLayoutTypes

type elementTypeEnum string

// Define the enum values as constants
const (
	InputElementType            elementTypeEnum = "input"
	SelectElementType           elementTypeEnum = "select"
	PromptElementType           elementTypeEnum = "prompt"
	ImageElementType            elementTypeEnum = "image"
	SearchableSelectElementType elementTypeEnum = "searchableSelect"
	LocationElementType         elementTypeEnum = "location"
	LocationStepperElementType  elementTypeEnum = "locationStepper"
)

// leaf element
type Element struct {
	Id       string          `json:"id"`
	Type     elementTypeEnum `json:"type"`
	Required bool            `json:"required"`
	Label    string          `json:"label"`
	IconUrl  string          `json:"iconUrl"`
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
	Options      []SelectOption `json:"options"`
	Placeholder  string         `json:"placeholder"`
	InitialValue string         `json:"initialValue"`
}

type SearchableSelectElement struct {
	Element
	Options          []SelectOption `json:"options"`
	Placeholder      string         `json:"placeholder"`
	DefaultOptionIds []string       `json:"defaultOptionIds"`
	InitialValue     string         `json:"initialValue"`
}

type PromptInput struct {
	Placeholder string `json:"placeholder"`
	InputType   string `json:"inputType"`
}

type Prompt struct {
	Element
	InputElement        PromptInput `json:"inputElement"`
	PromptFetchEndpoint string      `json:"promptFetchEndpoint"`
	Count               int         `json:"count"`
	UniquePrompts       bool        `json:"uniquePrompts"`
	PromptOptions       []string    `json:"staticPrompts"`
}

type Images struct {
	Element
	Count         int `json:"count"`
	RequiredCount int `json:"requiredCount"`
}

type DistanceStepper struct {
	Element
	MinDistance int    `json:"minDistance"`
	MaxDistance int    `json:"maxDistance"`
	Unit        string `json:"unit"`
}

type Location struct {
	Element
}

type LayoutElement interface{}

// element group (optional)

type ElementGroup struct {
	Id             string          `json:"id"`
	Label          string          `json:"label"`
	IsLogicalGroup bool            `json:"isLogicalGroup"`
	Elements       []LayoutElement `json:"elements"`
}

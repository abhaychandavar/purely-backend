package profileServiceTypes

type CreateProfileType struct {
	AuthId   *string  `json:"authId"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	Category *string  `json:"category"`
}

type GetProfileType struct {
	Category *string `json:"category"`
	AuthId   *string `json:"authId"`
}

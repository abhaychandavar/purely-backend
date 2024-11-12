package profileServiceTypes

type CreateProfileType struct {
	AuthId *string  `json:"authId"`
	Lat    *float64 `json:"lat"`
	Lng    *float64 `json:"lng"`
}

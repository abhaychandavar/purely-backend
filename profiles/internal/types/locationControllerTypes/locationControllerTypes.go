package locationControllerTypes

type GetLocationsType struct {
	Query     string `json:"place"`
	PageToken string `json:"pageToken"`
}

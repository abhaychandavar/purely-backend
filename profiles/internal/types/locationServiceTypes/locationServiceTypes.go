package locationServiceTypes

type GetLocationsType struct {
	Query     *string `json:"query"`
	PageToken *string `json:"pageToken"`
}

type GetLocationsResponseType struct {
	Results       []interface{} `json:"results"`
	NextPageToken *string       `json:"nextPageToken"`
}

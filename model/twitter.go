package model

//TwitterRequestBody .
type TwitterRequestBody struct {
	Data Twitter `json:"data" form:"data"`
}

//Twitter .
type Twitter struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
	// HashTag  string `json:"hash_tag" form:"hash_tag"`
	Sentence string `json:"sentence"`
	Key      string `json:"key"`
}

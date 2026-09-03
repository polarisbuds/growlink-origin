package inputmodel

type TopicInput struct {
	DisplayNo uint   `json:"displayNo"`
	Question  string `json:"question"`
	IsRequire bool   `json:"isRequire"`
}

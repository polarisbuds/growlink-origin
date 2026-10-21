package inputmodel

type TopicInput struct {
	DisplayNo   uint   `json:"displayNo"`
	Question    string `json:"question"`
	Description string `json:"description"`
	IsRequire   bool   `json:"isRequire"`
}

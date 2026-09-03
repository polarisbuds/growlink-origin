package outputmodel

type TopicOutput struct {
	ID        uint   `json:"id"`
	DisplayNo uint   `json:"displayNo"`
	Question  string `json:"question"`
	IsRequire bool   `json:"isRequire"`
}

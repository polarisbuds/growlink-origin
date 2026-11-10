package outputmodel

type TopicOutput struct {
	ID          uint   `json:"-"`
	UUID        string `json:"uuid"`
	DisplayNo   uint   `json:"displayNo"`
	Question    string `json:"question"`
	Description string `json:"description"`
	IsRequire   bool   `json:"isRequire"`
}

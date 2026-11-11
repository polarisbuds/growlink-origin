package outputmodel

type PairingOutput struct {
	ID     uint       `json:"-"`
	UUID   string     `json:"uuid"`
	Mentor UserOutput `json:"mentor"`
	Mentee UserOutput `json:"mentee"`
}

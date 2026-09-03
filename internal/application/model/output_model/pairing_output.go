package outputmodel

type PairingOutput struct {
	ID     uint       `json:"id"`
	Mentor UserOutput `json:"mentor"`
	Mentee UserOutput `json:"mentee"`
}

package inputmodel

type PairingTargetInput struct {
	Kind     string `json:"kind"`
	TargetID uint   `json:"targetId"`
}

type PairingInput struct {
	MentorID uint                 `json:"mentorId"`
	Targets  []PairingTargetInput `json:"targets"`
}

func (in PairingInput) ResolvedTargets() []PairingTargetInput {
	return in.Targets
}

package inputmodel

type PairingTargetInput struct {
	Kind       string `json:"kind"`
	TargetUuid string `json:"targetUuid"`
	TargetID   uint   `json:"-"`
}

type PairingInput struct {
	MentorUuid string               `json:"mentorUuid"`
	MentorID   uint                 `json:"-"`
	Targets    []PairingTargetInput `json:"targets"`
}

func (in PairingInput) ResolvedTargets() []PairingTargetInput {
	return in.Targets
}

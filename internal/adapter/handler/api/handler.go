package api

import (
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/application/usecase/planning_usecase"
	"growlink/internal/application/usecase/setting_usecase"
)

type Deps struct {
	Auth       *authenticate_usecase.AuthUsecase
	Meeting    meeting_usecase.Bundle
	Management management_usecase.Bundle
	Pairing    *pairing_usecase.PairingListCreatorUsecase
	PairingBrk *pairing_usecase.PairingBreakerUsecase
	Setting    *setting_usecase.IconRegisterUsecase
	SettingRm  *setting_usecase.IconRemoverUsecase
	Profile    *setting_usecase.ProfileUsecase
	Planning   *planning_usecase.PlanningAccessUsecase
}

type Handler struct {
	deps Deps
}

func New(deps Deps) *Handler {
	return &Handler{deps: deps}
}

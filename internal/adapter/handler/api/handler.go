package api

import (
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/application/usecase/setting_usecase"
	repository_interface "growlink/internal/domain/interface/repository"
)

type Deps struct {
	Auth       *authenticate_usecase.AuthUsecase
	Meeting    meeting_usecase.Bundle
	Management management_usecase.Bundle
	Pairing    *pairing_usecase.PairingListCreatorUsecase
	PairingBrk *pairing_usecase.PairingBreakerUsecase
	Setting    *setting_usecase.IconRegisterUsecase
	SettingRm  *setting_usecase.IconRemoverUsecase
	MeetingRepo     repository_interface.MeetingRepository
	PairingRepo     repository_interface.PairingRepository
	TemplateRepo    repository_interface.TemplateRepository
	UserRepo        repository_interface.UserRepository
	TeamRepo        repository_interface.TeamRepository
	TeamMemberRepo  repository_interface.TeamMemberRepository
	DepartmentRepo  repository_interface.DepartmentRepository
	RoleRepo        repository_interface.RoleRepository
}

type Handler struct {
	deps Deps
}

func New(deps Deps) *Handler {
	return &Handler{deps: deps}
}

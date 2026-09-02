package usecase_interface

import "growlink/internal/domain/entity"

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type ManagementUsecase interface {
	ProgressManagementUsecase
	MeetingManagementUsecase
	TemplateManagementUsecase
	PairingManagementUsecase
}

type (
	ProgressManagementUsecase interface {
		ProgressViewer
	}

	ProgressViewer interface {
		GetProgressByDepartmentID(departmentID uint) error // TODO; 実施状況の見た目が確定後に決定
		GetProgressByTeamID(teamID uint) error             // TODO; 実施状況の見た目が確定後に決定
		GetProgressByBossID(bossID uint) error             // TODO; 実施状況の見た目が確定後に決定
		GetProgressByMentorID(mentorID uint) error         // TODO; 実施状況の見た目が確定後に決定
	}
)

type (
	MeetingManagementUsecase interface {
		MeetingViewer
		MeetingCreator
		MeetingUpdater
		MeetingDeleter
	}

	MeetingViewer interface {
		ListMeetingsByMentor(mentorID uint) ([]entity.Meeting, error)
		ListMeetingsByBoss(bossID uint) ([]entity.Meeting, error)
		ListMeetingsByTeam(teamID uint) ([]entity.Meeting, error)
		ListMeetingsByDepartment(departmentID uint) ([]entity.Meeting, error)
	}
	MeetingCreator interface {
		CreateMeeting(pairingID, templateID uint) (*entity.Meeting, error)
		CreateMeetingByAll() ([]entity.Meeting, error)
	}
	MeetingUpdater interface {
		UpdateMeeting(*entity.Meeting) error
		ChangePairing(meetingID, newPairingID uint)
		ChangeTemplate(meetingID, newTemplateID uint)
	}
	MeetingDeleter interface {
		DeleteMeeting(meetingID uint) (bool, error)
	}
)

type (
	TemplateManagementUsecase interface {
		TemplateViewer
		TemplateCreator
		TemplateUpdater
		TemplateDeleter
	}

	TemplateViewer interface {
		ListTemplate() ([]entity.Template, error)
		GetTemplate(templateID uint) (*entity.Template, error)
	}
	TemplateCreator interface {
		CreateTemplate(newTemplate *entity.Template) error
	}
	TemplateUpdater interface {
		UpdateTemplate(newTemplate *entity.Template) error
	}
	TemplateDeleter interface {
		DeleteTemplate(templateID uint) (bool, error)
	}
)

type (
	PairingManagementUsecase interface {
		PairingViewer
		PairingCreator
		PairingUpdater
		PairingDeleter
	}

	PairingViewer interface {
		GetPairing(pairingID uint) (*entity.Pairing, error)
		ListPairingsByAll() ([]entity.Pairing, error)
		ListPairingsByMentor(mentorID uint) ([]entity.Pairing, error)
		ListPairingsByBoss(bossID uint) ([]entity.Pairing, error)
		ListPairingsByTeam(teamID uint) ([]entity.Pairing, error)
		ListPairingsByDepartment(departmentID uint) ([]entity.Pairing, error)
	}
	PairingCreator interface {
		CreatePairing(menteeID, mentorID uint) error
	}
	PairingUpdater interface {
		UpdatePairing(newPairing *entity.Pairing) error
	}
	PairingDeleter interface {
		BreakPairing(pairingID uint) error
		BreakPairings(pairingIDs ...uint) error
	}
)

package usecase_interface

import (
	"growlink/internal/domain/entity"
	"time"
)

//go:generate go tool mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type Manager interface {
	ProgressManager
	MeetingManager
	TemplateManager
	PairingManager
}

type (
	ProgressManager interface {
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
	MeetingManager interface {
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
		CreateMeetingEveryMonth(templateID uint, baseDate time.Time, week int, weekday time.Weekday) ([]entity.Meeting, error)
		CreateMeeting(pairingID, templateID uint, implementDate time.Time) (*entity.Meeting, error)
		CreateMeetingByAll(templateID uint, implementDate time.Time) ([]entity.Meeting, error)
	}
	MeetingUpdater interface {
		Update(meetingID uint, meeting *entity.Meeting) error
		ChangeImplementDate(meetingID uint, implementDate time.Time) error
		ChangePairing(meetingID, newPairingID uint) error
		ChangeTemplate(meetingID, newTemplateID uint) error
	}
	MeetingDeleter interface {
		DeleteMeeting(meetingIDs ...uint) (bool, error)
	}
)

type (
	TemplateManager interface {
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
	PairingManager interface {
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

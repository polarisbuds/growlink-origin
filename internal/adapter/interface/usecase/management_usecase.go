package usecase_interface

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
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
		// 		GetProgressByTeamID(teamID uint) error // TODO; 実施状況の見た目が確定後に決定
		// 		GetProgressByBossID(bossID uint) error // TODO; 実施状況の見た目が確定後に決定
		// 		GetProgressByMentorID(mentorID uint) error // TODO; 実施状況の見た目が確定後に決定
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
		ListMeetingsByMentor(mentorID uint) ([]outputmodel.MeetingOutput, error)
		ListMeetingsByBoss(bossID uint) ([]outputmodel.MeetingOutput, error)
		ListMeetingsByTeam(teamID uint) ([]outputmodel.MeetingOutput, error)
		ListMeetingsByDepartment(departmentID uint) ([]outputmodel.MeetingOutput, error)
	}
	MeetingCreator interface {
		CreateMeetingEveryMonth(templateID uint, baseDate time.Time, week int, weekday time.Weekday, createdByUserID uint, sharerUserIDs []uint) ([]outputmodel.MeetingOutput, error)
		CreateMeeting(pairingID, templateID uint, implementDate time.Time, createdByUserID uint, sharerUserIDs []uint) (*outputmodel.MeetingOutput, error)
		CreateMeetingByAll(templateID uint, implementDate time.Time, createdByUserID uint, sharerUserIDs []uint) ([]outputmodel.MeetingOutput, error)
	}
	MeetingUpdater interface {
		Update(meetingID uint, meeting inputmodel.MeetingInput) error
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
		ListTemplate() ([]outputmodel.TemplateOutput, error)
		GetTemplate(templateID uint) (*outputmodel.TemplateOutput, error)
	}
	TemplateCreator interface {
		CreateTemplate(input inputmodel.TemplateInput) error
	}
	TemplateUpdater interface {
		UpdateTemplate(templateID uint, input inputmodel.TemplateInput) error
	}
	TemplateDeleter interface {
		DeleteTemplate(templateID uint) error
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
		GetPairing(pairingID uint) (*outputmodel.PairingOutput, error)
		ListPairingsByAll() ([]outputmodel.PairingOutput, error)
		ListPairingsByMentor(mentorID uint) ([]outputmodel.PairingOutput, error)
		ListPairingsByBoss(bossID uint) ([]outputmodel.PairingOutput, error)
		ListPairingsByTeam(teamID uint) ([]outputmodel.PairingOutput, error)
		ListPairingsByDepartment(departmentID uint) ([]outputmodel.PairingOutput, error)
	}
	PairingCreator interface {
		CreatePairing(menteeID, mentorID uint) error
	}
	PairingUpdater interface {
		UpdatePairing(pairingID uint, input inputmodel.PairingInput) error
	}
	PairingDeleter interface {
		BreakPairing(pairingID uint) error
		BreakPairings(pairingIDs ...uint) error
	}
)

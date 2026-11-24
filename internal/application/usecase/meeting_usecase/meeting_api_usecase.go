package meeting_usecase

import (
	"context"
	"errors"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	logutils "growlink/internal/shared/logUtils"

	"github.com/google/uuid"
)

var (
	ErrMeetingForbidden         = errors.New("forbidden")
	ErrMeetingNotFound          = errors.New("not found")
	ErrRequireAnswerNotAnswered = errors.New("answered require")
)

type MeetingDetailFlags struct {
	CanEdit               bool
	ViewerRole            string
	CanWriteAnswers       bool
	CanWriteTopicComments bool
	CanEditSharers        bool
}

type MeetingDetailOutput struct {
	MeetingUUID string
	Meeting     outputmodel.MeetingOutput
	Answers     []outputmodel.AnswerOutput
	Topics      []outputmodel.TopicOutput
	Comments    []outputmodel.MeetingCommentOutput
	Sharers     []outputmodel.UserOutput
	Flags       MeetingDetailFlags
}

type MeetingAPIUsecase struct {
	meetingRepo      repository_interface.MeetingRepository
	userRepo         repository_interface.UserRepository
	answerRepo       repository_interface.AnswerRepository
	answerReader     *AnswerReaderUsecase
	draftSaver       *DraftSaverUsecase
	submitter        *SubmitterUsecase
	submitCanceller  *SubmitCancellerUsecase
	mentorCompleter  *MentorMeetingCompleterUsecase
	meetingCommenter *MeetingCommenterUsecase
	answerCommenter  *AnswerCommenterUsecase
	planning         PlanningMeetingAccess
}

func NewMeetingAPIUsecase(
	meetingRepo repository_interface.MeetingRepository,
	userRepo repository_interface.UserRepository,
	answerRepo repository_interface.AnswerRepository,
	answerReader *AnswerReaderUsecase,
	draftSaver *DraftSaverUsecase,
	submitter *SubmitterUsecase,
	submitCanceller *SubmitCancellerUsecase,
	mentorCompleter *MentorMeetingCompleterUsecase,
	meetingCommenter *MeetingCommenterUsecase,
	answerCommenter *AnswerCommenterUsecase,
) *MeetingAPIUsecase {
	return &MeetingAPIUsecase{
		meetingRepo:      meetingRepo,
		userRepo:         userRepo,
		answerRepo:       answerRepo,
		answerReader:     answerReader,
		draftSaver:       draftSaver,
		submitter:        submitter,
		submitCanceller:  submitCanceller,
		mentorCompleter:  mentorCompleter,
		meetingCommenter: meetingCommenter,
		answerCommenter:  answerCommenter,
	}
}

func (uc *MeetingAPIUsecase) SetPlanningAccess(access PlanningMeetingAccess) {
	uc.planning = access
}

func (uc *MeetingAPIUsecase) loadMeetingByUUID(id uuid.UUID) (*entity.Meeting, error) {
	meeting, err := uc.meetingRepo.GetByUUID(id)
	if err != nil {
		return nil, err
	}
	return meeting, nil
}

func (uc *MeetingAPIUsecase) canView(
	meeting *entity.Meeting,
	viewerID uint,
	authority valueobject.Authority,
	systemRole valueobject.SystemRole,
) (bool, error) {
	if viewerID == 0 {
		return true, nil
	}
	if uc.planning != nil {
		ok, err := uc.planning.CanViewPlanningMeeting(viewerID, authority, systemRole, meeting.ID)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return CanViewMeeting(meeting, viewerID), nil
}

func (uc *MeetingAPIUsecase) GetMeetingDetail(
	meetingUUID uuid.UUID,
	viewerID uint,
	authority any,
	systemRole any,
	canModify func(viewerID uint, systemRole valueobject.SystemRole, createdBy uint) bool,
) (*MeetingDetailOutput, error) {
	auth, _ := authority.(valueobject.Authority)
	role, _ := systemRole.(valueobject.SystemRole)
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	ok, err := uc.canView(meeting, viewerID, auth, role)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrMeetingNotFound
	}
	out := converter.MeetingEntityToOutput(*meeting)
	answers, err := uc.answerReader.ListAnswersByMeeting(meeting.ID)
	if err != nil {
		return nil, err
	}
	topics := out.Template.Topics
	if len(topics) == 0 && meeting.Template != nil {
		topics = converter.TopicEntitiesToOutputs(meeting.Template.Topics)
	}
	flags := MeetingDetailFlags{ViewerRole: "viewer"}
	if viewerID != 0 {
		flags.CanEdit = canModify(viewerID, role, meeting.CreatedByUserID)
		flags.ViewerRole = MeetingViewerRole(meeting, viewerID)
		flags.CanWriteAnswers = CanWriteMeetingAnswers(meeting, viewerID)
		flags.CanWriteTopicComments = CanWriteTopicComments(meeting, viewerID)
		flags.CanEditSharers = CanEditMeetingSharers(meeting, viewerID)
	}
	return &MeetingDetailOutput{
		MeetingUUID: meetingUUID.String(),
		Meeting:     out,
		Answers:     answers,
		Topics:      topics,
		Comments:    FilterMeetingComments(meeting, viewerID),
		Sharers:     MeetingSharersToOutputs(meeting),
		Flags:       flags,
	}, nil
}

func (uc *MeetingAPIUsecase) SaveDraftByUUID(meetingUUID uuid.UUID, viewerID uint, answers []outputmodel.AnswerOutput) error {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return err
	}
	if !CanWriteMeetingAnswers(meeting, viewerID) {
		return ErrMeetingForbidden
	}
	sanitized, err := sanitizeAnswerOutputs(meeting.ID, viewerID, answers, meeting, uc.answerReader.GetAnswer)
	if err != nil {
		return err
	}
	return uc.draftSaver.SaveDraft(meeting.ID, sanitized, false)
}

func (uc *MeetingAPIUsecase) SubmitByUUID(meetingUUID uuid.UUID, viewerID uint, answers []outputmodel.AnswerOutput) error {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return err
	}
	if !CanWriteMeetingAnswers(meeting, viewerID) {
		return ErrMeetingForbidden
	}
	if err := validateRequiredAnswers(meeting, answers); err != nil {
		return err
	}
	sanitized, err := sanitizeAnswerOutputs(meeting.ID, viewerID, answers, meeting, uc.answerReader.GetAnswer)
	if err != nil {
		return err
	}
	return uc.submitter.Submit(meeting.ID, sanitized, valueobject.MeetingStatusSubmitted)
}

func (uc *MeetingAPIUsecase) CompleteMentorByUUID(meetingUUID uuid.UUID, viewerID uint, comments []MentorTopicCommentInput) error {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return err
	}
	if !CanWriteTopicComments(meeting, viewerID) {
		return ErrMeetingForbidden
	}
	if meeting.Pairing == nil {
		return fmt.Errorf("pairing not found")
	}
	return uc.mentorCompleter.Complete(meeting.ID, viewerID, meeting.Pairing.MenteeID, comments)
}

func (uc *MeetingAPIUsecase) CancelSubmitByUUID(meetingUUID uuid.UUID, viewerID uint) error {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return err
	}
	if !CanCancelMeetingSubmit(meeting, viewerID) {
		return ErrMeetingForbidden
	}
	return uc.submitCanceller.CancelSubmit(meeting.ID)
}

func (uc *MeetingAPIUsecase) PostMeetingCommentByUUID(meetingUUID uuid.UUID, viewerID uint, content string, isLimited bool) error {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return err
	}
	if meeting.Status != valueobject.MeetingStatusDone {
		return ErrMeetingForbidden
	}
	return uc.meetingCommenter.PostComment(meeting.ID, viewerID, content, isLimited)
}

func (uc *MeetingAPIUsecase) PostAnswerCommentByAnswerUUID(answerUUID uuid.UUID, viewerID uint, content string) error {
	const op = "meeting_usecase.PostAnswerCommentByAnswerUUID"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	var err error
	defer func() { done(err) }()
	answer, err := uc.answerRepo.GetByUUID(answerUUID)
	if err != nil {
		return err
	}
	return uc.answerCommenter.PostComment(answer.ID, viewerID, content)
}

type SharerCandidate struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uc *MeetingAPIUsecase) ListSharerCandidates(meetingUUID uuid.UUID, viewerID uint, query string) ([]SharerCandidate, error) {
	meeting, err := uc.loadMeetingByUUID(meetingUUID)
	if err != nil {
		return nil, err
	}
	if !CanEditMeetingSharers(meeting, viewerID) {
		return nil, ErrMeetingForbidden
	}
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]SharerCandidate, 0, 32)
	for _, u := range users {
		if u.Suspended || !UserMatchesSharerQuery(u, query) {
			continue
		}
		out = append(out, SharerCandidate{UUID: u.UUID.String(), Name: u.Name, Email: u.Email})
		if len(out) >= 50 {
			break
		}
	}
	return out, nil
}

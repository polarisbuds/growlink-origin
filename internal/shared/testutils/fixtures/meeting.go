package fixtures

import (
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MeetingBuilder builds Meeting test data via GORM (not Repository).
type MeetingBuilder struct {
	t             *testing.T
	tx            *gorm.DB
	implementedAt *time.Time
	status        *valueobject.MeetingStatus
	template      *entity.Template
	pairing       *entity.Pairing
	meetingUUID   *uuid.UUID
	nestedPreload bool
}

// Meeting creates a MeetingBuilder with default dependencies.
func Meeting(t *testing.T, tx *gorm.DB) *MeetingBuilder {
	t.Helper()
	return &MeetingBuilder{t: t, tx: tx}
}

func (b *MeetingBuilder) WithImplementedAt(at time.Time) *MeetingBuilder {
	b.implementedAt = &at
	return b
}

func (b *MeetingBuilder) WithStatus(status valueobject.MeetingStatus) *MeetingBuilder {
	b.status = &status
	return b
}

func (b *MeetingBuilder) WithTemplate(template *entity.Template) *MeetingBuilder {
	b.template = template
	return b
}

func (b *MeetingBuilder) WithPairing(pairing *entity.Pairing) *MeetingBuilder {
	b.pairing = pairing
	return b
}

func (b *MeetingBuilder) WithUUID(id uuid.UUID) *MeetingBuilder {
	b.meetingUUID = &id
	return b
}

func (b *MeetingBuilder) WithNestedPreloadGraph() *MeetingBuilder {
	b.nestedPreload = true
	return b
}

// Create persists the meeting graph in ER dependency order:
// Layer2 Pairing/Template -> Layer4 Meeting -> Layer5/6 children.
func (b *MeetingBuilder) Create() *entity.Meeting {
	b.t.Helper()

	status := valueobject.MeetingStatusPending
	if b.status != nil {
		status = *b.status
	}

	pairing := b.pairing
	if pairing == nil {
		pairing = createPairing(b.t, b.tx)
	}

	template := b.template
	if template == nil {
		template = createTemplate(b.t, b.tx)
	}

	implementedAt := time.Now()
	if b.implementedAt != nil {
		implementedAt = *b.implementedAt
	}

	meeting := &entity.Meeting{
		ImplementedAt: implementedAt,
		Status:        status,
		TemplateID:    template.ID,
		PairingID:     pairing.ID,
	}
	if b.meetingUUID != nil {
		meeting.UUID = *b.meetingUUID
	}

	if err := b.tx.Create(meeting).Error; err != nil {
		b.t.Fatalf("failed to create meeting fixture: %v", err)
	}

	if b.nestedPreload {
		createMeetingChildren(b.t, b.tx, meeting, template, pairing)
	}

	return meeting
}

// createMeetingChildren inserts Answer/MeetingComment graph (Layer5/6).
// Topic must exist before Answer (topics -> answers FK).
func createMeetingChildren(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, template *entity.Template, pairing *entity.Pairing) {
	t.Helper()

	topic := createTopic(t, tx, template)
	answer := createAnswer(t, tx, meeting, topic, pairing.Mentor)
	createAnswerComment(t, tx, answer, pairing.Mentee)
	createMeetingComment(t, tx, meeting, pairing.Mentor)
}

// MeetingWithoutFK creates a meeting with invalid FK references for constraint tests.
func MeetingWithoutFK(t *testing.T, tx *gorm.DB) *entity.Meeting {
	t.Helper()

	meeting := &entity.Meeting{
		ImplementedAt: time.Now(),
		Status:        valueobject.MeetingStatusPending,
		TemplateID:    999999999,
		PairingID:     999999999,
	}
	return meeting
}

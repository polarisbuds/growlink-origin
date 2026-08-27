package fixtures

import (
	"testing"
	"time"

	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MeetingBuilder builds Meeting test data via GORM (not Repository).
type MeetingBuilder struct {
	t             *testing.T
	tx            *gorm.DB
	implementedAt *time.Time
	status        *entity.MeetingStatus
	template      *entity.Template
	peering       *entity.Peering
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

func (b *MeetingBuilder) WithStatus(status *entity.MeetingStatus) *MeetingBuilder {
	b.status = status
	return b
}

func (b *MeetingBuilder) WithTemplate(template *entity.Template) *MeetingBuilder {
	b.template = template
	return b
}

func (b *MeetingBuilder) WithPeering(peering *entity.Peering) *MeetingBuilder {
	b.peering = peering
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
// Layer0 Status -> Layer2 Peering/Template -> Layer4 Meeting -> Layer5/6 children.
func (b *MeetingBuilder) Create() *entity.Meeting {
	b.t.Helper()

	status := b.status
	if status == nil {
		status = createMeetingStatus(b.t, b.tx)
	}

	peering := b.peering
	if peering == nil {
		peering = createPeering(b.t, b.tx)
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
		StatusID:      int64(status.ID),
		TemplateID:    int64(template.ID),
		PeeringID:     int64(peering.ID),
	}
	if b.meetingUUID != nil {
		meeting.UUID = *b.meetingUUID
	}

	if err := b.tx.Create(meeting).Error; err != nil {
		b.t.Fatalf("failed to create meeting fixture: %v", err)
	}

	if b.nestedPreload {
		createMeetingChildren(b.t, b.tx, meeting, template, peering)
	}

	return meeting
}

// createMeetingChildren inserts Answer/MeetingComment graph (Layer5/6).
// Topic must exist before Answer (topics -> answers FK).
func createMeetingChildren(t *testing.T, tx *gorm.DB, meeting *entity.Meeting, template *entity.Template, peering *entity.Peering) {
	t.Helper()

	topic := createTopic(t, tx, template)
	answer := createAnswer(t, tx, meeting, topic, peering.Mentor)
	createAnswerComment(t, tx, answer, peering.Mentee)
	createMeetingComment(t, tx, meeting, peering.Mentor)
}

// MeetingWithoutFK creates a meeting with invalid FK references for constraint tests.
func MeetingWithoutFK(t *testing.T, tx *gorm.DB) *entity.Meeting {
	t.Helper()

	meeting := &entity.Meeting{
		ImplementedAt: time.Now(),
		StatusID:      999999999,
		TemplateID:    999999999,
		PeeringID:     999999999,
	}
	return meeting
}

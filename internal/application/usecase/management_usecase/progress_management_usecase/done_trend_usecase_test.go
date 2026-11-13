package progress_mng_usecase

import (
	"testing"
	"time"

	"growlink/internal/domain/entity"
	"growlink/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
)

func TestDoneTrendByMonth_ClipsRangePerMonth(t *testing.T) {
	loc := time.Local
	from := time.Date(2024, 10, 15, 0, 0, 0, 0, loc)
	to := time.Date(2024, 12, 12, 0, 0, 0, 0, loc).Add(24*time.Hour - time.Nanosecond)

	meetings := []entity.Meeting{
		meetingMentee(20, valueobject.MeetingStatusDone, time.Date(2024, 10, 14, 12, 0, 0, 0, loc)),
		meetingMentee(20, valueobject.MeetingStatusDone, time.Date(2024, 10, 20, 12, 0, 0, 0, loc)),
		meetingMentee(20, valueobject.MeetingStatusDone, time.Date(2024, 11, 15, 12, 0, 0, 0, loc)),
		meetingMentee(20, valueobject.MeetingStatusSubmitted, time.Date(2024, 11, 16, 12, 0, 0, 0, loc)),
		meetingMentee(20, valueobject.MeetingStatusDone, time.Date(2024, 12, 10, 12, 0, 0, 0, loc)),
	}
	pairings := []entity.Pairing{{MenteeID: 20}}
	uc := newDeptProgressUC(t, &stubReferable{}, meetings, pairings)

	got, err := uc.DoneTrendByMonth(ProgressFilter{
		From:       &from,
		To:         &to,
		SystemRole: valueobject.SystemRoleAdmin,
	})
	require.NoError(t, err)
	require.Len(t, got, 3)
	require.Equal(t, "2024年10月", got[0].Label)
	require.Equal(t, 1, got[0].Done)
	require.Equal(t, "2024年11月", got[1].Label)
	require.Equal(t, 1, got[1].Done)
	require.Equal(t, "2024年12月", got[2].Label)
	require.Equal(t, 1, got[2].Done)
}

func TestDoneTrendByMonth_EmptyWhenRangeIncomplete(t *testing.T) {
	uc := newDeptProgressUC(t, &stubReferable{}, nil, nil)
	got, err := uc.DoneTrendByMonth(ProgressFilter{SystemRole: valueobject.SystemRoleAdmin})
	require.NoError(t, err)
	require.Empty(t, got)
}

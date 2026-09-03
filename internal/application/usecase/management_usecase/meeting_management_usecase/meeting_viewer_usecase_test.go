package meeting_mng_usecase_test

import (
	"testing"

	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"go.uber.org/mock/gomock"
)

type meetingViewerMocks struct {
	meeting          *repomock.MockMeetingRepository
	pairing          *repomock.MockPairingRepository
	userRelationship *repomock.MockUserRelationshipRepository
	user             *repomock.MockUserRepository
	teamMember       *repomock.MockTeamMemberRepository
}

func newMeetingViewer(t *testing.T) (*meeting_mng_usecase.MeetingViewer, *meetingViewerMocks) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mocks := &meetingViewerMocks{
		meeting:          repomock.NewMockMeetingRepository(ctrl),
		pairing:          repomock.NewMockPairingRepository(ctrl),
		userRelationship: repomock.NewMockUserRelationshipRepository(ctrl),
		user:             repomock.NewMockUserRepository(ctrl),
		teamMember:       repomock.NewMockTeamMemberRepository(ctrl),
	}
	uc := meeting_mng_usecase.NewMeetingViewer(
		mocks.meeting,
		mocks.pairing,
		mocks.userRelationship,
		mocks.user,
		mocks.teamMember,
	)
	return uc, mocks
}

func TestMeetingViewer_ListMeetingsByMentor(t *testing.T) {
	tests := []struct {
		name        string
		mentorID    uint
		wantErr     bool
		errContains string
		setup       func(m *meetingViewerMocks)
		assertFn    func(t *testing.T, got []entity.Meeting)
	}{
		{
			name:     "通常系_複数ペアリングのミーティングをマージ",
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(m *meetingViewerMocks) {
				m.pairing.EXPECT().GetByMentorID(validMentorID).Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
					*newTestPairing(pairingID2, validMemberID, validMentorID),
				}, nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return([]entity.Meeting{newTestMeeting(meetingID1, pairingID1)}, nil)
				m.meeting.EXPECT().GetByPairingID(pairingID2).Return([]entity.Meeting{newTestMeeting(meetingID2, pairingID2)}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, []uint{meetingID1, meetingID2}, got)
			},
		},
		{
			name:     "通常系_ペアリングが空",
			mentorID: validMentorID,
			wantErr:  false,
			setup: func(m *meetingViewerMocks) {
				m.pairing.EXPECT().GetByMentorID(validMentorID).Return([]entity.Pairing{}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, nil, got)
			},
		},
		{
			name:        "異常系_GetByMentorID失敗",
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "ミーティングリストのメンターID指定取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.pairing.EXPECT().GetByMentorID(validMentorID).Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:        "異常系_GetByPairingID失敗",
			mentorID:    validMentorID,
			wantErr:     true,
			errContains: "メンターID指定でのミーティングリスト情報の取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.pairing.EXPECT().GetByMentorID(validMentorID).Return([]entity.Pairing{
					*newTestPairing(pairingID1, validMenteeID, validMentorID),
				}, nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return(nil, errGetMeetingsFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mocks := newMeetingViewer(t)
			if tt.setup != nil {
				tt.setup(mocks)
			}

			got, gotErr := uc.ListMeetingsByMentor(tt.mentorID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ListMeetingsByMentor")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

func TestMeetingViewer_ListMeetingsByBoss(t *testing.T) {
	tests := []struct {
		name        string
		bossID      uint
		wantErr     bool
		errContains string
		setup       func(m *meetingViewerMocks)
		assertFn    func(t *testing.T, got []entity.Meeting)
	}{
		{
			name:   "通常系_部下のミーティングを取得",
			bossID: validBossID,
			wantErr: false,
			setup: func(m *meetingViewerMocks) {
				m.userRelationship.EXPECT().GetByBossID(validBossID).Return([]entity.UserRelationship{
					newTestUserRelationship(1, validMenteeID, validBossID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(newTestPairing(pairingID1, validMenteeID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return([]entity.Meeting{newTestMeeting(meetingID1, pairingID1)}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, []uint{meetingID1}, got)
			},
		},
		{
			name:   "通常系_リレーションが空",
			bossID: validBossID,
			wantErr: false,
			setup: func(m *meetingViewerMocks) {
				m.userRelationship.EXPECT().GetByBossID(validBossID).Return([]entity.UserRelationship{}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, nil, got)
			},
		},
		{
			name:        "異常系_GetByBossID失敗",
			bossID:      validBossID,
			wantErr:     true,
			errContains: "上司ユーザID指定でのミーティングリスト作成のユーザリレーション情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.userRelationship.EXPECT().GetByBossID(validBossID).Return(nil, errGetRelationshipsFailed)
			},
		},
		{
			name:        "異常系_GetByMenteeID失敗",
			bossID:      validBossID,
			wantErr:     true,
			errContains: "部下ユーザID（メンティーユーザID）でのペアリング情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.userRelationship.EXPECT().GetByBossID(validBossID).Return([]entity.UserRelationship{
					newTestUserRelationship(1, validMenteeID, validBossID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:        "異常系_GetByPairingID失敗",
			bossID:      validBossID,
			wantErr:     true,
			errContains: "ペアリングユーザID指定でのミーティング情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.userRelationship.EXPECT().GetByBossID(validBossID).Return([]entity.UserRelationship{
					newTestUserRelationship(1, validMenteeID, validBossID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(newTestPairing(pairingID1, validMenteeID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return(nil, errGetMeetingsFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mocks := newMeetingViewer(t)
			if tt.setup != nil {
				tt.setup(mocks)
			}

			got, gotErr := uc.ListMeetingsByBoss(tt.bossID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ListMeetingsByBoss")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

func TestMeetingViewer_ListMeetingsByTeam(t *testing.T) {
	tests := []struct {
		name        string
		teamID      uint
		wantErr     bool
		errContains string
		setup       func(m *meetingViewerMocks)
		assertFn    func(t *testing.T, got []entity.Meeting)
	}{
		{
			name:   "通常系_チームメンバーのミーティングを取得",
			teamID: validTeamID,
			wantErr: false,
			setup: func(m *meetingViewerMocks) {
				m.teamMember.EXPECT().GetByTeamID(validTeamID).Return([]entity.TeamMember{
					newTestTeamMember(1, validTeamID, validMemberID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMemberID).Return(newTestPairing(pairingID1, validMemberID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return([]entity.Meeting{newTestMeeting(meetingID1, pairingID1)}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, []uint{meetingID1}, got)
			},
		},
		{
			name:   "通常系_チームメンバーが空",
			teamID: validTeamID,
			wantErr: false,
			setup: func(m *meetingViewerMocks) {
				m.teamMember.EXPECT().GetByTeamID(validTeamID).Return([]entity.TeamMember{}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, nil, got)
			},
		},
		{
			name:        "異常系_GetByTeamID失敗",
			teamID:      validTeamID,
			wantErr:     true,
			errContains: "チームID指定でのチームメンバー情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.teamMember.EXPECT().GetByTeamID(validTeamID).Return(nil, errGetTeamMembersFailed)
			},
		},
		{
			name:        "異常系_GetByMenteeID失敗",
			teamID:      validTeamID,
			wantErr:     true,
			errContains: "チームメンバーID指定でのペアリング情報の取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.teamMember.EXPECT().GetByTeamID(validTeamID).Return([]entity.TeamMember{
					newTestTeamMember(1, validTeamID, validMemberID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMemberID).Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:        "異常系_GetByPairingID失敗",
			teamID:      validTeamID,
			wantErr:     true,
			errContains: "ペアリングID指定でのミーティング情報の取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.teamMember.EXPECT().GetByTeamID(validTeamID).Return([]entity.TeamMember{
					newTestTeamMember(1, validTeamID, validMemberID),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMemberID).Return(newTestPairing(pairingID1, validMemberID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return(nil, errGetMeetingsFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mocks := newMeetingViewer(t)
			if tt.setup != nil {
				tt.setup(mocks)
			}

			got, gotErr := uc.ListMeetingsByTeam(tt.teamID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ListMeetingsByTeam")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

func TestMeetingViewer_ListMeetingsByDepartment(t *testing.T) {
	tests := []struct {
		name         string
		departmentID uint
		wantErr      bool
		errContains  string
		setup        func(m *meetingViewerMocks)
		assertFn     func(t *testing.T, got []entity.Meeting)
	}{
		{
			name:         "通常系_部署ユーザのミーティングを取得",
			departmentID: validDepartmentID,
			wantErr:      false,
			setup: func(m *meetingViewerMocks) {
				m.user.EXPECT().GetByDepartmentID(validDepartmentID).Return([]entity.User{
					*newTestUser(validMenteeID, "dept-user"),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(newTestPairing(pairingID1, validMenteeID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return([]entity.Meeting{newTestMeeting(meetingID1, pairingID1)}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, []uint{meetingID1}, got)
			},
		},
		{
			name:         "通常系_部署ユーザが空",
			departmentID: validDepartmentID,
			wantErr:      false,
			setup: func(m *meetingViewerMocks) {
				m.user.EXPECT().GetByDepartmentID(validDepartmentID).Return([]entity.User{}, nil)
			},
			assertFn: func(t *testing.T, got []entity.Meeting) {
				assertMeetingIDSet(t, nil, got)
			},
		},
		{
			name:         "異常系_GetByDepartmentID失敗",
			departmentID: validDepartmentID,
			wantErr:      true,
			errContains:  "部署ID指定でのユーザ取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.user.EXPECT().GetByDepartmentID(validDepartmentID).Return(nil, errGetUsersFailed)
			},
		},
		{
			name:         "異常系_GetByMenteeID失敗",
			departmentID: validDepartmentID,
			wantErr:      true,
			errContains:  "部署ID指定でのペアリング情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.user.EXPECT().GetByDepartmentID(validDepartmentID).Return([]entity.User{
					*newTestUser(validMenteeID, "dept-user"),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, errGetPairingsFailed)
			},
		},
		{
			name:         "異常系_GetByPairingID失敗",
			departmentID: validDepartmentID,
			wantErr:      true,
			errContains:  "ペアリングID指定でのミーティング情報取得に失敗しました",
			setup: func(m *meetingViewerMocks) {
				m.user.EXPECT().GetByDepartmentID(validDepartmentID).Return([]entity.User{
					*newTestUser(validMenteeID, "dept-user"),
				}, nil)
				m.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(newTestPairing(pairingID1, validMenteeID, validMentorID), nil)
				m.meeting.EXPECT().GetByPairingID(pairingID1).Return(nil, errGetMeetingsFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mocks := newMeetingViewer(t)
			if tt.setup != nil {
				tt.setup(mocks)
			}

			got, gotErr := uc.ListMeetingsByDepartment(tt.departmentID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "ListMeetingsByDepartment")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, got)
			}
		})
	}
}

package pairing_mng_usecase_test

import (
	"testing"

	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"go.uber.org/mock/gomock"
)

type pairingViewerMocks struct {
	pairing          *repomock.MockPairingRepository
	relationshipRepo *repomock.MockUserRelationshipRepository
	teamMember       *repomock.MockTeamMemberRepository
	user             *repomock.MockUserRepository
}

func newPairingViewer(t *testing.T) (*pairing_mng_usecase.PairingViewer, *pairingViewerMocks) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mocks := &pairingViewerMocks{
		pairing:          repomock.NewMockPairingRepository(ctrl),
		relationshipRepo: repomock.NewMockUserRelationshipRepository(ctrl),
		teamMember:       repomock.NewMockTeamMemberRepository(ctrl),
		user:             repomock.NewMockUserRepository(ctrl),
	}
	uc := pairing_mng_usecase.NewPairingViewer(
		mocks.user,
		mocks.relationshipRepo,
		mocks.pairing,
		mocks.teamMember,
	)
	return uc, mocks
}

func TestPairingViewer_ListPairingsByBoss_SkipsNilPairing(t *testing.T) {
	uc, mocks := newPairingViewer(t)
	mocks.relationshipRepo.EXPECT().GetByBossID(validBossID).Return([]entity.UserRelationship{
		newTestUserRelationship(1, validMenteeID, validBossID),
	}, nil)
	mocks.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, nil)

	got, err := uc.ListPairingsByBoss(validBossID)
	assertUsecaseErr(t, err, false, "", "ListPairingsByBoss")
	assertPairingOutputIDSet(t, nil, got)
}

func TestPairingViewer_ListPairingsByTeam_SkipsNilPairing(t *testing.T) {
	uc, mocks := newPairingViewer(t)
	mocks.teamMember.EXPECT().GetByTeamID(validTeamID).Return([]entity.TeamMember{
		newTestTeamMember(1, validTeamID, validMemberID),
	}, nil)
	mocks.pairing.EXPECT().GetByMenteeID(validMemberID).Return(nil, nil)

	got, err := uc.ListPairingsByTeam(validTeamID)
	assertUsecaseErr(t, err, false, "", "ListPairingsByTeam")
	assertPairingOutputIDSet(t, nil, got)
}

func TestPairingViewer_ListPairingsByDepartment_SkipsNilPairing(t *testing.T) {
	uc, mocks := newPairingViewer(t)
	mocks.user.EXPECT().GetByDepartmentID(validDepartmentID).Return([]entity.User{
		*newTestUser(validMenteeID, "dept-user"),
	}, nil)
	mocks.pairing.EXPECT().GetByMenteeID(validMenteeID).Return(nil, nil)

	got, err := uc.ListPairingsByDepartment(validDepartmentID)
	assertUsecaseErr(t, err, false, "", "ListPairingsByDepartment")
	assertPairingOutputIDSet(t, nil, got)
}

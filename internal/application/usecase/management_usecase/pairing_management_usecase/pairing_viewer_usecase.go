package pairing_mng_usecase

import (
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"sync"
)

type PairingViewer struct {
	userRepo         repository_interface.UserRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
	relationshipRepo repository_interface.UserRelationshipRepository
	pairingRepo      repository_interface.PairingRepository
}

func NewPairingViewer(
	userRepo repository_interface.UserRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
	pairingRepo repository_interface.PairingRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
) *PairingViewer {
	return &PairingViewer{
		pairingRepo:      pairingRepo,
		relationshipRepo: relationshipRepo,
		teamMemberRepo:   teamMemberRepo,
		userRepo:         userRepo,
	}
}

func (uc *PairingViewer) GetPairing(pairingID uint) (*outputmodel.PairingOutput, error) {
	pairing, err := uc.pairingRepo.Get(pairingID)
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報のID指定取得に失敗しました: %w", err)
	}
	outputPairing := converter.PairingEntityToOutput(*pairing)
	return &outputPairing, nil
}

func (uc *PairingViewer) ListPairingsByAll() ([]outputmodel.PairingOutput, error) {
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報の全取得に失敗しました: %w", err)
	}
	outputPairings := converter.PairingEntitiesToOutputs(pairings)
	return outputPairings, nil
}

func (uc *PairingViewer) ListPairingsByMentor(mentorID uint) ([]outputmodel.PairingOutput, error) {
	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング情報の取得に失敗しました: %w", err)
	}
	outputPairings := converter.PairingEntitiesToOutputs(pairings)
	return outputPairings, nil
}

func (uc *PairingViewer) ListPairingsByBoss(bossID uint) ([]outputmodel.PairingOutput, error) {
	relationships, err := uc.relationshipRepo.GetByBossID(bossID)
	if err != nil {
		return nil, fmt.Errorf("上司ユーザID指定でのミーティングリスト作成のユーザリレーション情報取得に失敗しました: %w", err)
	}
	var (
		result = []outputmodel.PairingOutput{}
		wg     = sync.WaitGroup{}
		mux    = sync.Mutex{}
	)
	errCh := make(chan error, len(relationships))
	for _, relationship := range relationships {
		wg.Add(1)
		relationship := relationship
		go func(menteeID uint, ec chan<- error) {
			defer wg.Done()
			pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
			if err != nil {
				errCh <- fmt.Errorf("部下ユーザID（メンティーユーザID）でのペアリング情報取得に失敗しました: %w", err)
				return
			}
			if pairing == nil {
				return
			}
			outputPairing := converter.PairingEntityToOutput(*pairing)
			mux.Lock()
			result = append(result, outputPairing)
			mux.Unlock()
		}(relationship.SourceUserID, errCh)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (uc *PairingViewer) ListPairingsByTeam(teamID uint) ([]outputmodel.PairingOutput, error) {
	teamMembers, err := uc.teamMemberRepo.GetByTeamID(teamID)
	if err != nil {
		return nil, fmt.Errorf("ペアリングリスト取得処理でチームメンバーの取得に失敗しました: %w", err)
	}

	var (
		result = []outputmodel.PairingOutput{}
		wg     = sync.WaitGroup{}
		mux    = sync.Mutex{}
	)
	errCh := make(chan error, len(teamMembers))

	for _, teamMember := range teamMembers {
		wg.Add(1)
		teamMember := teamMember

		go func(teamMemberID uint, ec chan<- error) {
			defer wg.Done()

			pairing, err := uc.pairingRepo.GetByMenteeID(teamMemberID)
			if err != nil {
				ec <- fmt.Errorf("メンバーID指定でのペアリング情報の取得に失敗しました: %w", err)
				return
			}
			if pairing == nil {
				return
			}

			outputPairing := converter.PairingEntityToOutput(*pairing)
			mux.Lock()
			result = append(result, outputPairing)
			mux.Unlock()
		}(teamMember.MemberID, errCh)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (uc *PairingViewer) ListPairingsByDepartment(departmentID uint) ([]outputmodel.PairingOutput, error) {
	users, err := uc.userRepo.GetByDepartmentID(departmentID)
	if err != nil {
		return nil, fmt.Errorf("部署所属ユーザの取得に失敗しました: %w", err)
	}
	var (
		result = []outputmodel.PairingOutput{}
		wg     = sync.WaitGroup{}
		mux    = sync.Mutex{}
	)
	errCh := make(chan error, len(users))
	for _, user := range users {
		wg.Add(1)
		go func(menteeID uint) {
			defer wg.Done()
			pairing, err := uc.pairingRepo.GetByMenteeID(menteeID)
			if err != nil {
				return
			}
			if pairing == nil {
				return
			}
			outputPairing := converter.PairingEntityToOutput(*pairing)
			mux.Lock()
			result = append(result, outputPairing)
			mux.Unlock()
		}(user.ID)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

package pairing_mng_usecase

import (
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"sync"

	"github.com/google/uuid"
)

type PairingViewer struct {
	userRepo         repository_interface.UserRepository
	teamMemberRepo   repository_interface.TeamMemberRepository
	teamRepo         repository_interface.TeamRepository
	deptRepo         repository_interface.DepartmentRepository
	relationshipRepo repository_interface.UserRelationshipRepository
	pairingRepo      repository_interface.PairingRepository
}

func NewPairingViewer(
	userRepo repository_interface.UserRepository,
	relationshipRepo repository_interface.UserRelationshipRepository,
	pairingRepo repository_interface.PairingRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	teamRepo repository_interface.TeamRepository,
	deptRepo repository_interface.DepartmentRepository,
) *PairingViewer {
	return &PairingViewer{
		pairingRepo:      pairingRepo,
		relationshipRepo: relationshipRepo,
		teamMemberRepo:   teamMemberRepo,
		teamRepo:         teamRepo,
		deptRepo:         deptRepo,
		userRepo:         userRepo,
	}
}

func (uc *PairingViewer) enrichOutputs(outs []outputmodel.PairingOutput) ([]outputmodel.PairingOutput, error) {
	users, err := uc.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	userNames := make(map[uint]string, len(users))
	for _, u := range users {
		userNames[u.ID] = u.Name
	}
	teams, err := uc.teamRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamNames := make(map[uint]string, len(teams))
	for _, t := range teams {
		teamNames[t.ID] = t.Name
	}
	depts, err := uc.deptRepo.GetAll()
	if err != nil {
		return nil, err
	}
	deptNames := make(map[uint]string, len(depts))
	for _, d := range depts {
		deptNames[d.ID] = d.Name
	}
	for i := range outs {
		for j := range outs[i].Targets {
			t := &outs[i].Targets[j]
			if t.Label != "" {
				continue
			}
			switch t.Kind {
			case entity.PairingTargetKindUser:
				t.Label = userNames[t.TargetID]
			case entity.PairingTargetKindTeam:
				t.Label = teamNames[t.TargetID]
			case entity.PairingTargetKindDepartment:
				t.Label = deptNames[t.TargetID]
			}
			if t.Label == "" {
				t.Label = fmt.Sprintf("#%d", t.TargetID)
			}
		}
	}
	return outs, nil
}

func (uc *PairingViewer) GetPairingByUUID(pairingUUID uuid.UUID) (*outputmodel.PairingOutput, error) {
	pairing, err := uc.pairingRepo.GetByUUID(pairingUUID)
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報のUUID指定取得に失敗しました: %w", err)
	}
	outputPairing := converter.PairingEntityToOutput(*pairing)
	outs, err := uc.enrichOutputs([]outputmodel.PairingOutput{outputPairing})
	if err != nil {
		return nil, err
	}
	return &outs[0], nil
}

func (uc *PairingViewer) GetPairing(pairingID uint) (*outputmodel.PairingOutput, error) {
	pairing, err := uc.pairingRepo.Get(pairingID)
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報のID指定取得に失敗しました: %w", err)
	}
	outputPairing := converter.PairingEntityToOutput(*pairing)
	outs, err := uc.enrichOutputs([]outputmodel.PairingOutput{outputPairing})
	if err != nil {
		return nil, err
	}
	return &outs[0], nil
}

func (uc *PairingViewer) ListPairingsByAll() ([]outputmodel.PairingOutput, error) {
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報の全取得に失敗しました: %w", err)
	}
	permanent := make([]entity.Pairing, 0, len(pairings))
	for _, p := range pairings {
		if !p.IsTemporary {
			permanent = append(permanent, p)
		}
	}
	outputPairings := converter.PairingEntitiesToOutputs(permanent)
	return uc.enrichOutputs(outputPairings)
}

func (uc *PairingViewer) ListPairingsByMentor(mentorID uint) ([]outputmodel.PairingOutput, error) {
	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング情報の取得に失敗しました: %w", err)
	}
	outputPairings := converter.PairingEntitiesToOutputs(pairings)
	return uc.enrichOutputs(outputPairings)
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
	return uc.enrichOutputs(result)
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
	return uc.enrichOutputs(result)
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
	return uc.enrichOutputs(result)
}

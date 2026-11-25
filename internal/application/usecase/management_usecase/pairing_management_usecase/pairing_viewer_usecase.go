package pairing_mng_usecase

import (
	"context"
	"fmt"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	logutils "growlink/internal/shared/logUtils"
	"sync"

	"github.com/google/uuid"
)

type PairingViewer struct {
	userRepo       repository_interface.UserRepository
	teamMemberRepo repository_interface.TeamMemberRepository
	teamRepo       repository_interface.TeamRepository
	deptRepo       repository_interface.DepartmentRepository
	pairingRepo    repository_interface.PairingRepository
}

func NewPairingViewer(
	userRepo repository_interface.UserRepository,
	pairingRepo repository_interface.PairingRepository,
	teamMemberRepo repository_interface.TeamMemberRepository,
	teamRepo repository_interface.TeamRepository,
	deptRepo repository_interface.DepartmentRepository,
) *PairingViewer {
	return &PairingViewer{
		pairingRepo:    pairingRepo,
		teamMemberRepo: teamMemberRepo,
		teamRepo:       teamRepo,
		deptRepo:       deptRepo,
		userRepo:       userRepo,
	}
}

func (uc *PairingViewer) GetPairingByUUID(pairingUUID uuid.UUID) (ret0 *outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.GetPairingByUUID"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairing, err := uc.pairingRepo.GetByUUID(pairingUUID)
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報のUUID指定取得に失敗しました: %w", err)
	}
	outputPairing := converter.PairingEntityToOutput(*pairing)
	return &outputPairing, nil
}

func (uc *PairingViewer) GetPairing(pairingID uint) (ret0 *outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.GetPairing"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairing, err := uc.pairingRepo.Get(pairingID)
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報のID指定取得に失敗しました: %w", err)
	}
	outputPairing := converter.PairingEntityToOutput(*pairing)
	return &outputPairing, nil
}

func (uc *PairingViewer) ListPairingsByAll() (ret0 []outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.ListPairingsByAll"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
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
	return converter.PairingEntitiesToOutputs(permanent), nil
}

func (uc *PairingViewer) ListPairingsByAllExTemporary() (ret0 []outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.ListPairingsByAll"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairings, err := uc.pairingRepo.GetAllExTemporary()
	if err != nil {
		return nil, fmt.Errorf("ペアリング情報の全取得に失敗しました: %w", err)
	}
	permanent := make([]entity.Pairing, 0, len(pairings))
	for _, p := range pairings {
		permanent = append(permanent, p)
	}
	return converter.PairingEntitiesToOutputs(permanent), nil
}

func (uc *PairingViewer) ListPairingsByMentor(mentorID uint) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.ListPairingsByMentor"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)
	if err != nil {
		return nil, fmt.Errorf("メンターID指定でのペアリング情報の取得に失敗しました: %w", err)
	}
	return converter.PairingEntitiesToOutputs(pairings), nil
}

func (uc *PairingViewer) ListPairingsByTeam(teamID uint) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.ListPairingsByTeam"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
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

			pairings, err := uc.pairingRepo.GetByMenteeID(teamMemberID)
			if err != nil {
				ec <- fmt.Errorf("メンバーID指定でのペアリング情報の取得に失敗しました: %w", err)
				return
			}
			mux.Lock()
			for _, pairing := range pairings {
				result = append(result, converter.PairingEntityToOutput(pairing))
			}
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

func (uc *PairingViewer) ListPairingsByDepartment(departmentID uint) (ret0 []outputmodel.PairingOutput, err error) {
	const op = "pairing_mng_usecase.ListPairingsByDepartment"
	done := logutils.Begin(context.Background(), logutils.LayerUsecase, op)
	defer func() {
		done(err)
	}()
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
			pairings, err := uc.pairingRepo.GetByMenteeID(menteeID)
			if err != nil {
				return
			}
			mux.Lock()
			for _, pairing := range pairings {
				result = append(result, converter.PairingEntityToOutput(pairing))
			}
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

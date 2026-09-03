package meeting_mng_usecase



import (
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	repository_interface "growlink/internal/domain/interface/repository"
	"sync"
)



type MeetingViewer struct {

	userRepo             repository_interface.UserRepository

	teamMemberRepo       repository_interface.TeamMemberRepository

	meetingRepo          repository_interface.MeetingRepository

	pairingRepo          repository_interface.PairingRepository

	userRelationshipRepo repository_interface.UserRelationshipRepository

}



func NewMeetingViewer(

	meetingRepo repository_interface.MeetingRepository, pairingRepo repository_interface.PairingRepository,

	userRelationshipRepo repository_interface.UserRelationshipRepository,

	userRepo repository_interface.UserRepository, teamMemberRepo repository_interface.TeamMemberRepository) *MeetingViewer {

	return &MeetingViewer{

		meetingRepo:          meetingRepo,

		pairingRepo:          pairingRepo,

		userRepo:             userRepo,

		teamMemberRepo:       teamMemberRepo,

		userRelationshipRepo: userRelationshipRepo,

	}

}



func (uc *MeetingViewer) ListMeetingsByMentor(mentorID uint) ([]outputmodel.MeetingOutput, error) {

	pairings, err := uc.pairingRepo.GetByMentorID(mentorID)

	if err != nil {

		return nil, fmt.Errorf("ミーティングリストのメンターID指定取得に失敗しました: %w", err)

	}



	var (

		resultMeetings = []outputmodel.MeetingOutput{}

		wg             = sync.WaitGroup{}

		mux            = sync.Mutex{}

	)



	errCh := make(chan error, len(pairings))

	for _, pairing := range pairings {

		wg.Add(1)

		pairing := pairing

		go func(pairingID uint, ec chan<- error) {

			defer wg.Done()



			meetings, err := uc.meetingRepo.GetByPairingID(pairingID)

			if err != nil {

				ec <- fmt.Errorf("メンターID指定でのミーティングリスト情報の取得に失敗しました: %w", err)

				return

			}



			mux.Lock()

			resultMeetings = append(resultMeetings, converter.MeetingEntitiesToOutputs(meetings)...)

			mux.Unlock()

		}(pairing.ID, errCh)

	}



	wg.Wait()

	close(errCh)



	for err := range errCh {

		if err != nil {

			return nil, err

		}

	}



	return resultMeetings, nil

}



func (uc *MeetingViewer) ListMeetingsByBoss(bossID uint) ([]outputmodel.MeetingOutput, error) {

	relationships, err := uc.userRelationshipRepo.GetByBossID(bossID)

	if err != nil {

		return nil, fmt.Errorf("上司ユーザID指定でのミーティングリスト作成のユーザリレーション情報取得に失敗しました: %w", err)

	}



	var (

		result = []outputmodel.MeetingOutput{}

		wg     = sync.WaitGroup{}

		mux    = sync.Mutex{}

	)

	errCh := make(chan error, len(relationships))



	for _, relationship := range relationships {

		wg.Add(1)

		relationship := relationship



		go func(menteeID uint, ec chan<- error) {

			defer wg.Done()



			pairings, err := uc.pairingRepo.GetByMenteeID(menteeID)

			if err != nil {

				errCh <- fmt.Errorf("部下ユーザID（メンティーユーザID）でのペアリング情報取得に失敗しました: %w", err)

				return

			}

			meeting, err := uc.meetingRepo.GetByPairingID(pairings.ID)

			if err != nil {

				errCh <- fmt.Errorf("ペアリングユーザID指定でのミーティング情報取得に失敗しました: %w", err)

				return

			}

			mux.Lock()

			result = append(result, converter.MeetingEntitiesToOutputs(meeting)...)

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



func (uc *MeetingViewer) ListMeetingsByTeam(teamID uint) ([]outputmodel.MeetingOutput, error) {

	teamMembers, err := uc.teamMemberRepo.GetByTeamID(teamID)

	if err != nil {

		return nil, fmt.Errorf("チームID指定でのチームメンバー情報取得に失敗しました: %w", err)

	}



	var (

		resultMeetings = []outputmodel.MeetingOutput{}

		wg             = sync.WaitGroup{}

		mux            = sync.Mutex{}

	)



	errCh := make(chan error, len(teamMembers))



	for _, teamMebmer := range teamMembers {

		wg.Add(1)

		teamMebmer := teamMebmer



		go func(memberID uint, ec chan<- error) {

			defer wg.Done()

			pairing, err := uc.pairingRepo.GetByMenteeID(memberID)

			if err != nil {

				ec <- fmt.Errorf("チームメンバーID指定でのペアリング情報の取得に失敗しました: %w", err)

				return

			}



			meetings, err := uc.meetingRepo.GetByPairingID(pairing.ID)

			if err != nil {

				ec <- fmt.Errorf("ペアリングID指定でのミーティング情報の取得に失敗しました: %w", err)

				return

			}



			mux.Lock()

			resultMeetings = append(resultMeetings, converter.MeetingEntitiesToOutputs(meetings)...)

			mux.Unlock()

		}(teamMebmer.MemberID, errCh)

	}



	wg.Wait()

	close(errCh)



	for err := range errCh {

		if err != nil {

			return nil, err

		}

	}



	return resultMeetings, nil

}



func (uc *MeetingViewer) ListMeetingsByDepartment(departmentID uint) ([]outputmodel.MeetingOutput, error) {

	users, err := uc.userRepo.GetByDepartmentID(departmentID)

	if err != nil {

		return nil, fmt.Errorf("部署ID指定でのユーザ取得に失敗しました: %w", err)

	}



	var (

		resultMeetings = []outputmodel.MeetingOutput{}

		wg             = sync.WaitGroup{}

		mux            = sync.Mutex{}

	)



	errCh := make(chan error, len(users))

	for _, user := range users {

		wg.Add(1)

		user := user



		go func(userID uint, ec chan<- error) {

			defer wg.Done()



			pairings, err := uc.pairingRepo.GetByMenteeID(userID)

			if err != nil {

				ec <- fmt.Errorf("部署ID指定でのペアリング情報取得に失敗しました: %w", err)

				return

			}

			meetings, err := uc.meetingRepo.GetByPairingID(pairings.ID)

			if err != nil {

				ec <- fmt.Errorf("ペアリングID指定でのミーティング情報取得に失敗しました: %w", err)

				return

			}



			mux.Lock()

			resultMeetings = append(resultMeetings, converter.MeetingEntitiesToOutputs(meetings)...)

			mux.Unlock()

		}(user.ID, errCh)

	}



	wg.Wait()

	close(errCh)



	for err := range errCh {

		if err != nil {

			return nil, err

		}

	}



	return resultMeetings, nil

}


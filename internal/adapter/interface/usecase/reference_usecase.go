package usecase_interface

type (
	ReferenceUsecase interface {
		ReferenceChecker
	}

	ReferenceChecker interface {
		CheckReferecneByUser(userID, targetUserID uint) (isOK bool, err error)
		CheckReferecneByTeam(userID, teamID uint) (isOK bool, err error)
		CheckReferecneByDepartment(userID, departmentID uint) (isOK bool, err error)
	}
)

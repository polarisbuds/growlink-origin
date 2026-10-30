package authctx

import "growlink/internal/domain/valueobject"

// CanModify returns true when the user may update/delete a resource they created, or is system admin.
func CanModify(userID uint, systemRole valueobject.SystemRole, createdByUserID uint) bool {
	if systemRole == valueobject.SystemRoleAdmin {
		return true
	}
	if createdByUserID == 0 {
		return false
	}
	return userID == createdByUserID
}

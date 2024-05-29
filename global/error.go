package global

import "errors"

var (
	ErrPermissionInUse     = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleAlreadyAssigned = errors.New("this role is already assigned to the user")
	ErrRoleInUse           = errors.New("cannot delete assigned role")
	ErrRoleNotFound        = errors.New("role not found")
	ErrGroupNotFound       = errors.New("group not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists		   = errors.New("user already exists")
	ErrOrgExists 		   = errors.New("organization already exists")
)

func CreateDuplicateErr(r string) (error) {
	return errors.New(r+" already exist")
}
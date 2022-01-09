package subordinates

import (
	"errors"
	"fmt"
)

var (
	ErrTargetUserNotFound = errors.New("Target user not found")
	ErrRoleNotFound       = errors.New("Role not found")
)

type Finder struct {
	roles map[int]Role
	users map[int]User
}

func NewFinder() *Finder {
	return &Finder{}
}

func (f *Finder) SetRoles(roles []Role) {
	f.roles = make(map[int]Role, len(roles))
	for _, role := range roles {
		f.roles[role.Id] = role
	}
}

func (f *Finder) SetUsers(users []User) {
	f.users = make(map[int]User, len(users))
	for _, user := range users {
		f.users[user.Id] = user
	}
}

func (f *Finder) GetSubordinates(userId int) ([]User, error) {
	// Get the target user
	targetUser, found := f.users[userId]
	if !found {
		return nil, fmt.Errorf("%w: %v", ErrTargetUserNotFound, userId)
	}

	// Check that the target user's role exists
	_, found = f.roles[targetUser.Role]
	if !found {
		return nil, fmt.Errorf("%w: %v", ErrRoleNotFound, targetUser.Role)
	}

	// Get all users with roles that are subordinate to the target user's role
	var results []User
	for _, user := range f.users {
		subordinate, err := f.isRoleSubordinate(user.Role, targetUser.Role)
		if err != nil {
			return nil, err
		}
		if subordinate {
			results = append(results, user)
		}
	}
	return results, nil
}

// Check if the role with checkRoleId is subordinate to the role with againstRoleId
// Returns an error (ErrRoleNotFound) if a role is not found
func (f *Finder) isRoleSubordinate(checkRoleId, againstRoleId int) (bool, error) {
	// A role cannot be subordinate to itself
	if checkRoleId == againstRoleId {
		return false, nil
	}

	// Get the role to check
	checkRole, found := f.roles[checkRoleId]
	if !found {
		return false, fmt.Errorf("%w: %v", ErrRoleNotFound, checkRoleId)
	}

	// A role is subordinate to its parent
	if checkRole.Parent == againstRoleId {
		return true, nil
	}

	// If there is a parent role, check against that
	if checkRole.Parent > 0 {
		return f.isRoleSubordinate(checkRole.Parent, againstRoleId)
	}

	// Otherwise it is not subordinate
	return false, nil
}

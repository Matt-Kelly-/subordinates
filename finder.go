/*
   Package subordinates provides a means for finding subordinates of a user.
   Each user has a role, and each role may have a parent role.
   User A is subordinate to user B if user B's role is an ancestor of user A's role.

   Usage
   1. Create a new subordinate finder with NewFinder()
   2. Call SetRoles() on the finder, passing a slice of roles
   3. Call SetUsers() on the finder, passing a slice of users
   4. Call GetSubordinates() on the finder, passing a user ID
   NB. Steps 2 and 3 can be done in any order

   GetSubordinates returns a slice of subordinate users, or a nil slice and an error if:
   - The target user is not found (ErrTargetUserNotFound)
   - The role of the target user, or any other user, is not found (ErrRoleNotFound)

   Subordinate finder instances are reusable - you can call any function on Finder as many times as you like.
*/

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
	memo := make(map[int]bool)
	for _, user := range f.users {
		subordinate, err := f.isRoleSubordinate(user.Role, targetUser.Role, memo)
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
func (f *Finder) isRoleSubordinate(checkRoleId, againstRoleId int, memo map[int]bool) (bool, error) {
	// A role cannot be subordinate to itself
	if checkRoleId == againstRoleId {
		return false, nil
	}

	// Return a memoized result if there is one
	memoResult, found := memo[checkRoleId]
	if found {
		return memoResult, nil
	}

	// Get the role to check
	checkRole, found := f.roles[checkRoleId]
	if !found {
		return false, fmt.Errorf("%w: %v", ErrRoleNotFound, checkRoleId)
	}

	var result bool
	if checkRole.Parent == againstRoleId {
		// A role is subordinate to its parent
		result = true
	} else if checkRole.Parent > 0 {
		// If there is a parent role, check against that
		var err error
		result, err = f.isRoleSubordinate(checkRole.Parent, againstRoleId, memo)
		if err != nil {
			return false, err
		}
	}
	// Otherwise it is not subordinate

	// Memoize result
	memo[checkRoleId] = result

	return result, nil
}

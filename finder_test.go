package subordinates_test

import (
	"github.com/Matt-Kelly-/subordinates"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSubordinatesWithSampleData(t *testing.T) {
	roles := []subordinates.Role{
		{
			Id:     1,
			Name:   "System Administrator",
			Parent: 0,
		},
		{
			Id:     2,
			Name:   "Location Manager",
			Parent: 1,
		},
		{
			Id:     3,
			Name:   "Supervisor",
			Parent: 2,
		},
		{
			Id:     4,
			Name:   "Employee",
			Parent: 3,
		},
		{
			Id:     5,
			Name:   "Trainer",
			Parent: 3,
		},
	}
	users := []subordinates.User{
		{
			Id:   1,
			Name: "Adam Admin",
			Role: 1,
		},
		{
			Id:   2,
			Name: "Emily Employee",
			Role: 4,
		},
		{
			Id:   3,
			Name: "Sam Supervisor",
			Role: 3,
		},
		{
			Id:   4,
			Name: "Mary Manager",
			Role: 2,
		},
		{
			Id:   5,
			Name: "Steve Trainer",
			Role: 5,
		},
	}

	finder := subordinates.NewFinder()
	finder.SetRoles(roles)
	finder.SetUsers(users)

	t.Run("sample case 1", func(t *testing.T) {
		results, err := finder.GetSubordinates(3)
		require.Nil(t, err)
		require.ElementsMatch(t, results, []subordinates.User{
			{
				Id:   2,
				Name: "Emily Employee",
				Role: 4,
			},
			{

				Id:   5,
				Name: "Steve Trainer",
				Role: 5,
			},
		})
	})

	t.Run("sample case 2", func(t *testing.T) {
		results, err := finder.GetSubordinates(1)
		require.Nil(t, err)
		require.ElementsMatch(t, results, []subordinates.User{
			{
				Id:   2,
				Name: "Emily Employee",
				Role: 4,
			},
			{
				Id:   3,
				Name: "Sam Supervisor",
				Role: 3,
			},
			{
				Id:   4,
				Name: "Mary Manager",
				Role: 2,
			},
			{
				Id:   5,
				Name: "Steve Trainer",
				Role: 5,
			},
		})
	})

	t.Run("user has no subordinates", func(t *testing.T) {
		results, err := finder.GetSubordinates(2)
		require.Nil(t, err)
		require.Empty(t, results)
	})
}

// Should return an error if the target user is not found
func TestGetSubordinatesWhenTargetUserNotFound(t *testing.T) {
	finder := subordinates.NewFinder()

	results, err := finder.GetSubordinates(1)

	require.NotNil(t, err)
	require.ErrorIs(t, err, subordinates.ErrTargetUserNotFound)
	require.Empty(t, results)
}

// Should return an error if the target user's role is not found
func TestGetSubordinatesWhenTargetUserRoleNotFound(t *testing.T) {
	finder := subordinates.NewFinder()
	finder.SetUsers([]subordinates.User{
		{
			Id:   1,
			Name: "Test user 1",
			Role: 1,
		},
	})

	results, err := finder.GetSubordinates(1)

	require.NotNil(t, err)
	require.ErrorIs(t, err, subordinates.ErrRoleNotFound)
	require.Empty(t, results)
}

// Should return an error if another user's role is not found
func TestGetSubordinatesWhenOtherUserRoleNotFound(t *testing.T) {
	finder := subordinates.NewFinder()
	finder.SetRoles([]subordinates.Role{
		{
			Id:     1,
			Name:   "Test role 1",
			Parent: 0,
		},
	})
	finder.SetUsers([]subordinates.User{
		{
			Id:   1,
			Name: "Test user 1",
			Role: 1,
		},
		{
			Id:   2,
			Name: "Test user 2",
			Role: 2,
		},
	})

	results, err := finder.GetSubordinates(1)

	require.NotNil(t, err)
	require.ErrorIs(t, err, subordinates.ErrRoleNotFound)
	require.Empty(t, results)
}

// Should be able to replace the roles with another set of roles and get subordinates again
func TestSwitchingRoles(t *testing.T) {
	roles1 := []subordinates.Role{
		{
			Id:     1,
			Name:   "Test role 1",
			Parent: 0,
		},
		{
			Id:     2,
			Name:   "Test role 2",
			Parent: 1,
		},
	}
	roles2 := []subordinates.Role{
		{
			Id:     1,
			Name:   "Test role 1",
			Parent: 1,
		},
		{
			Id:     2,
			Name:   "Test role 2",
			Parent: 0,
		},
	}
	users := []subordinates.User{
		{
			Id:   1,
			Name: "Test user 1",
			Role: 1,
		},
		{
			Id:   2,
			Name: "Test user 2",
			Role: 2,
		},
	}

	// First set of roles
	finder := subordinates.NewFinder()
	finder.SetRoles(roles1)
	finder.SetUsers(users)

	results, err := finder.GetSubordinates(1)

	require.Nil(t, err)
	require.ElementsMatch(t, results, []subordinates.User{
		{
			Id:   2,
			Name: "Test user 2",
			Role: 2,
		},
	})

	// Second set of roles
	finder.SetRoles(roles2)

	results, err = finder.GetSubordinates(1)

	require.Nil(t, err)
	require.Empty(t, results)
}

// Should be able to replace the users with another set of users and get subordinates again
func TestSwitchingUsers(t *testing.T) {
	roles := []subordinates.Role{
		{
			Id:     1,
			Name:   "Test role 1",
			Parent: 0,
		},
		{
			Id:     2,
			Name:   "Test role 2",
			Parent: 1,
		},
	}
	users1 := []subordinates.User{
		{
			Id:   1,
			Name: "Test user 1",
			Role: 1,
		},
		{
			Id:   2,
			Name: "Test user 2",
			Role: 2,
		},
	}
	users2 := []subordinates.User{
		{
			Id:   1,
			Name: "Test user 1",
			Role: 2,
		},
		{
			Id:   2,
			Name: "Test user 2",
			Role: 1,
		},
	}

	// First set of users
	finder := subordinates.NewFinder()
	finder.SetRoles(roles)
	finder.SetUsers(users1)

	results, err := finder.GetSubordinates(1)

	require.Nil(t, err)
	require.ElementsMatch(t, results, []subordinates.User{
		{
			Id:   2,
			Name: "Test user 2",
			Role: 2,
		},
	})

	// Second set of users
	finder.SetUsers(users2)

	results, err = finder.GetSubordinates(1)

	require.Nil(t, err)
	require.Empty(t, results)
}

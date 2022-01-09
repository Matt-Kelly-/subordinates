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

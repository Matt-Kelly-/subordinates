package subordinates_test

import (
	"fmt"
	"github.com/Matt-Kelly-/subordinates"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
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

// Benchmark for different numbers of roles and users, with a random unbalanced binary tree of roles
// This should be something like the average case
func BenchmarkGetSubordinatesWithRoleTree(b *testing.B) {
	inputSizes := []int{100, 200, 400, 800, 1600}
	rand.Seed(time.Now().UnixNano())

	for _, roleCount := range inputSizes {
		// Generate random tree of roles
		roles := generateTreeOfRoles(roleCount)

		for _, userCount := range inputSizes {

			// Generate random users and assign a random role to each user
			users := make([]subordinates.User, userCount)
			for i := 0; i < userCount; i++ {
				users[i] = subordinates.User{
					Id:   i + 1,
					Name: fmt.Sprintf("Test user %v", i+1),
					Role: roles[rand.Intn(roleCount)].Id,
				}
			}

			// Set up finder
			finder := subordinates.NewFinder()
			finder.SetRoles(roles)
			finder.SetUsers(users)

			// Benchmark GetSubordinates()
			b.Run(fmt.Sprintf("%v roles and %v users", roleCount, userCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, err := finder.GetSubordinates(users[0].Id)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Generate random unbalanced binary tree of roles
func generateTreeOfRoles(count int) []subordinates.Role {
	// Binary tree node for role
	type roleNode struct {
		role       subordinates.Role
		leftChild  *roleNode
		rightChild *roleNode
	}

	// Generate roles with random IDs and insert them into unbalanced binary tree
	var root *roleNode
	roles := make([]subordinates.Role, count)
	for i := 0; i < count; i++ {
		// Create node for with random ID and no parent
		roleId := rand.Int()
		node := &roleNode{
			role: subordinates.Role{
				Id:     roleId,
				Name:   fmt.Sprintf("Test role %v", roleId),
				Parent: 0,
			},
		}

		// Traverse the tree to find where to insert (tree is ordered by role ID)
		var parentNode *roleNode
		insertNode := root
		for insertNode != nil {
			parentNode = insertNode
			if roleId < insertNode.role.Id {
				insertNode = insertNode.leftChild
			} else {
				insertNode = insertNode.rightChild
			}
		}

		// Insert node into tree
		if parentNode == nil {
			root = node
		} else {
			node.role.Parent = parentNode.role.Id // Set role parent ID
			if roleId < parentNode.role.Id {
				parentNode.leftChild = node
			} else {
				parentNode.rightChild = node
			}
		}

		// Add role to list
		roles[i] = node.role
	}

	// Return roles
	return roles
}

// Benchmark for different numbers of roles and users, with tree of roles that forms a linked list
// This should be the worst case
func BenchmarkGetSubordinatesWithRoleList(b *testing.B) {
	inputSizes := []int{100, 200, 400, 800, 1600}

	for _, roleCount := range inputSizes {
		// Generate list of roles
		roles := make([]subordinates.Role, roleCount)
		for i := 0; i < roleCount; i++ {
			roles[i] = subordinates.Role{
				Id:     i + 1,
				Name:   fmt.Sprintf("Test role %v", i+1),
				Parent: i,
			}
		}

		for _, userCount := range inputSizes {

			// Generate users
			// First user has role at head of list
			// All other users have role at tail of list
			// Will need to check full role list for each user
			users := make([]subordinates.User, userCount)
			for i := 0; i < userCount; i++ {
				users[i] = subordinates.User{
					Id:   i + 1,
					Name: fmt.Sprintf("Test user %v", i+1),
					Role: roles[roleCount-1].Id,
				}
			}
			users[0].Role = roles[0].Id

			// Set up finder
			finder := subordinates.NewFinder()
			finder.SetRoles(roles)
			finder.SetUsers(users)

			// Benchmark GetSubordinates()
			b.Run(fmt.Sprintf("%v roles and %v users", roleCount, userCount), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, err := finder.GetSubordinates(users[0].Id)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// Benchmark for different numbers users, all with the same role
func BenchmarkGetSubordinatesWithSingleRole(b *testing.B) {
	userCounts := []int{100, 200, 400, 800, 1600}

	roles := []subordinates.Role{
		{
			Id:     1,
			Name:   "Test role 1",
			Parent: 0,
		},
	}

	for _, userCount := range userCounts {

		// Generate users, all with the same role
		users := make([]subordinates.User, userCount)
		for i := 0; i < userCount; i++ {
			users[i] = subordinates.User{
				Id:   i + 1,
				Name: fmt.Sprintf("Test user %v", i+1),
				Role: roles[0].Id,
			}
		}

		// Set up finder
		finder := subordinates.NewFinder()
		finder.SetRoles(roles)
		finder.SetUsers(users)

		// Benchmark GetSubordinates()
		b.Run(fmt.Sprintf("%v users", userCount), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := finder.GetSubordinates(users[0].Id)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

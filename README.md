# Subordinate Users

This is a demo for finding subordinate users in a system where each user has a role, and each role may have a parent role.
The roles form a forest (one or more trees). User A is subordinate to user B if user B's role is an ancestor of user A's role.

## Setup

1. `git clone https://github.com/Matt-Kelly-/subordinates`
2. `cd subordinates`
3. `go test -v ./...`

To use this code in another Go project

```
package main

import (
  "github.com/Matt-Kelly-/subordinates"
)

func main() {
  finder := subordinates.NewFinder()
  finder.SetRoles([]subordinates.Role{ ... })
  finder.SetUsers([]subordinates.User{ ... })
  results, err := finder.GetSubordinates(123)
  // Do something with results and err
}
```

## Algorithm

My solution is in the form of a struct `Finder` which stores the set of roles and the set of users that we are working with. The finder can then be queried to find subordinates of a user given the target user's ID.

Since we need to look up each user's role by ID, the finder stores roles in a map keyed by ID. This allows a constant time lookup of each user's role. Once we have retrieved a role, we can also look up the parent role by ID in the same way.

We might want to query the finder several times to find subordinates of different users in the same set of users and roles. The finder stores users in a map keyed by ID to allow for a constant time lookup of the target user for each query.

To find subordinate users, we first look up the target user to get their role, and check that their role exists. We then loop through all users, determine if each is a subordinate of the target user, and if so add them to the results. 

To determine if a user is subordinate to the target user, we look up their role, then:
- If their role is the same as the target user's role, they are not a subordinate
- If their role has a parent role that is the same as the target user's role, they are a subordinate
- If their role has a parent role, we recursively check if the parent role is subordinate
- Otherwise they are not a subordinate

### Complexity

Setting roles and setting users each copy the roles or users into a map. So these are both linear time complexity.

Finding subordinates must check each user, and for each user must check their role and potentially all its ancestor roles. Looping through the users is linear time complexity. If the role tree is reasonably balanced then checking if a user's role is subordinate will be logarithmic complexity. Overall, finding subordinates has linearithmic time complexity or O(U log R) where U is the number of users and R is the number of roles.

![tree-users](https://user-images.githubusercontent.com/6522959/148671156-50783500-8587-43d7-a476-77e1ff4e4a72.png)

![tree-roles](https://user-images.githubusercontent.com/6522959/148671175-453cf880-8c29-4729-bc21-92ab2a7ccc8a.png)

In the worst case, the tree of roles degenerates to a linked list so we must check potentially every role for every user. The complexity here is O(U R) where U is the number of users and R is the number of roles. 

![list-users](https://user-images.githubusercontent.com/6522959/148671253-a8df50ba-7570-4514-8294-a66c52bab3b2.png)

![list-roles](https://user-images.githubusercontent.com/6522959/148671268-1fa08104-9a83-42f9-a588-bc662bf148bb.png)

package subordinates

type Role struct {
	Id     int
	Name   string
	Parent int
}

type User struct {
	Id   int
	Name string
	Role int
}

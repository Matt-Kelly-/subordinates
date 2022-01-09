package subordinates

type Finder struct{}

func NewFinder() *Finder {
	return &Finder{}
}

func (f *Finder) SetRoles(roles []Role) {
}

func (f *Finder) SetUsers(users []User) {
}

func (f *Finder) GetSubordinates(userId int) ([]User, error) {
	return nil, nil
}

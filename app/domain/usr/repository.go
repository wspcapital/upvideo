package usr

type UserRepository interface {
	FindAll(dto UserSearchDto) ([]*User, error)
	Insert(item *User) error
	Update(item *User) error
	Delete(item *User) error
}

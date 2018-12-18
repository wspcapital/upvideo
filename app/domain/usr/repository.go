package usr

type UserRepository interface {
	FindAll(dto UserSearchDto) ([]*User, error)
	Insert(item *User) error
	Update(item *User) error
	Delete(item *User) error
	FindByForgotPasswordToken(dto *UserSearchDto) (*User, error)
	SetForgotPasswordToken(item *User) error
	RemoveForgotPasswordToken(item *User) error
}

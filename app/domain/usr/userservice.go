package usr

import (
	"bitbucket.org/marketingx/upvideo/app/cookie"
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type UserService interface {
	FindAll(dto UserSearchDto) ([]*User, error)
	Insert(item *User) error
	Update(item *User) error
	Delete(item *User) error
	PasswordHash(source string) string
	Login(email string, password string) (*User, error)
	FindById(id string) (*User, error)
	FindByKey(key string) (*User, error)
	FindByEmail(email string) (*User, error)
	CheckUserExists(user *User) (bool, error)
}

type userService struct {
	repo           UserRepository
	sessionService session.Service
}

var UserNotFound = errors.New("User not found")

func (this *userService) FindAll(dto UserSearchDto) ([]*User, error) {
	return this.repo.FindAll(dto)
}

func (this *userService) Login(email string, password string) (*User, error) {
	items, err := this.repo.FindAll(UserSearchDto{Email: email, PasswordHash: this.PasswordHash(password)})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, UserNotFound
	}
	return items[0], err
}

func (this *userService) CheckUserExists(user *User) (bool, error) {
	_, err := this.FindByEmail(user.Email)
	if err != nil {
		if err == UserNotFound {
			return false, nil
		}
	}

	return true, err
}

func (this *userService) FindByEmail(email string) (*User, error) {
	items, err := this.repo.FindAll(UserSearchDto{Email: email})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, UserNotFound
	}
	return items[0], err
}

func (this *userService) FindByKey(apiKey string) (*User, error) {
	items, err := this.repo.FindAll(UserSearchDto{APIKey: apiKey})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, UserNotFound
	}
	return items[0], err
}

func (this *userService) FindById(id string) (*User, error) {
	items, err := this.repo.FindAll(UserSearchDto{Id: id})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, UserNotFound
	}
	return items[0], err
}

func (this *userService) Insert(item *User) error {
	found, _ := this.FindAll(UserSearchDto{Email: item.Email})
	if len(found) == 0 {
		apikey, err := cookie.Generate(128)
		if err != nil {
			return err
		}
		item.APIKey = apikey
		return this.repo.Insert(item)
	} else {
		return errors.New("X001 - Could not create user")
	}

}

func (this *userService) Update(item *User) error {
	return this.repo.Update(item)
}

func (this *userService) Delete(item *User) error {
	return this.repo.Delete(item)
}

func (this *userService) PasswordHash(source string) string {
	hasher := md5.New()
	hasher.Write([]byte(source))
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

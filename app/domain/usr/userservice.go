package usr

import (
	"bitbucket.org/marketingx/upvideo/app/domain/session"
	"bitbucket.org/marketingx/upvideo/utils"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
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
	SetForgotPasswordToken(user *User) error
	SetNewForgottenPassword(user *User) error
	ResetApiKey(user *User) error
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
		item.APIKey = this.generateUniqueApiKey()
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

func (this *userService) SetForgotPasswordToken(user *User) error {
	user.ForgotPasswordToken = utils.SecureRandomString(36)
	err := this.repo.SetForgotPasswordToken(user)
	if err != nil {
		return err
	}
	return err
}

func (this *userService) SetNewForgottenPassword(user *User) error {
	item, err := this.repo.FindByForgotPasswordToken(&UserSearchDto{Email: user.Email, ForgotPasswordToken: user.ForgotPasswordToken})
	if err != nil {
		return err
	}

	item.PasswordHash = user.PasswordHash
	err = this.repo.Update(item)
	if err != nil {
		return err
	}

	err = this.repo.RemoveForgotPasswordToken(item)
	return err
}

func (this *userService) ResetApiKey(user *User) error {
	if user.Id == 0 {
		return UserNotFound
	}

	user.APIKey = this.generateUniqueApiKey()
	return this.repo.Update(user)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (this *userService) generateUniqueApiKey() (apikey string) {
	for {
		apikey = utils.SecureRandomString(36)
		if _, err := this.FindByKey(apikey); err != nil {
			if err != UserNotFound {
				fmt.Println("UserService -> generateUniqueApiKey", err)
			}
			break
		}
	}
	return
}

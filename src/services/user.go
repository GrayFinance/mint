package services

import (
	"fmt"
	"time"

	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) CreateUser() (models.User, error) {
	if len(u.Username) < 6 || len(u.Password) < 6 {
		err := fmt.Errorf("Username / Password invalid.")
		return models.User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		UserID:   uuid.New().String(),
		Username: u.Username,
		Password: string(password),
	}
	if storage.DB.Create(&user).Error != nil {
		err := fmt.Errorf("User already exists.")
		return models.User{}, err
	}
	return user, nil
}

func (u *User) GetUser() (models.User, error) {
	var user models.User

	if storage.DB.Where("user_id = ?", u.UserID).First(&user).Error != nil {
		err := fmt.Errorf("User does not exist.")
		return models.User{}, err
	}
	return user, nil
}

func (u *User) AuthUser() (map[string]string, error) {
	var user models.User

	if storage.DB.Model(user).Where("username = ?", u.Username).First(&user).Error != nil {
		err := fmt.Errorf("Username / Password invalid.")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		err = fmt.Errorf("Username / Password invalid.")
		return nil, err
	}

	t := jwt.New(jwt.SigningMethodHS256)
	c := t.Claims.(jwt.MapClaims)

	c["uid"] = user.UserID
	c["iat"] = time.Now().Add(time.Minute * (60 * 24)).Unix()

	token, err := t.SignedString([]byte(config.Config.SIGN_KEY))
	if err != nil {
		return nil, err
	}
	return map[string]string{"uid": user.UserID, "token": token}, nil
}

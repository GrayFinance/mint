package services

import (
	"fmt"
	"time"

	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/GrayFinance/mint/src/utils"
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
	if len(u.Password) < 6 || len(u.Password) > 16 || len(u.Username) < 3 || len(u.Username) > 16 {
		err := fmt.Errorf("Password / Username is invalid.")
		return models.User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		UserID:       uuid.New().String(),
		Usertag:      u.Username,
		Username:     u.Username,
		Password:     string(password),
		MasterAPIKey: utils.RandomHex(16),
	}

	if storage.DB.Create(&user).Error != nil {
		err := fmt.Errorf("User already exists.")
		return models.User{}, err
	}
	return user, nil
}

func (u *User) GetUser() (interface{}, error) {
	var user models.User

	if storage.DB.Where("user_id = ?", u.UserID).First(&user).Error != nil {
		err := fmt.Errorf("User does not exist.")
		return models.User{}, err
	}

	var balance int64

	storage.DB.Raw(`
	SELECT SUM(balance) 
	FROM wallets 
	WHERE user_id = ?
	`, u.UserID).Scan(&balance)
	return map[string]interface{}{"balance": balance, "master_api_key": user.MasterAPIKey, "user_id": user.UserID}, nil
}

func (u *User) AuthUser() (string, error) {
	if len(u.Username) < 6 || len(u.Password) < 6 {
		err := fmt.Errorf("Username and password must have a size of 6 character.")
		return "", err
	}

	var user models.User
	if storage.DB.Model(user).Where("username = ?", u.Username).First(&user).Error != nil {
		err := fmt.Errorf("Username / Password invalid.")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		err = fmt.Errorf("Username / Password invalid.")
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)

	claim["user_id"] = user.UserID
	claim["exp"] = time.Now().Add(time.Minute * (60 * 24)).Unix()
	return token.SignedString([]byte(config.Config.SIGN_KEY))
}

func (u *User) ChangePassword() (models.User, error) {
	if len(u.Password) < 6 || len(u.Password) > 16 {
		err := fmt.Errorf("New password is invalid.")
		return models.User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if storage.DB.Model(user).Where("user_id = ?", u.UserID).First(&user).Error != nil {
		err := fmt.Errorf("It was not possible to find user.")
		return user, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err == nil {
		err = fmt.Errorf("It was not possible to change the password.")
		return user, err
	}

	if storage.DB.Model(&user).Update("password", string(password)).Error != nil {
		err := fmt.Errorf("It was not possible to change the password.")
		return user, err
	}
	return user, nil
}

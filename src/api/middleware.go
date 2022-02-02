package api

import (
	"fmt"
	"net/http"

	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token_jwt, _, _ := r.BasicAuth()
		if token_jwt == "" {
			utils.SendJSONError(w, 401, "Authentication token not found.")
			return
		}

		token, err := jwt.Parse(token_jwt, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error.")
			}
			return []byte(config.Config.SIGN_KEY), nil
		})

		if err != nil {
			utils.SendJSONError(w, 401, err.Error())
			return
		}

		if token.Valid == true {
			claims := token.Claims.(jwt.MapClaims)
			context.Set(r, "user_id", claims["user_id"])
			next.ServeHTTP(w, r)
			return
		} else {
			utils.SendJSONError(w, 401, "Authentication token is valid.")
			return
		}
	}
}

func WalletMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		master_api_key, wallet_key, _ := r.BasicAuth()
		if len(master_api_key) < 16 || len(wallet_key) < 16 {
			utils.SendJSONError(w, 401, "Authentication token not found.")
			return
		}

		var user models.User
		if storage.DB.Model(user).Where("master_api_key = ?", master_api_key).First(&user).Error != nil {
			utils.SendJSONError(w, 401, "Master key is invalid.")
			return
		}

		var wallet models.Wallet
		if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ?", user.UserID, mux.Vars(r)["wallet_id"]).First(&wallet).Error != nil {
			utils.SendJSONError(w, 401, "Wallet id is invalid.")
			return
		}

		permission := ""
		if wallet.WalletAdminKey == wallet_key {
			permission = "admin"
		}

		if wallet.WalletReadKey == wallet_key {
			permission = "read"
		}

		if permission != "" {
			context.Set(r, "user_id", wallet.UserID)
			context.Set(r, "permission", permission)
			next.ServeHTTP(w, r)
			return
		} else {
			utils.SendJSONError(w, 401, "Wallet Key is invalid.")
			return
		}
	}
}

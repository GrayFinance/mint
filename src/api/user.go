package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GrayFinance/mint/src/services"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/gorilla/context"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	var user services.User
	json.Unmarshal(body, &user)

	data, err := user.CreateUser()
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}
	utils.SendJSONResponse(w, map[string]string{"user_id": data.UserID})
	return
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	var user services.User
	json.Unmarshal(body, &user)

	token, err := user.AuthUser()
	if err != nil {
		utils.SendJSONError(w, 401, err.Error())
		return
	}

	utils.SendJSONResponse(w, map[string]string{"token": token})
	return
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user := services.User{
		UserID: context.Get(r, "user_id").(string),
	}

	data, err := user.GetUser()
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	utils.SendJSONResponse(w, data)
	return
}

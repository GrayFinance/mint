package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/GrayFinance/mint/src/services"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	var params map[string]string
	json.Unmarshal(body, &params)

	wallet := services.Wallet{
		UserID: context.Get(r, "user_id").(string),
	}

	data, err := wallet.CreateWallet(params["label"])
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	utils.SendJSONResponse(w, data)
	return
}

func DeleteWallet(w http.ResponseWriter, r *http.Request) {
	wallet := services.Wallet{
		UserID: context.Get(r, "user_id").(string),
	}

	if _, err := wallet.DeleteWallet(mux.Vars(r)["wallet_id"]); err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}
}

func RenameWallet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	var data map[string]string

	if err := json.Unmarshal(body, &data); err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	if data["label"] == "" {
		utils.SendJSONError(w, 500, "Label not found.")
		return
	}

	wallet := services.Wallet{
		UserID: context.Get(r, "user_id").(string),
	}
	if _, err := wallet.RenameWallet(mux.Vars(r)["wallet_id"], data["label"]); err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}
}

func ListWallets(w http.ResponseWriter, r *http.Request) {
	wallet := services.Wallet{
		UserID: context.Get(r, "user_id").(string),
	}

	data, err := wallet.ListWallets()
	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	utils.SendJSONResponse(w, data)
	return
}

func GetWallet(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" {
		wallet := services.Wallet{
			UserID: context.Get(r, "user_id").(string),
		}

		data, err := wallet.GetWallet(mux.Vars(r)["wallet_id"])
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}

		utils.SendJSONResponse(w, data)
		return

	} else {
		err := "The key permission is not admin."
		utils.SendJSONError(w, 500, err)
		return
	}
}
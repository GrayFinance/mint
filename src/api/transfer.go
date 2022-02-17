package api

import (
	"io/ioutil"
	"net/http"

	"github.com/GrayFinance/mint/src/services"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

func Transfer(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}

		data := gjson.ParseBytes(body)

		destination := data.Get("destination").String()
		if destination == "" {
			utils.SendJSONError(w, 500, "")
			return
		}

		value := data.Get("value").Int()
		if value <= 0 {
			utils.SendJSONError(w, 500, "")
			return
		}

		description := data.Get("description").String()

		wallet := services.Wallet{
			UserID: context.Get(r, "user_id").(string),
		}

		transfer_wallet, err := wallet.Transfer(mux.Vars(r)["wallet_id"], destination, value, description)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		utils.SendJSONResponse(w, transfer_wallet)
		return
	} else {
		utils.SendJSONError(w, 500, "The key permission is not read.")
		return
	}
}

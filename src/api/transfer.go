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

type TransferParams struct {
	Value       uint64 `json:"value"`
	Destination string `json:"destination"`
	Description string `json:"description"`
}

func Transfer(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}

		var data TransferParams
		json.Unmarshal(body, &data)

		wallet := services.Wallet{
			UserID: context.Get(r, "user_id").(string),
		}

		transfer_wallet, err := wallet.Transfer(mux.Vars(r)["wallet_id"], data.Destination, data.Value, data.Description)
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

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

func Receive(w http.ResponseWriter, r *http.Request) {
	receive := services.Receive{
		UserID:   context.Get(r, "user_id").(string),
		WalletID: mux.Vars(r)["wallet_id"],
	}

	network, _ := r.URL.Query()["network"]
	if network[0] != "lightning" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}

		var params map[string]interface{}
		json.Unmarshal(body, &params)
	}

	address, err := receive.GetAddress(network[0])
	if err != nil {
		address, err = receive.GenerateAddress(network[0])
	}

	if err != nil {
		utils.SendJSONError(w, 500, err.Error())
		return
	}

	utils.SendJSONResponse(w, map[string]string{"address": address.Address})
	return
}

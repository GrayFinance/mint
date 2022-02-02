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
	permission := context.Get(r, "permission").(string)
	if permission == "admin" || permission == "read" {
		receive := services.Receive{
			UserID:   context.Get(r, "user_id").(string),
			WalletID: mux.Vars(r)["wallet_id"],
		}

		network, _ := r.URL.Query()["network"]
		if network[0] == "lightning" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				utils.SendJSONError(w, 500, err.Error())
				return
			}

			type params struct {
				Value int    `json:"value"`
				Memo  string `json:"memo"`
			}
			var data params
			json.Unmarshal(body, &data)

			invoice, err := receive.CreateInvoice(data.Value, data.Memo)
			if err != nil {
				utils.SendJSONError(w, 500, err.Error())
				return
			}

			utils.SendJSONResponse(w, map[string]string{"payment_request": invoice.Invoice})
			return
		} else {
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
	} else {
		utils.SendJSONError(w, 500, "The key permission is not read.")
		return
	}
}

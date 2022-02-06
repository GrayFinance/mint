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

			data := gjson.ParseBytes(body)
			value := data.Get("value").Uint()
			if value == 0 {
				utils.SendJSONError(w, 500, "")
				return
			}

			description := data.Get("description").String()
			invoice, err := receive.CreateInvoice(value, description)
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

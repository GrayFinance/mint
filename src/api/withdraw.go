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

func Withdraw(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		data := gjson.ParseBytes(body)

		address := data.Get("address").String()
		if address == "" {
			utils.SendJSONError(w, 500, "Address not found.")
			return
		}

		value := data.Get("value").Uint()
		if value == 0 {
			utils.SendJSONError(w, 500, "Value not found.")
			return
		}

		if value < 10_000 {
			utils.SendJSONError(w, 500, "Value is less than 10,000 satoshi.")
		}

		feerate := data.Get("feerate").Uint()
		if feerate == 0 {
			utils.SendJSONError(w, 500, "Feerate not found")
			return
		}

		description := data.Get("description").String()
		network := mux.Vars(r)["network"]
		if network == "bitcoin" {
			wallet := services.Wallet{
				UserID: context.Get(r, "user_id").(string),
			}

			wallet_withdraw, err := wallet.BitcoinWithdraw(
				mux.Vars(r)["wallet_id"], address, value, feerate, description)
			if err != nil {
				utils.SendJSONError(w, 500, err.Error())
				return
			}
			utils.SendJSONResponse(w, wallet_withdraw)
			return
		}
	} else {
		utils.SendJSONError(w, 500, "The key permission is not admin.")
		return
	}
}

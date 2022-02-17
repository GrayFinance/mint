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

func PayInvoice(w http.ResponseWriter, r *http.Request) {
	permission := context.Get(r, "permission").(string)
	if permission == "admin" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		data := gjson.ParseBytes(body)

		invoice := data.Get("invoice").String()
		if invoice == "" {
			utils.SendJSONError(w, 500, "Not found invoice.")
			return
		}

		wallet := services.Wallet{
			UserID: context.Get(r, "user_id").(string),
		}

		payinvoice, err := wallet.PayInvoice(mux.Vars(r)["wallet_id"], invoice)
		if err != nil {
			utils.SendJSONError(w, 500, err.Error())
			return
		}
		utils.SendJSONResponse(w, payinvoice)
		return
	} else {
		utils.SendJSONError(w, 500, "The key permission is not admin.")
		return
	}
}

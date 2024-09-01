package productcontroller

import (
	"go-jwt-mux/helper"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	data := []map[string]interface{}{
		{
			"id":           1,
			"name_product": "Mie Ayam",
			"price":        20000,
		},
		{
			"id":           2,
			"name_product": "Mie Ayam Bakso",
			"price":        25000,
		},
		{
			"id":           3,
			"name_product": "Mie Ayam Bakso",
			"price":        25000,
		},
	}
	helper.ResponseJSON(w, http.StatusOK, data)
}

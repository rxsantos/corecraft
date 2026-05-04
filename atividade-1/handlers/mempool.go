package handlers


import (
	// "fmt"
	"encoding/json"
	"net/http"

	"atividade-1/rpc"
	"atividade-1/services"
)

func MempoolSummary(btc rpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Chama getmempoolinfo
		info, err := btc.Call("getmempoolinfo", []interface{}{})
		if err != nil {
			http.Error(w, "Erro ao chamar getmempoolinfo: "+err.Error(), 500)
			return
		}

		//Chama getrawmempool com verbose true
		raw, err := btc.Call("getrawmempool", []interface{}{true})
		if err != nil {
			http.Error(w, "Erro ao chamar getrawmempool: "+err.Error(), 500)
			return
		}

		// fmt.Printf("RAW: %#v\n", raw)

		//Calcula estatísticas
		result := services.CalculateMempoolStats(info, raw)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
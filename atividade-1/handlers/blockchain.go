package handlers

import (
	"encoding/json"
	"net/http"

	"atividade-1/rpc"
)

func BlockchainLag(btc rpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Chama getblockchaininfo
		resp, err := btc.Call("getblockchaininfo", []interface{}{})
		if err != nil {
			http.Error(w, "Erro ao chamar getblockchaininfo: "+err.Error(), 500)
			return
		}		

		//Resultado da Chamada
		blocks, ok1 := resp["blocks"].(float64)
		headers, ok2 := resp["headers"].(float64)

		if !ok1 || !ok2 {
			http.Error(w, "Resposta inválida do Bitcoin Core", 500)
			return
		}
		data := map[string]interface{}{
			"blocks":  blocks,
			"headers": headers,
			"lag":     headers - blocks,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
package main

import (
	"fmt"
	"net/http"

	"atividade-1/handlers"
	"atividade-1/rpc"
)

func main() {
	btc := rpc.Client{
		URL:  "http://127.0.0.1:38443/",
		User: "teste",
		Pass: "teste",
	}

	//Exposição de Rodas (APIs)
	http.HandleFunc("/api/mempool/summary", handlers.MempoolSummary(btc))
	http.HandleFunc("/api/blockchain/lag", handlers.BlockchainLag(btc))
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	fmt.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar servidor: %v\n", err)
	}
}
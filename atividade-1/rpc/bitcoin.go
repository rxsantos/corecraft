package rpc

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	URL string
	User string
	Pass string
}

type RPCRequest struct {
	Jsonrpc	string			`json:"jsonrpc"`
	ID 		string			`json:"id"`
	Method	string			`json:"method"`
	Params	[]interface{}	`json:"params"`
}

func(c *Client) Call(method string, params []interface{}) (map[string]interface{}, error) {
	reqBody := RPCRequest{
		Jsonrpc: 	"1.0",
		ID: 		"corecraft",
		Method:		method,
		Params:		params,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}


	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.User, c.Pass)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	//tratar erro RPC corretamente
	if response["error"] != nil {
		return nil, fmt.Errorf("rpc error: %v", response["error"])
	}

	//retorna só o result
	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("resultado inválido")
	}

	return result, nil
}
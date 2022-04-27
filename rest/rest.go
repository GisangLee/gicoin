package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gisanglee/gicoin/blockchain"
	"github.com/gisanglee/gicoin/utils"
	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errResponse struct {
	ErrorMsg string `json:"errorMsg"`
}

type addTxPayload struct {
	To     string
	Amount int
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the blockchain status",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
		},
	}

	b, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Fprintf(rw, "%s", b)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())

	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)

	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	block_map := mux.Vars(r)
	block_hash := block_map["hash"]

	block, err := blockchain.FindBlock(block_hash)

	encoder := json.NewEncoder(rw)

	if err == blockchain.ErrNotFound {
		encoder.Encode(errResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	total := r.URL.Query().Get("total")

	switch total {
	case "true":
		amount := blockchain.Blockchain().BalanceByAddress(address)
		json.NewEncoder(rw).Encode(balanceResponse{Address: address, Balance: amount})
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blockchain().UTXOByAddress(address)))
	}

}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload

	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)

	if err != nil {
		json.NewEncoder(rw).Encode(errResponse{"not enough funds"})
	}

	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	//handler := http.NewServeMux()

	handler := mux.NewRouter()

	handler.Use(jsonContentTypeMiddleware)

	port = fmt.Sprintf(":%d", aPort)
	//explorer.Start()
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	handler.HandleFunc("/status", status).Methods("GET")
	handler.HandleFunc("/balance/{address}", balance).Methods("GET")
	handler.HandleFunc("/mempool", mempool).Methods("GET")
	handler.HandleFunc("/transactions", transactions).Methods("POST")

	fmt.Printf("REST API > Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

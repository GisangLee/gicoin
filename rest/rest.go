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

type addBlockBody struct {
	Message string
}

type errResponse struct {
	ErrorMsg string `json:"errorMsg"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See All Blocks",
			Payload:     "data:string",
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
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
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

func Start(aPort int) {
	//handler := http.NewServeMux()

	handler := mux.NewRouter()

	handler.Use(jsonContentTypeMiddleware)

	port = fmt.Sprintf(":%d", aPort)
	//explorer.Start()
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")

	fmt.Printf("REST API > Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

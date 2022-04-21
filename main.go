package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gisanglee/gicoin/utils"
)

//import "github.com/gisanglee/gicoin/explorer"

const PORT string = ":4000"

type URLDescription struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         "/block",
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Fprintf(rw, "%s", b)
}

func main() {
	//explorer.Start()
	http.HandleFunc("/", documentation)
	fmt.Printf("LIstening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

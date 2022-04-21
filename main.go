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
	URL         string
	Method      string
	Description string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
	}
	b, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Printf("%s", b)
}

func main() {
	//explorer.Start()
	http.HandleFunc("/", documentation)
	fmt.Printf("LIstening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

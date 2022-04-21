package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gisanglee/gicoin/blockchain"
)

const PORT string = ":4000"
const TEMPLATE_DIR string = "templates/"

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {

	data := homeData{PageTitle: "Home", Blocks: blockchain.AllBlocks()}

	templates.ExecuteTemplate(rw, "home", data)

}

func add(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)

	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func main() {
	templates = template.Must(template.ParseGlob(TEMPLATE_DIR + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(TEMPLATE_DIR + "partials/*.html"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

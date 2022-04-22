package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gisanglee/gicoin/blockchain"
)

const TEMPLATE_DIR string = "explorer/templates/"

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {

	data := homeData{PageTitle: "Home", Blocks: nil}

	templates.ExecuteTemplate(rw, "home", data)

}

func add(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)

	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.Blockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	handler := http.NewServeMux()

	templates = template.Must(template.ParseGlob(TEMPLATE_DIR + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(TEMPLATE_DIR + "partials/*.html"))

	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)

	fmt.Printf("HTML > Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

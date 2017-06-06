package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	serveWeb()
}

// struct to pass into the template
type defaultContext struct {
	Title      string
	ErrorMsg   string
	SuccessMsg string
	Usuario    string
	Senha      string
}

var themeName = getThemeName()
var staticPages = populateStaticPages()

func serveWeb() {
	gorillaRoute := mux.NewRouter()

	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{pageAlias}", serveContent).Methods("POST")
	gorillaRoute.HandleFunc("/{pageAlias}", serveContent).Methods("GET") //URL com parametros dinamicos

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	/*
		gorillaRoute.HandleFunc("/api/{uri}", api).Methods("GET")
		gorillaRoute.HandleFunc("/api/{uri}", api).Methods("POST")
	*/
	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":80", nil)
}

/*
func api(w http.ResponseWriter, r *http.Request) {
	log.Println("CHAMOU a API!!!!")
	login := r.FormValue("login")
	senha := r.FormValue("senha")

	pageAlias := "erro_autenticacao"
	if (senha=="y") {
		pageAlias = "autenticado"
	}
	staticPage := staticPages.Lookup(pageAlias + ".html")

	context := defaultContext{}
	context.Usuario = login
	context.Title = "ttt"

	staticPage.Execute(w, context)
}
*/

func serveContent(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	pageAlias := urlParams["pageAlias"]
	if pageAlias == "" {
		pageAlias = "sign"
	}

	staticPage := staticPages.Lookup(pageAlias + ".html")
	if staticPage == nil {
		log.Println("NAO ACHOU!!!!")
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}

	login := r.FormValue("login")
	senha := r.FormValue("senha")

	//Values to pass into the template
	context := defaultContext{}
	context.Title = pageAlias
	context.ErrorMsg = ""
	context.SuccessMsg = ""
	context.Usuario = login
	context.Senha = senha

	staticPage.Execute(w, context)
}

func getThemeName() string {
	return "bs4"
}

func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "pages"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	for _, pathinfo := range templatePathsRaw {
		log.Println(pathinfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathinfo.Name())
	}

	basePath = "themes/" + themeName
	templateFolder, _ = os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ = templateFolder.Readdir(-1)
	for _, pathinfo := range templatePathsRaw {
		log.Println(pathinfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathinfo.Name())
	}

	result.ParseFiles(*templatePaths...)
	return result
}

func serveResource(w http.ResponseWriter, req *http.Request) {

	path := "public/" + themeName + req.URL.Path
	var contentType string

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".	jpg") {
		contentType = "image/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	log.Println(path)
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

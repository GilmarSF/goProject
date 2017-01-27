package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"text/template"
	"os"
	"log"
	"bufio"
	"strings"
)

func main() {
	serveWeb()
}

// struct to pass into the template
type defaultContext struct{
	Title string
	ErrorMsg string
	SuccessMsg string
}

var themeName  = getThemeName()
var staticPages = populateStaticPages()

func serveWeb () {
	gorillaRoute := mux.NewRouter()

	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{page_alias}", serveContent)  //URL com parametros dinamicos

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":80", nil)
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]
	if page_alias == "" {
		page_alias = "home"
	}

	staticPage := staticPages.Lookup(page_alias + ".html")
	log.Println(staticPage)
	if staticPage == nil {
		log.Println("NAO ACHOU!!!!")
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}

	//Values to pass into the template
	context := defaultContext{}
	context.Title = page_alias
	context.ErrorMsg = ""
	context.SuccessMsg = ""

	staticPage.Execute(w, context)
}

func getThemeName() string {
	return "bs4"
}

func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "pages/"
	templateFolder, _:= os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	for _, pathinfo := range templatePathsRaw {
		log.Println(pathinfo.Name())
		*templatePaths = append(*templatePaths, basePath + "/" + pathinfo.Name())
	}

	basePath = "themes/" + themeName
	templateFolder, _= os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ = templateFolder.Readdir(-1)
	for _, pathinfo := range templatePathsRaw {
		log.Println(pathinfo.Name())
		*templatePaths = append(*templatePaths, basePath + "/" + pathinfo.Name())
	}

	result.ParseFiles (*templatePaths...)
	//log.Printf(result)
	return result
}

func serveResource(w http.ResponseWriter, req *http.Request){

	path := "public/" + themeName +req.URL.Path
	var contentType string

	if strings.HasSuffix(path, ".css"){
		contentType = "text/css; charset=utf-8"
	}else if strings.HasSuffix(path, ".png"){
		contentType = "image/png; charset=utf-8"
	}else if strings.HasSuffix(path, ".jpg"){
		contentType = "image/jpg; charset=utf-8"
	}else if strings.HasSuffix(path, ".js"){
		contentType = "application/javascript; charset=utf-8"
	}else {
		contentType = "text/plain; charset=utf-8"
	}

	log.Println(path)
	f, err := os.Open(path)
	if err == nil{
		defer f.Close()
		w.Header().Add("Content-type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w) 
	}else {
		w.WriteHeader(404)
	}
}







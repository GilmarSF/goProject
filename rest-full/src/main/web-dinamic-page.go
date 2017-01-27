package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"text/template"
	"os"
	"log"
	"bufio"
	"strings"
	"io/ioutil"
    "bytes"
    "encoding/json"
)

func main() {
	serveWeb()
}

// struct to pass into the template
type defaultContext struct{
	Title 	string
}

type CategoriaFull struct{
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Total   string  `json:"total,omitempty"`
}

type CategoriaEach struct{
    ID      string `json:"id,omitempty"`
    Nome    string `json:"nome,omitempty"`
    Regiao  string `json:"regiao,omitempty"`
    Total   string `json:"total,omitempty"`
}

var themeName  = getThemeName()
var staticPages = populateStaticPages()
var page_alias string

func attHTML (){
	themeName  = getThemeName()
	staticPages = populateStaticPages()
}

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

func atualizaJSON(){
	log.Printf("atualiza arquivo JSON") 
	respostaFull, err := http.Get("http://localhost:8080/categorias/")
	if err != nil {
		log.Println(err)
	}

	dadosRespostaFull, err := ioutil.ReadAll(respostaFull.Body)
	if err != nil {
		log.Println(err)
	}

	var categFull []CategoriaFull
	json.Unmarshal(dadosRespostaFull, &categFull) // coverte de json para struct

	for _, item := range categFull{
		log.Println(item.Nome, item.Total)	
	}	

	///////////////////////////////////// por categoria

	respostaEach, err := http.Get("http://localhost:8080/categorias/0")
	if err != nil {
		log.Println(err)
	}

	dadosRespostaEach, err := ioutil.ReadAll(respostaEach.Body)
	if err != nil {
		log.Println(err)
	}

	var categEach []CategoriaEach
	json.Unmarshal(dadosRespostaEach, &categEach)

	//for _, item := range categEach{
	//	log.Println(item.ID, item.Nome, item.Regiao, item.Total)	
	//}

	////////////////////////////////////////

	path := "B:/go/Github-faculdade/goProject/rest-full/bin/pages/"
	jsonOut, err := ioutil.ReadFile(path+"default.json")
    if err != nil {
        log.Println(err)
    }
   
    for _, item := range categFull { 

    	attCategoria := bytes.Replace(jsonOut, []byte("Categoria"), []byte(item.Nome), 1)
		if err = ioutil.WriteFile(path+"geral.json.html", attCategoria, 0666); err != nil {
            log.Println(err)
        }

        jsonOut, err = ioutil.ReadFile(path+"geral.json.html")
        if err != nil {
            log.Println(err)
        }  

        attTotal := bytes.Replace(jsonOut, []byte("00"), []byte(item.Total), 1)
        if err = ioutil.WriteFile(path+"geral.json.html", attTotal, 0666); err != nil {
            log.Println(err)
        }

        jsonOut, err = ioutil.ReadFile(path+"geral.json.html")
        if err != nil {
            log.Println(err)
        }   
    }

    jsonOut, err = ioutil.ReadFile(path+"default.json")
    if err != nil {
        log.Println(err)
    }

    for _, item := range categEach { 
        if strings.ToUpper(item.Nome) == strings.ToUpper(page_alias) { 
            //categoriaFound = append(categoriaFound,item)
            //log.Println(item.ID, item.Nome, item.Regiao, item.Total)
    
			attCategoria := bytes.Replace(jsonOut, []byte("Categoria"), []byte(item.Regiao), 1)
			if err = ioutil.WriteFile(path+"categoria.json.html", attCategoria, 0666); err != nil {
	            log.Println(err)
	        }

	        jsonOut, err = ioutil.ReadFile(path+"categoria.json.html")
	        if err != nil {
	            log.Println(err)
	        }  

	        attTotal := bytes.Replace(jsonOut, []byte("00"), []byte(item.Total), 1)
	        if err = ioutil.WriteFile(path+"categoria.json.html", attTotal, 0666); err != nil {
	            log.Println(err)
	        }

	        jsonOut, err = ioutil.ReadFile(path+"categoria.json.html")
	        if err != nil {
	            log.Println(err)
	        }   
	    }
    }
    attHTML()
}

func serveContent(w http.ResponseWriter, r *http.Request) {

	atualizaJSON()
	urlParams := mux.Vars(r)
	page_alias = urlParams["page_alias"]
	if page_alias == "" {
		page_alias = "geral"
	}

	staticPage := staticPages.Lookup(page_alias+".html")
	if staticPage == nil {
		log.Println("NAO ACHOU!!")
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}

	//Values to pass into the template   
	context := defaultContext{}
	context.Title = page_alias

	staticPage.Execute(w, context)
}

func getThemeName() string {
	return "bs4"
}

func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "pages"
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



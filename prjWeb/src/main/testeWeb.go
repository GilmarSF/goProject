package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"text/template"
	"os"
	"log"
	"strings"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//connectionDB()   // inicia conexão com o azure -Teste Git
	serveWeb()
}

type defaultContext struct {
	Title string
	ErrorMsgs string
	SuccessMsgs string
	Codigo    string
	Nome      string
	Email     string
}
var themeName  = getThemeName()
var staticPages = populateStaticPages()
func serveWeb () {
	log.Println(">>>>>>>serveWeb<<<<<<<")
	gorillaRoute := mux.NewRouter()

	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{page_alias}", serveContent)  //URL com parametros dinamicos

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}

func serveContent(w http.ResponseWriter, r *http.Request) {

	var (
		id_pessoa string
		nome     string
		email    string
	)
	log.Println(">>>>>>>      serveContent     <<<<<<<")
	urlParams   := mux.Vars(r)
	log.Println("Request: " + r.Method + "-" + r.RequestURI);
	codigo  := r.URL.Query()["Codigo"]

	log.Println("VER TAMANHO DO CODIGO")
	if len(codigo) > 0 {
		id_pessoa  = codigo[0]
	} else {
		id_pessoa  = "0"
	}

	log.Println("VER TAMANHO DO CODIGO OK: " + id_pessoa )
	page_alias := urlParams["page_alias"]

	log.Println("page_alias: " + page_alias + ".html")
	if page_alias == "" {
		page_alias = "home"
	}
	staticPage := staticPages.Lookup(page_alias + ".html")

	if staticPage == nil {
		log.Println("NAO ACHOU " + page_alias)
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	} else {
		log.Println("ACHOU " + page_alias)
	}

	context := defaultContext{}
	context.Title = "GilGil"//page_alias
	context.ErrorMsgs = ""
	context.SuccessMsgs = "Parse System..."


	if id_pessoa != "0" {
		// Banco no azure do Zé
		db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
		if err != nil { // se o retorno for 'nil' quer dizer que deu tudo certo
			log.Println("Open Failed: ", err.Error())
		}

		// db.Query usado para comandos no Banco
		//traz todas as linha do campo 'regiao' na tab_localidade
		rows, err := db.Query("select id, nome, email from [dbo].[tab_usuario] where id = ?", id_pessoa)
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close() // fecha o comando Query
		defer db.Close()   // fecha conexão com o Banco

		for rows.Next() {
			// percorre o objeto 'rows' salvando o que encontra na string 'dbDados'
			if err := rows.Scan(&id_pessoa, &nome, &email); err != nil {
				log.Fatal(err)
			}
			context.Nome = nome
			context.Email = email
			fmt.Printf("----------------------------------------------------------")
			fmt.Printf("id: %s \n", id_pessoa)
			fmt.Printf("Nome: %s \n", nome)
			fmt.Printf("email: %s \n", email)
		}
		context.Codigo = id_pessoa
	} else {

		context.Nome = ""
		context.Email = ""
		context.Codigo = "0"
	}

	staticPage.Execute(w,context)
}

func getThemeName() string {
	return "bs4"
}

func populateStaticPages() *template.Template {
	log.Println(">>>>>>>populateStaticPages<<<<<<<")
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
		log.Println(basePath + "/" + pathinfo.Name())
		*templatePaths = append(*templatePaths, basePath + "/" + pathinfo.Name())
	}

	result.ParseFiles (*templatePaths...)

	return result
}

func serveResource( w http.ResponseWriter, req *http.Request) {
	log.Println(">>>>>>>serveResource<<<<<<<")
	path := "public/" + themeName + req.URL.Path
	var contentType string

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg") {
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
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}


func connectionDB() {
	var (
		id_pessoa int
		nome     string
		email    string
	)
	// Banco no azure do Zé
	//db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	db, err := sql.Open("mssql", "server=elevainf1.database.windows.net;user id=gilmarsf;password=Canada01;database=DB_SQL;port=1433")
	if err != nil { // se o retorno for 'nil' quer dizer que deu tudo certo
		log.Println("Open Failed: ", err.Error())
	}

	// db.Query usado para comandos no Banco
	//traz todas as linha do campo 'regiao' na tab_localidade
	rows, err := db.Query("select id, nome, email from [dbo].[tbPessoa] where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close() // fecha o comando Query
	defer db.Close()   // fecha conexão com o Banco

	for rows.Next() {
		// percorre o objeto 'rows' salvando o que encontra na string 'dbDados'
		if err := rows.Scan(&id_pessoa, &nome, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("----------------------------------------------------------")
		fmt.Printf("id: %s \n", id_pessoa)
		fmt.Printf("Nome: %s \n", nome)
		fmt.Printf("email: %s \n", email)
	}

	fmt.Printf("\nbye")

}
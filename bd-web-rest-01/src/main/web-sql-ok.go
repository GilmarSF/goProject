package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	//"flag"
	"net/http"
)

var name string
var peso = 68

func main() {

	connectionHTTP() // chama função servidor http
	connectionDB()   // inicia conexão com o azure
}

func connectionDB() {

	// ABRINDO CONEXÃO E TESTANDO CONEXÃO COM PING
	log.Println("Main:")

	log.Println("Opening")

	// Banco Local
	//db, err := sql.Open("mssql", "server=WARRIOR\\SQLEXPRESS;Initial Catalog=dbo;user id=sa;password=123456;port=1433")

	// Banco no azure do Zé
	db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")

	if err != nil {
		log.Println("Open Failed: ", err.Error())
	}

	log.Println("Opened")

	log.Println("Pinging")

	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping: ", err.Error())
	}

	log.Println("Pinged")

	// db.Query usado para comandos no Banco
	rows, err := db.Query("select Nome from [dbo].[tbPessoa] where Peso=?", peso)
	//rows, err := db.Query("CREATE TABLE JoseDB.dbo.TestTable (ColA INT PRIMARY KEY, ColB INT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close() // fecha o comando Query
	defer db.Close()   // fecha conexão com o Banco
	//fmt.Println(rows)

	// rows.Next usado para varrer o objeto 'rows' e pegar os valores retornados da Query
	for rows.Next() {

		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s pesa %d\n", name, peso)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("bye\n")
}

func connectionHTTP() {
	http.HandleFunc("/", serveHome)           // Pagina home (localhost:80/)
	http.HandleFunc("/contato", serveContato) // (localhost:80/contato)

	http.ListenAndServe(":80", nil) // usando porta 80, pode ser 8080, apenas alterar
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(name)) //Exibe o nome buscado no select
}

func serveContato(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pagina de contato"))
}

// Sujeira =)

/*
	stmt, err := db.Prepare("select 1, 'abc'")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var somenumber int64
	var somechars string
	err = row.Scan(&somenumber, &somechars)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("somenumber:%d\n", somenumber)
	fmt.Printf("somechars:%s\n", somechars)
*/

/*connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
if *debug {
	fmt.Printf(" connString:%s\n", connString)
}*/

/*
	flag.Parse() // parse the command line args

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}
*/

/*

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "123abc!@#", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "pwbt.database.windows.net", "the database server")
var user = flag.String("user", "admin-jose", "the database user")
*/

/*
	// ABRE CONEXÃO COM BANCO DE DADOS SQL SERVER NO AZURE (BANCO DO JOSÉ)
	db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;port=1433")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}

	defer db.Close()

	err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
*/

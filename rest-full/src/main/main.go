package main
 
import (
    "encoding/json"
    "strconv" // usado para converter dados
    "log"
    "net/http"
	"fmt"
    "github.com/gorilla/mux"
    "database/sql"
    _ "github.com/denisenkom/go-mssqldb"
)
 
type Usuario struct {
    //"omitempty", quando usado omite este atributo caso seja 'null'
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Email   string  `json:"email,omitempty"`
    Senha   string  `json:"senha,omitempty"`
}

var cadastros []Usuario // array de Usuarios
var count int // usado para armazenar o ultimo 'id'
 
// Traz apenas um Usuario
func GetUsuario(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req) // salva o 'id' passado na requisição
    for _, item := range cadastros { // For sera executado de acordo com o tamanho do array 'cadastros'
        if item.ID == params["id"] { // Se encontrar uma Usuario com o ID enviado
			fmt.Fprintln(w, "Usuario: "+item.Nome+"\n")
            json.NewEncoder(w).Encode(item) // imprime os dados da Usuario
            return // sai da func
        }
    }
	fmt.Fprintln(w, "Usuario não encontrada!")
}
 
// traz a cadastros toda. Todas todas Usuarios do array 'cadastros'
func GetCadastros(w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(cadastros)
}
 
func PostUsuario(w http.ResponseWriter, req *http.Request) {
    //params := mux.Vars(req)
    var usuario Usuario
    
    _ = json.NewDecoder(req.Body).Decode(&usuario)
    fmt.Println(usuario)
    count++ // usado para sempre o numero do 'id' ser id+1
    //usuario.ID = strconv.Itoa(count)

    db, err := sql.Open("mssql", "server=WARRIOR\\SQLEXPRESS;user id=sa;password=123456;database=joseDB;port=1433")
    if err != nil {
        log.Println("Open Failed: ", err.Error())
    }

    rows, err := db.Query("INSERT INTO tab_usuario (id, nome, email, senha) VALUES (?1, ?2, ?3, ?4)", count, usuario.Nome, usuario.Email, usuario.Senha) 
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close() // fecha o comando Query
    defer db.Close()   // fecha conexão com o Banco
    
    connectionDB() // atualiza struct
    //cadastros = append(cadastros, usuario) // adiciona ao final de 'cadastros'
    json.NewEncoder(w).Encode(cadastros)
}
 
/* Quando o 'id' for deletado, sera preciso recriar o slice com o dados restantes, por isso o uso do 'indiceArray'
o for a cada rodada grava em 'indiceArray' o indice do array que ele esta.*/
func DeleteUsuario(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)

    db, err := sql.Open("mssql", "server=WARRIOR\\SQLEXPRESS;user id=sa;password=123456;database=joseDB;port=1433")
    if err != nil {
        log.Println("Open Failed: ", err.Error())
    }

    del, _ := strconv.Atoi(params["id"])
    rows, err := db.Query("DELETE FROM tab_usuario WHERE id = ?1", del) 
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close() // fecha o comando Query
    defer db.Close()   // fecha conexão com o Banco
    
    connectionDB() // atualiza struct

    /*for indiceArray, item := range cadastros {
        if item.ID == params["id"] {
            cadastros = append(cadastros[:indiceArray], cadastros[indiceArray+1:]...)
            break
        }
    }*/
    json.NewEncoder(w).Encode(cadastros)
}

func connectionDB() {

    // Banco no azure do Zé
    //db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
    db, err := sql.Open("mssql", "server=WARRIOR\\SQLEXPRESS;user id=sa;password=123456;database=joseDB;port=1433")
    if err != nil {
        log.Println("Open Failed: ", err.Error())
    }

    // db.Query usado para comandos no Banco
    //traz todas as linha do campo 'regiao' na tab_localidade
    rows, err := db.Query("select id, nome, email, senha from tab_usuario") 
    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close() // fecha o comando Query
    defer db.Close()   // fecha conexão com o Banco
    
    cadastros = cadastros[:0] // zera array antes de buscar novos dados no banco
    // rows.Next usado para varrer o objeto 'rows' e pegar os valores retornados da Query
    for rows.Next() {

        addUsuario := Usuario{}
        if err := rows.Scan(&addUsuario.ID, &addUsuario.Nome, &addUsuario.Email, &addUsuario.Senha); err != nil { 
            log.Fatal(err)
        }
        /*fmt.Printf(addUsuario.ID+"\n"); fmt.Printf(addUsuario.Nome+"\n"); fmt.Printf(addUsuario.Email+"\n")*/
        cadastros = append(cadastros, addUsuario) //salva o retorno neste slice "slice=array flexivel"
        count, _ = strconv.Atoi(addUsuario.ID) // pega o valor do ID do ultimo dado buscado no banco
    }

    fmt.Printf("\nbye\n")
}

func main() {

    connectionDB()   // inicia conexão com o azure    
    router := mux.NewRouter()
    router.HandleFunc("/cadastros", GetCadastros).Methods("GET")
    router.HandleFunc("/cadastros/{id}", GetUsuario).Methods("GET")
    router.HandleFunc("/cadastros/", PostUsuario).Methods("POST")
    router.HandleFunc("/cadastros/{id}", DeleteUsuario).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8080", router))
}
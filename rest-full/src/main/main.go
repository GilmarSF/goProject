// Todo codigo em Go precisa iniciar com este package
package main
 
 // Importa todas as bibliotecas que serão utilizadas no código
import (

    // bibliotecas Nativas do Golang
    "encoding/json"
    "strconv" 
    "log"
    "net/http"
	"fmt"
    "database/sql"

    // Pode ser usado bibliotecas não nativas 
    // /gorilla/mux para o servidor HTTP
    // /denisenkom/go-mssqldb para a comunicação com Banco SQLServer
    "github.com/gorilla/mux"
    _ "github.com/denisenkom/go-mssqldb"
)
 
// Struct Principal criada para armazenar os dados retornados do Banco de dados
// Structs são mais comuns para converter para JSON em Golang
type Usuario struct {

    // "omitempty", quando usado, omite este atributo caso seja 'null'
    ID      string  `json:"id,omitempty"`
    Nome    string  `json:"nome,omitempty"`
    Email   string  `json:"email,omitempty"`
    Senha   string  `json:"senha,omitempty"`
}

/* VARIAVEIS GLOBAIS PODEM SER USADAS EM QUALQUER PARTE DO CÓDiGO*/
// array para salvar os Usuarios
var cadastros []Usuario 
// Usado para armazenar o ultimo 'id' do banco de dados
var count int 
 
// Função utilizada enviar apenas um Usuario via Rest
func GetUsuario(w http.ResponseWriter, req *http.Request) {

    // Um aviso qualquer que será impresso no "prompt"
    fmt.Printf("\nGet Usuario\n") 

    // Salva o 'id' (valor) passado na requisição
    params := mux.Vars(req) 

    // For usado para varrer todos os Usuarios que estão salvos no Array 'cadastros'
    // 'range' significo por todo o array
    // Neste for esta sendo criado uma Struct 'item' com os dados do Usuario que possuir o mesmo ID da requisição
    // o sinal [ := ] é usado para criar uma nova variavel do tipo do valor que será adicionado a ela
    for _, item := range cadastros { 

        // Verifica em cada Usuario do array 'cadastros' se o ID é igual ao id passado na requisição
        if item.ID == params["id"] { 

            // Quando encontrado, converte a Struct 'item' para o 
            // formato 'JSON' e envia via HTTP na porta :8080
            json.NewEncoder(w).Encode(item) 
            return 
        }
    }

    // Caso não encontrar, envia um JSON vazio
    json.NewEncoder(w).Encode(&Usuario{}) 
}
 
// traz a cadastros toda. Todas todas Usuarios do array 'cadastros'
func GetCadastros(w http.ResponseWriter, req *http.Request) {
    fmt.Printf("\nGet cadastros\n") 
    json.NewEncoder(w).Encode(cadastros)
}

// Adiciona mais um Usuario na tab_usuario
func PostUsuario(w http.ResponseWriter, req *http.Request) {

    fmt.Printf("\nPost Usuario\n") 
    var usuario Usuario
    
    _ = json.NewDecoder(req.Body).Decode(&usuario)
    fmt.Println(usuario)
    count++ // usado para sempre o numero do 'id' ser id+1

    db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
    if err != nil {
        log.Println("Erro ao conectar com o Banco de dados:", err.Error())
    }

    rows, err := db.Query("INSERT INTO tab_usuario (id, nome, email, senha) VALUES (?1, ?2, ?3, ?4)", count, usuario.Nome, usuario.Email, usuario.Senha) 
    if err != nil {
        log.Println("Erro no INSERT:", err.Error())
    }
    defer rows.Close() // fecha o comando Query
    defer db.Close()   // fecha conexão com o Banco
    
    connectionDB() // atualiza struct
    json.NewEncoder(w).Encode(cadastros)
}

// Deleta o usuario da tab_usuario de acordo com o ID passado
func DeleteUsuario(w http.ResponseWriter, req *http.Request) {

    fmt.Printf("\nDelete Usuario\n") 
    params := mux.Vars(req) // salva em 'params' os dados passados na url

    db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
    if err != nil {
        log.Println("Erro ao conectar com o Banco de dados:", err.Error())
    }

    del, _ := strconv.Atoi(params["id"]) // converte de String para Int
    rows, err := db.Query("DELETE FROM tab_usuario WHERE id = ?1", del) 
    if err != nil {
        log.Println("Erro no DELETE:", err.Error())
    }
    defer rows.Close() // fecha o comando Query
    defer db.Close()   // fecha conexão com o Banco
    
    connectionDB() // atualiza struct

    json.NewEncoder(w).Encode(cadastros)
}

func connectionDB() {

    // Inicia conexão com Banco de dados no Azure 
    db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
    // Caso ocorrer algum erro na conexão, o ERRO será salvo em 'err'
    // Caso não houver ERRO o retorno será vazio 'nil'
    if err != nil {
        log.Println("Erro ao conectar com o Banco de dados:", err.Error())
    }

    // db.Query usado para executar comandos no Banco
    rows, err := db.Query("select id, nome, email, senha from tab_usuario") 
    if err != nil {
        log.Println("Erro no SELECT Principal:", err.Error())
    }

    // finaliza o comando para Query
    defer rows.Close() 
    // fecha conexão com o Banco
    defer db.Close()   
    
    // Zera o Array antes de buscar novos dados no banco.
    // Dessa forma evita dados repetidos
    cadastros = cadastros[:0]

    // rows.Next usado para varrer o objeto 'rows' e pegar os valores retornados da Query
    for rows.Next() {
        // Struct criada para receber os dados do banco de dados
        // Uma Struct permite receber dados de diferentes tipos: String, Int ...
        addUsuario := Usuario{} 

        // rows.Scan varre o objeto rows e salva os valores na variaveis citadas abaixo.
        // As variaveis são salvas na mesma ordem que são coletados do Banco de dados
        if err := rows.Scan(&addUsuario.ID, &addUsuario.Nome, &addUsuario.Email, &addUsuario.Senha); err != nil { 
            log.Println("Erro ao salvar os dados retornados do Banco:", err.Error())
        }
        
        // adiciona na struct principal os dados do banco
        // Sera está struct 'cadastros' que será convertida para o formato JSON
        cadastros = append(cadastros, addUsuario)

        // Pega o valor do ID do ultimo dado buscado no banco e converte de String para Int
        count, _ = strconv.Atoi(addUsuario.ID) 
    }

    fmt.Printf("\nAPI atualizada!\n")
}

func main() {

    connectionDB()   // inicia conexão com o azure    
    router := mux.NewRouter()
    router.HandleFunc("/cadastros", GetCadastros).Methods("GET")
    router.HandleFunc("/cadastros/{id}", GetUsuario).Methods("GET")
    router.HandleFunc("/cadastros/", PostUsuario).Methods("POST")
    router.HandleFunc("/cadastros/{id}", DeleteUsuario).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8080", router)) // Server na porta 8080 [ localhost:8080 ]
}

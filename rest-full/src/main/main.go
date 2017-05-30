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
    "io/ioutil"
    "bytes"

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

type Categoria struct{
    ID      string
    Nome    string
    Total   string

}
/* VARIAVEIS GLOBAIS PODEM SER USADAS EM QUALQUER PARTE DO CÓDiGO*/
// array para salvar os Usuarios
var cadastros []Usuario 

var categorias []Categoria
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
    
    AtualizaUsuarios() // atualiza struct no banco
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
    
    AtualizaUsuarios() // atualiza struct no banco

    json.NewEncoder(w).Encode(cadastros)
}

func CriaJSON(w http.ResponseWriter, req *http.Request){

    // Busca no banco de dados todas as categorias e os totais de reclamações de novo
    // Assim toda vez que chamar essa func os dados estarão sempre att 
    AtualizaCategorias()
    fmt.Printf("\nTabela.json atualizada\n") 
    
    // Lê todo conteudo do arquivo default.json e salva em 'jsonOut'
    jsonOut, err := ioutil.ReadFile("default.json")
    if err != nil {
            fmt.Println(err)
    }

    // for usado para varrer todos os dados Salvos no array d structs 'categorias'
    // a cada dado encontrado, é sobrescrito no arquivo default.json
    for _, item := range categorias { 

        // Salva em 'attCategoria' o conteudo copiado de 'default.json' 
        // já com a palavra "Categoria" substituida por 'item.Nome'
        // 'item' recebe do 'for' citado acima os dados da posição atual do array de struct
        // Obs.: O numero 1 no final significa que é para alterar apenas a primeira
        // palavra "Categoria encontrada", se estivesse 2 substituiria as 2 primeiras
        // palavras "Categoria" encontrada, se estivesse -1 substituiria todas as palavras
        // "Categoria" para a palavra salva no atual item.Nome
        attCategoria := bytes.Replace(jsonOut, []byte("Categoria"), []byte(item.Nome), 1)
        
        // Nesta parte sera gravado todo o conteudo de attCategoria com a alteração acima
        // no arquivo tabela.json
        if err = ioutil.WriteFile("tabela.json", attCategoria, 0666); err != nil {
            fmt.Println(err)
        }

        // Lê todo conteudo do arquivo tabela.json e salva em 'jsonOut'
        // a partir daqui cada alteração no arquivo 'tabela.json' o mesmo deve
        // ser relido e salvo em 'jsonOut' para as proximas alterações não serem
        // baseadas no arquivo 'default.json' 
        jsonOut, err = ioutil.ReadFile("tabela.json")
        if err != nil {
            fmt.Println(err)
        }  

        // Faz o mesmo descrito acima para o numero total de incidentes
        // Obs.: apesar de ser escrito no formato string
        // o Javascript na hora de ler entendera que é um numero INT pois não tem " "
        attTotal := bytes.Replace(jsonOut, []byte("0"), []byte(item.Total), 1)
        if err = ioutil.WriteFile("tabela.json", attTotal, 0666); err != nil {
            fmt.Println(err)
        }

        jsonOut, err = ioutil.ReadFile("tabela.json")
        if err != nil {
                fmt.Println(err)
        }   
    }
}

// Atualiza a struct Usuario com os dados do Banco
func AtualizaUsuarios() {

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
        log.Println("Erro no SELECT dos Usuarios:", err.Error())
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
            log.Println("Erro ao salvar os cadastros retornados do Banco:", err.Error())
        }
        
        // adiciona na struct principal os dados do banco
        // Sera está struct 'cadastros' que será convertida para o formato JSON
        cadastros = append(cadastros, addUsuario)

        // Pega o valor do ID do ultimo dado buscado no banco e converte de String para Int
        count, _ = strconv.Atoi(addUsuario.ID) 
    }

    fmt.Printf("\nUsuarios atualizados!\n")
}

// Atualiza a struct Categorias com os dados do Banco
func AtualizaCategorias() {
    
    // Inicia conexão com Banco de dados no Azure 
    db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
    // Caso ocorrer algum erro na conexão, o ERRO será salvo em 'err'
    // Caso não houver ERRO o retorno será vazio 'nil'
    if err != nil {
        log.Println("Erro ao conectar com o Banco de dados:", err.Error())
    }

    // Select traz todos os Nomes das categorias e as quantidades de denuncias de cada um
    rows, err := db.Query("SELECT d.id_categoria, c.categoria, COUNT(d.id_categoria) FROM tab_denuncia d JOIN tab_categoria c ON d.id_categoria = c.id GROUP BY d.id_categoria, c.categoria") 
    if err != nil {
        log.Println("Erro no SELECT das Categorias:", err.Error())
    }

    // finaliza o comando para Query
    defer rows.Close() 
    // fecha conexão com o Banco
    defer db.Close()   
    
    // Zera o Array antes de buscar novos dados no banco.
    // Dessa forma evita dados repetidos
    categorias = categorias[:0]

    // rows.Next usado para varrer o objeto 'rows' e pegar os valores retornados da Query
    for rows.Next() {
        // Struct criada para receber os dados do banco de dados
        // Uma Struct permite receber dados de diferentes tipos: String, Int ...
        addCategoria := Categoria{} 

        // rows.Scan varre o objeto rows e salva os valores na variaveis citadas abaixo.
        // As variaveis são salvas na mesma ordem que são coletados do Banco de dados
        if err := rows.Scan(&addCategoria.ID, &addCategoria.Nome, &addCategoria.Total); err != nil { 
            log.Println("Erro ao salvar as categorias retornados do Banco:", err.Error())
        }
        
        // adiciona na struct principal os dados do banco
        // Sera está struct 'cadastros' que será convertida para o formato JSON
        categorias = append(categorias, addCategoria)

        // Pega o valor do ID do ultimo dado buscado no banco e converte de String para Int
        count, _ = strconv.Atoi(addCategoria.ID) 
    }

    fmt.Printf("\nCategorias atualizadas!\n")
}

func main() {

    AtualizaUsuarios()   
    AtualizaCategorias() 
    router := mux.NewRouter()
    router.HandleFunc("/cadastros/", GetCadastros).Methods("GET")
    router.HandleFunc("/cadastros/{id}", GetUsuario).Methods("GET")
    router.HandleFunc("/cadastros/", PostUsuario).Methods("POST")
    router.HandleFunc("/cadastros/{id}", DeleteUsuario).Methods("DELETE")
    router.HandleFunc("/categorias/", CriaJSON) // Não retorna nem envia nada, apenas atualiza o arquivo tabela.json
    log.Fatal(http.ListenAndServe(":8080", router)) // Server na porta 8080 [ localhost:8080 ]
}

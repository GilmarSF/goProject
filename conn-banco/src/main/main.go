// código apenas faz um select no banco de dados e exibe.
// Feito para provar o funcionamento da comunicação com banco de dados Azure
// Feito por José Luis

package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/denisenkom/go-mssqldb"
)

var dbDados string // string criada para receber os valores do objeto
var slcDados = make([]string, 0) // slice criado para guardar as strings que o dbDados recebe

func main(){
	connectionDB()   // inicia conexão com o azure
}

func connectionDB() {

	// Banco no azure do Zé
	db, err := sql.Open("mssql", "server=pwbt.database.windows.net;user id=admin-jose;password=123abc!@#;database=PWBT;port=1433")
	if err != nil { // se o retorno for 'nil' quer dizer que deu tudo certo
		log.Println("Open Failed: ", err.Error())
	}

	// db.Query usado para comandos no Banco
	//traz todas as linha do campo 'regiao' na tab_localidade
	rows, err := db.Query("select regiao from [dbo].[tab_localidade]") 
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close() // fecha o comando Query
	defer db.Close()   // fecha conexão com o Banco
	
	// rows.Next usado para varrer o objeto 'rows' e pegar os valores retornados da Query
	for rows.Next() {
		// percorre o objeto 'rows' salvando o que encontra na string 'dbDados'
		if err := rows.Scan(&dbDados); err != nil { 
			log.Fatal(err)
		}
		fmt.Printf("Categoria: %s \n", dbDados)
		slcDados = append(slcDados, dbDados) //salva o retorno neste slice "slice=array flexivel" append=adicionar no final
	}

	fmt.Println(slcDados)  // exibe todo conteudo do slice

	for i := 0; i < len(slcDados); i++{ //for para exibir separadamento cada elemento do slice de acordo com seu tamanho
		fmt.Println(slcDados[i])	
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nbye")
}

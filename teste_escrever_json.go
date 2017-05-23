package main

import (
    "fmt"
    "os"
    "encoding/json"
)

type CategoriaStruct struct{
    Categoria string    `json:"v"`
    Quantidade int      `json:"v"`
}

type GrupoCategoriaStruct struct{
    GrupoCateg []CategoriaStruct `json:"c"`
}

type LinhasStruct struct{
    Linhas []GrupoCategoriaStruct `json:"rows"`
}

type DenunciaStruct struct{
    Denuncia string `json:"label"`
    Quantidade string `json:"type"`
}

type ColunaStruct struct{
    Coluna []DenunciaStruct `json:"cols"`
}

func main() {

    categoria := CategoriaStruct{Categoria:"Assalto", Quantidade:12}
    //quantidadeCategoria := CategoriaStruct{Quantidade:15}

    b, err := json.Marshal(categoria)

    userFile := "tabela.json"
    fout, err := os.Create(userFile)        
    if err != nil {
        fmt.Println(userFile, err)
        return
    }

    fmt.Print(b)
    defer fout.Close()
    //fout.WriteString(byteb)
    fout.Write([]byte(b))
}
/*
{
    "cols":[
        {"label":"Denuncia","type":"string"},
        {"label":"Quantidade","type":"number"}
    ],
    "rows":[
        {
            "c":[
                {"v":"Assalto"},
                {"v":15}
            ]
        },
    ]
}*/


{
	"cols":[
		{"label":"Denuncia","type":"string"},
        {"label":"Quantidade","type":"number"}
	],
	"rows":[
        {
			"c":[
				{"v":"Assalto"},
				{"v":15}
			]
		},
	]
}

        {"c":[{"v":"Abuso Sexual"},{"v":2}]},
        {"c":[{"v":"Transito"},{"v":10}]},
        {"c":[{"v":"Violencia"},{"v":20}]},
        {"c":[{"v":"Assassinato"},{"v":10}]}

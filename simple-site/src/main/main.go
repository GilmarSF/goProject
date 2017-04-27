package main
 
import (
  "log"
  "net/http"
  "html/template"
)

// struct que sera usada para exibir as 'strings' na pagina
type content struct{
	Nome string
	Idade string
}
 
func main() {
  http.HandleFunc("/", index) // chama a pagina principal
  http.HandleFunc("/greet", greeter) // chama a pagina com a mensagem
  log.Println("Listening...")
  err := http.ListenAndServe(":80", nil) // servidor http na porta 80 padrão web
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
   }
}
 
func index(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("public/index.html") // salva o html em 't'
  t.Execute(w, nil) // exibe a pagina index.html
}
 
func greeter(w http.ResponseWriter, r *http.Request) {
  dados := content{} // dados do tipo struct
  dados.Nome = r.FormValue("nome") // salva o nome na struct dados
  dados.Idade = r.FormValue("idade")
  t, _ := template.ParseFiles("public/greeter.html") // salva o html em 't'
  err := t.Execute(w, dados) // para exibir a pagina greeter.html
  // se tiver erro
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

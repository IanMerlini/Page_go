package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// Definir os handlers para cada rota
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/salvar", salvarHandler)
	http.HandleFunc("/dados", dadosHandler)
	http.HandleFunc("/excluir", excluirHandler)

	// Iniciar o servidor na porta 8000
	http.ListenAndServe(":8000", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Carregar o arquivo HTML de login
	html, err := ioutil.ReadFile("login.html")
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo HTML", http.StatusInternalServerError)
		return
	}
	w.Write(html)
}

// Salvar o nome e senha no documento usuarios.txt
func salvarHandler(w http.ResponseWriter, r *http.Request) {
	// Obter os valores do formulário
	nome := r.FormValue("nome")
	senha := r.FormValue("senha")

	// Abrir o arquivo em modo de anexação
	arquivo, err := os.OpenFile("usuarios.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo de usuários.", http.StatusInternalServerError)
		return
	}
	// Escrever os valores no final do arquivo
	texto := fmt.Sprintf("Nome: %s\nSenha: %s\n", nome, senha)
	if _, err := arquivo.WriteString(texto); err != nil {
		http.Error(w, "Erro ao salvar usuário.", http.StatusInternalServerError)
		return
	}
	defer arquivo.Close()

	// Redirecionar para a página de dados
	http.Redirect(w, r, "/dados", http.StatusSeeOther)
}

func dadosHandler(w http.ResponseWriter, r *http.Request) {
	// Ler o conteúdo do arquivo de usuários
	conteudo, err := ioutil.ReadFile("usuarios.txt")
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo de usuários.", http.StatusInternalServerError)
		return
	}

	// Carregar o template a partir do arquivo dados.html
	tmpl, err := template.ParseFiles("dados.html")
	if err != nil {
		http.Error(w, "Erro ao carregar o template.", http.StatusInternalServerError)
		return
	}

	// Executar o template com o conteúdo do arquivo de usuários
	err = tmpl.Execute(w, string(conteudo))
	if err != nil {
		http.Error(w, "Erro ao executar o template.", http.StatusInternalServerError)
		return
	}
}

// Exclui o conteúdo do arquivo de usuários
func excluirHandler(w http.ResponseWriter, r *http.Request) {
	// Truncar o arquivo para apagar o conteúdo
	arquivo, err := os.OpenFile("usuarios.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo de usuários.", http.StatusInternalServerError)
		return
	}
	defer arquivo.Close()

	// Redirecionar para a página de dados
	http.Redirect(w, r, "/dados", http.StatusSeeOther)
}

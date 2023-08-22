package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html charset=UTF-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

// útil quando retornamos um tipo e um erro e só retornamos o tipo se não houver erro
// sendo que faz sentido no contexto usar panic, pois a aplicação não pode executar
// se não acontecer. Indicativo de uso é na função main, onde erros no meio dela
// devem fazer o programa parar
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// Faz o parsing de templates buscando no sistema de arquivos (FS) gerando pelo
// embed (/templates/fs.go)
func ParseFS(fs fs.FS, pattern string) (Template, error) {
	// o pacte html/template possui uma função para buscar o template embutido no fs
	// e fazer o parsing como o parse comum na função abaixo
	tpl, err := template.ParseFS(fs, pattern)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}


func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

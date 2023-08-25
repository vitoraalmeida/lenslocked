package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/vitoraalmeida/lenslocked/context"
	"github.com/vitoraalmeida/lenslocked/models"
)

// usado para determinar se um erro é para usuários ou para ser mostrado internamente
type public interface {
	Public() string
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// Cada requisição que o servidor do go recebe é tratada numa nova goroutine,
	// porem estamos lidando com um recurso compartilhado (nosso campo htmlTpl é um ponteiro)
	// então no caso de muitas requisições acontecerem ao mesmo tempo, é possível que o
	// que a página especifica que foi gerada para um usuário tenha sido renderizada com
	// base no mesmo template de outra requisição, ou seja, 2 clientes terão o mesmo token
	// csrf. Para evitar isso, clonaremos o template, assim cada goroutine usará uma cópia
	// diferente do template
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	// atualiza a função de template csrfFiel adicionada no parseFS com o conteúdo correto
	// que é o token csrf gerado pelo gorilla csrf middleware
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	// executa num buffer antes para que o status code não seja definido para sucesso
	// ainda que tenha acontecido um erro na execução do template. O servidor do go
	// define como sucesso qualquer resposta que não tenha statuc code definido manualmente.
	// Quando escrevemos no response writer uma vez, na próxima vez que escrevermos (no caso
	// o tratamento do erro abaixo) a definição do status não vai ser levada em consideração
	// pois já foi definida antes
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	// escreve no http.ResponseWriter apenas se não houve erro
	io.Copy(w, &buf)
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
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// cria um template do zero
	tpl := template.New(patterns[0])
	// registra uma função personalizada para ser utilizada no template engine do go
	// deve ser feito antes de fazer o parsing do template para que a função possa ser reconhecida
	// durante o parsing
	tpl = tpl.Funcs(
		template.FuncMap{
			// cria uma função com o nome de csrfField que inclui um código no template
			// Inicialmente, a função csrfField não possui o conteúdo que queremos (que
			// é um campo de formulário contentdo o token csrf), pois queremos fazer
			// o parsing do template quando o servidor inicia para sabermos se o parsing
			// ocorreu corretamente. Essa função será atualizada posteriormente quando
			// existir uma request de fato, que é de onde vem o token csrf
			"csrfField": func() (template.HTML, error) {
				// a função retornará erro quando o template for executado sem a substituição
				// da função que de fato gera o csrf token
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				// a função retornará erro quando o template for executado sem a substituição
				// da função que de fato checa o usuário
				return "", fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)
	// o pacte html/template possui uma função para buscar o template embutido no fs
	// e fazer o parsing como o parse comum na função abaixo
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}

// func Parse(filepath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filepath)
// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}
// 	return Template{
// 		htmlTpl: tpl,
// 	}, nil
// }

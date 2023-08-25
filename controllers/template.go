package controllers

import "net/http"

// interface utilizada para desacoplar controllers de views
// não interessa a implementação do template, desde que tenha
// o método Execute
type Template interface {
	// o request é necessário para executar o template pois para a proteção contra
	// csrf, precisamos do token csrf que é gerado pelo middleware do gorilla csrf,
	// que só existe após um request ser feito
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}

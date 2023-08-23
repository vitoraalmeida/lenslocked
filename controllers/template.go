package controllers

import "net/http"

// interface utilizada para desacoplar controllers de views
// não interessa a implementação do template, desde que tenha
// o método Execute
type Template interface {
	Execute(w http.ResponseWriter, data interface{})
}

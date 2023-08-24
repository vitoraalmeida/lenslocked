package main

import (
	"context"
	"fmt"
)

// usar tipos customizados como keys em contexts, pois
// assim ele só poderá ser usado com valores que obedecem
// ao tipo definido, assim outros pacotes não conseguirão
// sobrescrever a chave por acidente
// não expor o tipo para que código fora do pacote não possa
// usar a mesma chave
type ctxKey string

const (
	// não expor o a chave para que código fora do pacote não possa
	// usar a mesma chave
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	fmt.Println(value)

	// type assertions
	strValue, ok := value.(int)
	if !ok {
		fmt.Println("não é string")
	}
	fmt.Printf("valor: %s", strValue)
}

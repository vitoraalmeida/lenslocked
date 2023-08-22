// o pacote embed possibilita que os templates façam parte do binário final,
// assim o servidor não precisa que os templates estejam em algum lugar no
// sistema de arquivos do host, pois estarão embutidos (embeded) no executável
// faz isso gerando um sistema de arquivos interno contendo os arquivos definidos
// declaramos uma variável que será o sistema de arquivos e acima da declaração
// colocamos um comentário explicitando um padrão de quais arquivos queremos que
// sejam embutidos, que estejam no diretório do pacote em questão
package templates

import "embed"

//go:embed *
var FS embed.FS

// o tipo embed.FS implementa a interface io.FS
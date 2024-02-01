package utils

import (
	"github.com/IgorCastilhos/toolkit/v2" // Módulo para gerar strings aleatórias
	"math/rand"
	"time"
)

var tools toolkit.Tools // Variável global que contém instâncias de ferramentas do toolkit

// init é chamada automaticamente antes da primeira execução de qualquer outra função.
// Aqui, ela é usada para inicializar a semente do gerador de números aleatórios, garantindo que os números gerados sejam diferentes a cada execução.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt gera um inteiro aleatório entre min e max (inclusive).
// A função rand.Int63n gera um número aleatório entre max-min (+1, pois se der 0 ele retornará panic)
// e então adiciona min ao resultado para ajustar o intervalo.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0->max-min
}

// RandomOwner gera uma string aleatória de comprimento fixo.
// A função utiliza tools.RandomString do toolkit importado para criar uma string aleatória de 10 caracteres.
func RandomOwner() string {
	length := 10
	return tools.RandomString(length)
}

// RandomMoney gera um valor int64 aleatório.
// A função utiliza RandomInt para gerar um valor inteiro aleatório entre 0 e 1000.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency gera uma string representando uma moeda aleatória.
// Atualmente, a função retorna "BRL" como única opção, mas pode ser expandida para incluir mais moedas.
func RandomCurrency() string {
	currencies := []string{"BRL"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

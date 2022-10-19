package main

import "fmt"

type message_type int

const (
	// Election message
	ELECTION message_type = iota
	// Election response message
	// ELECTION_RESPONSE
	// Election winner message
	ELECTION_WINNER
	NEW_NODE
)

func (s message_type) String() string {
	switch s {
	case ELECTION:
		return "Eleicão"
	// case ELECTION_RESPONSE:
	// return "Resposta Eleição"
	case ELECTION_WINNER:
		return "Resultado da eleição"
	case NEW_NODE:
		return "Novo nó"
	}
	return "Desconhecido"
}

type mensagem struct {
	tipo  message_type // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo map[int]int  // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
	owner int          // id do processo que enviou a mensagem
}

func (m mensagem) String() string {
	corpo := "("
	for _, val := range m.corpo {
		corpo += fmt.Sprintf("%d", val)
	}

	corpo += ")"
	return fmt.Sprintf("%s: %v", m.tipo, corpo)
}

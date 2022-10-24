package main

import "fmt"

type message_type int

const (
	// Election message
	ELECTION message_type = iota
	// Election winner message
	ELECTION_WINNER
)

func (s message_type) String() string {
	switch s {
	case ELECTION:
		return "Eleicão"
	// case ELECTION_RESPONSE:
	// return "Resposta Eleição"
	case ELECTION_WINNER:
		return "Resultado da eleição"
		// case NEW_NODE:
		// 	return "Novo nó"
	}
	return "Desconhecido"
}

type mensagem struct {
	tipo  message_type // tipo da mensagem para fazer o controle do que fazer (eleição, confirmacao da eleicao)
	corpo interface{}  // conteudo da mensagem para colocar os ids (usar um tamanho ocmpativel com o numero de processos no anel)
	owner int          // id do processo que enviou a mensagem
}

func (m mensagem) String() string {
	switch m.tipo {
	case ELECTION:
		body := m.corpo.(map[int]int)
		corpo := "[ "
		for _, val := range body {
			corpo += fmt.Sprintf("%d ", val)
		}

		corpo += "]"
		return fmt.Sprintf("%s: %v", m.tipo, corpo)

	case ELECTION_WINNER:
		// case NEW_NODE:
		body := m.corpo.(*Node)
		return fmt.Sprintf("%s: %v", m.tipo, body.TaskId)
	}
	return "Desconhecido"
}

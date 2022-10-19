package main

import (
	"log"
	"math/rand"
	"time"
)

type Node struct {
	TaskId  int
	Next    *Node
	Message chan mensagem

	IsMaster bool
	Master   int  // define index do master
	IsAlive  bool // Define se o node está vivo ou morto
}

func (n *Node) Handle(done <-chan bool) {
	timeBetweenChecks := time.Duration(rand.Intn(500)+500) * time.Millisecond // 100ms - 200ms
	checkMaster := time.NewTicker(timeBetweenChecks)
	for {
		select {
		case msg := <-n.Message:
			n.ReceiveMessage(msg)
		case <-done:
			return
		case <-checkMaster.C:
			if !n.MasterIsAlive() {
				log.Printf("Node %d: Master(%d) está morto, iniciando eleição", n.TaskId, n.Master)
				n.InitElection()
			}
		}
	}
}

func (n *Node) InitElection() {
	n.SendMessageToNext(mensagem{
		owner: n.TaskId,
		tipo:  ELECTION,
		corpo: map[int]int{
			n.TaskId: n.TaskId,
		},
	})
}

func (n *Node) SendMessageToNext(msg mensagem) {
	next := n.Next
	if next != nil {
		if n.IsAlive {
			log.Printf("Node %d: Enviando para o Node %d a msg: %s", n.TaskId, n.Next.TaskId, msg)
		}
		next.Message <- msg
	} else {
		log.Printf("Node %d: Não existe ninguem vivo para receber a mensagem: %s", n.TaskId, msg)
	}
}

func (n *Node) ReceiveMessage(msg mensagem) {
	// log.Printf("Node %d: recebeu - %s", n.TaskId, msg)
	switch msg.tipo {
	case ELECTION:
		n.receiveElectionMessage(msg)
	case ELECTION_WINNER:
		n.ReceiveElectionResponseMessage(msg)
		// case NEW_NODE:
		// 	n.ReceiveNewNodeMessage(msg)
	}

}

func (n *Node) receiveElectionMessage(msg mensagem) {
	corpo, ok := msg.corpo.(map[int]int)
	if !ok {
		log.Printf("Node %d: Erro ao ler corpo da mensagem: %s", n.TaskId, msg)
		return
	}
	if msg.owner != n.TaskId {
		if n.IsAlive {
			corpo[n.TaskId] = n.TaskId
		}
		n.SendMessageToNext(msg)
	} else {
		maior := 0

		for _, val := range corpo {
			if val > maior { // TODO: Criar um criterio melhor para decidir o vencedor
				maior = val
			}
		}

		n.Master = maior
		n.SendMessageToNext(mensagem{
			owner: n.TaskId,
			tipo:  ELECTION_WINNER,
			corpo: rede.Get(maior),
		})
	}
}

func (n *Node) ReceiveElectionResponseMessage(msg mensagem) {
	corpo, ok := msg.corpo.(*Node)
	if !ok {
		log.Printf("Node %d: Erro ao ler corpo da mensagem: %s", n.TaskId, msg)
		return
	}

	if corpo.TaskId == n.TaskId {
		n.IsMaster = true
	} else {
		n.IsMaster = false
	}
	n.Master = corpo.TaskId

	if msg.owner != n.TaskId {
		n.SendMessageToNext(msg)
	} else {
		log.Printf("Node %d: Recebi de volta a mensagem: %s", n.TaskId, msg)
	}
}

// func (n *Node) ReceiveNewNodeMessage(msg mensagem) {
// 	if msg.owner != n.TaskId {
// 		corpo, ok := msg.corpo.(*Node)
// 		if !ok {
// 			log.Printf("Node %d: Erro ao ler corpo da mensagem: %s", n.TaskId, msg)
// 			return
// 		}
// 		n.rede.Nodes = append(n.rede.Nodes, corpo)
// 		n.SendMessageToNext(msg)
// 	} else {
// 		log.Printf("Node %d: Recebi de volta a mensagem: %s", n.TaskId, msg)
// 	}
// }

func (n *Node) HealthCheck() bool {
	return n.IsAlive
}

func (n *Node) MasterIsAlive() bool {
	master := rede.Get(n.Master)
	if master == nil {
		return false
	}
	return (n.IsMaster) || master.HealthCheck()

}

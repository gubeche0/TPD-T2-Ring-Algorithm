package main

import "log"

type Rede struct {
	Nodes []*Node
}

func (r *Rede) GetMaster() *Node {
	for _, node := range r.Nodes {
		if node.IsMaster {
			return node
		}
	}
	return nil
}

func (r *Rede) InsereNode(node *Node) {
	if len(r.Nodes) == 0 {
		node.Master = 0
		node.IsMaster = true
		node.Next = nil

	} else {
		r.Nodes[len(r.Nodes)-1].Next = node
		node.Next = r.Nodes[0]

		node.IsMaster = false
		node.Master = r.GetMaster().TaskId
	}

	r.Nodes = append(r.Nodes, node)

	node.rede = *r

	node.SendMessageToNext(mensagem{
		tipo:  NEW_NODE,
		owner: node.TaskId,
		corpo: map[int]int{0: node.TaskId},
	})
}

type Node struct {
	TaskId  int
	Next    *Node
	Message chan mensagem

	IsMaster bool
	Master   int // define index do master
	IsAlive  bool

	rede Rede
}

func (n *Node) Handle(done <-chan bool) {
	for {
		select {
		case msg := <-n.Message:
			n.ReceiveMessage(msg)
		case <-done:
			return
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
		next.Message <- msg
	} else {
		log.Printf("%d - Não existe ninguem vivo para receber a mensagem: %s", n.TaskId, msg)
	}
}

func (n *Node) ReceiveMessage(msg mensagem) {
	log.Printf("%d recebeu - %s", n.TaskId, msg)
	switch msg.tipo {
	case ELECTION:
		n.receiveElectionMessage(msg)
	case ELECTION_WINNER:
		n.ReceiveElectionResponseMessage(msg)
	case NEW_NODE:
		n.ReceiveNewNodeMessage(msg)
	}

}

func (n *Node) receiveElectionMessage(msg mensagem) {
	if msg.owner != n.TaskId {
		if n.IsAlive {
			msg.corpo[n.TaskId] = n.TaskId
		}
		n.SendMessageToNext(msg)
	} else {
		maior := 0

		for _, val := range msg.corpo {
			if val > maior { // TODO: Criar um criterio melhor para decidir o vencedor
				maior = val
			}
		}

		n.SendMessageToNext(mensagem{
			owner: n.TaskId,
			tipo:  ELECTION_WINNER,
			corpo: map[int]int{0: maior},
		})
	}
}

func (n *Node) ReceiveElectionResponseMessage(msg mensagem) {
	if msg.owner != n.TaskId {
		if msg.corpo[0] == n.TaskId {
			n.IsMaster = true
		} else {
			n.IsMaster = false
		}
		n.Master = msg.corpo[0]

		n.SendMessageToNext(msg)
	} else {
		log.Printf("%d, Recebi de volta a mensagem: %s", n.TaskId, msg)
	}
}

func (n *Node) ReceiveNewNodeMessage(msg mensagem) {
	if msg.owner != n.TaskId {
		n.Master = msg.corpo[0]
		if n.Master == n.TaskId {
			n.IsMaster = true
		} else {
			n.IsMaster = false
		}
		n.SendMessageToNext(msg)
	} else {
		log.Printf("%d, Recebi de volta a mensagem: %s", n.TaskId, msg)
	}
}

func (n *Node) HealthCheck() bool {

	return n.IsAlive
}

func (n *Node) MasterIsAlive() bool {
	return (n.IsMaster) || n.rede.GetMaster().HealthCheck()

}

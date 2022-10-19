package main

import (
	"fmt"
	"log"
)

type Rede struct {
	Nodes []*Node
}

// Metodo interno
func (r *Rede) getMaster() *Node {
	for _, node := range r.Nodes {
		if node.IsMaster {
			return node
		}
	}
	return nil
}

func (r *Rede) Get(nodeId int) *Node {
	for _, node := range r.Nodes {
		if node.TaskId == nodeId {
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
		node.Master = r.getMaster().TaskId
	}

	r.Nodes = append(r.Nodes, node)

	// node.SendMessageToNext(mensagem{
	// 	tipo:  NEW_NODE,
	// 	owner: node.TaskId,
	// 	corpo: node,
	// })
}

func (r *Rede) Debug() {
	d := ""
	for _, node := range r.Nodes {
		d += fmt.Sprintf("ID: %d, IsMaster: %v, IsVivo:%v, MasterConhecido: %d\n", node.TaskId, node.IsMaster, node.IsAlive, node.Master)
	}
	log.Printf("\n\n[DEBUG] Rede:\n%s \n", d)

}

package main

import "time"

func main() {

	rede := &Rede{}

	nodes := []*Node{
		{
			TaskId:  0,
			IsAlive: true,
			Message: make(chan mensagem),
		},
		{
			TaskId:  1,
			IsAlive: true,
			Message: make(chan mensagem),
		},
		{
			TaskId:  2,
			IsAlive: true,
			Message: make(chan mensagem),
		},
	}

	done := make(chan bool) // channel to stop the nodes

	for _, node := range nodes {
		go node.Handle(done) // Inicializa os processos dos nós
	}

	for _, node := range nodes {
		rede.InsereNode(node) // Insere os nós na rede
	}

	time.Sleep(time.Second * 5)

	done <- true // stop the nodes

}

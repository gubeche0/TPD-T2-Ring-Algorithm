package main

import (
	"log"
	"time"
)

var (
	rede = &Rede{}
)

func main() {

	// Utiliza seed padrao para gerar sempre a mesma sequencia de numeros aleatorios para teste
	// rand.Seed(time.Now().UnixNano())

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
		{
			TaskId:  3,
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

	time.Sleep(time.Second)
	log.Println("Rede iniciada\n ")

	log.Println("Matando o node 0")
	nodes[0].IsAlive = false
	time.Sleep(time.Second)

	rede.Debug()

	log.Println("Matando o node 2")
	nodes[2].IsAlive = false
	time.Sleep(time.Second * 2)

	rede.Debug()

	log.Println("Matando o node 3")
	nodes[3].IsAlive = false
	time.Sleep(time.Second)

	rede.Debug()

	log.Println("Revivendo o node 0")
	log.Println("Revivendo o node 2")
	log.Println("Revivendo o node 3")
	nodes[0].IsAlive = true
	nodes[2].IsAlive = true
	nodes[3].IsAlive = true
	time.Sleep(time.Second)
	rede.Debug()

	log.Println("Matando o node 1")
	nodes[1].IsAlive = false
	time.Sleep(time.Second)

	rede.Debug()

	log.Println("Fim")
	time.Sleep(time.Second)

	done <- true // stop the nodes

}

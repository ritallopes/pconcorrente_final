package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var POINT_TO_WIN = 5 //quantidade de pontos para vencder
/* Um jogador*/
type person struct {
	name      string
	lostBalls int //bola perdida e ponto adversario
}

func newPerson(name string) *person {
	p := person{name: name}
	p.lostBalls = 0
	return &p
}
func (p *person) lost1Ball() int {
	p.lostBalls = p.lostBalls + 1
	return p.lostBalls
}

func play(waitGp *sync.WaitGroup, tennisCourt chan int, playerCurrent *person) {
	defer waitGp.Done() //para finalizar
	for {               //loop infinito
		jackpot, open := <-tennisCourt
		if !open {
			fmt.Println(playerCurrent.name, " ganhou!!")
			return
		}
		fmt.Print("Jogada ", jackpot)
		fmt.Println(" - ", playerCurrent.name, " recebeu a bola")
		num := rand.Intn(100)
		fmt.Println("Numero randomico", num)
		if num <= 50 {
			playerCurrent.lost1Ball()
			fmt.Println(playerCurrent.name, " não rebateu")
			if playerCurrent.lostBalls >= POINT_TO_WIN {
				fmt.Println(playerCurrent.name, " perdeu!")
				fmt.Println()
				close(tennisCourt)
				return
			}

		} else {
			fmt.Println(playerCurrent.name, " rebateu!")
		}
		fmt.Println("")

		jackpot++
		tennisCourt <- jackpot

	}

}

/*
Executa assim que inicia com uma semente
objetivo e nao repetir o numero randomico, mas usar
o momento como base
*/
func init() {
	rand.Seed(time.Now().UnixNano())
}

func declareVictorius(p1 *person, p2 *person) {
	fmt.Println(p1.name, " fez ", p2.lostBalls)
	fmt.Println(p2.name, " fez ", p1.lostBalls)
}

func main() {
	var sets int
	flag.IntVar(&sets, "sets", 5, "numero de sets no game")
	var games int
	flag.IntVar(&games, "games", 5, "numero de games no jogo")
	var points int
	flag.IntVar(&points, "points", 5, "numero de points por set")
	flag.Parse()

	var waitGp sync.WaitGroup // wait group passado como ponteiro porque vou usar em várias funcoes

	fmt.Println("Iniciando jogo!!!")
	player1 := person{name: "Player 01", lostBalls: 0}
	player2 := person{name: "Player 02", lostBalls: 0}

	tennisCourt := make(chan int) //recurso compartilhado

	waitGp.Add(2) //Adicionar com dois jogadores

	go play(&waitGp, tennisCourt, &player1)
	go play(&waitGp, tennisCourt, &player2)

	tennisCourt <- 1
	waitGp.Wait()
	declareVictorius(&player1, &player2)
}

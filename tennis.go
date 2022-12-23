// Pacote contendo toda a implementação do
// Problema do Jogo de Tenis
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var POINT_TO_WIN = 4 //quantidade de pontos para vencder

/* Um jogador*/
type person struct {
	name      string
	lostBalls int //bola perdida e ponto adversario
	wonGames  int
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

func playGame(waitGp *sync.WaitGroup, tennisCourt chan int, playerCurrent *person) {
	defer waitGp.Done() //para finalizar
	for {               //loop infinito
		jackpot, open := <-tennisCourt
		if !open { //verifica se o jogo ja finalizou
			fmt.Println(playerCurrent.name, " ganhou!!") //Caso sim, o jogador atual venceu
			return
		}

		fmt.Print("Jogada ", jackpot)
		fmt.Println(" - ", playerCurrent.name, " recebeu a bola")
		num := rand.Intn(100) //gera um numero randomico menor que 100
		if num <= 50 {        //menor que 50 perdeu a bola
			playerCurrent.lost1Ball() //perdeu a bola, ponto para o adversario
			fmt.Println(playerCurrent.name, " não rebateu")
			if playerCurrent.lostBalls >= POINT_TO_WIN {
				fmt.Println(playerCurrent.name, " perdeu!") //o adversario atingiu a pontuacao
				fmt.Println()
				close(tennisCourt)
				return
			}

		} else {
			fmt.Println(playerCurrent.name, " rebateu e esperando a bola!") //passa para o adversario sem perder
		}
		fmt.Println("")

		jackpot++ //Proxima jogada
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

// Funcao para contagem dos pontos finais
func declareVictorius(p1 *person, p2 *person) {
	fmt.Println(p1.name, " fez ", p2.lostBalls)
	fmt.Println(p2.name, " fez ", p1.lostBalls)
}

/*funcao principal*/
func main() {
	//Leitura da quantidade de pontos definidos
	var points int
	flag.IntVar(&points, "points", 4, "numero de points por game")
	flag.Parse()
	//RF: os pontos precisam ser maiores que quatro
	if points >= 4 {
		POINT_TO_WIN = points
	}

	var waitGp sync.WaitGroup // wait group passado como ponteiro porque vou usar em várias funcoes

	fmt.Println("Iniciando jogo!!!")
	player1 := person{name: "Player 01", lostBalls: 0}
	player2 := person{name: "Player 02", lostBalls: 0}

	tennisCourt := make(chan int) //recurso compartilhado

	waitGp.Add(2) //Adicionar com dois jogadores

	go playGame(&waitGp, tennisCourt, &player1)
	go playGame(&waitGp, tennisCourt, &player2)

	tennisCourt <- 1
	waitGp.Wait()
	declareVictorius(&player1, &player2)
}

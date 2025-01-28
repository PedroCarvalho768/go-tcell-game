package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"time"
)

const tick = time.Second

var ticker *time.Ticker

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func drawString(screen tcell.Screen, x, y int, msg string) {
	for i, char := range msg {
		screen.SetContent(x+i, y, char, nil, tcell.StyleDefault)
	}
}

func setupCoins(level int) []*Sprite {
	coins := make([]*Sprite, level+2)
	for i := range level + 2 {
		coins[i] = NewSprite('O', rand.Intn(20), randRange(4, 30))
	}
	return coins
}

func main() {
	// Creating a screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	defer screen.Fini() // Destroying it at the end of main

	// initializing the screen
	err = screen.Init()
	if err != nil {
		log.Fatal(err)
	}

	// game init section
	player := NewSprite('@', 10, 10)
	coins := setupCoins(1)

	timeLeft := 10
	score := 0
	level := 1

	ticker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				timeLeft--
			}
		}
	}()

	// game loop
	running := true
	for running {
		// draw logic
		screen.Clear()
		//__ __ __ __ __ __ drawing __ __ __ __ __ __

		player.Draw(screen)

		for _, coin := range coins {
			coin.Draw(screen)
		}

		// UI
		drawString(screen, 1, 1, "GRAB ALL THE COINS")
		drawString(screen, 1, 2, fmt.Sprintf("TIME LEFT: %d", timeLeft))
		drawString(screen, 1, 3, fmt.Sprintf("Score: %d", score))
		drawString(screen, 1, 4, fmt.Sprintf("Level: %d", level))

		//___________________________________________
		screen.Show()

		// update logic

		playerMoved := false

		// getting the event
		ev := screen.PollEvent()
		// checking event type
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// checking the Event Key
			switch ev.Rune() {
			case 'q':
				running = false
			// Player Movement
			case 'w':
				player.Y -= 1
				playerMoved = true
			case 'a':
				player.X -= 1
				playerMoved = true
			case 's':
				player.Y += 1
				playerMoved = true
			case 'd':
				player.X += 1
				playerMoved = true
			}
		}

		// end
		if timeLeft == 0 {
			fmt.Println("YOU LOSE")
			fmt.Printf("Your score was: %d", score)
			running = false
		}

		// check for coin collisions
		if playerMoved {
			coinCollectedIndex := -1
			for i, coin := range coins {
				if coin.X == player.X && coin.Y == player.Y {
					// Collect coin
					coinCollectedIndex = i
					score++
				}
			}

			// handle coin collisions
			if coinCollectedIndex > -1 {
				coins[coinCollectedIndex] = coins[len(coins)-1]
				coins = coins[0 : len(coins)-1]

				if len(coins) == 0 {
					level++
					coins = setupCoins(level)
				}
			}
		}
	}
}

package main

import (
	"log"

	"github.com/sago35/koebiten"
	tensecgo "github.com/sago35/koebiten/games/10secGo/10secGo"
)

func main() {
	koebiten.SetWindowSize(128, 64)
	koebiten.SetWindowTitle("10sec Go")
	game := tensecgo.NewGame()

	if err := koebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package tensecgo

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/sago35/koebiten"
)

type Game struct {
}

func NewGame() *Game {
	game := &Game{}
	return game
}

// Game update process
func (g *Game) Update() error {
	return nil
}

// Screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return 128, 64
}

func (g *Game) Draw(screen *koebiten.Image) {
	draw()
}

var (
	scene   = "title"
	score   = 0
	counts  = 90
	games   = 300
	ranking = []int{}
	write   = false
	read    = false

	buf bytes.Buffer
)

func draw() {
	switch scene {
	case "title":
		drawTitle()
	case "countstart":
		drawGameStart()
	case "countend":
		drawGameEnd()
	case "ranking":
		drawRanking()
	}
}

func drawTitle() {
	koebiten.Println("click Rotate to")
	koebiten.Println("start 10 Sec Go")
	koebiten.Println("rotating to Ranking")
	if onStartBtn() {
		scene = "countstart"
	}
	if onRankingBtn() {
		scene = "ranking"
	}
}

func drawGameStart() {
	// ３秒カウントダウンする
	if counts > 0 {
		koebiten.Println(fmt.Sprintf("%d", counts/30+1))
		counts--
		return
	}

	// 10秒間ゲームをプレイ
	if counts == 0 {
		if games > 0 {
			koebiten.Println(fmt.Sprintf("%d", games/30+1))
			games--

			if onClickBtn() {
				score++
			}

			if games == 0 {
				scene = "countend"
				return
			}
		}
	}
}

func drawGameEnd() {
	koebiten.Println("Game End")
	koebiten.Println(fmt.Sprintf("Score: %d", score))

	// 1回だけ得点を書き込む
	if !write {
		ranking = append(ranking, score)
		buf = writeScoresToMemory(ranking)
		write = true
	}

	if onStartBtn() {
		scene = "title"
		score = 0
		counts = 90
		games = 300
		write = false
		return
	}
	if onRankingBtn() {
		scene = "ranking"
	}
}

func drawRanking() {
	koebiten.Println("Ranking")

	if !read {
		ranking = readScoresFromMemory(&buf)
		read = true
	}

	for i := 0; i < len(ranking) && i < 3; i++ {
		if i == 0 {
			koebiten.Println(fmt.Sprintf("%dst: %d", i+1, ranking[i]))
		} else if i == 1 {
			koebiten.Println(fmt.Sprintf("%dnd: %d", i+1, ranking[i]))
		} else if i == 2 {
			koebiten.Println(fmt.Sprintf("%drd: %d", i+1, ranking[i]))
		}
	}
	koebiten.Println("click or rotarte to Title")

	if onStartBtn() || onRankingBtn() {
		scene = "title"
		score = 0
		counts = 90
		games = 300
		write = false
		return
	}
}

func onStartBtn() bool {
	if len(koebiten.AppendJustPressedKeys(nil)) > 0 {
		if koebiten.AppendJustPressedKeys(nil)[0] == koebiten.KeyRotaryButton {
			return true
		}
	}
	return false
}

func onRankingBtn() bool {
	if len(koebiten.AppendJustPressedKeys(nil)) > 0 {
		if koebiten.AppendJustPressedKeys(nil)[0] == koebiten.KeyRotaryLeft || koebiten.AppendJustPressedKeys(nil)[0] == koebiten.KeyRotaryRight {
			return true
		}
	}
	return false
}

func onClickBtn() bool {
	return len(koebiten.AppendJustPressedKeys(nil)) > 0
}

// スコアデータの書き込み
func writeScoresToMemory(scores []int) bytes.Buffer {
	sort.Slice(scores, func(i, j int) bool { return scores[i] > scores[j] })
	for _, score := range scores {
		fmt.Fprintln(&buf, score)
	}

	return buf
}

// スコアデータの読み込み
func readScoresFromMemory(buf *bytes.Buffer) []int {
	var scores []int
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		score, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			log.Fatal(err)
		}
		scores = append(scores, score)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return scores
}

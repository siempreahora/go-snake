package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// const PaddleSymbol = '1'
const SnakeSymbol = 0x2588
const AppleSymbol = 0x25CF

const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = "||"

var VelocityRow = 1
var VelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var isGamePaused bool
var debugLog string

var gameObjects = []*GameObject{}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	initGameState()
	inputChan := initUserInput()

	for {
		handleInput(readInput(inputChan))
		UpdateState()
		drawState()

		time.Sleep(70 * time.Millisecond)
	}

	// screenWidth, screenHeight := screen.Size()
	// winner := getWinner()
	// pringStringCentered(screenHeight/2-1, screenWidth/2, "Game Over!")
	// pringStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins!", winner))

	// screen.Show()

	time.Sleep(3 * time.Second)
	screen.Fini()
	os.Exit(0)
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

}

func initGameState() {
	gameObjects = []*GameObject{}
}

func UpdateState() {
	if isGamePaused {
		return
	}

	// screenWidth, screenHeight := screen.Size()
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
}

func drawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	printString(0, 0, debugLog)
	for _, obj := range gameObjects {
		drawObject(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

// func gameEnd() bool {
// 	return getWinner() != ""
// }

// func getWinner() string {
// 	screenWidth, _ := screen.Size()
// 	if ball.col < 0 {
// 		return "Player 2"
// 	} else if ball.col >= screenWidth {
// 		return "Player 1"
// 	} else {
// 		return ""
// 	}
// }

func pringStringCentered(row, col int, str string) {
	col = col - len(str)/2
	printString(row, col, str)
}

func printString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func drawObject(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func initUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func readInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}

func handleInput(key string) {
	// _, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	}
	// } else if key == "Rune[w]" && playerLeft.row > 0 {
	// 	playerLeft.row--
	// } else if key == "Rune[s]" && playerLeft.row+playerLeft.height < screenHeight {
	// 	playerLeft.row++
	// } else if key == "Up" && playerRight.row > 0 {
	// 	playerRight.row--
	// } else if key == "Down" && playerRight.row+playerRight.height < screenHeight {
	// 	playerRight.row++
	// }
}

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
const GameFrameSymbol = '║'

type Point struct {
	row, col int
}

type Snake struct {
	body           []*Point
	velRow, velCol int
	symbol         rune
}

type Apple struct {
	point  *Point
	symbol rune
}

var screen tcell.Screen
var snake *Snake
var apple *Apple
var isGamePaused bool
var debugLog string

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

/* ****************************************** */
// INITIALIZING/UPDATING THE GAME //
/* ****************************************** */

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
	snake = &Snake{
		body: []*Point{
			{row: 9, col: 3},
			{row: 8, col: 3},
			{row: 7, col: 3},
			{row: 6, col: 3},
			{row: 5, col: 3},
		},
		velRow: -1,
		velCol: 0,
		symbol: SnakeSymbol,
	}
	apple = &Apple{
		point:  &Point{row: 10, col: 10},
		symbol: AppleSymbol,
	}
}

func UpdateState() {
	if isGamePaused {
		return
	}

	UpdateSnake()
}

func UpdateSnake() {
	head := snake.body[len(snake.body)-1]
	snake.body = append(snake.body,
		&Point{
			row: head.row + snake.velRow,
			col: head.col + snake.velCol,
		})

	snake.body = snake.body[1:]
}

/* ****************************************** */
// DRAWING THE GAMESTATES //
/* ****************************************** */

func drawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	printString(0, 0, debugLog)
	drawGameFrame()
	drawSnake()
	drawApple()

	screen.Show()
}

func drawGameFrame() {

	gameTLRow, gameTLCol := getTopLeft()
	row, col := gameTLRow-1, gameTLCol-1
	width, height := GameFrameWidth+2, GameFrameHeight+2

	drawRectPerimeter(row, col, width, height, GameFrameSymbol)
	// drawRectPerimeter(row+1, col+1, width, height, '*')
	// drawRectPerimeter(row+1, col+1, GameFrameWidth, GameFrameHeight, '*')

}

func drawSnake() {
	for _, p := range snake.body {
		drawInsideRectArea(p.row, p.col, 1, 1, snake.symbol)
	}
}

func drawApple() {
	drawInsideRectArea(apple.point.row, apple.point.col, 1, 1, apple.symbol)
}

func drawRectPerimeter(row, col, width, height int, ch rune) {

	// ╔ ╗ ╝ ╚ ═

	// Top wall of the box
	screen.SetContent(col, row, '╔', nil, tcell.StyleDefault)
	for c := 2; c < width; c++ {
		screen.SetContent(col+c-1, row, '═', nil, tcell.StyleDefault)
	}
	screen.SetContent(col+width-1, row, '╗', nil, tcell.StyleDefault)

	// Middle/Side walls of the box
	for r := 1; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, ch, nil, tcell.StyleDefault)
	}

	// Bottom wall of the box
	screen.SetContent(col, row+height-1, '╚', nil, tcell.StyleDefault)
	for c := 2; c < width; c++ {
		screen.SetContent(col+c-1, row+height-1, '═', nil, tcell.StyleDefault)
	}
	screen.SetContent(col+width-1, row+height-1, '╝', nil, tcell.StyleDefault)

}

func drawRectArea(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func drawInsideRectArea(row, col, width, height int, ch rune) {
	r, c := getTopLeft()
	drawRectArea(row+r, col+c, width, height, ch)
}

/* ****************************************** */
// USER INPUT AND HANDLING //
/* ****************************************** */

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

/* ****************************************** */
// CALCULATIONS AND PRINTING STRINGS //
/* ****************************************** */

func getTopLeft() (int, int) {
	screenWidth, screenHeight := screen.Size()

	return (screenHeight/2 - GameFrameHeight/2) - 1, (screenWidth/2 - GameFrameWidth/2) - 1
}

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

/* ****************************************** */
// END OF GAME //
/* ****************************************** */

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

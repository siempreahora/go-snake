package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

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
var score int
var pointsToClear []*Point
var isGameOver bool
var isGamePaused bool
var debugLog string

func main() {
	initScreen()
	initGameState()
	screen.HideCursor()
	inputChan := initUserInput()

	for !isGameOver {
		handleInput(readInput(inputChan))
		UpdateState()
		drawState()

		time.Sleep(75 * time.Millisecond)
	}

	gameEnd()
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
		point:  &Point{row: rand.Intn(GameFrameHeight), col: rand.Intn(GameFrameWidth)},
		symbol: AppleSymbol,
	}
}

func UpdateState() {
	if isGamePaused {
		drawGameFrame()
		return
	}

	UpdateSnake()
	UpdateAppple()
}

func UpdateSnake() {
	head := getSnakeHead()
	snake.body = append(snake.body,
		&Point{
			row: head.row + snake.velRow,
			col: head.col + snake.velCol,
		})

	if !appleInsideSnake() {
		snake.body = snake.body[1:]
	} else if apple.symbol == '■' {
		score += 5
	} else {
		score++

	}

	if wallCollide() || bitesItself() {
		isGameOver = true
	}
}

func UpdateAppple() {
	rand.Seed(time.Now().UnixNano())
	for appleInsideSnake() {
		apple.point.row, apple.point.col =
			rand.Intn(GameFrameHeight),
			rand.Intn(GameFrameWidth)

		if rand.Intn(5) == 1 {
			apple.symbol = '■'
		} else {
			apple.symbol = AppleSymbol
		}
	}
}

/* ****************************************** */
// DRAWING THE GAMESTATES //
/* ****************************************** */

func clearScreen() {
	for _, p := range pointsToClear {
		drawInsideRectArea(p.row, p.col, 1, 1, ' ')
	}

	pointsToClear = []*Point{}

}

func drawState() {
	if isGamePaused {
		return
	}

	clearScreen()

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
}

func drawSnake() {
	for _, b := range snake.body {
		drawInsideRectArea(b.row, b.col, 1, 1, snake.symbol)

		pointsToClear = append(pointsToClear, b)
	}
}

func drawApple() {
	drawInsideRectArea(apple.point.row, apple.point.col, 1, 1, apple.symbol)
	pointsToClear = append(pointsToClear, apple.point)

}

func drawRectPerimeter(row, col, width, height int, ch rune) {

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
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused

	} else if snake.velRow != 1 && key == "Rune[w]" || snake.velRow != 1 && key == "Up" {
		snake.velRow = -1
		snake.velCol = 0

	} else if snake.velCol != 1 && key == "Rune[a]" || snake.velCol != 1 && key == "Left" {
		snake.velRow = 0
		snake.velCol = -1

	} else if snake.velRow != -1 && key == "Rune[s]" || snake.velRow != -1 && key == "Down" {
		snake.velRow = 1
		snake.velCol = 0

	} else if snake.velCol != -1 && key == "Rune[d]" || snake.velCol != -1 && key == "Right" {
		snake.velRow = 0
		snake.velCol = 1
	}
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

func appleInsideSnake() bool {
	for _, b := range snake.body {
		if b.row == apple.point.row && b.col == apple.point.col {
			return true
		}
	}
	return false
}

func wallCollide() bool {
	head := getSnakeHead()
	return head.row < 0 ||
		head.row >= GameFrameHeight ||
		head.col < 0 ||
		head.col >= GameFrameWidth
}

func bitesItself() bool {
	head := getSnakeHead()
	for _, b := range snake.body[:snakeHeadIndex()] {
		if b.row == head.row && b.col == head.col {
			return true
		}
	}
	return false
}

func snakeHeadIndex() int {
	return len(snake.body) - 1
}

func getSnakeHead() *Point {
	return snake.body[snakeHeadIndex()]
}

/* ****************************************** */
// END OF GAME //
/* ****************************************** */

func gameEnd() {
	screenWidth, screenHeight := screen.Size()
	pringStringCentered(screenHeight/2-3, screenWidth/2-1, "Game Over!")
	pringStringCentered(screenHeight/2-2, screenWidth/2-2, fmt.Sprintf("Your score is %d", score))

	screen.Show()

	time.Sleep(3 * time.Second)
	screen.Fini()
	os.Exit(0)
}

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Item struct {
	x int
	y int
}

type Map struct {
	height     int
	width      int
	pacman     Item
	ghosts     []Item
	walls      [][]bool
	coins      [][]bool
	totalCoins int
	direction  int
}

const UP = int(tcell.KeyUp)
const DOWN = int(tcell.KeyDown)
const LEFT = int(tcell.KeyLeft)
const RIGHT = int(tcell.KeyRight)

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var wallStyle = tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorWhite)
var ghostStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
var pacmanStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
var deathStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
var coinStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBeige)
var winStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlue)

var totalGhosts = 5
var totalPoints = 0
var allData = Map{}
var endgame = false
var screen tcell.Screen = nil

func main() {

	var wg sync.WaitGroup

	log_file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(log_file)

	argc := len(os.Args)

	if argc == 1 {
		totalGhosts = 5
	} else if argc == 2 {
		ghsts, err := strconv.Atoi(os.Args[1])
		if totalGhosts < 0 || err != nil {
			fmt.Printf("Please input total number of ghosts ./program 5\n")
			return
		}
		totalGhosts = ghsts
	} else {
		fmt.Printf("Please input total number of ghosts ./program 5\n")
		return
	}

	file_map, err := os.ReadFile("./maps/map.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	var coinsCount = 0
	var wallsCount = 0
	var ghostSpaces = 0
	var width = 0
	var height = 0

	for _, element := range file_map {
		switch element {
		case 'W':
			wallsCount++
			break
		case '*':
			coinsCount++
			break
		case 'g':
			ghostSpaces++
			break
		case '\n':
		case '\r':
			height++
			continue
		}
		if height == 0 {
			width++
		}
	}
	height++ // EOF last line

	limitedGhosts := 0
	if totalGhosts < ghostSpaces {
		limitedGhosts = totalGhosts
	} else {
		limitedGhosts = ghostSpaces
	}

	walls := make([][]bool, width)
	for i := range walls {
		walls[i] = make([]bool, height)
	}
	coins := make([][]bool, width)
	for i := range coins {
		coins[i] = make([]bool, height)
	}
	ghosts := make([]Item, limitedGhosts)
	var pacman Item
	coinsCount = 0
	ghostsCount := 0
	line := 0
	col := 0

	for _, element := range file_map {
		switch element {
		case 'W':
			walls[col][line] = true
			break
		case '*':
			coins[col][line] = true
			coinsCount++
			break
		case 'g':
			if ghostsCount < limitedGhosts {
				ghosts[ghostsCount].x = col
				ghosts[ghostsCount].y = line
				ghostsCount++
			}
			break
		case 'C':
			pacman.x = col
			pacman.y = line
			break
		case '\n':
			line++
			col = -1
			break
		}
		col++
	}

	allData.height = height
	allData.width = width
	allData.coins = coins
	allData.walls = walls
	allData.ghosts = ghosts
	allData.pacman = pacman
	allData.totalCoins = coinsCount

	screen = initUI()

	wg.Add(1)
	go manageUI(&wg)

	wg.Add(1)
	go readInput(&wg)

	i := 0
	for i < limitedGhosts {
		wg.Add(1)
		go controlGhost(&wg, i)
		i++
	}
	log.Println("Waiting...")
	wg.Wait()
	log.Println("Threads ended")
	exit()

}

func readInput(wg *sync.WaitGroup) {
	defer wg.Done()
	for !endgame {
		// Poll event
		ev := screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				endgame = true
			} else if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyDown || ev.Key() == tcell.KeyLeft || ev.Key() == tcell.KeyRight {
				allData.direction = int(ev.Key())
			}
		}
	}
	log.Println("Ended Input")
}

func controlGhost(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	time.Sleep(time.Duration((rand.Int() % (2 * (id + 1)))) * time.Second)

	currentDirection := UP

	for !endgame {
		prevX := allData.ghosts[id].x
		prevY := allData.ghosts[id].y

		var possiblePath []int
		if !allData.walls[prevX][prevY-1] && !(currentDirection == DOWN) {
			possiblePath = append(possiblePath, UP)
		}
		if !allData.walls[prevX][prevY+1] && !(currentDirection == UP) {
			possiblePath = append(possiblePath, DOWN)
		}
		if !allData.walls[prevX+1][prevY] && !(currentDirection == LEFT) {
			possiblePath = append(possiblePath, RIGHT)
		}
		if !allData.walls[prevX-1][prevY] && !(currentDirection == RIGHT) {
			possiblePath = append(possiblePath, LEFT)
		}

		if currentDirection == UP {
			allData.ghosts[id].y--
		} else if currentDirection == DOWN {
			allData.ghosts[id].y++
		} else if currentDirection == LEFT {
			allData.ghosts[id].x--
		} else if currentDirection == RIGHT {
			allData.ghosts[id].x++
		}

		if allData.walls[allData.ghosts[id].x][allData.ghosts[id].y] {
			allData.ghosts[id].x = prevX
			allData.ghosts[id].y = prevY
		}

		currentDirection = possiblePath[(rand.Int() % len(possiblePath))]
		sleepCycle()
		if allData.ghosts[id].x == allData.pacman.x && allData.ghosts[id].y == allData.pacman.y {
			endgame = true
		}
	}
	log.Println("Ended ghosts " + strconv.Itoa(id))

}

func initUI() tcell.Screen {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()
	return s
}

func manageUI(wg *sync.WaitGroup) {
	defer wg.Done()
	for !endgame {

		prevX := allData.pacman.x
		prevY := allData.pacman.y
		if allData.direction == UP {
			allData.pacman.y--
		} else if allData.direction == DOWN {
			allData.pacman.y++
		} else if allData.direction == LEFT {
			allData.pacman.x--
		} else if allData.direction == RIGHT {
			allData.pacman.x++
		}

		if allData.walls[allData.pacman.x][allData.pacman.y] {
			allData.pacman.x = prevX
			allData.pacman.y = prevY
		}

		if allData.coins[allData.pacman.x][allData.pacman.y] {
			totalPoints++
			if totalPoints >= allData.totalCoins {
				endgame = true
			}
			allData.coins[allData.pacman.x][allData.pacman.y] = false
		}

		draw()

	}
	draw()
	log.Println("Ended UI")

}

func draw() {

	for i := range allData.walls {
		for j := range allData.walls[i] {
			if allData.walls[i][j] {
				screen.SetContent(i, j+1, ' ', nil, wallStyle)
			}
		}
	}

	for i := range allData.coins {
		for j := range allData.coins[i] {
			if allData.coins[i][j] {
				screen.SetContent(i, j+1, '*', nil, coinStyle)
			}
		}
	}

	for _, ghost := range allData.ghosts {
		screen.SetContent(ghost.x, ghost.y+1, 'g', nil, ghostStyle)
	}

	screen.SetContent(allData.pacman.x, allData.pacman.y+1, 'C', nil, pacmanStyle)

	if endgame {
		screen.SetContent(allData.pacman.x, allData.pacman.y+1, 'X', nil, deathStyle)
	}
	screen.SetContent(0, 0, ' ', nil, defStyle)
	drawText(screen, 0, 0, 1000, 1000, defStyle, "Points: "+strconv.Itoa(totalPoints))

	screen.Show()
	sleepCycle()
	screen.Clear()
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func exit() {
	time.Sleep(2 * time.Second)
	screen.Clear()

	drawText(screen, 0, 0, 1000, 1000, defStyle, "Points: "+strconv.Itoa(totalPoints))
	drawText(screen, 0, 1, 1000, 1000, defStyle, "END GAME")
	if totalPoints >= allData.totalCoins {
		drawText(screen, 0, 2, 1000, 1000, winStyle, "You WON")
	}
	screen.Show()
	time.Sleep(5 * time.Second)
	screen.Fini()
	os.Exit(0)
}

func sleepCycle() {
	time.Sleep(50 * time.Millisecond)
}

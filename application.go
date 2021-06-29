package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// const game_width int = 40 / 2
// const game_height int = 40 / 2
const GAME_TOP_OFFSET = 1

type Tile struct {
	x          int
	y          int
	value      int
	isEmpty    bool
	mergedFrom []Tile
}

func (t *Tile) updatePosition(pos *Tile) {
	t.x = pos.x
	t.y = pos.y
}

func drawMessage(x, y int, str string, color termbox.Attribute) {
	backgroundColor := termbox.ColorDefault
	runes := []rune(str)

	for i := 0; i < len(runes); i += 1 {
		termbox.SetCell(x+i, y, runes[i], color, backgroundColor)
	}
}

func drawLine(x, y int, str string) {
	drawMessage(x, y, str, termbox.ColorDefault)
}

func defaultColorFill(x, y, w, h int, cell termbox.Cell) {
	fill(x, y, w, h, cell, termbox.ColorDefault)
}

func fill(x, y, w, h int, cell termbox.Cell, color termbox.Attribute) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, color, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func drawCell(tile Tile, left int, top int, cellWidth int, cellHeight int) {
	const coldef = termbox.ColorDefault

	defaultColorFill(left, top, cellWidth, 1, termbox.Cell{Ch: '─'})
	defaultColorFill(left, top, 1, cellHeight, termbox.Cell{Ch: '|'})
	defaultColorFill(left, top+cellHeight, cellWidth, 1, termbox.Cell{Ch: '─'})
	defaultColorFill(left+cellWidth, top, 1, cellHeight, termbox.Cell{Ch: '|'})

	if !tile.isEmpty {
		tbprint(left+cellWidth/2, top+cellHeight/2, getColor(tile.value), coldef, strconv.Itoa(tile.value))
	}
}

func handleKeyEvent() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeySpace:
				return
			case termbox.KeyArrowDown:
				dispatch("down", &Message{data: "push left"})
				break
			case termbox.KeyArrowLeft:
				dispatch("left", &Message{data: "push left"})
				break
			case termbox.KeyArrowRight:
				dispatch("right", &Message{data: "push left"})
				break
			case termbox.KeyArrowUp:
				dispatch("up", &Message{data: "push left"})
				break
			default:
			}
		default:
		}
	}
}

type Drawer struct{}

func (d *Drawer) redraw(grid *Grid, score int, highScore int, isOver bool) {
	if debugRun {
		dumpCell(grid, score, highScore, isOver)
	} else {
		gridDraw(grid, score, highScore, isOver)
	}

	var info GameInfo

	if isOver {
		info = GameInfo{HighScore: score}
	} else {
		info = GameInfo{HighScore: score, CurrentScore: score, TileState: tileToPrimitive(grid.cells)}
	}

	var f *os.File
	var err error

	if f, err = os.Create(getFilePath()); err != nil {
		panic(err)
	}

	if err = info.save(f); err != nil {
		panic(err)
	}
}

func tileToPrimitive(t [][]Tile) [][][]int {
	p := [][][]int{}
	for ly := 0; ly < len(t); ly++ {
		r := [][]int{}
		for lx := 0; lx < len(t[ly]); lx++ {
			tt := t[ly][lx]
			var v int

			if tt.isEmpty {
				v = 0
			} else {
				v = tt.value
			}
			r = append(r, []int{tt.x, tt.y, v})
		}
		p = append(p, r)
	}

	return p
}

func gridDraw(grid *Grid, score int, highScore int, isOver bool) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	drawCellNumber(grid)

	//draw score
	drawScore(score)
	drawHighScore(highScore)

	if isOver {
		drawOver()
	}

	termbox.Flush()
}

func drawOver() {
	gameOver := "Game Over"
	gameOverLen := len(gameOver)

	lastMessage := "If you quit it, please press ESC or SPACE"

	top := game_height / 2
	width := (game_width - gameOverLen) / 2
	color := termbox.ColorRed

	fill(0, top, width, 1, termbox.Cell{Ch: '='}, color)
	drawMessage(width, top, gameOver, color)
	fill(width+gameOverLen, top, width, 1, termbox.Cell{Ch: '='}, color)

	drawMessage((game_width-len(lastMessage))/2, top+1, lastMessage, color)
}

func drawScore(score int) {
	drawLine(0, 0, fmt.Sprintf("Score: %d", score))
}

func drawHighScore(score int) {
	highScoreMsg := fmt.Sprintf("High Score: %d", score)
	drawLine(game_width-len(highScoreMsg), 0, highScoreMsg)
}

func drawCellNumber(grid *Grid) {
	cellWidth := game_width / grid.size
	cellHeight := game_height / grid.size

	for ly := 0; ly < grid.size; ly++ {
		for lx := 0; lx < grid.size; lx++ {
			tile := grid.cells[lx][ly]
			drawCell(tile, lx*cellWidth, GAME_TOP_OFFSET+ly*cellHeight, cellWidth, cellHeight)
		}
	}

}

func dumpCell(grid *Grid, score int, highScore int, isOver bool) {
	sumValue := 0
	countIsNotEmpty := 0
	for ly := 0; ly < grid.size; ly++ {
		for lx := 0; lx < grid.size; lx++ {
			if !grid.cells[lx][ly].isEmpty {
				sumValue += grid.cells[lx][ly].value
				fmt.Println("==========================================", grid.cells[ly][lx].x, grid.cells[ly][lx].y, grid.cells[ly][lx].value)
				countIsNotEmpty += 1
			}
		}
	}
	fmt.Println("==================isOver================", isOver)
	fmt.Println("==================sumValue================", sumValue)
	fmt.Println("================countIsNotEmpty===========", (16 - countIsNotEmpty))
}

type GameInfo struct {
	HighScore    int
	CurrentScore int
	TileState    [][][]int
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (g *GameInfo) load() error {
	dir := getDirPath()
	file := getFilePath()

	var err error
	var f *os.File

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if fileExists(file) {
		if _, err := toml.DecodeFile(file, g); err != nil {
			return err
		}
		return nil
	}

	if f, err = os.Create(file); err != nil {
		return err
	}

	g.HighScore = 0
	g.CurrentScore = 0

	if err != g.save(f) {
		return err
	}

	return nil
}

func (g *GameInfo) getTiles() [][]Tile {
	res := [][]Tile{}

	for ly := 0; ly < len(g.TileState); ly++ {
		r := []Tile{}
		for lx := 0; lx < len(g.TileState[ly]); lx++ {
			v := g.TileState[ly][lx]
			var e bool
			if v[2] == 0 {
				e = true
			} else {
				e = false
			}
			r = append(r, Tile{x: v[0], y: v[1], value: v[2], isEmpty: e})
		}
		res = append(res, r)
	}

	return res
}

func getDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		
		exe, err := os.Executable()
		
		if err != nil {
			panic("Couldn't find user's home directory or the executable")
		}
		
		home = filepath.Dir(exe)
	}
	return filepath.Join(home, ".config", "2048")
}

func getFilePath() string {
	return filepath.Join(getDirPath(), "game.toml")
}

func (g *GameInfo) save(f *os.File) error {
	if err := toml.NewEncoder(f).Encode(g); err != nil {
		return err
	}

	return nil
}

func controlFromCommand() error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		switch scanner.Text() {
		case "d":
			dispatch("down", &Message{data: "push left"})
			break
		case "l":
			dispatch("left", &Message{data: "push left"})
			break
		case "r":
			dispatch("right", &Message{data: "push left"})
			break
		case "u":
			dispatch("up", &Message{data: "push left"})
			break
		default:
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var (
	debugRun bool
	rowLength int
	tileSize int
	game_width int
	game_height int
)

func main() {
	flag.BoolVar(&debugRun, "debug", false, "debugRun flag")
	flag.IntVar(&rowLength, "rowLength", 4, "set the number of tiles per row")
	flag.IntVar(&tileSize, "tileSize", 8, "set the size of the square tiles")
	flag.Parse()
	
	game_width = tileSize * 10
	game_height = game_width / 2
	
	var gInfo GameInfo
	err := gInfo.load()

	drawer := Drawer{}
	gameState := Game{gridSize: rowLength, drawer: &drawer}
	gameState.setup(gInfo)

	if debugRun {
		err = controlFromCommand()
	} else {
		err = termbox.Init()

		gridDraw(gameState.grid, gameState.score, gameState.highScore, false)

		defer termbox.Close()

		handleKeyEvent()
	}

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

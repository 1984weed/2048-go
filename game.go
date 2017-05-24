package main

import (
	"math/rand"
	"fmt"
)

type Game struct {
	gridSize int
	startTiles int
	score int
	over bool
	won bool
	grid *Grid
	drawer *Drawer
}

type Vector struct{
	x int
	y int
}

type PositionTraversal struct{
	x []int
	y []int
}

func (g *Game)setup(){
	g.score = 0
	g.startTiles = 2
	g.grid = &Grid{size: g.gridSize}
	g.grid.setup()
	g.addStartTiles()

	add("up", func(message *Message){
		fmt.Println("up")
		g.move(0)
	})
	add("right", func(message *Message){
		fmt.Println("right")
		g.move(1)
	})
	add("down", func(message *Message){
		fmt.Println("down")
		g.move(2)
	})
	add("left", func(message *Message){
		fmt.Println("left")
		g.move(3)
	})
}

func (g *Game) addStartTiles(){
	defaultTiles := [2]Tile{Tile{x: 0, y: 0, value: 2, isEmpty: false}, Tile{x: 0, y: 3, value: 2, isEmpty: false}}
	for i := 0; i < g.startTiles; i++{
		//g.addRandomTile()
		g.grid.insertTile(defaultTiles[i])
	}
}

func (g *Game) addRandomTile(){
	if g.grid.cellsAvailable() {
		value := 2
		if rand.Float32() < 0.9 {
			value = 4
		}
		tile := g.grid.randomAvailableCell()
		newTile := Tile{x: tile.x, y: tile.y, value: value, isEmpty: false}

		g.grid.insertTile(newTile)
	}
}

func (g *Game) GetVector(direction int) Vector{
	res := make(map[int]Vector)

	res[0] = Vector{x: 0, y: -1} // Up
	res[1] = Vector{x: 1, y: 0} // Right
	res[2] = Vector{x: 0, y: 1} // Down
	res[3] = Vector{x: -1, y: 0} // Left

	return res[direction]
}

/**
ここを変えたい。。。。updatePositionをどうするか
 */
func (g *Game) moveTile(tile *Tile, farPos *Tile) Tile{
	g.grid.removeTile(tile)
	g.grid.cells[farPos.x][farPos.y] = Tile{x: farPos.x, y: farPos.y, value: tile.value, mergedFrom: tile.mergedFrom, isEmpty: false}
	return g.grid.cells[farPos.x][farPos.y]
}

func (g *Game) IsGameTerminated() bool{
	return false
}

func (g *Game) BuildTraversals(vec Vector) PositionTraversal{
	traversals := PositionTraversal{x: make([]int, g.gridSize), y: make([]int, g.gridSize)}

	for i := 0; i < g.gridSize; i++ {
		traversals.x[i] = i
		traversals.y[i] = i
	}

	if(vec.x == 1){
		traversals.x = ReverseList(traversals.x)
	}

	if(vec.y == 1){
		traversals.y = ReverseList(traversals.y)
	}

	return traversals
}

func (g *Game) FindFarthestPosition(cell Tile, vector Vector) (*Tile, *Tile){
	previous := cell
	isFirst := true

	for isFirst || (g.grid.WithinBounds(&cell) && g.grid.CellAvailable(&cell)) {
		previous = cell
		cell = Tile{x: previous.x + vector.x, y: previous.y + vector.y}

		isFirst = false
	}

	return &previous, &cell
}

func (g *Game) positionsEqual(first *Tile, second *Tile) bool {
	return first.x == second.x && first.y == second.y
}

//Todo, I have to implement later...
func (g *Game) tileMatchesAvailable() bool{
	return true
}

func (g *Game) movesAvailable() bool {
	return g.grid.cellsAvailable()// || g.tileMatchesAvailable()
}

func (g *Game) move(direction int){


	if(g.IsGameTerminated()){
		return
	}

	moved := false
	vector := g.GetVector(direction)
	traversals := g.BuildTraversals(vector)

	for i := 0; i < len(traversals.x); i++ {
		for j := 0; j < len(traversals.y); j++{
			cell := Tile{x: traversals.x[i], y: traversals.y[j]}
			tile := g.grid.CellContent(&cell)
			if(tile != nil && !tile.isEmpty){
				farPos, nextPos := g.FindFarthestPosition(cell, vector)
				next := g.grid.CellContent(nextPos)
				if( next != nil && next.value == tile.value /*&& !next.mergedFrom*/){
					merged := Tile{x: nextPos.x, y: nextPos.y, value: tile.value * 2}
					tiles := make([]Tile, 2)
					tiles[0] = g.copyTile(tile)
					tiles[1] = g.copyTile(nextPos)
					merged.mergedFrom = tiles

					g.grid.insertTile(merged)
					g.grid.removeTile(tile)

					temp := g.copyTile(tile)
					tile = &temp
					tile.updatePosition(nextPos)

					g.score += merged.value
				} else {
					temp := g.moveTile(tile, farPos)
					tile = &temp
				}

				if(!g.positionsEqual(&cell, tile)){
					moved = true
				}
			}
		}
	}
	if !g.movesAvailable() {
		fmt.Println("game is over")
	}

	if moved {
		g.addRandomTile()
		g.actuate()
	}
}

func (g *Game) copyTile(tile *Tile) Tile{
	return Tile{x: tile.x, y: tile.y, value: tile.value, isEmpty: tile.isEmpty}
}

func (g *Game) actuate(){
	g.drawer.redraw(g.grid, g.score)
}

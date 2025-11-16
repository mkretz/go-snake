package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your Battlesnake username
		Color:      "#888888", // TODO: Choose color
		Head:       "default", // TODO: Choose head
		Tail:       "default", // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

func contains(coords []Coord, coord Coord) int {
	for index, c := range coords {
		if c.X == coord.X && c.Y == coord.Y {
			return index
		}
	}
	return -1
}

func snakeCoords(snake Battlesnake) []Coord {
	return append(snake.Body, snake.Head)
}

func distance(from Coord, to Coord) int {
	return int(math.Abs(float64(to.X-from.X))) + int(math.Abs(float64(to.Y-from.Y)))
}

func distanceToClosestFood(from Coord, gameState GameState) int {
	result := gameState.Board.Width + gameState.Board.Height
	for _, coord := range gameState.Board.Food {
		if distance(from, coord) < result {
			result = distance(from, coord)
		}
	}
	return result
}

func resultingCoord(battlesnake Battlesnake, move string) Coord {
	switch move {
	case "right":
		return Coord{X: battlesnake.Head.X + 1, Y: battlesnake.Head.Y}
	case "left":
		return Coord{X: battlesnake.Head.X - 1, Y: battlesnake.Head.Y}
	case "up":
		return Coord{X: battlesnake.Head.X, Y: battlesnake.Head.Y + 1}
	case "down":
		return Coord{X: battlesnake.Head.X, Y: battlesnake.Head.Y - 1}
	default:
		return battlesnake.Head
	}
}

func foodDistances(snake Battlesnake, gameState GameState) map[string]int {
	foodDistances := map[string]int{
		"right": gameState.Board.Width + gameState.Board.Height,
		"left":  gameState.Board.Width + gameState.Board.Height,
		"up":    gameState.Board.Width + gameState.Board.Height,
		"down":  gameState.Board.Width + gameState.Board.Height,
	}

	// search for food to the right
	for i := snake.Head.X; i < gameState.Board.Width; i++ {
		coord := Coord{X: i, Y: snake.Head.Y}
		if contains(gameState.Board.Food, coord) >= 0 {
			foodDistances["right"] = i
			break
		}
	}

	// search for food to the left
	for i := snake.Head.X; i >= 0; i-- {
		coord := Coord{X: i, Y: snake.Head.Y}
		if contains(gameState.Board.Food, coord) >= 0 {
			foodDistances["left"] = snake.Head.X - i
			break
		}
	}

	// search for food to the top
	for i := snake.Head.Y; i < gameState.Board.Height; i++ {
		coord := Coord{X: snake.Head.X, Y: i}
		if contains(gameState.Board.Food, coord) >= 0 {
			foodDistances["up"] = i
			break
		}
	}

	// search for food to the bottom
	for i := snake.Head.Y; i >= 0; i-- {
		coord := Coord{X: snake.Head.X, Y: i}
		if contains(gameState.Board.Food, coord) >= 0 {
			foodDistances["down"] = snake.Head.Y - i
			break
		}
	}
	return foodDistances
}

func bestDirectionForFood(snake Battlesnake, gameState GameState, safeMoves []string) string {
	result := "down"
	minDistance := gameState.Board.Width + gameState.Board.Height
	for _, move := range safeMoves {
		foodDistance := distanceToClosestFood(resultingCoord(snake, move), gameState)
		if foodDistance < minDistance {
			minDistance = foodDistance
			result = move
		}
	}
	return result
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {
	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false
	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false
	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false
	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	if myHead.X == 0 { // at the left edge
		isMoveSafe["left"] = false
	}
	if myHead.X == state.Board.Width-1 { // at the right edge
		isMoveSafe["right"] = false
	}
	if myHead.Y == 0 { // at bottom edge
		isMoveSafe["down"] = false
	}
	if myHead.Y == state.Board.Height-1 { // at top edge
		isMoveSafe["up"] = false
	}

	// do not crash into yourself
	if contains(state.You.Body, Coord{X: myHead.X + 1, Y: myHead.Y}) >= 0 {
		isMoveSafe["right"] = false
	}
	if contains(state.You.Body, Coord{X: myHead.X - 1, Y: myHead.Y}) >= 0 {
		isMoveSafe["left"] = false
	}
	if contains(state.You.Body, Coord{X: myHead.X, Y: myHead.Y + 1}) >= 0 {
		isMoveSafe["up"] = false
	}
	if contains(state.You.Body, Coord{X: myHead.X, Y: myHead.Y - 1}) >= 0 {
		isMoveSafe["down"] = false
	}

	// do not crash into the other snakes
	for _, snake := range state.Board.Snakes {
		if contains(snakeCoords(snake), Coord{X: myHead.X + 1, Y: myHead.Y}) >= 0 {
			isMoveSafe["right"] = false
		}
		if contains(snakeCoords(snake), Coord{X: myHead.X - 1, Y: myHead.Y}) >= 0 {
			isMoveSafe["left"] = false
		}
		if contains(snakeCoords(snake), Coord{X: myHead.X, Y: myHead.Y + 1}) >= 0 {
			isMoveSafe["up"] = false
		}
		if contains(snakeCoords(snake), Coord{X: myHead.X, Y: myHead.Y - 1}) >= 0 {
			isMoveSafe["down"] = false
		}
	}
	//
	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	nextMove := bestDirectionForFood(state.You, state, safeMoves)

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func main() {
	RunServer()
}

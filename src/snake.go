package main

import (
	"honnef.co/go/js/dom"
	"time"
	"math/rand"
	"strconv"
	"fmt"
)

const colNumber = 10
const rowNumber = 10
const snakeInitialLength = 3
var goDir string
type position struct {
	row	int
	col	int
}
var snakeBody []position // This is a QUEUE
const speed = time.Millisecond * 900
var foodPosition position


func randomDirection(random *rand.Rand) string {
	dirNum := random.Intn(4)
	if dirNum == 0 { return "←" }
	if dirNum == 1 { return "↑" }
	if dirNum == 2 { return "→" }
	return "↓"
}


func isFoodCollideWithSnakeBody(fpos position, sb []position) bool {
	for _, pos := range sb {
		if fpos.row == pos.row && fpos.col == pos.col { return true }
	}
	return false
}


func randomFood(sb []position) position {
	// input: snakeBody
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var pos position
	for {
		// food position: row
		pos.row = r.Intn(rowNumber)
		// food position: column
		pos.col = r.Intn(colNumber)
		if !isFoodCollideWithSnakeBody(pos, sb) { return pos }
	}
}


func nextPosition(direction string, pos position) position {
	if direction == "←" {
		newCol := (pos.col + colNumber - 1) % colNumber
		return position{row: pos.row, col: newCol}
	}

	if direction == "↑" {
		newRow := (pos.row + rowNumber - 1) % rowNumber
		return position{row: newRow, col: pos.col}
	}

	if direction == "→" {
		newCol := (pos.col + 1) % colNumber;
		return position{row: pos.row, col: newCol}
	}

	if direction == "↓" {
		newRow := (pos.row + 1) % rowNumber;
		return position{row: newRow, col: pos.col}
	}

	panic("End of func nextPosition")
}


func startGame() {
	d := dom.GetWindow().Document()
	btn := d.GetElementByID("start").(*dom.HTMLButtonElement)
	btn.Style().SetProperty("display", "none", "")

	main := d.GetElementByID("main").(*dom.HTMLDivElement)

	// plot grids
	for i:=0; i<rowNumber; i++ {
		for j:=0; j<colNumber; j++ {
			elm := d.CreateElement("div").(*dom.HTMLDivElement)
			elm.Class().Add("square")
			// not working, why???
			//elm.SetAttribute("data-row", string(i))
			//elm.SetAttribute("data-col", string(j))
			elm.SetAttribute("data-row", strconv.Itoa(i))
			elm.SetAttribute("data-col", strconv.Itoa(j))
			elm.SetTextContent(" ")
			main.AppendChild(elm)
		}
		br := d.CreateElement("br")
		main.AppendChild(br)
	}

	// generate random snake
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var snakeHeadPosition = position{
		row: r.Intn(rowNumber),
		col: r.Intn(colNumber),
	}

	// snake initial go direction
	goDir = randomDirection(r);

	sel := fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
		snakeHeadPosition.row, snakeHeadPosition.col)
	main.QuerySelector(sel).Class().Add("select")

	// add snake head position to snakeBody(Queue)
	snakeBody = append(snakeBody, snakeHeadPosition)

	// snake body direction
	snakeBodyDir := randomDirection(r)
	for snakeBodyDir == goDir {
		snakeBodyDir = randomDirection(r)
	}

	// generate snake body
	pos := snakeBody[0]
	for i:=0; i<(snakeInitialLength-1); i++ {
		pos = nextPosition(snakeBodyDir, pos)
		snakeBody = append(snakeBody, pos)
		sel2 := fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
			pos.row, pos.col)
		main.QuerySelector(sel2).Class().Add("select")
	}

	// generate food
	foodPosition = randomFood(snakeBody)
	sel3 := fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
		foodPosition.row, foodPosition.col)
	main.QuerySelector(sel3).Class().Add("food")

	// start to move snake
	ticker := time.NewTicker(speed)
	quit := make(chan struct{})
	go func() {
		for {
			select {
				case <- ticker.C:
					snakeMove(quit, d)
				case <- quit:
					ticker.Stop()
					return
			}
		}
	}()
}


func snakeMove(quit chan struct{}, d dom.Document) {
	// move snake: get next position
	npos := nextPosition(goDir, snakeBody[0])

	// check if snake eat itself
	for _, pos := range snakeBody {
		if pos.row == npos.row && pos.col == npos.col {
			close(quit)
			d.GetElementByID("info").SetTextContent("Game Over!")
			return
		}
	}

	// add next position to snake body
	// FIXME: how to prepend npos to snakeBody efficiently???
	var hpos []position
	hpos = append(hpos, npos)
	for idx, _ := range snakeBody {
		hpos = append(hpos, snakeBody[idx])
	}
	snakeBody = hpos
	sel := fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
		npos.row, npos.col)
	d.QuerySelector(sel).Class().Add("select")
	if foodPosition.row == npos.row && foodPosition.col == npos.col {
		// snake eat food
		d.QuerySelector(sel).Class().Remove("food")
		// generate food
		foodPosition = randomFood(snakeBody)
		sel = fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
			foodPosition.row, foodPosition.col)
		d.QuerySelector(sel).Class().Add("food")
	} else {
		// snake does not eat food
		// move snake: remove last position
		pos := snakeBody[len(snakeBody)-1]
		snakeBody = snakeBody[:len(snakeBody)-1]
		sel = fmt.Sprintf("div[data-row=\"%d\"][data-col=\"%d\"]",
			pos.row, pos.col)
		d.QuerySelector(sel).Class().Remove("select")
	}
}


func handleArrowKey(event dom.Event) {
	ke := event.(*dom.KeyboardEvent)

	switch ke.KeyCode {
	case 37: //left
		if goDir != "→" { goDir = "←" }
	case 38: //up
		if goDir != "↓" { goDir = "↑" }
	case 39: //right
		if goDir != "←" { goDir = "→" }
	case 40: //down
		if goDir != "↑" { goDir = "↓" }
	default:
	}
}

func main() {
	startGame()

	dom.GetWindow().AddEventListener("keydown", false, handleArrowKey)
}

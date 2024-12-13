package main

import (
	"fmt"

	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/platform"
	"github.com/bennicholls/tyumi/platform/platform_sdl"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func main() {
	engine.InitConsole(40, 20)
	platform.Set(platform_sdl.New())
	err := engine.SetupRenderer("res/curses24x24.bmp", "res/font12x24.bmp", "TEST WINDOW")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	state := TestState{}
	state.Setup()

	engine.InitMainState(&state)

	engine.Run()
}

type TestState struct {
	engine.StatePrototype

	text ui.Textbox
	list ui.List

	tick int
}

func (ts *TestState) Setup() {
	ts.Init(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	ts.Window().SetDefaultColours(col.Pair{col.RED, col.LIME})
	ts.Window().SetBorderStyle(ui.BORDER_STYLE_CUSTOM, ui.BorderStyles["Thick"])

	container := ui.ElementPrototype{}
	container.Init(20, 15, vec.Coord{1, 1}, 1)
	container.EnableBorder("Container", "")

	ts.text = ui.NewTextbox(ui.FIT_TEXT, ui.FIT_TEXT, vec.Coord{1, 1}, 0, "TEST STRING DO NOT UPVOTE", true)
	ts.text.SetDefaultColours(col.Pair{col.CYAN, col.FUSCHIA})
	ts.text.EnableBorder("TEST TITLE", "hint")
	ts.text.SetBorderStyle(ui.BORDER_STYLE_INHERIT)

	text2 := ui.NewTextbox(10, ui.FIT_TEXT, vec.Coord{7, 3}, 0, util.LoremIpsum(31), true)
	text2.EnableBorder("", "")
	container.AddChildren(&ts.text, &text2)
	ts.AddInputHandler(ts.HandleInputs)

	ts.Window().AddChild(&container)

	inputbox := ui.NewInputbox(10, 1, vec.Coord{25, 18}, 10)
	inputbox.EnableBorder("inputs!", "do the input")
	inputbox.SetBorderStyle(ui.BORDER_STYLE_CUSTOM, ui.BorderStyles["Block"])
	ts.Window().AddChild(&inputbox)

	ts.list = ui.NewList(15, 10, vec.Coord{8, 8}, 1)
	for i := 0; i < 20; i++ {
		item := ui.NewTextbox(15, i%3+1, vec.ZERO_COORD, 1, "List item "+fmt.Sprint(i)+"/n", false)
		ts.list.AddChild(&item)
	}

	ts.list.EnableBorder("LIST", "")
	custom_border := ui.BorderStyles["Thin"]
	custom_border.TextDecorationR = 'A'
	custom_border.DisableLink = true
	ts.list.SetBorderStyle(ui.BORDER_STYLE_CUSTOM, custom_border)
	ts.list.ToggleHighlight()
	ts.list.ToggleScrollbar()
	ts.list.SetDefaultColours(col.Pair{col.BLUE, col.WHITE})

	ts.Window().AddChild(&ts.list)
}

func (ts *TestState) Update() {
	ts.tick++
}

func (ts *TestState) UpdateUI() {
	return
}

func (ts *TestState) HandleInputs(e event.Event) {
	switch e.ID() {
	case input.EV_KEYBOARD:
		ev := e.(input.KeyboardEvent)
		switch ev.Key {
		case input.K_a:
			// item := ui.NewTextbox(15, 1, vec.ZERO_COORD, 1, "new item", false)
			// ts.list.AddChild(&item)
		case input.K_RIGHT:
			ts.text.Move(1, 0)
		case input.K_LEFT:
			ts.text.Move(-1, 0)
		case input.K_UP:
			ts.text.Move(0, -1)
		case input.K_DOWN:
			ts.text.Move(0, 1)
		case input.K_RETURN:
			ts.text.ToggleVisible()
		}
	}
}

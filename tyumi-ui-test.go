package main

import (
	"fmt"

	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/platform"
	"github.com/bennicholls/tyumi/platform/platform_sdl"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

func main() {
	engine.InitConsole(40, 30)
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

	tabs        *ui.PageContainer
	moving_text *ui.Textbox
	list        *ui.List

	tick int
}

func (ts *TestState) Setup() {
	ts.Init(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	ts.AddInputHandler(ts.HandleInputs)
	ts.Window().SetDefaultColours(col.Pair{col.RED, col.LIME})
	ts.Window().SetBorderStyle(ui.BORDER_STYLE_CUSTOM, ui.BorderStyles["Thick"])

	ts.tabs = ui.NewPageContainer(ts.Window().Bounds().W-2, ts.Window().Bounds().H-2, vec.Coord{1, 1}, 0)
	ts.tabs.SetupBorder("UI TESTS!!", "tab to swap")
	ts.tabs.SetDefaultColours(col.Pair{col.MAROON, col.DARKGREY})
	ts.Window().AddChild(ts.tabs)

	//
	//some basic elements with input handlers
	//
	page1 := ts.tabs.CreatePage("input")

	inputbox := ui.NewInputbox(10, 1, vec.Coord{25, 18}, 2)
	inputbox.SetupBorder("type some text", "use your hands")
	inputbox.SetBorderStyle(ui.BORDER_STYLE_CUSTOM, ui.BorderStyles["Block"])
	page1.AddChild(inputbox)

	ts.list = ui.NewList(15, 10, vec.Coord{8, 8}, 1)
	for i := 0; i < 20; i++ {
		item := ui.NewTextbox(15, i%3+1, vec.ZERO_COORD, 1, "List item "+fmt.Sprint(i)+"/n", false)
		ts.list.AddChild(item)
	}

	ts.list.SetupBorder("LIST", "UP/DOWN/INSERT")
	custom_border := ui.BorderStyles["Thin"]
	custom_border.TextDecorationR = 'A'
	custom_border.DisableLink = true
	ts.list.SetBorderStyle(ui.BORDER_STYLE_CUSTOM, custom_border)
	ts.list.ToggleHighlight()
	ts.list.ToggleScrollbar()
	ts.list.SetDefaultColours(col.Pair{col.BLUE, col.WHITE})
	page1.AddChild(ts.list)

	choices := ui.NewChoiceBox(10, 1, vec.Coord{25, 10}, 5, "choice 1", "2", "really long choice")
	choices.SetDefaultColours(col.Pair{col.GREY, col.MAROON})
	page1.AddChild(choices)

	//
	// container test. multiple elements, including one that can move with the arrow keys
	//
	page2 := ts.tabs.CreatePage("container")
	container := ui.ElementPrototype{}
	container.Init(20, 15, vec.Coord{1, 1}, 1)
	container.SetupBorder("Container", "")

	ts.moving_text = ui.NewTextbox(ui.FIT_TEXT, ui.FIT_TEXT, vec.Coord{1, 1}, 0, "TEST STRING DO NOT UPVOTE", true)
	ts.moving_text.SetDefaultColours(col.Pair{col.CYAN, col.FUSCHIA})
	ts.moving_text.SetupBorder("TEST TITLE", "hint")
	ts.moving_text.SetBorderStyle(ui.BORDER_STYLE_INHERIT)

	text2 := ui.NewTextbox(10, ui.FIT_TEXT, vec.Coord{7, 3}, 0, util.LoremIpsum(31), true)
	text2.SetupBorder("", "")
	container.AddChildren(ts.moving_text, text2)
	page2.AddChild(&container)
}

func (ts *TestState) Update() {
	ts.tick++
}

func (ts *TestState) UpdateUI() {
	return
}

func (ts *TestState) GiveUserSeizure() {
	flash := gfx.NewFlashAnimation(ts.Window().Bounds(), 0, col.Pair{col.WHITE, col.WHITE}, 10)
	flash.OneShot = true
	flash.Start()
	ts.Window().AddAnimation(flash)
}

func (ts *TestState) HandleInputs(e event.Event) {
	switch e.ID() {
	case input.EV_KEYBOARD:
		ev := e.(input.KeyboardEvent)
		switch ev.Key {
		case input.K_INSERT:
			item := ui.NewTextbox(15, 1, vec.ZERO_COORD, 1, "new item", false)
			ts.list.AddChild(item)
		case input.K_RIGHT:
			ts.moving_text.Move(1, 0)
		case input.K_LEFT:
			ts.moving_text.Move(-1, 0)
		case input.K_UP:
			ts.moving_text.Move(0, -1)
		case input.K_DOWN:
			ts.moving_text.Move(0, 1)
		case input.K_RETURN:
			ts.moving_text.ToggleVisible()
		case input.K_F1:
			ts.GiveUserSeizure()
		}
	}
}

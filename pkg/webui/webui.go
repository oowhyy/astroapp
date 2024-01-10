package webui

import (
	"syscall/js"
)

type UserInterface interface {
	// TODO: if window resize listener is needed - implement callback instead
	WindowSize() (float64, float64)
	IsPaused() bool
	IsAddMode() bool

	OnClearTrail(callback func())
	OnSpeedUp(callback func() int)
	OnSlowDown(callback func() int)
}

type WebInterface struct {
	document     js.Value
	pauseButton  *button
	addButton    *button
	clearButton  *button
	speedControl js.Value
}

type button struct {
	stateOn bool
	elem    js.Value
}

func newButton(elem js.Value) *button {
	b := &button{
		stateOn: false,
		elem:    elem,
	}
	return b
}

func NewWebInterface() *WebInterface {
	wi := &WebInterface{}
	wi.document = js.Global().Get("document")
	// play button
	wi.pauseButton = newButton(wi.document.Call("getElementById", "play-pause"))
	wi.pauseButton.elem.Call("addEventListener", "click", js.FuncOf(wi.togglePlay))
	// add button
	wi.addButton = newButton(wi.document.Call("getElementById", "add"))
	wi.addButton.elem.Call("addEventListener", "click", js.FuncOf(wi.toggleAdd))
	// clearTrail button
	wi.clearButton = newButton(wi.document.Call("getElementById", "clearTrailDots"))

	// speed controller
	wi.speedControl = wi.document.Call("getElementById", "speedControl")
	return wi
}

func (wi *WebInterface) WindowSize() (float64, float64) {
	w := wi.document.Get("body").Get("clientWidth").Float()
	h := wi.document.Get("body").Get("clientHeight").Float()
	return w, h

}

func (wi *WebInterface) OnClearTrail(callback func()) {
	wi.clearButton.elem.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		callback()
		return nil
	}))
}

func (wi *WebInterface) OnSpeedUp(callback func() int) {
	wi.speedControl.Call("addEventListener", "wheel", js.FuncOf(func(this js.Value, args []js.Value) any {
		if args[0].Get("deltaY").Float() < 0 {
			curspeed := callback()
			wi.speedControl.Set("innerHTML", js.ValueOf(curspeed))
		}
		// fmt.Println("ARGS", args[0].Get("deltaY"))
		return nil
	}))
}

func (wi *WebInterface) OnSlowDown(callback func() int) {
	wi.speedControl.Call("addEventListener", "wheel", js.FuncOf(func(this js.Value, args []js.Value) any {
		if args[0].Get("deltaY").Float() > 0 {
			curspeed := callback()
			wi.speedControl.Set("innerHTML", js.ValueOf(curspeed))
		}
		// fmt.Println("ARGS", args[0].Get("deltaY"))
		return nil
	}))
}

func (wi *WebInterface) togglePlay(this js.Value, args []js.Value) any {
	wi.pauseButton.stateOn = !wi.pauseButton.stateOn
	wi.pauseButton.elem.Get("classList").Call("toggle", "pause")
	return nil
}

func (wi *WebInterface) toggleAdd(this js.Value, args []js.Value) any {
	wi.addButton.stateOn = !wi.addButton.stateOn
	wi.addButton.elem.Get("classList").Call("toggle", "active")
	return nil
}

func (wi *WebInterface) IsPaused() bool {
	return wi.pauseButton.stateOn
}

func (wi *WebInterface) IsAddMode() bool {
	return wi.addButton.stateOn
}

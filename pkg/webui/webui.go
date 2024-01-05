package webui

import (
	"syscall/js"
)

type UserInterface interface {
	// TODO: if window resize listener is needed - implement callback instead
	WindowSize() (float64, float64)
	IsPaused() bool
	IsAddMode() bool
}

type WebInterface struct {
	document    js.Value
	pauseButton *button
	addButton   *button
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
	return wi
}

func (wi *WebInterface) WindowSize() (float64, float64) {
	w := wi.document.Get("body").Get("clientWidth").Float()
	h := wi.document.Get("body").Get("clientHeight").Float()
	return w, h

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

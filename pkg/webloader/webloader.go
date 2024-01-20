package webloader

import (
	"syscall/js"
)

func LoadFile(name string) string {
	promise := js.Global().Call("fetch", name)
	resp, err := await(promise)
	if err != nil {
		panic(err)
	}
	text, err := await(resp[0].Call("text"))
	if err != nil {
		panic(err)
	}
	return text[0].String()

}

func await(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}

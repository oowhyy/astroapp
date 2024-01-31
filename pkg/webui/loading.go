package webui

import (
	"fmt"
	"time"
)

func (wi *WebInterface) SetLoadingMessage(message string) {
	if !wi.isLoading {
		wi.isLoading = true
		go func() {
			for {
				select {
				case msg := <-wi.messagesChan:
					if wi.loadingMessage != "" {
						oldFull := fmt.Sprintf("%s %.2f", wi.loadingMessage, wi.loadingTimer.Seconds())
						wi.oldMsgs = oldFull + "<br />" + wi.oldMsgs
					}
					wi.loadingMessage = msg
					wi.loadingTimer = 0
				case <-time.After(time.Millisecond * 10):
					wi.loadingTimer += time.Millisecond * 10
					s := fmt.Sprintf("%s %.2f<br />%s", wi.loadingMessage, wi.loadingTimer.Seconds(), wi.oldMsgs)
					wi.loadingDiv.Set("innerHTML", s)
				case <-wi.done:
					return
				}
			}
		}()
	}
	fmt.Println("GOT MESSAGE", message)
	wi.messagesChan <- message
}

func (wi *WebInterface) DoneLoading() {
	if !wi.isLoading {
		return
	}
	fmt.Println(wi.oldMsgs)
	wi.addButton.elem.Get("style").Set("visibility", "visible")
	wi.clearButton.elem.Get("style").Set("visibility", "visible")
	wi.pauseButton.elem.Get("style").Set("visibility", "visible")
	wi.speedControl.Get("style").Set("visibility", "visible")

	wi.loadingDiv.Get("style").Set("display", "none")
	wi.done <- true
	wi.isLoading = false
}

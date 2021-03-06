package main

import "fmt"

type hub struct {
    listener *connection
    listening bool
    command chan string
    connect chan *connection
    disconnect chan *connection
    currentsong string
}

var h = hub{
    command: make(chan string),
    connect: make(chan *connection),
    disconnect: make(chan *connection),
    listening: false,
}

func (h *hub) run() {
    for {
        select {
        case c := <- h.disconnect:
            if c == h.listener && h.listening {
                close(h.listener.send)
                h.listener.ws.Close()
                h.listening = false
            }
        case c := <- h.connect:
            if !h.listening {
                h.listener = c
                h.listening = true
            } else {
                c.send <- "close"
                close(c.send)
                c.ws.Close()
            }
        case m := <- h.command:
            if m[:2] == "np" {
                fmt.Println(m)
                h.currentsong = m[2:]
            } else if h.listening {
                h.listener.send <- m 
            }
        }
    }
}

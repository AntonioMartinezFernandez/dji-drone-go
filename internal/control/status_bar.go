package control

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

type StatusBar struct {
	mu  sync.Mutex
	Msg string
	Col tcell.Color
}

func (sb *StatusBar) Set(msg string, col tcell.Color) {
	sb.mu.Lock()
	sb.Msg = msg
	sb.Col = col
	sb.mu.Unlock()
}

func (sb *StatusBar) Get() (string, tcell.Color) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Msg, sb.Col
}

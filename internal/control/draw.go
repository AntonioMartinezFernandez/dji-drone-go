package control

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// ─── Drawing helpers ──────────────────────────────────────────────────────────

func DrawAll(screen tcell.Screen, st *Sticks, sb *StatusBar) {
	screen.Fill(' ', tcell.StyleDefault)

	lr, fb, ud, y := st.Get()

	lines := []struct {
		text string
		col  tcell.Color
	}{
		{"╔══════════════════════════════════════════╗", tcell.ColorYellow},
		{"║       DJI DRONE — REALTIME CONTROL       ║", tcell.ColorYellow},
		{"╚══════════════════════════════════════════╝", tcell.ColorYellow},
		{"", tcell.ColorWhite},
		{"  COMMANDS                                  ", tcell.ColorWhite},
		{"  [T] Take off      [L] Land                ", tcell.ColorLightCyan},
		{"  [E] Emergency stop.                       ", tcell.ColorRed},
		{"", tcell.ColorWhite},
		{"  REALTIME MOVEMENT  (hold key)             ", tcell.ColorWhite},
		{"  [W/S]  Forward / Backward  (FB axis)      ", tcell.ColorGreen},
		{"  [A/D]  Yaw left / right    (Yaw axis)     ", tcell.ColorGreen},
		{"  [Q/E]  Strafe left / right (LR axis)      ", tcell.ColorGreen},
		{"  [R/F]  Up / Down           (UD axis)      ", tcell.ColorGreen},
		{"", tcell.ColorWhite},
		{"  FLIPS  (single press)                     ", tcell.ColorWhite},
		{"  [I] Forward  [K] Backward                 ", tcell.ColorDarkMagenta},
		{"  [J] Left     [O] Right                    ", tcell.ColorDarkMagenta},
		{"", tcell.ColorWhite},
		{"  [ESC] Land and quit                       ", tcell.ColorGray},
		{"", tcell.ColorWhite},
		{"─────────────────────────────────────────────", tcell.ColorDarkGray},
		{"  STICK VALUES                              ", tcell.ColorWhite},
		{fmt.Sprintf("  LR: %+4d   FB: %+4d   UD: %+4d   Yaw: %+4d", lr, fb, ud, y), tcell.ColorLightCyan},
		{"─────────────────────────────────────────────", tcell.ColorDarkGray},
	}

	for row, line := range lines {
		drawText(screen, 2, row+1, line.text, line.col)
	}

	// Status bar at the bottom
	_, h := screen.Size()
	msg, col := sb.Get()
	drawText(screen, 2, h-1, fmt.Sprintf("  %-60s", msg), col)

	screen.Show()
}

func drawText(screen tcell.Screen, x, y int, text string, color tcell.Color) {
	style := tcell.StyleDefault.Foreground(color)
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

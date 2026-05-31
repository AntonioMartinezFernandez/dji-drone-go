package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AntonioMartinezFernandez/dji-drone-go/internal/control"
	astitello "github.com/asticode/go-astitello"
	"github.com/gdamore/tcell/v2"
)

// Stick speed values sent to the drone (-100 to 100).
const stickSpeed = 60

// keyHoldTimeout: if no repeat event arrives within this duration, the key is
// considered released and its axis is zeroed. OS key-repeat fires every ~30ms,
// so 150ms gives comfortable margin.
const keyHoldTimeout = 150 * time.Millisecond

// rcTicker controls how often SetSticks is sent to the drone (~20 Hz).
const rcTickInterval = 50 * time.Millisecond

func main() {
	logFile, err := os.OpenFile("dji-drone.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()
	l := log.New(logFile, "", log.LstdFlags)

	drone := astitello.New(l)

	l.Println("Connecting to Tello…")
	if err := drone.Start(); err != nil {
		l.Fatalf("start: %v", err)
	}
	defer drone.Close()
	l.Println("Connected!")

	drone.On(astitello.TakeOffEvent, func(_ interface{}) { l.Println("took off") })
	drone.On(astitello.LandEvent, func(_ interface{}) { l.Println("landed") })

	// Graceful shutdown on SIGINT/SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = drone.Land()
		drone.Close()
		os.Exit(0)
	}()

	screen, err := tcell.NewScreen()
	if err != nil {
		l.Fatalf("screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		l.Fatalf("screen init: %v", err)
	}
	defer screen.Fini()
	screen.EnableMouse() // not needed but harmless

	st := &control.Sticks{}
	ks := control.NewKeyState(keyHoldTimeout)
	sb := &control.StatusBar{Msg: "Ready — press T to take off", Col: tcell.ColorGreen}

	control.DrawAll(screen, st, sb)

	// ── RC ticker: send sticks to drone at ~20 Hz ─────────────────────────────
	go func() {
		ticker := time.NewTicker(rcTickInterval)
		defer ticker.Stop()
		for range ticker.C {
			lr, fb, ud, y := st.Get()
			if err := drone.SetSticks(lr, fb, ud, y); err != nil {
				sb.Set(fmt.Sprintf("SetSticks error: %v", err), tcell.ColorRed)
			}
		}
	}()

	// ── UI refresh ticker ─────────────────────────────────────────────────────
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			control.DrawAll(screen, st, sb)
		}
	}()

	// ── Event loop ────────────────────────────────────────────────────────────
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {

		case *tcell.EventResize:
			screen.Sync()
			control.DrawAll(screen, st, sb)

		case *tcell.EventKey:
			// Single-shot commands (take-off, land, flip) are sent in a goroutine
			// so they don't block the input loop.
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				sb.Set("Landing…", tcell.ColorYellow)
				control.DrawAll(screen, st, sb)
				_ = drone.Land()
				screen.Fini()
				fmt.Println("Drone landed. Bye!")
				return
			}

			neg := func(v int) int { return -v }
			_ = neg

			switch ev.Rune() {
			// ── Take-off / Land ───────────────────────────────────────────────
			case 't', 'T':
				sb.Set("Taking off…", tcell.ColorLightCyan)
				go func() {
					if err := drone.TakeOff(); err != nil {
						sb.Set(fmt.Sprintf("takeoff: %v", err), tcell.ColorRed)
					} else {
						sb.Set("Airborne — use WASD/QZ to fly", tcell.ColorGreen)
					}
				}()
			case 'l', 'L':
				sb.Set("Landing…", tcell.ColorYellow)
				go func() {
					if err := drone.Land(); err != nil {
						sb.Set(fmt.Sprintf("land: %v", err), tcell.ColorRed)
					} else {
						sb.Set("Landed", tcell.ColorGreen)
					}
				}()
			case '1':
				sb.Set("EMERGENCY STOP", tcell.ColorRed)
				go drone.Emergency()

			// ── Flips (single-shot) ───────────────────────────────────────────
			case 'i', 'I':
				sb.Set("Flip forward", tcell.ColorDarkMagenta)
				go drone.Flip(astitello.FlipForward)
			case 'k', 'K':
				sb.Set("Flip backward", tcell.ColorDarkMagenta)
				go drone.Flip(astitello.FlipBack)
			case 'j', 'J':
				sb.Set("Flip left", tcell.ColorDarkMagenta)
				go drone.Flip(astitello.FlipLeft)
			case 'o', 'O':
				sb.Set("Flip right", tcell.ColorDarkMagenta)
				go drone.Flip(astitello.FlipRight)

			// ── Realtime movement ─────────────────────────────────────────────
			// Forward / backward  →  FB axis
			case 'w', 'W':
				ks.Press(st, control.AxisFB, stickSpeed)
			case 's', 'S':
				ks.Press(st, control.AxisFB, -stickSpeed)

			// Strafe left / right  →  LR axis
			case 'q', 'Q':
				ks.Press(st, control.AxisLR, -stickSpeed)
			case 'e', 'E':
				ks.Press(st, control.AxisLR, stickSpeed)

			// Up / down  →  UD axis
			case 'r', 'R':
				ks.Press(st, control.AxisUD, stickSpeed)
			case 'f', 'F':
				ks.Press(st, control.AxisUD, -stickSpeed)

			// Yaw left / right  →  Yaw axis
			case 'a', 'A':
				ks.Press(st, control.AxisYaw, -stickSpeed)
			case 'd', 'D':
				ks.Press(st, control.AxisYaw, stickSpeed)
			}
		}
	}
}

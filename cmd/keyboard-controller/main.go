package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astitello"
	"github.com/eiannone/keyboard"
)

func main() {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	worker := astikit.NewWorker(astikit.WorkerOptions{Logger: l})
	drone := astitello.New(l)
	worker.HandleSignals(astikit.TermSignalHandler(func() {
		if err := drone.Land(); err != nil {
			l.Println(fmt.Errorf("main: landing failed: %w\n", err))
			return
		}
	}))

	drone.On(astitello.TakeOffEvent, func(any) { l.Println("main: drone has took off!") })
	if err := drone.Start(); err != nil {
		l.Println(fmt.Errorf("main: starting to the drone failed: %w", err))
		return
	}
	defer drone.Close()

	err := keyboard.Open()
	if err != nil {
		l.Printf("Failed to open keyboard: %v", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	fmt.Println("------------------------------")
	fmt.Println("           CONTROLS")
	fmt.Println("------------------------------")
	fmt.Println(" T - Take off")
	fmt.Println(" L - Land")
	fmt.Println("------------------------------")
	fmt.Println(" P - Flip right")
	fmt.Println(" O - Flip left")
	fmt.Println("------------------------------")
	fmt.Println(" W - Move forward")
	fmt.Println(" S - Move backward")
	fmt.Println(" A - Move left")
	fmt.Println(" D - Move right")
	fmt.Println("------------------------------")
	fmt.Println(" U - Move up")
	fmt.Println(" J - Move down")
	fmt.Println("------------------------------")
	fmt.Println(" ESC - Exit")
	fmt.Println("------------------------------")

	for {
		char, key, err := keyboard.GetKey()

		if err != nil {
			l.Printf("Failed to get keyboard input: %v\n", err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		switch char {

		case 't', 'T':
			if err := drone.TakeOff(); err != nil {
				l.Printf("Failed to take off: %v\n", err)
			}

		case 'l', 'L':
			if err := drone.Land(); err != nil {
				l.Printf("Failed to land: %v\n", err)
			}

		case 'p', 'P':
			if err := drone.Flip(astitello.FlipRight); err != nil {
				l.Printf("Failed to flip: %v\n", err)
			}

		case 'o', 'O':
			if err := drone.Flip(astitello.FlipLeft); err != nil {
				l.Printf("Failed to flip: %v\n", err)
			}

		case 'w', 'W':
			if err := drone.Forward(5); err != nil {
				l.Printf("Failed to move forward: %v\n", err)
			}

		case 's', 'S':
			if err := drone.Back(5); err != nil {
				l.Printf("Failed to move backward: %v\n", err)
			}

		case 'a', 'A':
			if err := drone.Left(5); err != nil {
				l.Printf("Failed to move left: %v\n", err)
			}

		case 'd', 'D':
			if err := drone.Right(5); err != nil {
				l.Printf("Failed to move right: %v\n", err)
			}

		case 'u', 'U':
			if err := drone.Up(5); err != nil {
				l.Printf("Failed to move up: %v\n", err)
			}

		case 'j', 'J':
			if err := drone.Down(5); err != nil {
				l.Printf("Failed to move down: %v\n", err)
			}

		}

		time.Sleep(10 * time.Millisecond)
	}
}

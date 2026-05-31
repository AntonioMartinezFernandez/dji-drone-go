package main

import (
	"fmt"
	"log"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astitello"
)

func main() {
	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create worker
	worker := astikit.NewWorker(astikit.WorkerOptions{Logger: l})

	// Create the drone
	drone := astitello.New(l)

	// Handle signals
	worker.HandleSignals(astikit.TermSignalHandler(func() {
		// Make sure to land on term signal
		if err := drone.Land(); err != nil {
			l.Println(fmt.Errorf("main: landing failed: %w", err))
			return
		}
	}))

	// Handle take off event
	drone.On(astitello.TakeOffEvent, func(any) { l.Println("main: drone has took off!") })

	// Start the drone
	if err := drone.Start(); err != nil {
		l.Println(fmt.Errorf("main: starting to the drone failed: %w", err))
		return
	}
	defer drone.Close()

	// Execute in a task
	worker.NewTask().Do(func() {
		// Take off
		if err := drone.TakeOff(); err != nil {
			l.Println(fmt.Errorf("main: taking off failed: %w", err))
			return
		}

		// Flip
		if err := drone.Flip(astitello.FlipRight); err != nil {
			l.Println(fmt.Errorf("main: flipping failed: %w", err))
			return
		}

		// Log state
		l.Printf("main: state is: %+v\n", drone.State())

		// Land
		if err := drone.Land(); err != nil {
			l.Println(fmt.Errorf("main: landing failed: %w", err))
			return
		}

		// Stop worker
		worker.Stop()
	})

	// Wait
	worker.Wait()
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func StartNewProcess(binary string) {
	cmd := exec.Command(binary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting new process:", err)
		return
	}

	// Store the new process instance
	runningVersion = cmd.Process

	// Wait for a brief moment to allow the new process to start
	time.Sleep(3 * time.Second)
}

func StopPreviousProcess() {
	if runningVersion != nil {
		// Send a termination signal to the new process
		err := runningVersion.Kill()
		if err != nil {
			fmt.Println("Error stopping new process:", err)
			return
		}

		processState, err := runningVersion.Wait()

		if err != nil {
			panic(err)
		}

		fmt.Println("Waiting for program to finish..")
		time.Sleep(5 * time.Second)

		if processState.Exited() {
			panic("program could not be exited within timeframe")
		}

		err = runningVersion.Release()

		if err != nil {
			panic(err)
		}

		// Reset the newProcess variable
		runningVersion = nil
	}
}

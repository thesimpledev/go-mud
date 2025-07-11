package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	fmt.Println("Server Started.")

	ctx, cancel := context.WithCancel(context.Background())
	commands := make(chan string)
	errors := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go getServerCommands(cancel, commands, errors)
	go commandLoop(&wg, ctx, commands, errors)

	wg.Wait()

}

func commandLoop(wg *sync.WaitGroup, ctx context.Context, commands <-chan string, errors <-chan error) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errors:
			if !ok {
				return
			}
			//TODO: As error handling becomes more complex we will need to do more than print
			//and return from the errors.
			fmt.Println(err)
			return
		case command := <-commands:
			process(command)
		}
	}
}

func process(command string) {
	fmt.Println("Command Received: ", command)
}

func getServerCommands(cancel context.CancelFunc, commands chan<- string, errors chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	//Note: this is the only place this channel signal happens so this is safe.
	defer cancel()
	defer close(commands)
	defer close(errors)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			errors <- err
			return
		}
		input = strings.TrimSpace(input)

		switch input {
		case "quit":
			fallthrough
		case "exit":
			return
		default:
			commands <- input
		}
	}

}

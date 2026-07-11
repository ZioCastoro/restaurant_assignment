package main

import (
	"fmt"
	"os"

	"github.com/ZioCastoro/restaurant_assignment/deps"
)

func main() {
	app, err := deps.InjectApplication()
	if err != nil {
		fmt.Printf("error while injecting application: %v", err)

		os.Exit(1)
	}

	if err = app.Routes().Listen(":8080"); err != nil {
		fmt.Printf("error while serving application: %v", err)

		os.Exit(1)
	}
}

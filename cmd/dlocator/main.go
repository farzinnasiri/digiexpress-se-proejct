package main

import (
	"context"

	"github.com/digiexpress/dlocator/internal/app/dlocator"
)

func main() {
	app, err := dlocator.CreateApp()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	app.Run(ctx)
}

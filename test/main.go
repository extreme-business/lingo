package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	g.Go(func() error {
		<-ctx.Done()
		fmt.Printf("done1\n")
		return errors.New("error1")
	})

	g.Go(func() error {
		<-ctx.Done()
		fmt.Printf("done2\n")
		return errors.New("error2")
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("done\n")
}

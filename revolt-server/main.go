package main

import (
	"fmt"
	"os"
)

func run() error {
	RunServer()

	game := NewGame()
	game.AddPlayer()
	game.AddPlayer()
	game.AddPlayer()

	// Do a round of tax so everyone has cash.
	for range len(game.Players) {
		// Turn 1
		tax := Action{
			Type: Tax,
		}
		err := game.AttemptAction(tax)
		if err != nil {
			return err
		}
		err = game.CommitTurn()
		if err != nil {
			return err
		}
		fmt.Println(game)
		err = game.EndTurn()
		if err != nil {
			return err
		}
	}

	assassinate := Action{
		Type:         Assassinate,
		TargetPlayer: 1,
	}
	err := game.AttemptAction(assassinate)

	if err != nil {
		return err
	}

	err = game.CommitTurn()
	if err != nil {
		return err
	}
	err = game.ResolveDeath(0)
	if err != nil {
		return err
	}
	fmt.Println(game)

	err = game.EndTurn()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

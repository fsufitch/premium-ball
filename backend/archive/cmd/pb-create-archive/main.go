package main

import (
	"fmt"
	"os"

	"github.com/fsufitch/premium-ball/archive"
)

func main() {
	nCards := 50
	count := 0
	for result := range archive.GetAllPokemonTCGCards() {
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "error fetching cards: %s\n", result.Error)
			continue
		}
		pbc := archive.ToPremiumBallCard(*result.Card)

		fmt.Printf("- %s -- %s \n", pbc.CardData.Name, pbc.MechanicsHash)
		count += 1
		if count >= nCards {
			break
		}
	}
}

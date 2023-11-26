package main

import "fmt"

func printGame(game FreeGameEntry) {
	fmt.Println("EpicId: ", game.EpicId)
	fmt.Println("EpicStoreLink: ", game.EpicStoreLink)
	fmt.Println("FreeDate: ", game.FreeDate)
	fmt.Println("GameTitle: ", game.GameTitle)
	fmt.Println("MappingSlug: ", game.MappingSlug)
	fmt.Println("ProductSlug: ", game.ProductSlug)
	fmt.Println("SandboxId: ", game.SandboxId)
	fmt.Println("UrlSlug: ", game.UrlSlug)
	fmt.Println("------------------------------------")
}

func CliHandlerFree() {
	freeGames, _ := GetFreeGames()

	fmt.Println("Free This Week:")
	for _, game := range freeGames.ThisWeek {
		printGame(game)
	}
	fmt.Println("Free Next Week:")
	for _, game := range freeGames.NextWeek {
		printGame(game)
	}

}

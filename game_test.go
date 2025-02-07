package main

import (
	"fmt"
	"testing"
)

func TestGameRollDiceAndCards(t *testing.T) {

	game := NewGame()
	gs := *new(GameSetting)
	gs.Map = 0
	gs.NumberOfPlayers = 2
	gs.TurnTimeOut = false
	game.UpdateGameSetting(gs)
	game.Start()

	//fmt.Println(game.Board())
	availableIntersections := game.context.board.GetAvailableIntersections([]int{})
	if len(availableIntersections) != 54 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	//player 0
	game.PutSettlement(14)

	availableRoads := game.context.board.GetNeighborIntersections1(14)
	if len(availableRoads) != 3 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutRoad([2]int{14, 15})

	//player 1

	availableIntersections = game.context.board.GetAvailableIntersections([]int{14})
	if len(availableIntersections) != 50 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutSettlement(26)

	availableRoads = game.context.board.GetNeighborIntersections1(26)
	if len(availableRoads) != 3 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutRoad([2]int{25, 26})

	//player 1

	availableIntersections = game.context.board.GetAvailableIntersections([]int{14, 26})
	if len(availableIntersections) != 46 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutSettlement(41)

	availableRoads = game.context.board.GetNeighborIntersections1(41)
	if len(availableRoads) != 3 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutRoad([2]int{41, 32})

	//player 0
	availableIntersections = game.context.board.GetAvailableIntersections([]int{14, 26, 41})
	if len(availableIntersections) != 42 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}

	game.PutSettlement(13)
	availableRoads = game.context.board.GetNeighborIntersections1(13)
	if len(availableRoads) != 3 {
		t.Log("expected to have 54 intersections, failed.")
		t.Fail()
	}
	game.PutRoad([2]int{13, 20})

	player0 := game.getPlayer(0)
	if len(player0.settlements) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}
	if len(player0.roads) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}

	player1 := game.getPlayer(1)
	if len(player1.settlements) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}
	if len(player1.roads) != 2 {
		t.Log("expecting 2 settlements, failed")
		t.Fail()
	}

	if player0.settlements[0].Intersection != 14 {
		t.Log("expecting settlement in 14, failed")
		t.Fail()
	}
	if player0.settlements[1].Intersection != 13 {
		t.Log("expecting settlement in 13, failed")
		t.Fail()
	}

	if player1.settlements[0].Intersection != 26 {
		t.Log("expecting settlement in 14, failed")
		t.Fail()
	}
	if player1.settlements[1].Intersection != 41 {
		t.Log("expecting settlement in 13, failed")
		t.Fail()
	}

	game.context.handleDice(12)
	game.context.handleDice(6)
	game.context.handleDice(8)
	game.context.handleDice(4)

	//if game.getPlayer(0).cards[3] != 1 {
	//	t.Log("expecting 1 grain, failed")
	//	t.Fail()
	//}
	//if game.getPlayer(0).cards[1] != 1 {
	//	t.Log("expecting 1 brick, failed")
	//	t.Fail()
	//}
	//if game.getPlayer(0).cards[2] != 1 {
	//	t.Log("expecting 1 wool, failed")
	//	t.Fail()
	//}
	//
	//if game.getPlayer(1).cards[0] != 1 {
	//	t.Log("expecting 1 tree, failed")
	//	t.Fail()
	//}
	//if game.getPlayer(1).cards[3] != 2 {
	//	t.Log("expecting 1 grain, failed")
	//	t.Fail()
	//}
	//
	//if game.context.Bank.cards[0] != 18 {
	//	t.Log("expecting 18 tree remaining, failed")
	//	t.Fail()
	//}
	//
	//if game.context.Bank.cards[1] != 18 {
	//	t.Log("expecting 18 brick remaining, failed")
	//	t.Fail()
	//}
	//
	//if game.context.Bank.cards[3] != 16 {
	//	t.Log("expecting 18 grain remaining, failed")
	//	t.Fail()
	//}

	fmt.Println("")
}

package game

import (
	"catans/board"
	"catans/utils"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type gameAction struct {
	Name    string
	Timeout time.Time
}

type gameTrade struct {
	ID          int
	From        int
	To          int
	Gives       [][2]int
	Wants       [][2]int
	HasAllCards bool
	OK          int
}

type GameContext struct {
	tiles        []Tile
	board        board.Board
	tradeCounter int
	trades       []*gameTrade

	GameSetting
	GameState
}

func (context GameContext) getCurrentPlayer() *Player {
	return context.Players[context.CurrentPlayer]
}

func (context GameContext) getGamePhase() string {
	return context.Phase
}

func (context *GameContext) putSettlement(validate bool, intersection int) error {
	if validate {
		availableIntersections, _ := context.getPossibleSettlementLocations()
		matched := false
		for _, availableIntersection := range availableIntersections {
			if availableIntersection == intersection {
				matched = true
			}
		}
		if !matched {
			return errors.New("invalid action")
		}
	}

	var settlement Settlement

	if Phase2 == context.Phase || Phase3 == context.Phase {
		nextAction := context.getAction()
		if nextAction != nil && nextAction.Name != ActionPlaceSettlement {
			return errors.New("invalid action")
		}
		tileIndices := context.board.GetTiles(intersection)
		tokens := make([]int, len(tileIndices))
		for i, idx := range tileIndices {
			tokens[i] = context.tiles[idx].Token
		}
		settlement = Settlement{Indices: tileIndices, Tokens: tokens, Intersection: intersection}
	}

	return context.getCurrentPlayer().putSettlement(settlement)
}

func (context *GameContext) getPossibleRoads() ([][2]int, error) {

	if Phase4 == context.Phase {
		currentPlayer := context.getCurrentPlayer()

		noOfCoords := len(currentPlayer.Settlements) + len(currentPlayer.Roads)
		var occupiedIns = make([]int, noOfCoords)
		i := 0
		for _, settlement := range currentPlayer.Settlements {
			occupiedIns[i] = settlement.Intersection
			i++
		}
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				occupiedIns[i] = road[0]
			}
			if !utils.Contains(occupiedIns, road[1]) {
				occupiedIns[i] = road[1]
			}
			i++
		}

		i = 0
		var roads = make([][2]int, noOfCoords*3)
		for _, ins := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections2(ins)
			for _, ni := range neighborIntersections {
				roads[i] = ni
				i++
			}
		}

		var allRoads [][2]int
		for _, player := range context.Players {
			allRoads = append(allRoads, player.Roads...)
		}

		var availableRoads [][2]int
		for _, newRoad := range roads {
			found := false
			for _, currentRoad := range allRoads {
				if newRoad[0] == currentRoad[0] && newRoad[1] == currentRoad[1] {
					found = true
					break
				}
			}
			if !found {
				availableRoads = append(availableRoads, newRoad)
			}
		}
		return availableRoads, nil
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		getRoadsForIntersection := func(settlement *Settlement) [][2]int {
			var roads [][2]int
			if settlement != nil {
				neighborIntersections := context.board.GetNeighborIntersections2(settlement.Intersection)
				roads = append(roads, neighborIntersections...)
			}
			return roads
		}

		if Phase2 == context.Phase {
			nextAction := context.getAction()
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
				return nil, errors.New("invalid action")
			}
			currentPlayer := context.getCurrentPlayer()

			var firstSettlement *Settlement
			if len(currentPlayer.Settlements) > 0 {
				firstSettlement = &currentPlayer.Settlements[0]
			}
			return getRoadsForIntersection(firstSettlement), nil
		}

		if Phase3 == context.Phase {
			nextAction := context.getAction()
			if nextAction != nil && nextAction.Name != ActionPlaceRoad {
				return nil, errors.New("invalid action")
			}

			currentPlayer := context.getCurrentPlayer()
			var (
				settlementCounter = 1
				secondSettlement  *Settlement
			)

			if len(currentPlayer.Settlements) > 1 {
				settlementCounter++
				secondSettlement = &currentPlayer.Settlements[1]
			}

			return getRoadsForIntersection(secondSettlement), nil
		}
	}

	return nil, errors.New("invalid action")
}

func (context *GameContext) getPossibleSettlementLocations() ([]int, error) {
	if Phase4 == context.Phase {

		var occupiedIns []int
		for _, player := range context.Players {
			for _, settlement := range player.Settlements {
				occupiedIns = append(occupiedIns, settlement.Intersection)
			}
		}

		for _, v := range occupiedIns {
			neighborIntersections := context.board.GetNeighborIntersections1(v)
			for _, ni := range neighborIntersections {
				if !utils.Contains(occupiedIns, ni) {
					occupiedIns = append(occupiedIns, ni)
				}
			}
		}

		currentPlayer := context.getCurrentPlayer()
		var roadsIntersections []int
		for _, road := range currentPlayer.Roads {
			if !utils.Contains(occupiedIns, road[0]) {
				roadsIntersections = append(roadsIntersections, road[0])
			}
			if !utils.Contains(occupiedIns, road[1]) {
				roadsIntersections = append(roadsIntersections, road[1])
			}
		}

		var availableIntersections []int
		for _, ins := range roadsIntersections {
			if !utils.Contains(occupiedIns, ins) {
				availableIntersections = append(availableIntersections, ins)
			}
		}
		return availableIntersections, nil
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		nextAction := context.getAction()
		if nextAction != nil && nextAction.Name != ActionPlaceSettlement {
			return nil, errors.New("invalid action")
		}
		occupied := make([]int, 0)
		for _, player := range context.Players {
			for _, s := range player.Settlements {
				occupied = append(occupied, s.Intersection)
			}
		}
		return context.board.GetAvailableIntersections(occupied), nil
	}

	return nil, nil
}

func (context *GameContext) putRoad(validate bool, road [2]int) error {
	if validate {
		availableRoads, _ := context.getPossibleRoads()
		matched := false
		for _, availableRoad := range availableRoads {
			if availableRoad == road {
				matched = true
			}
		}
		if !matched {
			return errors.New("invalid action")
		}
	}

	if Phase2 == context.Phase || Phase3 == context.Phase {
		nextAction := context.getAction()
		if nextAction != nil && nextAction.Name != ActionPlaceRoad {
			return errors.New("invalid action")
		}
	}

	return context.getCurrentPlayer().putRoad(road)
}

func (context *GameContext) handOverCards(player *Player, cardType int, count int) {
	player.Cards[cardType] = player.Cards[cardType] + count
}

func (context *GameContext) updateGameSetting(gs GameSetting) error {
	if context.GameState.Phase != Phase1 {
		return errors.New("invalid action")
	}
	context.GameSetting = gs
	for i := 0; i < gs.NumberOfPlayers; i++ {
		player := NewPlayer()
		player.ID = i
		context.Players = append(context.Players, player)
	}
	context.tiles = GenerateTiles("10MO,2PA,9FO,12FI,6HI,4PA,10HI,9FI,11FO,0DE,3FO,8MO,8FO,3MO,4FI,5PA,5HI,6FI,11PA")
	context.board = board.NewBoard(gs.Map)
	return nil
}

func (context *GameContext) isInitialSettlementDone() bool {
	settlementCount := 0
	for _, player := range context.Players {
		settlementCount = settlementCount + len(player.Settlements)
	}
	return settlementCount == (context.GameSetting.NumberOfPlayers * 2)
}

func (context GameContext) getAction() *gameAction {
	return &context.Action
}

func (context GameContext) getActionString() string {
	return context.Action.Name
}

func (context *GameContext) scheduleAction(action string) {
	context.Action = gameAction{Name: action, Timeout: time.Now()}
}

func (context *GameContext) bankTrade(gives, wants [][2]int) error {
	if len(gives) == 1 && len(wants) == 1 && wants[0][1] == 1 && gives[0][1] > 1 {
		currentPlayer := context.getCurrentPlayer()
		if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
			return errors.New("invalid operation")
		}

		bank := context.Bank
		bank.Begin()
		defer bank.Commit()

		wantCardType := wants[0][0]
		wantTradeCount := wants[0][1]
		giveCardType := gives[0][0]
		giveTradeCount := gives[0][1]

		if currentPlayer.has21 || currentPlayer.has31 {
			if currentPlayer.has21 && currentPlayer.cards21[giveCardType] == 1 && giveTradeCount == 2 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -2
				currentPlayer.Cards[wantCardType]++
			} else if currentPlayer.has31 && giveTradeCount == 3 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -3
				currentPlayer.Cards[wantCardType]++
			}
		} else {
			if giveTradeCount == 4 {
				if err := bank.Return(giveCardType, giveTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				if _, err := bank.Give(wantCardType, wantTradeCount); err != nil {
					bank.Rollback()
					return err
				}
				currentPlayer.Cards[giveCardType] = -4
				currentPlayer.Cards[wantCardType]++
			}
		}

	}

	return errors.New("invalid operation")
}

func (context *GameContext) getTrade(tradeID int) *gameTrade {
	for _, t := range context.trades {
		if t.ID == tradeID {
			return t
		}
	}
	return nil
}

func (context *GameContext) setupTrade(gives [][2]int, wants [][2]int) error {

	currentPlayer := context.getCurrentPlayer()
	if !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(currentPlayer.ID, gives) {
		return errors.New("invalid operation")
	}

	for _, otherPlayer := range context.Players {
		if otherPlayer.ID != currentPlayer.ID {

			hasAllCards := false
			if context.isPlayerHasAllCards(otherPlayer.ID, wants) {
				hasAllCards = true
			}

			var trade = new(gameTrade)
			trade.ID = context.tradeCounter
			trade.From = currentPlayer.ID
			trade.To = otherPlayer.ID
			trade.Gives = gives
			trade.Wants = wants
			trade.HasAllCards = hasAllCards
			trade.OK = 0
			context.trades = append(context.trades, trade)

			//race condition
			context.tradeCounter++
		}
	}

	return nil
}

func (context *GameContext) overrideTrade(tradeID int, gives [][2]int, wants [][2]int) error {
	var trade = context.getTrade(tradeID)
	if trade == nil || !context.isSafeTrade(gives, wants) || !context.isPlayerHasAllCards(trade.To, gives) {
		return errors.New("invalid operation")
	}

	currentPlayer := context.getCurrentPlayer()

	hasAllCards := false
	if context.isPlayerHasAllCards(currentPlayer.ID, wants) {
		hasAllCards = true
	}
	trade.From = trade.To
	trade.To = context.CurrentPlayer
	trade.Gives = gives
	trade.Wants = wants
	trade.HasAllCards = hasAllCards
	trade.OK = 0

	return nil
}

func (context *GameContext) acceptTrade(tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 {
		return errors.New("invalid operation")
	}
	trade.OK = 1
	return nil
}

func (context *GameContext) rejectTrade(tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 0 {
		return errors.New("invalid operation")
	}
	trade.OK = -1
	return nil
}

func (context *GameContext) completeTrade(tradeID int) error {
	trade := context.getTrade(tradeID)
	if trade == nil || trade.OK != 1 {
		return errors.New("invalid operation")
	}

	trade.OK = 2

	player1 := context.Players[trade.From]
	player2 := context.Players[trade.To]

	for _, card := range trade.Gives {
		player1.Cards[card[0]] -= card[1]
		player2.Cards[card[0]] += card[1]
	}

	for _, card := range trade.Wants {
		player1.Cards[card[0]] += card[1]
		player2.Cards[card[0]] -= card[1]
	}

	return nil
}

func (context GameContext) isSafeTrade(gives [][2]int, wants [][2]int) bool {
	if context.Phase != Phase4 {
		return false
	}
	gl := len(gives)
	wl := len(wants)
	if gl == 0 || wl == 0 || wl > 5 || gl > 5 {
		return false
	}
	return true
}

func (context GameContext) isPlayerHasAllCards(playerID int, cards [][2]int) bool {
	player := context.Players[playerID]
	for _, giveCard := range cards {
		giveCardType := giveCard[0]
		giveTradeCount := giveCard[1]
		if giveTradeCount <= 0 || player.Cards[giveCardType] <= giveTradeCount {
			return false
		}
	}
	return true
}

func (context *GameContext) startPhase2() error {
	if context.Phase != Phase1 {
		return errors.New("invalid operation")
	}
	context.Phase = Phase2
	context.CurrentPlayer = 0
	context.scheduleAction(ActionPlaceSettlement)
	return nil
}

func (context *GameContext) startPhase3() {
	context.Phase = Phase3
	context.scheduleAction(ActionPlaceSettlement)
}

func (context *GameContext) startPhase4() {
	context.Phase = Phase4
	context.CurrentPlayer = 0
	context.scheduleAction(ActionRollDice)
}

func (context *GameContext) phase2GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 0 && roadCount == 0 {
		return ActionPlaceSettlement
	}
	if settlementCount == 1 && roadCount == 0 {
		return ActionPlaceRoad
	}
	if settlementCount == 1 && roadCount == 1 {
		return ""
	}
	return ""
}

func (context *GameContext) phase3GetNextAction() string {
	currentPlayer := context.getCurrentPlayer()
	settlementCount := len(currentPlayer.Settlements)
	roadCount := len(currentPlayer.Roads)

	if settlementCount == 1 && roadCount == 1 {
		return ActionPlaceSettlement
	}
	if settlementCount == 2 && roadCount == 1 {
		return ActionPlaceRoad
	}
	return ""
}

func (context *GameContext) endAction() error {
	fmt.Println("END", context.getActionString(), context.CurrentPlayer)

	NumberOfPlayers := context.GameSetting.NumberOfPlayers - 1

	if Phase4 == context.Phase {
		//clean up trades
		context.trades = []*gameTrade{}

		lastAction := context.getActionString()

		if lastAction == ActionDiscardCards {
			context.scheduleAction(ActionPlaceRobber)
			return nil
		}

		if lastAction == ActionPlaceRobber {
			context.scheduleAction(ActionSelectToRob)
			return nil
		}

		if lastAction == ActionSelectToRob || lastAction == ActionRollDice {
			context.scheduleAction(ActionTurn)
			return nil
		}

		if lastAction == ActionTurn {
			if context.CurrentPlayer < NumberOfPlayers {
				context.CurrentPlayer++
			} else {
				context.CurrentPlayer = 0
			}
			context.scheduleAction(ActionRollDice)
		}
		return nil
	}

	if Phase2 == context.Phase {
		nextAction := context.phase2GetNextAction()
		if nextAction == "" && context.CurrentPlayer < NumberOfPlayers {
			context.CurrentPlayer++
			nextAction = context.phase2GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayer == NumberOfPlayers {
			context.startPhase3()
		} else {
			context.scheduleAction(nextAction)
		}
	}

	if Phase3 == context.Phase {
		nextAction := context.phase3GetNextAction()
		if nextAction == "" && context.CurrentPlayer > 0 {
			context.CurrentPlayer--
			nextAction = context.phase3GetNextAction()
		}
		if nextAction == "" && context.CurrentPlayer == 0 {
			context.startPhase4()
		} else {
			context.scheduleAction(nextAction)
		}
	}
	return nil
}

func (context *GameContext) isActionTimeout() bool {
	action := context.Action

	timeout := 0
	if context.GameSetting.Speed == 0 {
		switch action.Name {
		case ActionTurn:
			timeout = 30
		case ActionRollDice:
			timeout = 10
		case ActionPlaceSettlement:
			if context.Phase == Phase2 || context.Phase == Phase3 {
				timeout = 12
			}
		case ActionPlaceRoad:
			if context.Phase == Phase2 || context.Phase == Phase3 {
				timeout = 15
			}
		case ActionDiscardCards:
			timeout = 15
		case ActionPlaceRobber:
			timeout = 10
		case ActionSelectToRob:
			timeout = 10
		}
	}
	durationSeconds := time.Now().Sub(action.Timeout).Seconds()
	if int(durationSeconds) < timeout {
		return false
	}
	return true
}

func (context *GameContext) handleDice(dice int) error {
	if dice == 7 {
		context.scheduleAction(ActionDiscardCards)
		return nil
	}

	bank := context.Bank
	players := context.Players

	bank.Begin()
	defer bank.Commit()
	for _, player := range players {
		for _, settlement := range player.Settlements {
			for idx, token := range settlement.Tokens {
				if token == dice {
					terrain := context.tiles[settlement.Indices[idx]].Terrain
					count, err := context.Bank.Give(terrain, 1)
					if err != nil {
						bank.Rollback()
						return err
					}
					context.handOverCards(player, terrain, count)
				}
			}
		}
	}
	return nil
}

func (context *GameContext) randomPlaceInitialSettlement() {
	availableIntersections, _ := context.getPossibleSettlementLocations()
	selectedIntersection := availableIntersections[rand.Intn(len(availableIntersections))]
	context.putSettlement(false, selectedIntersection)
}

func (context *GameContext) randomPlaceInitialRoad() {
	availableRoads, _ := context.getPossibleRoads()
	selectedRoad := availableRoads[rand.Intn(len(availableRoads))]
	context.putRoad(false, selectedRoad)
}

func (context *GameContext) randomDiscardCards() {
	for _, player := range context.Players {
		if len(player.Cards) > 7 {
			fmt.Println(fmt.Sprintf("Player %d looses %d cards.", player.ID, len(player.Cards)/2))
		}
	}
}

func NewGameContext() *GameContext {
	gc := new(GameContext)
	gc.Phase = Phase1
	gc.Bank = NewBank()
	return gc
}

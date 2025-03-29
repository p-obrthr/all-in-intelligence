package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	gamePlay "server/gameplay"
	"sync"

	"github.com/gorilla/websocket"
)

type BaseMessage struct {
	Type string `json:"type"`
}

type Config struct {
	Type             string `json:"type"`
	Players          int    `json:"players"`
	BBAmount         int    `json:"bb_amount"`
	StartMoneyAmount int    `json:"start_money_amount"`
}

type ConfigSet struct {
	Type  string `json:"type"`
	IsSet bool   `json:"is_set"`
}

type GetId struct {
	Type     string `json:"type"`
	PlayerId int    `json:"player_id"`
}

type Move struct {
	Type     string `json:"type"`
	PlayerId int    `json:"player_id"`
	Action   string `json:"action"`
}

type RoundState struct {
	Type            string          `json:"type"`
	Pot             int             `json:"pot"`
	Board           []gamePlay.Card `json:"board"`
	Cards           []gamePlay.Card `json:"cards"`
	StatusName      string          `json:"status_name"`
	PlayerId        int             `json:"player_id"`
	CurrentPlayerId int             `json:"current_player_id"`
	Money           int             `json:"money"`
	MsgLog          []string        `json:"msg_log"`
}

type RoundSend struct {
	Type     string         `json:"type"`
	Round    gamePlay.Round `json:"round"`
	BBAmount int            `json:"bb_amount"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients    = make(map[*websocket.Conn]int)
	mu         sync.Mutex
	gameConfig *Config
	game       *gamePlay.Game
	playerId   int
)

func newConfigSet(isSet bool) ConfigSet {
	return ConfigSet{
		Type:  "config_set",
		IsSet: isSet,
	}
}

func newRoundState(round *gamePlay.Round, player *gamePlay.Player, playerId, currentPlayerId int, currentPlayerMoney int) RoundState {
	cards := []gamePlay.Card{}
	statusName := ""

	if player != nil {
		cards = player.Cards
		statusName = player.Status.TypeName
	}

	return RoundState{
		Type:            "RoundState",
		Pot:             round.Pot,
		Board:           round.Board,
		Cards:           cards,
		StatusName:      statusName,
		PlayerId:        playerId,
		CurrentPlayerId: currentPlayerId,
		Money:           currentPlayerMoney,
		MsgLog:          round.MsgLog,
	}
}

func runLLMPlayers() bool {
	gameChanged := false

	for {
		game.CheckPlay()
		round := &game.Rounds[game.CurrentRound]
		player := round.GetPlayerById(round.CurrentPlayerId)

		broadcastRound()
		gameChanged = true

		if !player.IsLLM {
			break
		}
	}

	return gameChanged
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("err ws upgrade:", err)
		return
	}
	defer ws.Close()

	mu.Lock()
	if game == nil {
		playerId = 1
	} else {
		playerId = game.AppendNewPlayer()
	}
	clients[ws] = playerId
	fmt.Printf("Player %d connected\n", playerId)
	mu.Unlock()

	configSet := newConfigSet(gameConfig != nil)
	configSetJSON, _ := json.Marshal(configSet)
	ws.WriteMessage(websocket.TextMessage, []byte(configSetJSON))

	if game != nil {
		round := &game.Rounds[game.CurrentRound]
		round.MsgLog = append(round.MsgLog, fmt.Sprintf("Player %d connected", playerId))
		// player := round.GetPlayerById(round.CurrentPlayerId)
		// fmt.Println("round player current: " + player.Name)
	}

	for {
		var msg json.RawMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			fmt.Printf("Player %d disconnected\n", playerId)
			break
		}

		fmt.Println("received msg:", string(msg))

		var base BaseMessage
		if err := json.Unmarshal(msg, &base); err != nil {
			continue
		}

		handleMessage(base.Type, msg)
	}
}

func handleMessage(msgType string, msg json.RawMessage) {
	switch msgType {
	case "config":
		handleConfigMessage(msg)
	case "move":
		handleMoveMessage(msg)
	default:
		fmt.Println("err msg unknown type:", msgType)
	}
}

func handleConfigMessage(rawMsg json.RawMessage) {
	var configMsg Config
	if err := json.Unmarshal(rawMsg, &configMsg); err != nil {
		fmt.Println("Error parsing config message:", err)
		return
	}

	mu.Lock()
	gameConfig = &configMsg
	mu.Unlock()

	if game == nil {
		game = gamePlay.NewGame(gameConfig.Players, gameConfig.StartMoneyAmount, gameConfig.BBAmount)
		runLLMPlayers()
	}

	broadcastRound()
}

func handleMoveMessage(rawMsg json.RawMessage) {
	var moveMsg Move
	if err := json.Unmarshal(rawMsg, &moveMsg); err != nil {
		fmt.Println("err parsing move message:", err)
		return
	}

	fmt.Printf("received: %s\n", moveMsg.Action)

	round := &game.Rounds[game.CurrentRound]
	player := round.GetPlayerById(round.CurrentPlayerId)

	if moveMsg.PlayerId != player.Id {
		fmt.Println("not current player")
		return
	}

	if game.ProcessPlayerAction(moveMsg.Action) {
		broadcastRound()
		runLLMPlayers()
	}

	broadcastRound()
}

func broadcastRound() {
	mu.Lock()
	defer mu.Unlock()

	round := game.Rounds[game.CurrentRound]

	for client, clientPlayerId := range clients {
		player := round.GetPlayerById(clientPlayerId)
		state := newRoundState(&round, player, clientPlayerId, round.CurrentPlayerId, player.Money)

		responseJSON, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			fmt.Println("err serializing JSON:", err)
			continue
		}

		if err := client.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("--> server runs on 8080...")
	http.ListenAndServe(":8080", nil)
}

package message

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

type Move struct {
	Type     string `json:"type"`
	PlayerId int    `json:"player_id"`
	Action   string `json:"action"`
}

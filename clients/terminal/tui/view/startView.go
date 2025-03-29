package view

import (
	"encoding/json"
	"fmt"
	"terminal/tui/message"

	tea "github.com/charmbracelet/bubbletea"
)

type StartView struct {
	gameStarted bool
}

func NewStartView() *StartView {
	return &StartView{}
}

func (v *StartView) Init() tea.Cmd {
	return nil
}

func (v StartView) Render() string {

	cardArt := `
       _____      _____      _____      _____
      |A .  |    |K .  |    |Q .  |    |J .  |
      | /.\ ||   | /.\ ||   | /.\ ||   | /.\ ||
      |(_._)||   |(_._)||   |(_._)||   |(_._)||
      |  |  ||   |  |  ||   |  |  ||   |  |  ||
      |____V||   |____V||   |____V||   |____V||
             |____V||   |____V||   |____V||   
             |____V||   |____V||   |____V||   
           .__  .__    .__         .__        __         .__  .__  .__                                   
    _____  |  | |  |   |__| ____   |__| _____/  |_  ____ |  | |  | |  | ___    ____   ____   ____  ____  
    \__  \ |  | |  |   |  |/    \  |  |/    \   __\/ __ \|  | |  | |  |/ ___\_/ __ \ /    \_/ ___\/ __ \ 
     / __ \|  |_|  |__ |  |   |  \ |  |   |  \  | \  ___/|  |_|  |_|  / /_/  >  ___/|   |  \  \__\  ___/ 
    (____  /____/____/ |__|___|  / |__|___|  /__|  \___  >____/____/__\___  / \___  >___|  /\___  >___  >
         \/                    \/          \/          \/            /_____/      \/     \/     \/    \/  
	`

	return "            Welcome to" + "\n" + cardArt + "\n" + "            Press enter to start\n\n"

}

func (v *StartView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			// fmt.Println(v.gameStarted)
			if v.gameStarted {
				config := message.Config{
					Type: "config",
				}
				return NewGameView(), SendConfigCmd(config)
			} else {
				return NewConfigView(), NewConfigView().Init()
			}

		case "esc":
			return NewStartView(), nil
		}
	}
	var cmd tea.Cmd
	return v, cmd
}

func (v *StartView) UpdateMessage(msg string) {
	var configSet message.ConfigSet
	err := json.Unmarshal([]byte(msg), &configSet)
	if err != nil {
		fmt.Println(fmt.Sprintf("err: %v", err))
		return
	}

	v.gameStarted = configSet.IsSet

	// fmt.Println("gameStarted:", v.gameStarted)
}

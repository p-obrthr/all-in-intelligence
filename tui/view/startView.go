package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StartView struct{}

func NewStartView() *StartView {
	return &StartView{}
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
    _____  |  | |  |   |__| ____   |__| _____/  |_  ____ |  | |  | |  | |__| ____   ____   ____   ____  
    \__  \ |  | |  |   |  |/    \  |  |/    \   __\/ __ \|  | |  | |  |/ ___\_/ __ \ /    \_/ ___\/ __ \ 
     / __ \|  |_|  |__ |  |   |  \ |  |   |  \  | \  ___/|  |_|  |_|  / /_/  >  ___/|   |  \  \__\  ___/ 
    (____  /____/____/ |__|___|  / |__|___|  /__|  \___  >____/____/__\___  / \___  >___|  /\___  >___  >
         \/                    \/          \/          \/            /_____/      \/     \/     \/    \/  
	`

	return "\n" + "            Welcome to" + "\n" + cardArt + "\n" + "            Press enter to start\n\n"

}

func (v StartView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			return NewConfigView(), NewConfigView().Init()
		case "esc":
			return NewStartView(), nil
		}
	}
	var cmd tea.Cmd
	return v, cmd
}

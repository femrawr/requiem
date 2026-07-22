package bot

import (
	"requiem/bot/buttons"
	"requiem/store"
	"strings"
)

var buttonssList = make(map[string]store.Button)

func registerButtons() {
	buttons := []store.Button{
		&buttons.TestButton{},
	}

	for _, button := range buttons {
		buttonssList[strings.Split(button.Iden(), ".")[1]] = button
	}
}

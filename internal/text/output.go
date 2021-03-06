package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/RasmusLindroth/i3keys/internal/i3parse"
	"github.com/RasmusLindroth/i3keys/internal/keyboard"
	"github.com/RasmusLindroth/i3keys/internal/xlib"
)

func printKeyboards(keyboards []keyboard.Keyboard, groups []i3parse.ModifierGroup, prefix string) {
	for kbIndex, kb := range keyboards {
		dots := "-"
		for i := 0; i < len(kb.Name); i++ {
			dots = dots + "-"
		}

		fmt.Printf("%s%s:\n%s%s\n", prefix, kb.Name, prefix, dots)

		for _, keyRow := range kb.Keys {
			var unused []string
			for _, key := range keyRow {
				if key.InUse == false {
					unused = append(unused, key.Symbol)
				}
			}
			unusedStr := strings.Join(unused, ", ")
			fmt.Printf("%s%s\n", prefix, unusedStr)
		}
		if kbIndex+1 < len(groups) {
			fmt.Printf("\n\n")
		}
	}
}

//Output prints the keyboards to os.Stdout
func Output(layout string) {
	modes, keys, err := i3parse.ParseFromRunning()

	if err != nil {
		log.Fatalln(err)
	}

	layout = strings.ToUpper(layout)

	keys = i3parse.KeysToSymbol(keys)

	groups := i3parse.GetModifierGroups(keys)
	modifiers := xlib.GetModifiers()

	var keyboards []keyboard.Keyboard
	for _, group := range groups {
		kb, err := keyboard.MapKeyboard(layout, group, modifiers)
		if err == nil {
			keyboards = append(keyboards, kb)
		}
	}

	fmt.Printf("Available keybindings per modifier group\n\n")
	printKeyboards(keyboards, groups, "")

	for _, mode := range modes {
		keys := i3parse.KeysToSymbol(mode.Bindings)
		groups := i3parse.GetModifierGroups(keys)
		var mKeyboards []keyboard.Keyboard

		for _, group := range groups {
			kb, err := keyboard.MapKeyboard(layout, group, modifiers)
			if err == nil {
				mKeyboards = append(mKeyboards, kb)
			}
		}

		fmt.Printf("\n\nMode: %s\n", mode.Name)
		printKeyboards(mKeyboards, groups, "\t")
	}
}

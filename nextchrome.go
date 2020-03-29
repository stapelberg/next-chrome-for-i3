// Program next-chrome-for-i3 focuses the chrome window on the current
// workspace or starts a new chrome instance.
package main

import (
	"fmt"
	"log"
	"strings"

	"go.i3wm.org/i3/v4"
)

func logic() error {
	tree, err := i3.GetTree()
	if err != nil {
		return err
	}

	ws := tree.Root.FindFocused(func(n *i3.Node) bool { return n.Type == i3.WorkspaceNode })
	if ws == nil {
		return fmt.Errorf("could not locate workspace")
	}

	if chrome := ws.FindChild(func(n *i3.Node) bool { return strings.HasSuffix(n.Name, "- Google Chrome") }); chrome != nil {
		_, err = i3.RunCommand(fmt.Sprintf(`[con_id="%d"] focus`, chrome.ID))
	} else {
		_, err = i3.RunCommand(`exec google-chrome`)
	}

	return err
}

func main() {
	if err := logic(); err != nil {
		log.Fatal(err)
	}
}

// Program next-chrome-for-i3 focuses the chrome window on the current
// workspace or starts a new chrome instance.
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"go.i3wm.org/i3/v4"
)

func logic() error {
	var (
		titleExpr = flag.String(
			"title_regexp",
			"- Google Chrome$",
			"Go regular expression (https://golang.org/pkg/regexp) that will be matched on the window title")

		cmd = flag.String(
			"not_found_command",
			"exec google-chrome",
			"i3 command to run if no window matching -title_regexp is found")

		scope = flag.String(
			"scope",
			"workspace",
			"workspace or root, specifies which child windows to match")
	)
	flag.Parse()

	titleRe, err := regexp.Compile(*titleExpr)
	if err != nil {
		return err
	}

	tree, err := i3.GetTree()
	if err != nil {
		return err
	}

	var parent *i3.Node
	if *scope == "workspace" {
		parent = tree.Root.FindFocused(func(n *i3.Node) bool { return n.Type == i3.WorkspaceNode })
		if parent == nil {
			return fmt.Errorf("could not locate workspace")
		}
	} else {
		parent = tree.Root
	}

	if chrome := parent.FindChild(func(n *i3.Node) bool { return titleRe.MatchString(n.Name) }); chrome != nil {
		_, err = i3.RunCommand(fmt.Sprintf(`[con_id="%d"] focus`, chrome.ID))
	} else {
		_, err = i3.RunCommand(*cmd)
	}

	return err
}

func main() {
	if err := logic(); err != nil {
		log.Fatal(err)
	}
}

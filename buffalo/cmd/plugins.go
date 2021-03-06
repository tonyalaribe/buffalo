package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

var plugx = &sync.Mutex{}
var _plugs plugins.List

func plugs() plugins.List {
	plugx.Lock()
	defer plugx.Unlock()
	if _plugs == nil {
		var err error
		_plugs, err = plugins.Available()
		if err != nil {
			_plugs = plugins.List{}
			log.Printf("error loading plugins %s\n", err)
		}
	}
	return _plugs
}

func decorate(name string, cmd *cobra.Command) {
	for _, c := range plugs()[name] {
		func(c plugins.Command) {
			cc := &cobra.Command{
				Use:   c.Name,
				Short: fmt.Sprintf("[PLUGIN] %s", c.Description),
				RunE: func(cmd *cobra.Command, args []string) error {
					ax := []string{c.Name}
					ax = append(ax, args...)
					ex := exec.Command(c.Binary, ax...)
					ex.Stdin = os.Stdin
					ex.Stdout = os.Stdout
					ex.Stderr = os.Stderr
					return ex.Run()
				},
			}
			cc.DisableFlagParsing = true
			cmd.AddCommand(cc)
		}(c)
	}
}

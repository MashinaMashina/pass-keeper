package clibell

import (
	"github.com/chzyer/readline"
	"io"
	"sync"
)

type noBellStdout struct {
	Stdout io.WriteCloser
}

func (n *noBellStdout) Write(p []byte) (int, error) {
	readline.Stdout = n.Stdout
	if len(p) == 1 && p[0] == readline.CharBell {
		return 0, nil
	}
	return readline.Stdout.Write(p)
}

func (n *noBellStdout) Close() error {
	return readline.Stdout.Close()
}

var once sync.Once
var bell *noBellStdout

func Instance(stdout io.WriteCloser) *noBellStdout {
	once.Do(func() {
		bell = &noBellStdout{
			stdout,
		}
	})

	return bell
}

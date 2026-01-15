package flag

import (
	"flag"
	"os"
	"sync"
)

var (
	FS       = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	once     sync.Once
	errParse error
	visited  map[string]bool
)

// Parse парсит os.Args один раз и запоминает, какие флаги реально передали
func Parse() error {
	once.Do(func() {
		errParse = FS.Parse(os.Args[1:])
		visited = map[string]bool{}
		FS.Visit(func(f *flag.Flag) {
			visited[f.Name] = true
		})
	})
	return errParse
}

// WasSet возвращает true, если флаг реально был передан при запуске
func WasSet(name string) bool {
	if visited == nil {
		return false
	}
	return visited[name]
}

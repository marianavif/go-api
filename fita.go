package poker

import (
	"os"
)

type fita struct {
	arquivo *os.File
}

func (f *fita) Write(p []byte) (n int, err error) {
	f.arquivo.Truncate(0)
	f.arquivo.Seek(0, 0)
	return f.arquivo.Write(p)
}

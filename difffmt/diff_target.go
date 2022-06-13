package difffmt

import (
	"io/ioutil"
	"os"
	"time"
)

type DiffTarget struct {
	Path         string
	ModifiedTime time.Time
}

func NewDiffTarget(path string) *DiffTarget {
	return &DiffTarget{
		Path: path,
	}
}

func (t *DiffTarget) LoadStats() error {
	file, err := os.Stat(t.Path)
	if err != nil {
		return err
	}
	t.ModifiedTime = file.ModTime()
	return nil
}

func (t *DiffTarget) ReadText() (string, error) {
	text, err := ioutil.ReadFile(t.Path)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

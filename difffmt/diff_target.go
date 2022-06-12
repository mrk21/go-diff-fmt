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

func NewDiffTarget(path string) (*DiffTarget, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &DiffTarget{
		Path:         path,
		ModifiedTime: file.ModTime(),
	}, nil
}

func (i *DiffTarget) ReadText() (string, error) {
	text, err := ioutil.ReadFile(i.Path)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

package gen

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Generator interface {
	Render(*sqlx.DB) ([]string, error)
	GetFileName() string
}

func WriteToFile(g Generator, db *sqlx.DB) error {
	p, err := os.Getwd()
	if err != nil {
		return err
	}
	ouf := p + "/" + os.Getenv("QGEN_PATH") + "/" + g.GetFileName()
	if strings.HasPrefix(os.Getenv("QGEN_PATH"), "/") {
		ouf = os.Getenv("QGEN_PATH") + "/" + g.GetFileName()
	}
	f, err := os.Create(path.Clean(ouf))
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := g.Render(db)
	if err != nil {
		return err
	}
	for i := range r {
		// это лечит левую кодировку которую возвращает база mysql
		s := string([]rune(r[i]))
		n, err2 := f.Write([]byte(s + "\n"))
		if err2 != nil {
			return err2
		}
		if n != len(s)+1 {
			return fmt.Errorf("short write")
		}
	}
	return err
}

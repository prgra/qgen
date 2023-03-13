package gen

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type Generator interface {
	Render(*sqlx.DB) ([]string, error)
	GetFileName() string
}

func WriteToFile(g Generator, db *sqlx.DB) error {
	f, err := os.Create(g.GetFileName())
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

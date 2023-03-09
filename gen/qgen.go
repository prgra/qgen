package gen

import (
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
		_, err = f.Write([]byte(r[i] + "\n"))
	}
	return err
}

package gen

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"golang.org/x/text/encoding/charmap"
)

type Generator interface {
	Render(*sqlx.DB, config.Config) ([]string, error)
	GetFileName() string
	GetRemoteDir() string
}

func WriteToFile(g Generator, cfg config.Config, db *sqlx.DB) error {
	p, err := os.Getwd()
	if err != nil {
		return err
	}
	ouf := p + "/" + cfg.Path + "/" + g.GetFileName()
	if strings.HasPrefix(cfg.Path, "/") {
		ouf = cfg.Path + "/" + g.GetFileName()
	}
	f, err := os.Create(path.Clean(ouf))
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := g.Render(db, cfg)
	if err != nil {
		return err
	}
	for i := range r {
		// это лечит левую кодировку которую возвращает база mysql
		s := string([]rune(r[i]))
		if cfg.CSVChatset == "windows-1251" {
			enc := charmap.Windows1251.NewEncoder()
			s, err = enc.String(s)
			if err != nil {
				fmt.Println("encode", err, r[i])
			}
		}
		n, err2 := f.WriteString(s + "\n")
		if err2 != nil {
			return err2
		}
		if n != len(s)+1 {
			return fmt.Errorf("short write")
		}
	}
	return err
}

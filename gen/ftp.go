package gen

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"golang.org/x/text/encoding/charmap"
)

func UploadToFTP(g Generator, cfg config.Config, db *sqlx.DB) error {
	r, err := g.Render(db, cfg)
	if err != nil {
		return err
	}
	c, err := ftp.Dial(cfg.FTPHost, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(cfg.FTPUser, cfg.FTPPass)
	if err != nil {
		return err
	}
	defer c.Quit()
	log.Println("uploading", g.GetRemoteDir()+"/"+g.GetFileName())
	reader, writer := io.Pipe()
	defer reader.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		err = c.Stor(g.GetRemoteDir()+"/"+g.GetFileName(), reader)
		if err != nil {
			log.Println("upload", err)
		}
		wg.Done()
	}(&wg)

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
		if writer == nil {
			n, err2 := writer.Write([]byte(s + "\n"))
			if err2 != nil {
				return err2
			}
			if n != len(s)+1 {
				return fmt.Errorf("short write")
			}
		}
	}

	writer.Close()
	wg.Wait()
	return nil
}

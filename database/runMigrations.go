package database

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RunMigrations(db *Database, directories ...string) error {
	for _, dir := range directories {
		err := runMigrations(db, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMigrations(db *Database, dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if len(path) > 4 && path[len(path)-4:] == ".sql" {
			readFile, err := os.Open(path)
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(readFile)

			scanner.Split(bufio.ScanLines)

			q := ""
			for scanner.Scan() {
				text := scanner.Text()
				if text == "" {
					continue
				}

				if text[0:2] == "--" {
					continue
				}

				re, err := regexp.Compile(`/[*].*?[*]/`)
				if err != nil {
					return err
				}

				text = string(re.ReplaceAll([]byte(text), []byte("")))
				if text == ";" {
					continue
				}

				q += " " + text
				if strings.Contains(text, ";") {
					_, err = db.Db.Exec(q)
					if err != nil {
						return err
					}
					q = ""
				}

			}
		}
		return nil
	})
}

package brunchfeed

import (
	"log"
	"os"
	"path/filepath"

	"github.com/linterpreteur/brunchfeed/lib"
	"github.com/vmihailenco/msgpack"
)

const dateFormat = "2006-01-02T15:04:05-07:00"

type FetchParams struct {
	Id  string
	Src string
}

func Fetch(params FetchParams) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	src := params.Src
	id := params.Id

	if len(id) == 0 {
		log.Fatal("Command line argument -id should not be empty")
	}
	if len(src) == 0 {
		log.Fatal("Command line argument -src should not be empty")
	}

	raw, err := brunchfeed.Fetch(id)
	if err != nil {
		log.Fatal(err)
	}

	for _, raw := range raw.Items {
		item, err := brunchfeed.FullItem(raw)
		if err != nil {
			log.Fatal(err)
		}

		f := path(src, item.Meta.Slug)
		if _, err := os.Stat(f); os.IsNotExist(err) {
			save(f, item)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func path(base string, slug string) string {
	absPath, err := filepath.Abs(base)
	if err != nil {
		log.Fatal(err)
	}

	f := filepath.Join(absPath, slug+".md")
	return f
}

func save(path string, item brunchfeed.Item) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := msgpack.Marshal(&item)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
}

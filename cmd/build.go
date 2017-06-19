package brunchfeed

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/linterpreteur/brunchfeed/lib"
	"github.com/vmihailenco/msgpack"
)

type BuildParams struct {
	Src      string
	Dest     string
	Template string
}

func Build(params BuildParams) {
	src := params.Src
	dest := params.Dest
	template := params.Template
	flag.Parse()

	if len(src) == 0 {
		log.Fatal("Command line argument -src should not be empty")
	}
	if len(dest) == 0 {
		log.Fatal("Command line argument -dest should not be empty")
	}
	if len(template) == 0 {
		log.Fatal("Command line argument -template should not be empty")
	}

	root, err := filepath.Abs(src)
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if root == path || f.Name() == ".gitkeep" {
			return nil
		}

		item, err := load(path)
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(dest)
		if err != nil {
			log.Fatal(err)
		}

		file := filepath.Join(absPath, item.Meta.Slug+".md")
		err = savePost(file, template, item)
		if err != nil {
			return err
		}

		return nil
	})
}

const hugoDateFormat = "2006-01-02"

func savePost(path string, template string, item brunchfeed.Item) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	content := inject(template, item)

	defer f.Close()
	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func inject(template string, item brunchfeed.Item) string {
	return replaceAll(template, map[string]string{
		"title":     item.Body.Title,
		"date":      item.Meta.Date.Format(hugoDateFormat),
		"thumbnail": item.Meta.Thumbnail,
		"category":  item.Meta.Category,
		"tags":      "\"" + strings.Join(item.Meta.Tags, "\", \"") + "\"",
		"summary":   item.Body.Summary,
		"content":   item.Body.Content,
		"link":      item.Meta.Link,
	})
}

func replaceAll(s string, m map[string]string) string {
	out := s
	for k, v := range m {
		out = replace(out, k, v)
	}
	return out
}

func replace(s string, src string, repl string) string {
	return regexp.MustCompile("<%\\s*"+src+"\\s*%>").ReplaceAllString(s, repl)
}

func load(path string) (brunchfeed.Item, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return brunchfeed.Item{}, err
	}

	var item brunchfeed.Item
	err = msgpack.Unmarshal(data, &item)
	if err != nil {
		return brunchfeed.Item{}, err
	}

	return item, nil
}

package brunchfeed

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/SlyMarbo/rss"
)

type Meta struct {
	Thumbnail string
	Category  string
	Tags      []string
	Slug      string
	Link      string
	Date      time.Time
}

type Body struct {
	Title   string
	Summary string
	Content string
}

type Item struct {
	Meta Meta
	Body Body
}

func Fetch(id string) (*rss.Feed, error) {
	raw, err := rss.Fetch("https://brunch.co.kr/rss/@@" + id)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func FullItem(raw *rss.Item) (Item, error) {
	link := raw.Link
	doc, err := goquery.NewDocument(link)
	if err != nil {
		return Item{}, err
	}

	thumb := thumbnail(doc)
	content, tags, err := fullContent(link, doc)

	if err != nil {
		return Item{}, err
	}

	return Item{
		Meta: Meta{
			Thumbnail: thumb,
			Category:  tags[0],
			Tags:      tags[1:],
			Slug:      slugify(raw.Title),
			Link:      link,
			Date:      raw.Date,
		},
		Body: Body{
			Title:   raw.Title,
			Content: content,
		},
	}, nil
}

func thumbnail(doc *goquery.Document) string {
	cover := doc.Find(".cover_image")
	style, _ := cover.Attr("style")
	url := "http:" + r(style, "background-image:", "url(", ")", ";")
	return url
}

func r(s string, rs ...string) string {
	if len(rs) == 0 {
		return s
	}
	return r(strings.Replace(s, rs[0], "", 1), rs[1:]...)
}

func fullContent(sourceURL string, doc *goquery.Document) (string, []string, error) {
	body := doc.Find(".wrap_article .wrap_body").First()
	article, err := body.Html()
	if err != nil {
		return "", nil, err
	}

	tags := doc.Find(".list_keyword a, .title_magazine").Map(func(_ int, v *goquery.Selection) string {
		return v.Text()
	})

	r := regexp.MustCompile("data-[a-z]+=\"[^\"]+\"")
	article = string(r.ReplaceAllString(article, ""))
	article = strings.Replace(article, "\"//", "\"https://", -1)

	return article, tags, nil
}

func slugify(raw string) string {
	nonletters := regexp.MustCompile("[^ㄱ-ㅣ가-힣0-9a-zA-Z ]")
	onlyLetters := string(nonletters.ReplaceAllString(raw, ""))
	whitespaces := regexp.MustCompile("\\s+")
	normalized := string(whitespaces.ReplaceAllString(onlyLetters, "-"))
	return normalized
}

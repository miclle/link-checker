package checker

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Link element struct
type Link struct {
	Href       string
	Text       string
	Referer    string
	Status     string
	StatusCode int
}

// Checker struct
type Checker struct {
	*http.Client
	TargetURL string
	host      string
	Depth     int
	queue     map[string][]*Link
}

// NewChecker return Checker
func NewChecker(targetURL string, depth int) (*Checker, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return &Checker{
		Client:    http.DefaultClient,
		TargetURL: targetURL,
		host:      u.Host,
		Depth:     depth,
		queue:     map[string][]*Link{},
	}, nil
}

// Checking the url
func (c *Checker) Checking() (err error) {
	return c.walk(c.TargetURL)
}

// walk the url
func (c *Checker) walk(url string) (err error) {

	resp, err := c.Get(url)
	if err != nil {
		return
	}

	log.Println("resp.Status:", resp.StatusCode, resp.Status)

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	log.Println("title:", doc.Find("title").Text())

	if c.host == resp.Request.Host {

		// Find the a elements
		doc.Find("a").Each(func(i int, a *goquery.Selection) {
			if href, exists := a.Attr("href"); exists {

				link := &Link{
					Href:    href,
					Text:    a.Text(),
					Referer: url,
				}

				c.queue[link.Href] = append(c.queue[link.Href], link)
			}
		})
	}
	return
}

// Check the url
func Check(address string, depth int) (err error) {
	checker, err := NewChecker(address, depth)
	if err != nil {
		return err
	}
	return checker.Checking()
}

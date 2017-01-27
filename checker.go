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
	Referers   LinkDictionary
	Status     string
	StatusCode int
}

// LinkDictionary linkddictionary
type LinkDictionary map[string]*Link

// AddReferer add referer to link
func (l *Link) AddReferer(link *Link) {
	if l.Referers == nil {
		l.Referers = LinkDictionary{}
	}
	l.Referers[link.Href] = link
}

// Checker struct
type Checker struct {
	*http.Client
	TargetURL string
	host      string
	Depth     int
	queue     LinkDictionary
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
		queue:     LinkDictionary{},
	}, nil
}

// Checking the url
func (c *Checker) Checking() (err error) {
	link := &Link{
		Href: c.TargetURL,
	}

	return c.walk(link)
}

// walk the url
func (c *Checker) walk(link *Link) (err error) {

	resp, err := c.Get(link.Href)
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

				internalLink := &Link{
					Href: href,
					Text: a.Text(),
				}

				internalLink.AddReferer(link)
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

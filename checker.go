package checker

import (
	"net/http"
	"net/url"

	"log"

	"github.com/PuerkitoBio/goquery"
)

// Link element struct
type Link struct {
	Href       string
	URL        *url.URL
	Text       string
	Referers   LinkDictionary
	Status     string
	StatusCode int
	PageTitle  string
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
	base      *url.URL
	Depth     int
	queue     LinkDictionary
}

// NewChecker return Checker
func NewChecker(targetURL string, depth int) (checker *Checker, err error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return &Checker{
		Client:    http.DefaultClient,
		TargetURL: targetURL,
		base:      u,
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
	u, err := url.Parse(link.Href)
	if err != nil {
		return err
	}

	link.URL = c.base.ResolveReference(u)

	resp, err := c.Get(link.URL.String())
	if err != nil {
		log.Println("http get error", link.Href, err.Error())
		return
	}

	link.Status = resp.Status
	link.StatusCode = resp.StatusCode

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	link.PageTitle = doc.Find("title").Text()

	log.Printf("link: %#v", link)

	if c.base.Host == resp.Request.Host {
		// Find the a elements
		doc.Find("a").Each(func(i int, a *goquery.Selection) {
			if href, exists := a.Attr("href"); exists {
				internalLink := &Link{
					Href: href,
					Text: a.Text(),
				}
				internalLink.AddReferer(link)

				if _, e := c.queue[internalLink.Href]; e == false {
					c.queue[internalLink.Href] = internalLink
					c.walk(internalLink)
				}
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

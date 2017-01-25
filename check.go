package checker

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Checker the url
func Checker(url string, depth int) (err error) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	// Find the review items
	doc.Find("a").Each(func(i int, a *goquery.Selection) {

		fmt.Printf("Review %s \t: %s\n", a.Text(), a.AttrOr("href", ""))
	})

	return
}

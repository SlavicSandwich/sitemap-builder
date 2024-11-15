package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	parser "sitemap_builder"
	"strings"
)

type loc struct {
	Loc string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func main() {
	urlLink := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the depth of the sitemap")
	flag.Parse()

	pages := bfs(*maxDepth, *urlLink)
	//fmt.Println(pages)
	var locs []loc
	for _, page := range pages {
		locs = append(locs, loc{page})
	}
	fmt.Print(xml.Header)
	set := urlset{locs, xmlns}
	output, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		fmt.Println("error, %v\n", err)
	}
	os.Stdout.Write(output)
	/*
		1. GET the webpage
		2. Parse the webpage for links
		3. Build proper urls for other pages
		4. filter out links that are diff domain
		5. Find all pages(BFS)
		6. print out XML
	*/

}

func bfs(depth int, urlStr string) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= depth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	res, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	reqUrl := res.Request.URL
	baseURL := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseURL.String()
	return filter(hrefs(res.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := parser.Parse(r)
	var h []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			h = append(h, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			h = append(h, l.Href)
		}
	}
	return h
}

func filter(links []string, keepFN func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFN(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

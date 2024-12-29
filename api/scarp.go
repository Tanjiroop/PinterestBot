package api

import (
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func FetchWallpapers(query string) []map[string]string {
    imagesData := []map[string]string{}

    url := "https://wallpapers.com/search/"
    if query != "" {
        url = "https://wallpapers.com/search/" + query
    }

    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    totalPages := 1
    pageCounter := doc.Find(".page-counter.mobi").First().Text()
    if pageCounter != "" {
        fmt.Sscanf(strings.Split(pageCounter, " ")[len(strings.Split(pageCounter, " "))-1], "%d", &totalPages)
    }

    page := rand.Intn(totalPages) + 1

    url = "https://wallpapers.com/search/"
    if query != "" {
        url = "https://wallpapers.com/search/" + query
    }
    url += fmt.Sprintf("?p=%d", page)

    response, err = http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    doc, err = goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    doc.Find("li.content-card").Each(func(i int, s *goquery.Selection) {
        title := s.Find("a").AttrOr("title", "")
        imgURL := s.Find("img").AttrOr("data-src", "")
        if title != "" && imgURL != "" {
            imageURL := strings.Join(strings.Split(url, "/")[:len(strings.Split(url, "/"))-1], "/") + imgURL
            imagesData = append(imagesData, map[string]string{"title": title, "url": imageURL})
        }
    })

    return imagesData
}

package api

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "math/rand"
    "net/http"
    "strings"
)

type ImageData struct {
    Title string
    URL   string
}

func FetchWallpapers(query string) []ImageData {
    imagesData := []ImageData{}
    
    var url string
    if query == "" {
        url = "https://wallpapers.com/search/anime"
    } else {
        url = "https://wallpapers.com/search/" + quote(query)
    }

    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching the URL:", err)
        return imagesData
    }
    defer response.Body.Close()

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        fmt.Println("Error loading the HTML:", err)
        return imagesData
    }

    totalPages := 1
    pageCounter := doc.Find(".page-counter.mobi")
    if len(pageCounter.Nodes) > 0 {
        totalPages, _ = strconv.Atoi(strings.Fields(pageCounter.Text())[len(strings.Fields(pageCounter.Text()))-1])
    }

    page := rand.Intn(totalPages) + 1
    pageURL := url + "?p=" + strconv.Itoa(page)

    response, err = http.Get(pageURL)
    if err != nil {
        fmt.Println("Error fetching the page URL:", err)
        return imagesData
    }
    defer response.Body.Close()

    doc, err = goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        fmt.Println("Error loading the HTML:", err)
        return imagesData
    }

    doc.Find("li.content-card").Each(func(i int, s *goquery.Selection) {
        aTag := s.Find("a")
        imgTag := s.Find("img")
        title := aTag.AttrOr("title", "")
        imgURL := imgTag.AttrOr("data-src", "")
        if title != "" && imgURL != "" {
            imageURL := strings.Join(strings.Split(pageURL, "/")[:len(strings.Split(pageURL, "/"))-1], "/") + imgURL
            imagesData = append(imagesData, ImageData{Title: title, URL: imageURL})
        }
    })

    return imagesData
}


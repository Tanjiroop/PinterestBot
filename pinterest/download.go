package pinterest

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
	
	"golang.org/x/net/html"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func ExtractImageUrl(pinterestUrl string) (string, error) {
	resp, err := http.Get(pinterestUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error fetching content from URL: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var imageTags []*html.Node
	var found bool

	finden:
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "img" {
			for _, a := range c.Attr {
				if a.Key == "class" && (a.Val == "h-image-fit" || a.Val == "h-unsplash-img") {
					imageTags = append(imageTags, c)
					found = true
					break finden
				}
			}
		}
	}

	if !found {
		for c := doc.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "img" {
				for _, a := range c.Attr {
					if a.Key == "src" && strings.HasPrefix(a.Val, "https://i.pinimg.com/") {
						imageTags = append(imageTags, c)
						break
					}
				}
			}
		}
	}

	if len(imageTags) > 0 {
		return imageTags[0].Attr[0].Val, nil
	}

	return "", nil
}

func DownloadSend(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	chk := message.Text
	if strings.HasPrefix(chk, "/") {
		return nil
	}
	pattern := regexp.MustCompile(`https://pin\.it/(?P<url>[\w]+)`)
	if !pattern.MatchString(chk) {
		return nil
	}

	link, _ := ExtractImageUrl(chk)
	photo := gotgbot.InputMediaPhoto{			
		Media: gotgbot.InputFileByURL(link),
	}
	_, err := b.SendPhoto(ctx.EffectiveChat.Id, photo.Media, &gotgbot.SendPhotoOpts{})
	if err != nil {
		message.Reply(b, "Failed to Send Photo", &gotgbot.SendMessageOpts{})
		return err
	}
	return nil
}

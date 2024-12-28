package pinterest

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type PinterestResponse struct {
    Data []struct {
        URL string `json:"url"` 
    } `json:"data"` 
}

func searchPinterest(query string) (PinterestResponse, error) {
    url := fmt.Sprintf("https://horridapi.onrender.com/pinterest?query=%s", query)
    resp, err := http.Get(url)
    if err != nil {
        return PinterestResponse{}, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return PinterestResponse{}, err
    }

    var result PinterestResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        return PinterestResponse{}, err
    }

    return result, nil
}


func FindImage(b *gotgbot.Bot, ctx *ext.Context) error {
    message := ctx.Message
    query := strings.TrimSpace(strings.Replace(message.Text, "/h", "", -1))
    urls, err := searchPinterest(query)
    if err != nil {
        fmt.Println(err)
        message.Reply(b, "Image not found", &gotgbot.SendMessageOpts{})
        return err
    }

    media := make([]gotgbot.InputMedia, 0)
    for _, item := range urls.Data {
        fmt.Printf("Found image URL: %s\n", item.URL) // Debugging line to print the URL
        if item.URL != "" { // Check if URL is not empty
            media = append(media, gotgbot.InputMediaPhoto{
                Media: gotgbot.InputFileByURL(item.URL), // Change to InputFileByURL
            })
        } else {
            fmt.Println("Skipped empty URL") // Warn about empty URLs
        }
    }

    if len(media) == 0 {
        message.Reply(b, "No media to send", &gotgbot.SendMessageOpts{})
        return fmt.Errorf("no valid media found to send")
    }

    _, err = b.SendMediaGroup(
        ctx.EffectiveUser.Id,
        media,
        &gotgbot.SendMediaGroupOpts{},
    )
    if err != nil {
        fmt.Printf("Error sending media group: %s\n", err) // Log the error
        return err
    }

    return nil
}

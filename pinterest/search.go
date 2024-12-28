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
    query := strings.SplitN(message.GetText(), " ", 2)
    message.Reply(b, query, &gotgbot.SendMessageOpts{})
	if len(query) < 2 {     
        message.Reply(b, "No Query Provied So i can't send Photo, so Please Provide Query Eg: /pinterest Iron man", &gotgbot.SendMessageOpts{})
        return fmt.Errorf("no query provided")
    }

    urls, err := searchPinterest(query)
    if err != nil {
        fmt.Println(err)
        message.Reply(b, "Image not found", &gotgbot.SendMessageOpts{})
        return err
    }
    
    media := make([]gotgbot.InputMedia, 0)
    for _, item := range urls.Data {
        fmt.Printf("Found image URL: %s\n", item.URL)
        if item.URL != "" { 
            media = append(media, gotgbot.InputMediaPhoto{
                Media: gotgbot.InputFileByURL(item.URL),
            })
        } else {
            fmt.Println("Skipped empty URL")
        }
    }
   
    if len(media) == 0 {
        message.Reply(b, "No media to send", &gotgbot.SendMessageOpts{})
        return fmt.Errorf("no valid media found to send")
    }
 
    for i := 0; i < len(media); i += 10 {
        end := i + 10
        if end > len(media) {
            end = len(media)
        }

        batch := media[i:end]        
        
        _, err = b.SendMediaGroup(
            ctx.EffectiveUser.Id,
            batch,
            &gotgbot.SendMediaGroupOpts{},
        )
        if err != nil {
            fmt.Printf("Error sending media group: %s\n", err)
            message.Reply(b, "Error sending media group, please try again later.", &gotgbot.SendMessageOpts{})
            return err
        }
      
    }

    return nil
}

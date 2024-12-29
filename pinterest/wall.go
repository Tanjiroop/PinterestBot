package pinterest

import (    
    "fmt"    
    "strings"    
    
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/Mishel-07/PinterestBot/api"
)

func WallSearch(b *gotgbot.Bot, ctx *ext.Context) error {
    message := ctx.Message
    split := strings.SplitN(message.GetText(), " ", 2)            
    if len(split) < 2 {     
        message.Reply(b, "No Query Provied So i can't send Photo, so Please Provide Query Eg: <code>/pinterest Iron man</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
        return fmt.Errorf("no query provided")
    }

    query := split[1]
    quotequery := strings.Replace(query, " ", "+", -1)
    images := api.FetchWallpapers(quotequery)
    
    media := make([]gotgbot.InputMedia, 0)
    for _, item := range images {
        fmt.Printf("Found image URL: %s\n", item)
        media = append(media, gotgbot.InputMediaPhoto{
            Media: gotgbot.InputFileByURL(item),
            })        
    }
   
    if len(media) == 0 {
        message.Reply(b, "No Image found", &gotgbot.SendMessageOpts{})
        return fmt.Errorf("no valid media found to send")
    }
 
    for i := 0; i < len(media); i += 10 {
        end := i + 10
        if end > len(media) {
            end = len(media)
        }

        batch := media[i:end]        
        
        b.SendMediaGroup(
            ctx.EffectiveUser.Id,
            batch,
            &gotgbot.SendMediaGroupOpts{},
        )        
      
    }

    return nil
}

package app

import (
	"convert-song-tg-bot/pkg/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const pathRequest = "links"

type responceLink struct {
	EntityUniqueId     string `json:"entityUniqueId"`
	UserCountry        string `json:"userCountry"`
	PageUrl            string `json:"pageUrl"`
	EntitiesByUniqueId map[string]entitiesByUniqueId
	LinksByPlatform    map[string]linksByPlatform
}

type entitiesByUniqueId struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	Title           string   `json:"title"`
	ArtistName      string   `json:"artistName"`
	ThumbnailUrl    string   `json:"thumbnailUrl"`
	ThumbnailWidth  int      `json:"thumbnailWidth"`
	ThumbnailHeight int      `json:"thumbnailHeight"`
	ApiProvider     string   `json:"apiProvider"`
	Platforms       []string `json:"platforms"`
}

type linksByPlatform struct {
	Country             string `json:"country"`
	Url                 string `json:"url"`
	NativeAppUriMobile  string `json:"nativeAppUriMobile"`
	NativeAppUriDesktop string `json:"nativeAppUriDesktop"`
	EntityUniqueId      string `json:"entityUniqueId"`
}

func getSong(link string) (*Song, error) {
	hostName := config.GetEnv("SONGLINK_API_URL", "https://api.song.link/v1-alpha.1")
	respVars := &Song{}
	var structResp responceLink
	requestLink := url.Values{}
	requestLink.Add("url", link)
	urlRequest := fmt.Sprintf("%s/%s?%s", hostName, pathRequest, requestLink.Encode())
	//userCountry
	//fmt.Printf("Request: %#v\n",urlRequest)
	cl := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cl.Get(urlRequest)
	fmt.Printf("response is: %#v\n", resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error in responce:\nStatus code %v\nStatus \"%s\"", resp.StatusCode, resp.Status))
	}
	byteResp, err := io.ReadAll(resp.Body)
	//fmt.Println(string(byteResp))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteResp, &structResp)
	if err != nil {
		return nil, err
	}

	respVars.Web = structResp.PageUrl
	for key, entity := range structResp.EntitiesByUniqueId {
		if strings.Contains(key, "ITUNES") {
			respVars.Type = entity.Type
			respVars.Title = entity.Title
			respVars.Artist = entity.ArtistName
			respVars.Pic = entity.ThumbnailUrl
			break
		} else {
			respVars.Type = entity.Type
			respVars.Title = entity.Title
			respVars.Artist = entity.ArtistName
			respVars.Pic = entity.ThumbnailUrl
		}
		continue
	}
	respVars.Links = make([]LinkSong, 0, 10)
	for key, val := range structResp.LinksByPlatform {
		respVars.Links = append(respVars.Links, LinkSong{
			Name: key,
			Url:  val.Url,
		})
	}

	//message := fmt.Sprintf("*%s*\n🎤%s\n[pic](%s)\nЯ нашел эту песню на следующих платформах:\n%s",respVars.Title,respVars.Artist,respVars.Pic,respVars.StrList)
	//message := respVars.Title+"\n"+respVars.Artist+"\nЯ нашел эту песню на следующих платформах:\n"+respVars.StrList
	//fmt.Printf("%v\n",message)
	return respVars, nil
}

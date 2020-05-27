package cogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"
)

type Hastebin struct {
	Key string `json:"string"`
}

type GeniusAPI struct {
	Response struct {
		Hits []struct {
			Result struct {
				Title   string `json:"full_title"`
				URL     string `json:"url"`
				SongArt string `json:"song_art_image_thumbnail_url"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

func tinyurl(ctx *commands.Context, link string) (err error) {
	url := "http://tinyurl.com/api-create.php?url=" + link

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error when getting tinyurl response")
	}

	defer resp.Body.Close()
	ShortenedURL, _ := ioutil.ReadAll(resp.Body)

	_, err = ctx.Send("Here is your shortened URL: <" + string(ShortenedURL) + ">")
	return
}

func lyrics(ctx *commands.Context, query ...string) (err error) {
	q := url.Values{}
	q.Set("q", strings.Join(query, " "))
	req, _ := http.NewRequest("GET", "https://api.genius.com/search?"+q.Encode(), nil)
	req.Header.Set("Authorization", "Bearer "+utils.GetConfig("geniusapi"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Genius returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	data := &GeniusAPI{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return
	}

	if len(data.Response.Hits) < 1 {
		ctx.Send("No results found.")
		return
	}

	song := data.Response.Hits[0].Result
	LyricsResp, err := http.Get(song.URL)
	if err != nil {
		return
	}
	defer LyricsResp.Body.Close()
	if LyricsResp.StatusCode != 200 {
		return fmt.Errorf("Genius returned status code %d", LyricsResp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(LyricsResp.Body)
	if err != nil {
		return
	}
	lyrics := strings.TrimSpace(doc.Find("div.lyrics").First().Text())
	lines := strings.Split(lyrics, "\n")

	chars := 0
	split := 0
	for i, line := range lines {
		chars += len(line) + 1
		if chars > 2048 {
			split = i - 1
			break
		}
	}
	if split == 0 {
		split = len(lines)
	}

	em := commands.NewEmbed().
		SetTitle(song.Title).
		SetURL(song.URL).
		SetDescription(strings.Join(lines[:split], "\n")).
		SetThumbnail(song.SongArt).
		SetColor(0x2ecc71)

	if split != len(lines) {
		em.SetFooter("View full lyrics by clicking the title of the embed.")
	}

	_, err = ctx.SendComplex("", em.MessageEmbed)
	return
}

func hastebin(ctx *commands.Context, code ...string) (err error) {
	pushData := utils.CleanupCode(strings.Join(code, " "))
	fmt.Println(pushData)

	jsonData := map[string]string{"data": pushData}
	jsonValue, _ := json.Marshal(jsonData)

	resp, err := http.Post("https://hastebin.com/documents", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("error making POST request")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := &Hastebin{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("error unparsing returned JSON")
	}

	_, err = ctx.Send("Hastebin-ified! Here is your code: https://hastebin.com/" + data.Key + ".go")
	return

}

func init() {
	cog := commands.NewCog("Utility", "Useful commands to help you out")
	cog.AddCommand("tinyurl", "Shorten a URL with the tinyurl API", "<link>", tinyurl)
	// cog.AddCommand("hastebin", "Hastebin-ify your code!", "<code>", hastebin)
	cog.AddCommand("lyrics", "Find the lyrics for a song", "<query>", lyrics)
	cog.Load()
}

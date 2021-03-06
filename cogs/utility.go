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

	var data struct {
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

	err = json.Unmarshal(body, &data)
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
	cleanData := utils.CleanupCode(strings.Join(code, " "))

	resp, err := http.Post("https://hastebin.com/documents", "application/json", bytes.NewBuffer([]byte(cleanData)))
	if err != nil {
		return fmt.Errorf("error making POST request")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Key string `json:"key"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("error unparsing returned JSON")
	}

	_, err = ctx.Send("Hastebin-ified! Here is your code: https://hastebin.com/" + data.Key)
	return

}

func wikipedia(ctx *commands.Context, article ...string) (err error) {
	resp, err := http.Get("https://en.wikipedia.org/api/rest_v1/page/summary/" + url.PathEscape(strings.Join(article, "_")))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ctx.Send("Wikipedia article was not found.")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var data struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
		Thumbnail struct {
			Source string `json:"source"`
		} `json:"thumbnail"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	if data.ContentURLs.Desktop.Page == "" {
		ctx.Send("Wikipedia article was not found.")
		return
	}

	em := commands.NewEmbed().
		SetTitle(data.Title).
		SetURL(data.ContentURLs.Desktop.Page).
		SetColor(0x2ecc71).
		SetDescription(data.Extract).
		SetThumbnail(data.Thumbnail.Source).
		MessageEmbed

	_, err = ctx.SendComplex("", em)
	return
}

func wolfram(ctx *commands.Context, query ...string) (err error) {
	q := url.Values{}
	q.Set("appid", utils.GetConfig("wolframapi"))
	q.Set("i", strings.Join(query, " "))
	req, _ := http.NewRequest("GET", "http://api.wolframalpha.com/v1/result?"+q.Encode(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	answer := string(body)
	_, err = ctx.Send(answer)
	return
}

func init() {
	cog := commands.NewCog("Utility", "Useful commands to help you out")
	cog.AddCommand("tinyurl", "Shorten a URL with the tinyurl API", "<link>", tinyurl)
	cog.AddCommand("hastebin", "Hastebin-ify your code!", "<code>", hastebin)
	cog.AddCommand("lyrics", "Find the lyrics for a song", "<query>", lyrics)
	cog.AddCommand("wikipedia", "Find a wikipedia article", "<article>", wikipedia).
		SetAliases("wiki")
	cog.AddCommand("wolfram", "Search wolfram for an answer to anything", "<query>", wolfram)
	cog.Load()
}

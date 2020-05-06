package cogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SharpBit/go-enigma/utils"
)

// Hastebin responses
type Hastebin struct {
	Key string `json:"string"`
}

func tinyurl(ctx *Context, link string) (err error) {
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

func hastebin(ctx *Context, code ...string) (err error) {
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
	cog := NewCog("Utility", "Useful commands to help you out", false)
	cog.AddCommand("tinyurl", "Shorten a URL with the tinyurl API", "<link>", tinyurl)
	// cog.AddCommand("hastebin", "Hastebin-ify your code!", "<code>", hastebin)
	cog.Load()
}

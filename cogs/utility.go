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

func tinyurl(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.Send("Please pass in a URL")
		return
	}
	link := ctx.Args[0]
	url := "http://tinyurl.com/api-create.php?url=" + link

	resp, err := http.Get(url)
	if err != nil {
		ctx.Send("An error has occured.")
		return
	}

	defer resp.Body.Close()
	ShortenedURL, _ := ioutil.ReadAll(resp.Body)

	ctx.Send("Here is your shortened URL: <" + string(ShortenedURL) + ">")

}

func hastebin(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.Send("Please send some code")
		return
	}
	code := utils.CleanupCode(strings.Join(ctx.Args, " "))
	fmt.Println(code)

	jsonData := map[string]string{"data": code}
	jsonValue, _ := json.Marshal(jsonData)

	resp, err := http.Post("https://hastebin.com/documents", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		ctx.Send("error making POST request")
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := &Hastebin{}
	err = json.Unmarshal(body, data)
	if err != nil {
		ctx.Send("error unparsing returned JSON")
		return
	}

	ctx.Send("Hastebin-ified! Here is your code: https://hastebin.com/" + data.Key + ".go")

}

func init() {
	cog := NewCog("Utility", "Useful commands to help you out", false)
	cog.AddCommand("tinyurl", "Shorten a URL with the tinyurl API", []string{}, tinyurl)
	// cog.AddCommand("hastebin", "Hastebin-ify your code!", []string{}, hastebin)
	cog.Load()
}

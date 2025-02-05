package main

import (
	"flag"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"good.site/pkg/config"
	"good.site/pkg/services"
	"good.site/pkg/utils/convert"
	"good.site/pkg/utils/file"
)

func main() {
	percent := flag.Int("percent", 20, "A number from 0-100 denoting the percentage of nonsense to create")
	flag.Parse()

	reg, err := regexp.Compile(`index\.html$`)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	cfg := config.GetConfig()

	if err != nil {
		panic(err)
	}

	files := file.FindFiles(cfg.Build.Public, reg)

	if len(files) == 0 {
		return
	}

	tokens := []string{}

	for _, f := range files {
		htmlStr := file.ReadFile(f)
		tokens = append(tokens, convert.TokenizeHtml(string(htmlStr))...)
	}

	c := services.NewMarkov(tokens)

	for _, f := range files {
		htmlStr := file.ReadFile(f)

		doc, err := convert.ParseHtml(string(htmlStr))

		if err != nil {
			panic(err)
		}

		body := convert.FindNodes(doc, "body", 1)

		if len(body) > 0 {
			convert.ReplaceText(body[0], func(input string) string {
				words := strings.Split(input, " ")
				newWords := []string{}

				for i := 0; i < len(words); i++ {
					num := r.Intn(100)

					if num < *percent {
						newWords = append(newWords, c.Next())
					} else {
						newWords = append(newWords, words[i])
					}
				}

				out := strings.Join(newWords, " ")

				if input[len(input)-1:] == " " {
					out = out + " "
				}

				if input[0:1] == " " {
					out = " " + out
				}

				return out
			})
		}

		file.WriteFile(strings.Replace(f, cfg.Build.Public, cfg.Build.Nonsense, 1), convert.RenderNode(doc))
	}
}

package processor

import (
	"context"
	"goprocess/contentretriever"
	"goprocess/persistor"
	"strings"

	"golang.org/x/net/html"
)

var (
	script = "script"
	style = "style"
	batchSize = 10
)

func getTermCounts(htmlContent string) map[string]int {
	termCounts := make(map[string]int)
	reader := strings.NewReader(htmlContent)
	tokenizer := html.NewTokenizer(reader)
	var inScriptOrStyle int

	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			return termCounts
		case html.StartTagToken, html.SelfClosingTagToken:
			name, _ := tokenizer.TagName()
			if string(name) == script || string(name) == style {
				inScriptOrStyle++
			}
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			if string(name) == script || string(name) == style {
				inScriptOrStyle--
			}
		case html.TextToken:
			if inScriptOrStyle == 0 {
				text := strings.TrimSpace(string(tokenizer.Text()))
				if text != "" {
					for word := range strings.FieldsSeq(text) {
						termCounts[strings.ToLower(word)]++
					}
				}
			}
		}
	}
}

func Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		urlContents := contentretriever.GetUnprocessed(batchSize)
		if len(urlContents) == 0 {
			break
		}

		urlTermCounts := make([]persistor.UrlTermCount, 0, len(urlContents))
		for _, urlContent := range urlContents {
			termCounts := getTermCounts(urlContent.Content)
			urlTermCounts = append(urlTermCounts, persistor.UrlTermCount{Url: urlContent.Url, TermCount: termCounts})
		}

		persistor.SaveUrlTermCounts(urlTermCounts)

		urls := make([]string, len(urlContents))
		for i, uc := range urlContents {
			urls[i] = uc.Url
		}
		contentretriever.MarkProcessed(urls)
	}
}

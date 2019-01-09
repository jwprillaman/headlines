package main

import (
	"fmt"
	"io"
	"golang.org/x/net/html"
	"strings"
	"github.com/tebeka/selenium"
	"os"
)

//Attempt to use https://sourcegraph.com/github.com/sourcegraph/go-webkit2@master/-/blob/README.md

//var sources = [...]string{"http://www.cnn.com", "http://www.foxnews.com"}
var sources = [...]string{"http://www.cnn.com"}
var userAgentHeader = "Mozilla/5.0 (X11; Linux x86_64; rv:64.0) Gecko/l20100101 Firefox/64.0"

func askSources() {

	const (
		seleniumPath = "vendor/selenium-server-standalone-3.14.0.jar"
		gekoDriverPath = "vendor/geckodriver-v0.23.0-linux64"
		port = 8080
	)

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.GeckoDriver(gekoDriverPath),
		selenium.Output(os.Stderr),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()


	capabilities := selenium.Capabilities{"browserName":"firefox"}
	for _, source := range sources {

		wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err != nil {
			panic(err)
		} else {
			err := wd.Get(source)
			if err != nil {
				panic(err)
			}
		}


		//old

	}
}

func getHeadlines(reader io.Reader) []string {
	headlines := make([]string, 0)
	tokenizer := html.NewTokenizer(reader)
	isSearching := false
	//currentSearchType := ""
	for {
		tokenType := tokenizer.Next()
		if tokenType != html.ErrorToken {
			token := tokenizer.Token()

			/*fmt.Printf("%v\n", token.String())
			fmt.Printf("\t%v\n", token.Data)
			fmt.Printf("\t\t%v\n", string(tokenizer.Text()))
			*/
			if !isSearching {
				if string(token.Data) == "span" {
							for _, attr := range token.Attr {
								fmt.Printf("%v = %v\n", attr.Key, attr.Val)
							}
				}


				for _, attr := range token.Attr {
					if strings.Contains(attr.Val, "headline-text") {
						fmt.Printf("Has headline-text")
						//currentSearchType = token.String()
						isSearching = true
						break
					}
				}
			} else {
				if token.String() != "strong" {
					headlines = append(headlines, string(tokenizer.Text()))
					isSearching = false
				}
			}
		} else {
			fmt.Printf("Token was error. Breaking....")
			break
		}

	}

	return headlines
}

func main() {
	fmt.Printf("Start\n")
	askSources()
	fmt.Printf("End\n")
}

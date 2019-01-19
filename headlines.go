package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
)

//var sources = [...]string{"http://www.cnn.com", "http://www.foxnews.com"}
var sources = [...]string{"http://www.cnn.com"}
var userAgentHeader = "Mozilla/5.0 (X11; Linux x86_64; rv:64.0) Gecko/l20100101 Firefox/64.0"
var headlineClassNames = [...]string{"cd__headline-text"}
var headlineElementNames = [...]string{"span"}

func askSources() {

	const (
		seleniumPath   = "vendor/selenium-server-standalone-3.14.0.jar"
		gekoDriverPath = "vendor/geckodriver-v0.23.0-linux64"
		port           = 8080
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

	capabilities := selenium.Capabilities{"browserName": "firefox"}
	for _, source := range sources {

		//TODO handle selenium warnings
		wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		fmt.Printf("%v\n", wd)

		err = wd.Get(source)
		if err != nil {
			panic(err)
		}
		headlines := getHeadlines(wd)
		for _, headline := range headlines {
			fmt.Printf("%v\n", headline)
		}
	}
}

func getHeadlines(wd selenium.WebDriver) []string {
	headlines := make([]string, 0)
	for _, className := range headlineClassNames {
		elements, _ := wd.FindElements(selenium.ByCSSSelector, fmt.Sprintf(".%v", className))
		for _, e := range elements {
			currentText, _ := e.Text()
			if len(currentText) > 0 {
				headlines = append(headlines, currentText)
			} else { //now find acceptable child
				children, _ := e.FindElements(selenium.ByTagName, "*")
				if len(children) > 0 {
					for _, child := range children {
						childText, _ := child.Text()
						if len(childText) > 0 {
							headlines = append(headlines, childText)
						}
					}
				}
			}
		}
	}

	return headlines
}

func main() {
	fmt.Printf("Start\n")
	askSources()
	fmt.Printf("End\n")
}

package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
)

//var sources = [...]string{"http://www.cnn.com", "http://www.foxnews.com"}
var sources = [...]string{"http://www.cnn.com", "http://www.foxnews.com"}
var headlineClassNames = [...]string{"cd__headline-text", "title"}

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
	//TODO handle selenium warnings
	capabilities := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	for _, source := range sources {
		fmt.Printf("Source : %v\n", source)
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
				headlines = append(headlines, extractHeadlinesFromChildren(e)...)
			}
		}
	}

	return headlines
}

func extractHeadlinesFromChildren(element selenium.WebElement)([]string){
	children, _ := element.FindElements(selenium.ByTagName, "*")
	childHeadlines := make([]string, 0)
	if len(children) > 0 {
		for _, child := range children {
			childText, _ := child.Text()
			if len(childText) > 0 {
				childHeadlines = append(childHeadlines, childText)
			}
			childHeadlines = append(childHeadlines, extractHeadlinesFromChildren(child)...)
		}
	}
	return childHeadlines
}

func main() {
	fmt.Printf("Start\n")
	askSources()
	fmt.Printf("End\n")
}

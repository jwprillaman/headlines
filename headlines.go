package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
	"strings"
)

type SourceSummary struct {
	source            string
	headlines         []string
	filteredHeadlines []string
}

type CommonHeadlines struct {
	sources   string
	topic     string
	headlines []string
}

var sources = [...]string{"http://www.cnn.com"}

//TODO maybe store this somewhere else
var stopWords = [...]string{"a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "as", "at", "be", "because", "been", "before", "being", "below", "between", "both", "but", "by", "could", "did", "do", "does", "doing", "down", "during", "each", "few", "for", "from", "further", "had", "has", "have", "having", "he", "he'd", "he'll", "he's", "her", "here", "here's", "hers", "herself", "him", "himself", "his", "how", "how's", "i", "i'd", "i'll", "i'm", "i've", "if", "in", "into", "is", "it", "it's", "its", "itself", "let's", "me", "more", "most", "my", "myself", "nor", "of", "on", "once", "only", "or", "other", "ought", "our", "ours", "ourselves", "out", "over", "own", "same", "she", "she'd", "she'll", "she's", "should", "so", "some", "such", "than", "that", "that's", "the", "their", "theirs", "them", "themselves", "then", "there", "there's", "these", "they", "they'd", "they'll", "they're", "they've", "this", "those", "through", "to", "too", "under", "until", "up", "very", "was", "we", "we'd", "we'll", "we're", "we've", "were", "what", "what's", "when", "when's", "where", "where's", "which", "while", "who", "who's", "whom", "why", "why's", "with", "would", "you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves"}

var headlineClassNames = [...]string{"cd__headline-text", "title"}

func getSourceSummaries() []SourceSummary {
	output := make([]SourceSummary, 0)
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

	//create stopword set for efficiency
	stopWordSet := make(map[string]struct{})
	for _, word := range stopWords {
		stopWordSet[word] = struct{}{}
	}

	//TODO handle selenium warnings
	capabilities := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	for _, source := range sources {
		err = wd.Get(source)
		if err != nil {
			panic(err)
		}
		headlines := getHeadlines(wd)
		filteredHeadlines := make([]string, len(headlines))
		for i, headline := range headlines {
			filteredHeadlines[i] = filterStopWords(headline, stopWordSet)
		}

		currentSourceSummary := SourceSummary{source, headlines, filteredHeadlines}
		output = append(output, currentSourceSummary)
	}
	return output
}

func filterStopWords(original string, stopWordSet map[string]struct{}) string {
	builder := strings.Builder{}
	tokens := strings.Split(original, " ")
	for i, token := range tokens {
		if _,isStopWord := stopWordSet[token]; !isStopWord {
			builder.WriteString(token)
			if i != len(tokens) -1 {
				builder.WriteString(" ")
			}
		}
	}
	return builder.String()
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

func extractHeadlinesFromChildren(element selenium.WebElement) []string {
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

func printSourceSummary(sourceSummary SourceSummary) {
	fmt.Printf("%v\n", sourceSummary.source)
	for i,v := range sourceSummary.headlines{
		fmt.Printf("\t%v\n\t\t%v\n%v\n", v, sourceSummary.filteredHeadlines[i], i)
	}
}

func main() {
	summaries := getSourceSummaries()
	for _,summary := range summaries {
		printSourceSummary(summary)
	}
}

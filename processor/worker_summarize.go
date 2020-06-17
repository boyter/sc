// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"encoding/json"
	"fmt"
	str "github.com/boyter/cs/str"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"io/ioutil"
	"os"
)

type ResultSummarizer struct {
	input            chan *fileJob
	ResultLimit      int64
	FileReaderWorker *FileReaderWorker
	SnippetCount     int64
	NoColor          bool
	Format           string
	FileOutput       string
}

func NewResultSummarizer(input chan *fileJob) ResultSummarizer {
	return ResultSummarizer{
		input:        input,
		ResultLimit:  -1,
		SnippetCount: 1,
		NoColor:      os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())),
		Format:       Format,
		FileOutput:   FileOutput,
	}
}

func (f *ResultSummarizer) Start() {
	// First step is to collect results so we can rank them
	results := []*fileJob{}
	for res := range f.input {
		results = append(results, res)
	}

	// Consider moving this check into processor to save on CPU burn there at some point in
	// the future
	if f.ResultLimit != -1 {
		if int64(len(results)) > f.ResultLimit {
			results = results[:f.ResultLimit]
		}
	}

	rankResults(int(f.FileReaderWorker.GetFileCount()), results)

	switch f.Format {
	case "json":
		f.formatJson(results)
	default:
		f.formatDefault(results)
	}
}

func (f *ResultSummarizer) formatJson(results []*fileJob) {
	var jsonResults []jsonResult

	documentFrequency := calculateDocumentTermFrequency(results)

	for _, res := range results {
		v3 := extractRelevantV3(res, documentFrequency, int(SnippetLength), "…")[0]

		// We have the snippet so now we need to highlight it
		// we get all the locations that fall in the snippet length
		// and then remove the length of the snippet cut which
		// makes out location line up with the snippet size
		var l [][]int
		for _, value := range res.MatchLocations {
			for _, s := range value {
				if s[0] >= v3.StartPos && s[1] <= v3.EndPos {
					s[0] = s[0] - v3.StartPos
					s[1] = s[1] - v3.StartPos
					l = append(l, s)
				}
			}
		}

		jsonResults = append(jsonResults, jsonResult{
			Filename:       res.Filename,
			Location:       res.Location,
			Content:        v3.Content,
			Score:          res.Score,
			MatchLocations: l,
		})
	}

	jsonString, _ := json.Marshal(jsonResults)
	if f.FileOutput == "" {
		fmt.Println(string(jsonString))
	} else {
		_ = ioutil.WriteFile(FileOutput, []byte(jsonString), 0600)
		fmt.Println("results written to " + FileOutput)
	}
}

func (f *ResultSummarizer) formatDefault(results []*fileJob) {
	fmtBegin := "\033[1;31m"
	fmtEnd := "\033[0m"
	if f.NoColor {
		fmtBegin = ""
		fmtEnd = ""
	}

	documentFrequency := calculateDocumentTermFrequency(results)

	for _, res := range results {
		color.Magenta(fmt.Sprintf("%s (%.3f)", res.Location, res.Score))

		snippets := extractRelevantV3(res, documentFrequency, int(SnippetLength), "…")

		if int64(len(snippets)) > f.SnippetCount {
			snippets = snippets[:f.SnippetCount]
		}


		for i:= 0; i< len(snippets); i++ {

			// We have the snippet so now we need to highlight it
			// we get all the locations that fall in the snippet length
			// and then remove the length of the snippet cut which
			// makes out location line up with the snippet size
			var l [][]int
			for _, value := range res.MatchLocations {
				for _, s := range value {
					if s[0] >= snippets[i].StartPos && s[1] <= snippets[i].EndPos {
						s[0] = s[0] - snippets[i].StartPos
						s[1] = s[1] - snippets[i].StartPos
						l = append(l, s)
					}
				}
			}

			displayContent := snippets[i].Content

			// If the start and end pos are 0 then we don't need to highlight because there is
			// nothing to do so, which means its likely to be a filename match with no content
			if !(snippets[i].StartPos == 0 && snippets[i].EndPos == 0) {
				displayContent = str.HighlightString(snippets[i].Content, l, fmtBegin, fmtEnd)
			}

			fmt.Println(displayContent)
			if i == len(snippets)-1 {
				fmt.Println("")
			} else {
				fmt.Println("")
				fmt.Println("----------")
				fmt.Println("")
			}
		}
	}
}

type jsonResult struct {
	Filename       string  `json:"filename"`
	Location       string  `json:"location"`
	Content        string  `json:"content"`
	Score          float64 `json:"score"`
	MatchLocations [][]int `json:"matchlocations"`
}

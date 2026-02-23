package docs

import _ "embed"

//go:embed russian-history.md
var russianHistory string

//go:embed analysis-framework.md
var analysisFramework string

//go:embed agents/results-extractor.md
var extractPrompt string

func RussianHistory() string {
	return russianHistory
}

func AnalysisFramework() string {
	return analysisFramework
}

func ExtractPrompt() string {
	return extractPrompt
}

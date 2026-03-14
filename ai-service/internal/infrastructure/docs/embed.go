package docs

import _ "embed"

//go:embed russian-history.md
var russianHistory string

//go:embed analysis-framework.md
var analysisFramework string

//go:embed agents/results-extractor.md
var extractPrompt string

//go:embed agents/raw-data-extractor.md
var rawDataAgentPrompt string

//go:embed agents/news-collector.md
var newsCollectorAgent string

func RussianHistory() string {
	return russianHistory
}

func AnalysisFramework() string {
	return analysisFramework
}

func ExtractPrompt() string {
	return extractPrompt
}

func NewsCollectorAgent() string {
	return newsCollectorAgent
}

//go:embed agents/business-researcher.md
var businessResearcherPrompt string

func RawDataAgentPrompt() string {
	return rawDataAgentPrompt
}

func BusinessResearcherPrompt() string {
	return businessResearcherPrompt
}

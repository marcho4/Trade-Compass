package docs

import _ "embed"

//go:embed russian-history.md
var russianHistory string

//go:embed analysis-framework.md
var analysisFramework string

//go:embed agents/results-extractor.md
var extractPrompt string

//go:embed agents/raw-data-extractor-v2.md
var rawDataAgentPrompt string

//go:embed agents/raw-data-validator.md
var rawDataValidatorPrompt string

//go:embed agents/news-collector.md
var newsCollectorAgent string

//go:embed agents/business-researcher.md
var businessResearcherPrompt string

//go:embed agents/risk-n-growth-finder.md
var riskAndGrowthPrompt string

//go:embed agents/scenarios-generator.md
var scenarioGeneratorPrompt string

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

func RawDataAgentPrompt() string {
	return rawDataAgentPrompt
}

func BusinessResearcherPrompt() string {
	return businessResearcherPrompt
}

func RiskAndGrowthPrompt() string {
	return riskAndGrowthPrompt
}

func ScenarioGeneratorPrompt() string {
	return scenarioGeneratorPrompt
}

func RawDataValidatorPrompt() string {
	return rawDataValidatorPrompt
}

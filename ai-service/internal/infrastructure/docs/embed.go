package docs

import _ "embed"

//go:embed russian-history.md
var russianHistory string

//go:embed analysis-framework.md
var analysisFramework string

func RussianHistory() string {
	return russianHistory
}

func AnalysisFramework() string {
	return analysisFramework
}

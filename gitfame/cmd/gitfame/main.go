//go:build !solution
// +build !solution

package main

import (
	"flag"
	"gitfame/internal/statistics"
	"gitfame/internal/writers"
	"os"
	"strings"
)

func main() {
	repo := flag.String("repository", ".", "path to repo")
	revision := flag.String("revision", "HEAD", "revision")
	ord := flag.String("order-by", "lines", "order to sort (by lines, commits, files)")
	useCommiter := flag.Bool("use-committer", false, "change author to committer")
	format := flag.String("format", "tabular", "format to print statistics (tabular, csv, json, json-lines)")
	extensions := flag.String("extensions", "", "extensions to calculate")
	langs := flag.String("languages", "", "languages to calculate")
	excludePatterns := flag.String("exclude", "", "patterns to exclude")
	restrictPatterns := flag.String("restrict-to", "", "patterns to calculate")

	flag.Parse()

	var listLangs, listExtensions, listExclude, listRestict []string
	if len(*langs) > 0 {
		listLangs = strings.Split(*langs, ",")
	}
	if len(*extensions) > 0 {
		listExtensions = strings.Split(*extensions, ",")
	}

	if len(*excludePatterns) > 0 {
		listExclude = strings.Split(*excludePatterns, ",")
	}
	if len(*restrictPatterns) > 0 {
		listRestict = strings.Split(*restrictPatterns, ",")
	}

	files := statistics.GetInterestingFiles(*repo, *revision, listExtensions, listLangs, listExclude, listRestict)
	stats := statistics.CollectStatistics(*repo, *revision, files, *useCommiter)
	statistics.SortStatistics(stats, *ord)

	var writer statistics.FrameWriter
	switch *format {
	case "tabular":
		writer = writers.CreateTabularWriter(os.Stdout)
	case "csv":
		writer = writers.CreateCSVWriter(os.Stdout)
	case "json":
		writer = writers.CreateJSONWriter(os.Stdout)
	case "json-lines":
		writer = writers.CreateJSONLinesWriter(os.Stdout)
	}
	writer.WriteFrames(stats)
}

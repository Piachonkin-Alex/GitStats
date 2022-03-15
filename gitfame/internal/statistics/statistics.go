package statistics

import (
	"gitfame/internal/clearing"
	"gitfame/pkg/gitapi"
	"sort"
	"strings"
)

func GetInterestingFiles(repo, revision string, goodExtensions, langsList, excludePaths, restrictPaths []string) []string {
	files := gitapi.GetFileNames(repo, revision)
	files = clearing.TakeFilesToExtensions(files, langsList, goodExtensions)
	files = clearing.TakeFilesToPaths(files, excludePaths, restrictPaths)
	return files
}

func CollectStatistics(repo, revision string, files []string, useCommitter bool) []Frame {
	authorToCommits := make(map[string]map[string]struct{})
	authorToLines := make(map[string]int)
	authorToFiles := make(map[string]int)
	for _, file := range files {
		if file == "" {
			continue
		}
		fileAuthorToCommits, fileAuthorToLines := gitapi.BlameFile(repo, file, revision, useCommitter)
		for author, commits := range fileAuthorToCommits {
			if _, ok := authorToCommits[author]; !ok {
				authorToCommits[author] = make(map[string]struct{})
			}
			authorToFiles[author] += 1
			for commit := range commits {
				authorToCommits[author][commit] = struct{}{}
			}
		}
		for author, numLines := range fileAuthorToLines {
			authorToLines[author] += numLines
		}
	}

	var result []Frame
	for author, commits := range authorToCommits {
		frame := Frame{AuthorName: author, NumLines: authorToLines[author],
			NumCommits: len(commits), NumFiles: authorToFiles[author]}
		result = append(result, frame)
	}
	return result
}

func SortStatistics(frames []Frame, ord string) {
	comparator := func(i, j int) bool {
		switch ord {
		case "lines":
			if frames[i].NumLines != frames[j].NumLines {
				return frames[i].NumLines > frames[j].NumLines
			}
			if frames[i].NumCommits != frames[j].NumCommits {
				return frames[i].NumCommits > frames[j].NumCommits
			}
			if frames[i].NumFiles != frames[j].NumFiles {
				return frames[i].NumFiles > frames[j].NumFiles
			}
			return strings.Compare(frames[i].AuthorName, frames[j].AuthorName) == -1
		case "commits":
			if frames[i].NumCommits != frames[j].NumCommits {
				return frames[i].NumCommits > frames[j].NumCommits
			}
			if frames[i].NumLines != frames[j].NumLines {
				return frames[i].NumLines > frames[j].NumLines
			}
			if frames[i].NumFiles != frames[j].NumFiles {
				return frames[i].NumFiles > frames[j].NumFiles
			}
			return strings.Compare(frames[i].AuthorName, frames[j].AuthorName) == -1
		case "files":
			if frames[i].NumFiles != frames[j].NumFiles {
				return frames[i].NumFiles > frames[j].NumFiles
			}
			if frames[i].NumLines != frames[j].NumLines {
				return frames[i].NumLines > frames[j].NumLines
			}
			if frames[i].NumCommits != frames[j].NumCommits {
				return frames[i].NumCommits > frames[j].NumCommits
			}
			return strings.Compare(frames[i].AuthorName, frames[j].AuthorName) == -1
		}
		panic("No such comparator!")
	}
	sort.Slice(frames, comparator)
}

package gitapi

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetFileNames(dir, commit string) []string {
	command := exec.Command("git", "ls-tree", commit, "--name-only", "-r")
	command.Dir = dir
	output, err2 := command.Output()
	if err2 != nil {
		panic("ls-tree fails!!!")
	}
	return strings.Split(string(output), "\n")
}

func LogFile(dir, filename, revision string) (string, string) {
	command := exec.Command("git", "log", "--pretty=format:%an%n%H", revision, "--", filename)
	command.Dir = dir
	output, _ := command.Output()
	result := strings.Split(string(output), "\n")
	return result[0], result[1]
}

func BlameFile(dir, filename, revision string, useCommitter bool) (map[string]map[string]struct{}, map[string]int) {
	authorToCommits := make(map[string]map[string]struct{})
	authorTonumLines := make(map[string]int)
	command := exec.Command("git", "blame", "--porcelain", filename, revision)
	command.Dir = dir
	output, _ := command.Output()

	if len(output) == 0 {
		author, commit := LogFile(dir, filename, revision)
		authorToCommits[author] = make(map[string]struct{})
		authorToCommits[author][commit] = struct{}{}
		authorTonumLines[author] = 0
		return authorToCommits, authorTonumLines
	}

	commitToNumLines := make(map[string]int)

	lines := strings.Split(string(output), "\n")
	readAuthor, numReadLines := false, 0
	curHash := ""

	for _, line := range lines {
		if line == "" {
			continue
		}
		if numReadLines == 0 {
			args := strings.Split(line, " ")
			curHash = args[0]
			numReadLines, _ = strconv.Atoi(args[len(args)-1])
			commitToNumLines[curHash] += numReadLines
			readAuthor = true
		} else if line[0] == '\t' {
			numReadLines--
		} else if readAuthor {
			args := strings.Split(line, " ")
			if (!useCommitter && args[0] == "author") || (useCommitter && args[0] == "committer") {
				readAuthor = false
				author := strings.Join(args[1:], " ")
				if _, ok := authorToCommits[author]; !ok {
					authorToCommits[author] = make(map[string]struct{})
				}
				authorToCommits[author][curHash] = struct{}{}
			}
		}
	}
	for author, commits := range authorToCommits {
		for commit, num := range commitToNumLines {
			if _, ok := commits[commit]; ok {
				authorTonumLines[author] += num
			}
		}
	}
	return authorToCommits, authorTonumLines
}

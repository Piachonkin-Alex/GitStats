package writers

import (
	"encoding/csv"
	"gitfame/internal/statistics"
	"io"
	"strconv"
)

type CSVWriter struct {
	writer *csv.Writer
}

func CreateCSVWriter(w io.Writer) *CSVWriter {
	_, _ = w.Write([]byte("Name,Lines,Commits,Files\n"))
	return &CSVWriter{writer: csv.NewWriter(w)}
}

func (w *CSVWriter) WriteFrames(frames []statistics.Frame) {
	for _, frame := range frames {
		lines := strconv.Itoa(frame.NumLines)
		commits := strconv.Itoa(frame.NumCommits)
		files := strconv.Itoa(frame.NumFiles)
		_ = w.writer.Write([]string{frame.AuthorName, lines, commits, files})
	}
	w.writer.Flush()
}

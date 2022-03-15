package writers

import (
	"fmt"
	"gitfame/internal/statistics"
	"io"
	"text/tabwriter"
)

type TabularWriter struct {
	writer *tabwriter.Writer
}

func CreateTabularWriter(w io.Writer) *TabularWriter {
	res := &TabularWriter{writer: tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)}
	_, _ = fmt.Fprintln(res.writer, "Name\tLines\tCommits\tFiles")
	return res
}

func (w *TabularWriter) WriteFrames(frames []statistics.Frame) {
	for _, frame := range frames {
		_, _ = fmt.Fprintf(w.writer, "%s\t%d\t%d\t%d\n", frame.AuthorName, frame.NumLines, frame.NumCommits, frame.NumFiles)
	}
	_ = w.writer.Flush()
}

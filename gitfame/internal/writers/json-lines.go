package writers

import (
	"encoding/json"
	"fmt"
	"gitfame/internal/statistics"
	"io"
)

type JSONLinesWriter struct {
	writer io.Writer
}

func CreateJSONLinesWriter(w io.Writer) *JSONLinesWriter {
	return &JSONLinesWriter{writer: w}
}

func (w *JSONLinesWriter) WriteFrames(frames []statistics.Frame) {
	for _, frame := range frames {
		bytes, _ := json.Marshal(frame)
		_, _ = fmt.Fprintf(w.writer, "%s\n", string(bytes))
	}
}

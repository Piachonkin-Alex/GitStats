package writers

import (
	"encoding/json"
	"gitfame/internal/statistics"
	"io"
)

type JSONWriter struct {
	writer io.Writer
}

func CreateJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{writer: w}
}

func (w *JSONWriter) WriteFrames(frames []statistics.Frame) {
	bytes, _ := json.Marshal(frames)
	_, _ = w.writer.Write(bytes)
}

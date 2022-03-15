package statistics

type Frame struct {
	AuthorName string `json:"name"`
	NumLines   int    `json:"lines"`
	NumCommits int    `json:"commits"`
	NumFiles   int    `json:"files"`
}

type FrameWriter interface {
	WriteFrames(frame []Frame)
}

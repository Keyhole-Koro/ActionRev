package domain

import "time"

type DocumentStatus string

const (
	DocumentStatusUploaded             DocumentStatus = "uploaded"
	DocumentStatusPendingNormalization DocumentStatus = "pending_normalization"
	DocumentStatusProcessing           DocumentStatus = "processing"
	DocumentStatusCompleted            DocumentStatus = "completed"
	DocumentStatusFailed               DocumentStatus = "failed"
)

type ExtractionDepth string

const (
	ExtractionDepthFull    ExtractionDepth = "full"
	ExtractionDepthSummary ExtractionDepth = "summary"
)

type Document struct {
	DocumentID      string
	WorkspaceID     string
	UploadedBy      string // Firebase Auth UID
	Filename        string
	GCSURI          string
	MIMEType        string
	FileSize        int64
	Status          DocumentStatus
	ExtractionDepth ExtractionDepth
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DocumentChunk struct {
	DocumentID        string
	ChunkID           string
	ChunkIndex        int
	Text              string
	SourceFilename    string // zip 展開後ファイル名
	SourcePage        int
	SourceOffsetStart int
	SourceOffsetEnd   int
}

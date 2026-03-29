package domain

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrStorageQuotaExceeded = errors.New("storage quota exceeded")
	ErrFileTooLarge         = errors.New("file too large")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrInvalidInput         = errors.New("invalid input")
	ErrPendingNormalization = errors.New("document pending normalization approval")
)

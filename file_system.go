package scl

import (
	"io"
	"time"
)

type FileSystem interface {
	Glob(pattern string) ([]string, error)
	ReadCloser(path string) (io.ReadCloser, time.Time, error)
}

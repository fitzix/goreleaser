// Package targz implements the Archive interface providing tar.gz archiving
// and compression.
package targz

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Archive as tar.gz
type Archive struct {
	gw *gzip.Writer
	tw *tar.Writer
}

// Close all closeables
func (a Archive) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.gw.Close()
}

// New tar.gz archive
func New(target io.Writer) Archive {
	gw := gzip.NewWriter(target)
	tw := tar.NewWriter(gw)
	return Archive{
		gw: gw,
		tw: tw,
	}
}

// Add file to the archive
func (a Archive) Add(name, src string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}
		filename := strings.Replace(path, src, name, 1)
		if filename == "." && !info.IsDir() {
			for i := len(path) - 1; i >= 0; i-- {
				if path[i] == '/' {
					filename = path[i:]
					break
				}
			}
		}
		header.Name = filename
		if err := a.tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(a.tw, f)
		return err
	})
}

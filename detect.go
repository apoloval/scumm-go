package scumm

import (
	"io"
)

// ResourceFileType is the type of a resource file.
type ResourceFileType string

const (
	// ResourceFileUknown is an unknown resource file.
	ResourceFileUknown ResourceFileType = "unknown resource file"
)

// DetectResourceFile detects the type of a resource file.
func DetectResourceFile(r io.ReadSeeker) ResourceFileType {
	defer r.Seek(0, io.SeekStart)

	if isFileIndexv4(r) {
		return ResourceFileIndexV4
	}
	if isCharsetV4(r) {
		return ResourceFileCharsetV4
	}
	if isResourceBundleV4(r) {
		return ResourceFileBundleV4
	}
	return ResourceFileUknown
}

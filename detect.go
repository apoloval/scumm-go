package scumm

import (
	"io"

	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4"
)

// DetectResourceFile detects the type of a resource file.
func DetectResourceFile(r io.ReadSeeker) vm.ResourceFileType {
	defer r.Seek(0, io.SeekStart)

	if vm4.IsFileIndex(r) {
		return vm4.ResourceFileIndex
	}
	if vm4.IsCharset(r) {
		return vm4.ResourceFileCharset
	}
	if vm4.IsResourceBundle(r) {
		return vm4.ResourceFileBundle
	}
	return vm.ResourceFileUknown
}

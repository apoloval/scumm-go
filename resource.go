package scumm

import (
	"fmt"
	"os"
	"path"

	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4"
)

// FromIndex creates a resource manager from an index file.
func FromIndexFile(f string) (vm.ResourceManager, error) {
	indexFile, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}

	rt := DetectResourceFile(indexFile)
	switch rt {
	case vm4.ResourceFileIndex:
		index, err := vm4.DecodeIndex(indexFile)
		if err != nil {
			return nil, err
		}
		return vm4.NewResourceManager(path.Dir(f), index), nil
	default:
		return nil, fmt.Errorf("invalid index resource: unexpected %s", rt)
	}
}

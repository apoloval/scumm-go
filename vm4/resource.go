package vm4

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/apoloval/scumm-go/ioutils"
	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4/inst"
)

// ResourceBundleKey is the key used to decrypt data resource files for SCUMM v4.
const ResourceBundleKey = 0x69

// ChunkType is the type of a chunk in a resource file for SCUMM v4.
type ChunkType [2]byte

var (
	ChunkTypeLE = ChunkType{'L', 'E'}
	ChunkTypeFO = ChunkType{'F', 'O'}
	ChunkTypeLF = ChunkType{'L', 'F'}
	ChunkTypeRO = ChunkType{'R', 'O'}
	ChunkTypeHD = ChunkType{'H', 'D'}
	ChunkTypeCC = ChunkType{'C', 'C'}
	ChunkTypeSP = ChunkType{'S', 'P'}
	ChunkTypeBX = ChunkType{'B', 'X'}
	ChunkTypePA = ChunkType{'P', 'A'}
	ChunkTypeSA = ChunkType{'S', 'A'}
	ChunkTypeBM = ChunkType{'B', 'M'}
	ChunkTypeOI = ChunkType{'O', 'I'}
	ChunkTypeNL = ChunkType{'N', 'L'}
	ChunkTypeSL = ChunkType{'S', 'L'}
	ChunkTypeOC = ChunkType{'O', 'C'}
	ChunkTypeEX = ChunkType{'E', 'X'}
	ChunkTypeEN = ChunkType{'E', 'N'}
	ChunkTypeLC = ChunkType{'L', 'C'}
	ChunkTypeLS = ChunkType{'L', 'S'}
	ChunkTypeSC = ChunkType{'S', 'C'} // SC: Global vm.Script
)

// String implements the Stringer interface.
func (b ChunkType) String() string {
	return string(b[:])
}

// ChunkHeader is the header of a chunk in a resource file for SCUMM v4.
type ChunkHeader struct {
	Size uint32
	Type ChunkType
}

// ChunkHeaderSize is the size of a chunk header in a resource file for SCUMM v4.
const ChunkHeaderSize = 6

// Decode decodes a chunk header from a reader.
func (h *ChunkHeader) Decode(r io.ReadSeeker, rem *uint32) error {
	if rem != nil && *rem < ChunkHeaderSize {
		return fmt.Errorf("invalid input: chunk header size exceeds remaining bytes")
	}
	err := binary.Read(r, binary.LittleEndian, h)
	if rem != nil {
		*rem -= ChunkHeaderSize
	}
	return err
}

// DecodeAs decodes a chunk header from a reader and checks that the chunk type matches the expected.
func (h *ChunkHeader) DecodeAs(r io.ReadSeeker, t ChunkType, rem *uint32) error {
	offset, _ := r.Seek(0, io.SeekCurrent)
	if err := h.Decode(r, rem); err != nil {
		return err
	}
	if h.Type != t {
		return fmt.Errorf(
			"invalid input: unexpected chunk type %s while decoding %s at offset %d",
			h.Type, t, offset)
	}
	return nil
}

// BodyLen returns the length of the chunk body.
func (h *ChunkHeader) BodyLen() uint32 {
	return h.Size - ChunkHeaderSize
}

// ResourceBundle is a resource bundle for SCUMM v4. This is what is stored in the DISKxx.LEC
// files.
type ResourceBundle struct {
	r io.ReadSeeker

	indexLF map[vm.RoomID]vm.ChunkOffset
}

// NewResourceBundle creates a new resource bundle for SCUMM v4.
func NewResourceBundle(r io.ReadSeeker) *ResourceBundle {
	return &ResourceBundle{
		r: ioutils.NewXorReader(r, ResourceBundleKey),
	}
}

// GetRoom returns the room r from the resource bundle.
func (b *ResourceBundle) GetRoom(r vm.IndexedRoom) (*vm.Room, error) {
	rem, err := b.seekLF(r.ID)
	if err != nil {
		return nil, err
	}
	room := &vm.Room{ID: r.ID, Name: r.Name}
	if err := b.decodeRO(room, &rem); err != nil {
		return nil, err
	}
	return room, nil
}

// GetScript returns the global script r from the resource bundle.
func (b *ResourceBundle) GetScript(r vm.IndexedScript) (*vm.Script, error) {
	_, err := b.seekLF(r.Room)
	if err != nil {
		return nil, err
	}

	rem, err := b.seekChunk(ChunkTypeSC, r.Offset, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bytecode := make([]byte, rem)
	if err := b.decode(binary.LittleEndian, &bytecode, nil); err != nil {
		return nil, err
	}
	return &vm.Script{ID: r.ID, Bytecode: bytecode}, nil
}

func (b *ResourceBundle) decodeRO(r *vm.Room, lfrem *uint32) error {
	var roh ChunkHeader
	if err := roh.DecodeAs(b.r, ChunkTypeRO, lfrem); err != nil {
		return err
	}

	if lfrem != nil && *lfrem < roh.BodyLen() {
		return fmt.Errorf("invalid input: RO chunk size exceeds remaining bytes")
	}
	rorem := roh.BodyLen()

	if err := b.decodeHD(r, &rorem); err != nil {
		return err
	}

	if err := b.decodeAndSkipBlock(ChunkTypeCC, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeSP, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeBX, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypePA, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeSA, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeBM, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfObjects); i++ {
		if err := b.decodeAndSkipBlock(ChunkTypeOI, &rorem); err != nil {
			return err
		}
	}
	if err := b.decodeAndSkipBlock(ChunkTypeNL, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeSL, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfObjects); i++ {
		if err := b.decodeAndSkipBlock(ChunkTypeOC, &rorem); err != nil {
			return err
		}
	}
	if err := b.decodeAndSkipBlock(ChunkTypeEX, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeEN, &rorem); err != nil {
		return err
	}
	if err := b.decodeLC(r, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfLocalScripts); i++ {
		var lsh ChunkHeader
		if err := lsh.DecodeAs(b.r, ChunkTypeLS, &rorem); err != nil {
			return err
		}
		var id uint8
		if err := b.decode(binary.LittleEndian, &id, &rorem); err != nil {
			return err
		}
		bytecode := make([]byte, lsh.BodyLen()-1)
		if err := b.decode(binary.LittleEndian, &bytecode, &rorem); err != nil {
			return err
		}
		r.LocalScripts = append(r.LocalScripts, vm.Script{
			ID:       vm.ScriptID(id),
			Bytecode: bytecode,
		})
	}

	if rorem > 0 {
		return fmt.Errorf(
			"invalid input: %d bytes remaining after decoding entire RO chunk body", rorem)
	}
	if lfrem != nil {
		*lfrem = roh.BodyLen()
	}
	return nil
}

func (b *ResourceBundle) decodeHD(r *vm.Room, rorem *uint32) error {
	var hdh ChunkHeader
	if err := hdh.DecodeAs(b.r, ChunkTypeHD, rorem); err != nil {
		return err
	}

	var hd struct {
		Width           uint16
		Height          uint16
		NumberOfObjects uint16
	}
	if err := b.decode(binary.LittleEndian, &hd, rorem); err != nil {
		return err
	}
	r.Width = hd.Width
	r.Height = hd.Height
	r.NumberOfObjects = hd.NumberOfObjects
	return nil
}

func (b *ResourceBundle) decodeLC(r *vm.Room, rem *uint32) error {
	var lch ChunkHeader
	if err := lch.DecodeAs(b.r, ChunkTypeLC, rem); err != nil {
		return err
	}

	var lc struct {
		NumberOfLocalScripts uint8
		_                    uint8
	}
	if err := b.decode(binary.LittleEndian, &lc, rem); err != nil {
		return err
	}
	r.NumberOfLocalScripts = lc.NumberOfLocalScripts
	return nil
}

func (b *ResourceBundle) seekLF(r vm.RoomID) (size uint32, err error) {
	if err := b.ensureIndexLF(r); err != nil {
		return 0, err
	}
	rem, err := b.seekChunk(ChunkTypeLF, b.indexLF[r], io.SeekStart)
	if err != nil {
		return 0, err
	}

	var id uint16
	if err := b.decode(binary.LittleEndian, &id, &rem); err != nil {
		return 0, err
	}
	if vm.RoomID(id) != r {
		return 0, fmt.Errorf("invalid input: unexpected room ID %d while seeking room %d", id, r)
	}
	return rem, nil
}

func (b *ResourceBundle) ensureIndexLF(r vm.RoomID) error {
	if b.indexLF == nil {
		return b.readFO(r)
	}
	return nil
}

func (b *ResourceBundle) readFO(r vm.RoomID) error {
	rem, err := b.seekChunk(ChunkTypeLE, 0, io.SeekStart)
	if err != nil {
		return err
	}

	var foh ChunkHeader
	if err := foh.DecodeAs(b.r, ChunkTypeFO, &rem); err != nil {
		return err
	}

	var fo struct {
		NumberOfBundles uint8
	}
	if err := b.decode(binary.LittleEndian, &fo, &rem); err != nil {
		return err
	}

	b.indexLF = make(map[vm.RoomID]vm.ChunkOffset, fo.NumberOfBundles)
	for i := uint8(0); i < fo.NumberOfBundles; i++ {
		var loc struct {
			LF     uint8
			Offset vm.ChunkOffset
		}
		if err := b.decode(binary.LittleEndian, &loc, &rem); err != nil {
			return err
		}
		b.indexLF[vm.RoomID(loc.LF)] = loc.Offset
	}

	return nil
}

func (b *ResourceBundle) seek(offset vm.ChunkOffset, whence int) error {
	_, err := b.r.Seek(int64(offset), whence)
	return err
}

func (b *ResourceBundle) seekChunk(t ChunkType, offset vm.ChunkOffset, whence int) (size uint32, err error) {
	if err := b.seek(offset, whence); err != nil {
		return 0, err
	}
	var h ChunkHeader
	if err := h.DecodeAs(b.r, t, nil); err != nil {
		return 0, err
	}
	return h.BodyLen(), nil
}

func (b *ResourceBundle) decode(bo binary.ByteOrder, data any, rem *uint32) error {
	from, _ := b.r.Seek(0, io.SeekCurrent)
	if err := binary.Read(b.r, bo, data); err != nil {
		return err
	}
	to, _ := b.r.Seek(0, io.SeekCurrent)
	len := uint32(to - from)
	if rem != nil {
		if *rem < len {
			return fmt.Errorf("invalid input: decoded object exceeds remaining bytes")
		}
		*rem -= len
	}
	return nil
}

func (b *ResourceBundle) decodeAndSkipBlock(t ChunkType, rem *uint32) error {
	var h ChunkHeader
	if err := h.DecodeAs(b.r, t, rem); err != nil {
		return err
	}
	return b.skip(h.BodyLen(), rem)
}

func (b *ResourceBundle) skip(n uint32, rem *uint32) error {
	if rem != nil && *rem < n {
		return fmt.Errorf("skip failed: not enough remaining bytes")
	}
	_, err := b.r.Seek(int64(n), io.SeekCurrent)
	if rem != nil {
		*rem -= n
	}
	return err
}

// ResourceManager is a resource manager for SCUMM v4.
type ResourceManager struct {
	basePath string
	index    vm.Index
	bundles  map[int]*ResourceBundle
}

// NewResourceManager creates a new resource manager for SCUMM v4.
func NewResourceManager(basePath string, index vm.Index) *ResourceManager {
	return &ResourceManager{
		basePath: basePath,
		index:    index,
		bundles:  make(map[int]*ResourceBundle),
	}
}

// GetRoom implements the ResourceManager interface.
func (m *ResourceManager) GetRoom(id vm.RoomID) (*vm.Room, error) {
	r, ok := m.index.Rooms[id]
	if !ok {
		return nil, fmt.Errorf("unknown room ID %d", id)
	}
	bundle, err := m.getBundle(int(r.FileNumber))
	if err != nil {
		return nil, err
	}
	return bundle.GetRoom(r)
}

// GetRoomByName implements the ResourceManager interface.
func (b *ResourceManager) GetRoomByName(name vm.RoomName) (*vm.Room, error) {
	for _, r := range b.index.Rooms {
		if r.Name == name {
			return b.GetRoom(r.ID)
		}
	}
	return nil, fmt.Errorf("unknown room %s", name)
}

// GetScript implements the ResourceManager interface.
func (m *ResourceManager) GetScript(id vm.ScriptID, decode bool) (*vm.Script, error) {
	s, ok := m.index.Scripts[id]
	if !ok {
		return nil, fmt.Errorf("unknown script ID %d", id)
	}
	r, ok := m.index.Rooms[s.Room]
	if !ok {
		return nil, fmt.Errorf("unknown room ID %d", s.Room)
	}
	bundle, err := m.getBundle(int(r.FileNumber))
	if err != nil {
		return nil, err
	}
	script, err := bundle.GetScript(s)
	if err != nil {
		return nil, err
	}

	if decode {
		err = script.Decode(inst.Decode)
	}
	return script, err
}

func (m *ResourceManager) getBundle(id int) (*ResourceBundle, error) {
	bundle, ok := m.bundles[id]
	if !ok {
		bundle, err := m.openBundle(id)
		if err != nil {
			return nil, err
		}
		m.bundles[id] = bundle
		return bundle, nil
	}
	return bundle, nil
}

func (m *ResourceManager) openBundle(id int) (bundle *ResourceBundle, err error) {
	var file *os.File

	// Because in case-sensitive file systems, shit happens.
	paths := []string{
		path.Join(m.basePath, fmt.Sprintf("DISK%02d.LEC", id)),
		path.Join(m.basePath, fmt.Sprintf("disk%02d.lec", id)),
		path.Join(m.basePath, fmt.Sprintf("DISK%02d.lec", id)),
		path.Join(m.basePath, fmt.Sprintf("disk%02d.LEC", id)),
		path.Join(m.basePath, fmt.Sprintf("Disk%02d.lec", id)),
		path.Join(m.basePath, fmt.Sprintf("Disk%02d.LEC", id)),
	}

	for _, p := range paths {
		file, err = os.Open(p)
		if err == nil {
			return NewResourceBundle(file), nil
		}
	}
	return nil, fmt.Errorf("failed to open bundle %d file", id)
}

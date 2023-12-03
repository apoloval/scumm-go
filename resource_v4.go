package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/apoloval/scumm-go/ioutils"
)

// ResourceBundleV4Key is the key used to decrypt data resource files for SCUMM v4.
const ResourceBundleV4Key = 0x69

// ChunkTypeV4 is the type of a chunk in a resource file for SCUMM v4.
type ChunkTypeV4 [2]byte

var (
	ChunkTypeV4LE = ChunkTypeV4{'L', 'E'}
	ChunkTypeV4FO = ChunkTypeV4{'F', 'O'}
	ChunkTypeV4LF = ChunkTypeV4{'L', 'F'}
	ChunkTypeV4RO = ChunkTypeV4{'R', 'O'}
	ChunkTypeV4HD = ChunkTypeV4{'H', 'D'}
	ChunkTypeV4CC = ChunkTypeV4{'C', 'C'}
	ChunkTypeV4SP = ChunkTypeV4{'S', 'P'}
	ChunkTypeV4BX = ChunkTypeV4{'B', 'X'}
	ChunkTypeV4PA = ChunkTypeV4{'P', 'A'}
	ChunkTypeV4SA = ChunkTypeV4{'S', 'A'}
	ChunkTypeV4BM = ChunkTypeV4{'B', 'M'}
	ChunkTypeV4OI = ChunkTypeV4{'O', 'I'}
	ChunkTypeV4NL = ChunkTypeV4{'N', 'L'}
	ChunkTypeV4SL = ChunkTypeV4{'S', 'L'}
	ChunkTypeV4OC = ChunkTypeV4{'O', 'C'}
	ChunkTypeV4EX = ChunkTypeV4{'E', 'X'}
	ChunkTypeV4EN = ChunkTypeV4{'E', 'N'}
	ChunkTypeV4LC = ChunkTypeV4{'L', 'C'}
	ChunkTypeV4LS = ChunkTypeV4{'L', 'S'}
	ChunkTypeV4SC = ChunkTypeV4{'S', 'C'} // SC: Global Script
)

// String implements the Stringer interface.
func (b ChunkTypeV4) String() string {
	return string(b[:])
}

// ChunkHeaderV4 is the header of a chunk in a resource file for SCUMM v4.
type ChunkHeaderV4 struct {
	Size uint32
	Type ChunkTypeV4
}

// ChunkHeaderV4Size is the size of a chunk header in a resource file for SCUMM v4.
const ChunkHeaderV4Size = 6

// Decode decodes a chunk header from a reader.
func (h *ChunkHeaderV4) Decode(r io.ReadSeeker, rem *uint32) error {
	if rem != nil && *rem < ChunkHeaderV4Size {
		return fmt.Errorf("invalid input: chunk header size exceeds remaining bytes")
	}
	err := binary.Read(r, binary.LittleEndian, h)
	if rem != nil {
		*rem -= ChunkHeaderV4Size
	}
	return err
}

// DecodeAs decodes a chunk header from a reader and checks that the chunk type matches the expected.
func (h *ChunkHeaderV4) DecodeAs(r io.ReadSeeker, t ChunkTypeV4, rem *uint32) error {
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
func (h *ChunkHeaderV4) BodyLen() uint32 {
	return h.Size - ChunkHeaderV4Size
}

// ResourceBundleV4 is a resource bundle for SCUMM v4. This is what is stored in the DISKxx.LEC
// files.
type ResourceBundleV4 struct {
	r io.ReadSeeker

	indexLF map[RoomID]ChunkOffset
}

// NewResourceBundleV4 creates a new resource bundle for SCUMM v4.
func NewResourceBundleV4(r io.ReadSeeker) *ResourceBundleV4 {
	return &ResourceBundleV4{
		r: ioutils.NewXorReader(r, ResourceBundleV4Key),
	}
}

// GetRoom returns the room r from the resource bundle.
func (b *ResourceBundleV4) GetRoom(r IndexedRoom) (*Room, error) {
	rem, err := b.seekLF(r.ID)
	if err != nil {
		return nil, err
	}
	room := &Room{ID: r.ID, Name: r.Name}
	if err := b.decodeRO(room, &rem); err != nil {
		return nil, err
	}
	return room, nil
}

// GetScript returns the global script r from the resource bundle.
func (b *ResourceBundleV4) GetScript(r IndexedScript) (*Script, error) {
	_, err := b.seekLF(r.Room)
	if err != nil {
		return nil, err
	}

	rem, err := b.seekChunk(ChunkTypeV4SC, r.Offset, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bytecode := make([]byte, rem)
	return &Script{ID: r.ID, Bytecode: bytecode}, nil
}

func (b *ResourceBundleV4) decodeRO(r *Room, lfrem *uint32) error {
	var roh ChunkHeaderV4
	if err := roh.DecodeAs(b.r, ChunkTypeV4RO, lfrem); err != nil {
		return err
	}

	if lfrem != nil && *lfrem < roh.BodyLen() {
		return fmt.Errorf("invalid input: RO chunk size exceeds remaining bytes")
	}
	rorem := roh.BodyLen()

	if err := b.decodeHD(r, &rorem); err != nil {
		return err
	}

	if err := b.decodeAndSkipBlock(ChunkTypeV4CC, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4SP, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4BX, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4PA, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4SA, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4BM, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfObjects); i++ {
		if err := b.decodeAndSkipBlock(ChunkTypeV4OI, &rorem); err != nil {
			return err
		}
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4NL, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4SL, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfObjects); i++ {
		if err := b.decodeAndSkipBlock(ChunkTypeV4OC, &rorem); err != nil {
			return err
		}
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4EX, &rorem); err != nil {
		return err
	}
	if err := b.decodeAndSkipBlock(ChunkTypeV4EN, &rorem); err != nil {
		return err
	}
	if err := b.decodeLC(r, &rorem); err != nil {
		return err
	}
	for i := 0; i < int(r.NumberOfLocalScripts); i++ {
		var lsh ChunkHeaderV4
		if err := lsh.DecodeAs(b.r, ChunkTypeV4LS, &rorem); err != nil {
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
		r.LocalScripts = append(r.LocalScripts, Script{
			ID:       ScriptID(id),
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

func (b *ResourceBundleV4) decodeHD(r *Room, rorem *uint32) error {
	var hdh ChunkHeaderV4
	if err := hdh.DecodeAs(b.r, ChunkTypeV4HD, rorem); err != nil {
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

func (b *ResourceBundleV4) decodeLC(r *Room, rem *uint32) error {
	var lch ChunkHeaderV4
	if err := lch.DecodeAs(b.r, ChunkTypeV4LC, rem); err != nil {
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

func (b *ResourceBundleV4) seekLF(r RoomID) (size uint32, err error) {
	if err := b.ensureIndexLF(r); err != nil {
		return 0, err
	}
	rem, err := b.seekChunk(ChunkTypeV4LF, b.indexLF[r], io.SeekStart)
	if err != nil {
		return 0, err
	}

	var id uint16
	if err := b.decode(binary.LittleEndian, &id, &rem); err != nil {
		return 0, err
	}
	if RoomID(id) != r {
		return 0, fmt.Errorf("invalid input: unexpected room ID %d while seeking room %d", id, r)
	}
	return rem, nil
}

func (b *ResourceBundleV4) ensureIndexLF(r RoomID) error {
	if b.indexLF == nil {
		return b.readFO(r)
	}
	return nil
}

func (b *ResourceBundleV4) readFO(r RoomID) error {
	rem, err := b.seekChunk(ChunkTypeV4LE, 0, io.SeekStart)
	if err != nil {
		return err
	}

	var foh ChunkHeaderV4
	if err := foh.DecodeAs(b.r, ChunkTypeV4FO, &rem); err != nil {
		return err
	}

	var fo struct {
		NumberOfBundles uint8
	}
	if err := b.decode(binary.LittleEndian, &fo, &rem); err != nil {
		return err
	}

	b.indexLF = make(map[RoomID]ChunkOffset, fo.NumberOfBundles)
	for i := uint8(0); i < fo.NumberOfBundles; i++ {
		var loc struct {
			LF     uint8
			Offset ChunkOffset
		}
		if err := b.decode(binary.LittleEndian, &loc, &rem); err != nil {
			return err
		}
		b.indexLF[RoomID(loc.LF)] = loc.Offset
	}

	return nil
}

func (b *ResourceBundleV4) seek(offset ChunkOffset, whence int) error {
	_, err := b.r.Seek(int64(offset), whence)
	return err
}

func (b *ResourceBundleV4) seekChunk(t ChunkTypeV4, offset ChunkOffset, whence int) (size uint32, err error) {
	if err := b.seek(offset, whence); err != nil {
		return 0, err
	}
	var h ChunkHeaderV4
	if err := h.DecodeAs(b.r, t, nil); err != nil {
		return 0, err
	}
	return h.BodyLen(), nil
}

func (b *ResourceBundleV4) decode(bo binary.ByteOrder, data any, rem *uint32) error {
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

func (b *ResourceBundleV4) decodeAndSkipBlock(t ChunkTypeV4, rem *uint32) error {
	var h ChunkHeaderV4
	if err := h.DecodeAs(b.r, t, rem); err != nil {
		return err
	}
	return b.skip(h.BodyLen(), rem)
}

func (b *ResourceBundleV4) skip(n uint32, rem *uint32) error {
	if rem != nil && *rem < n {
		return fmt.Errorf("skip failed: not enough remaining bytes")
	}
	_, err := b.r.Seek(int64(n), io.SeekCurrent)
	if rem != nil {
		*rem -= n
	}
	return err
}

// ResourceManagerV4 is a resource manager for SCUMM v4.
type ResourceManagerV4 struct {
	basePath string
	index    Index
	bundles  map[int]*ResourceBundleV4
}

// NewResourceManagerV4 creates a new resource manager for SCUMM v4.
func NewResourceManagerV4(basePath string, index Index) *ResourceManagerV4 {
	return &ResourceManagerV4{
		basePath: basePath,
		index:    index,
		bundles:  make(map[int]*ResourceBundleV4),
	}
}

// GetRoom implements the ResourceManager interface.
func (m *ResourceManagerV4) GetRoom(id RoomID) (*Room, error) {
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
func (b *ResourceManagerV4) GetRoomByName(name RoomName) (*Room, error) {
	for _, r := range b.index.Rooms {
		if r.Name == name {
			return b.GetRoom(r.ID)
		}
	}
	return nil, fmt.Errorf("unknown room %s", name)
}

// GetScript implements the ResourceManager interface.
func (m *ResourceManagerV4) GetScript(id ScriptID) (*Script, error) {
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
	return bundle.GetScript(s)
}

func (m *ResourceManagerV4) getBundle(id int) (*ResourceBundleV4, error) {
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

func (m *ResourceManagerV4) openBundle(id int) (bundle *ResourceBundleV4, err error) {
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
			return NewResourceBundleV4(file), nil
		}
	}
	return nil, fmt.Errorf("failed to open bundle %d file", id)
}

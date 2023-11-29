package scumm_test

import (
	"bytes"
	"testing"

	"github.com/apoloval/scumm-go"
	"github.com/apoloval/scumm-go/scummtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeIndexFile(t *testing.T) {
	input := bytes.NewReader(scummtest.MonkeyIsland["000.LFL"])
	index, err := scumm.DecodeIndexFile(input)

	require.NoError(t, err)
	assert.Equal(t, 83, len(index.Rooms))

	// Check room names. Some rooms have no names. Here the named ones are checked.
	var expectedRoomNames = map[scumm.RoomNumber]string{
		90: "copycrap", 96: "part1", 10: "logo", 38: "lookout", 33: "dock", 63: "map",
		44: "masters-h", 61: "sword-mas", 83: "cu-dock", 42: "underwate", 28: "bar", 41: "kitchen",
		81: "cu-bar-2", 82: "cu-bar-3", 79: "cu-bar-1", 35: "low-stree", 29: "fortune",
		23: "cu-gov", 53: "foyer", 32: "alley", 34: "high-stre", 31: "jail", 30: "store",
		78: "church", 45: "cu-church", 57: "bridge", 89: "cu-works", 36: "mansion-e", 59: "stans",
		85: "melee", 58: "damnfores", 62: "cu-sword", 49: "road", 43: "trainers-", 60: "gym",
		76: "cu-traine", 88: "fighting", 51: "circus-te", 52: "circus-gr", 48: "crossing",
		64: "treasure-", 37: "meats-hou", 50: "cu-brush", 97: "part2", 7: "sh-captai", 8: "sh-hull",
		9: "sh-crews", 14: "sh-galley", 19: "sh-deck", 17: "sh-crows", 84: "recipe", 98: "part3",
		87: "seamonkey", 12: "monkey-he", 69: "bighead", 65: "hellhall", 70: "hellcliff",
		39: "hellmaze", 71: "gh-room", 72: "gh-captai", 73: "gh-crews", 74: "gh-storag",
		75: "gh-hull", 77: "gh-deck", 95: "part4", 20: "main-beac", 1: "beach", 2: "monkey-1",
		3: "monkey-2", 4: "monkey-3", 5: "monkey-4", 6: "monkey-5", 21: "gen-jungl", 18: "crack",
		11: "vista-1", 16: "seesaw", 40: "pond", 25: "village", 27: "hut", 80: "fort",
		86: "cu-naviga", 15: "fork", 94: "roland",
	}
	for roomNumber, roomName := range expectedRoomNames {
		assert.Equal(t, roomName, index.Rooms[roomNumber].Name.String())
	}

	// Check directory of rooms
	var expectedFileNumbers = map[scumm.RoomNumber]uint8{
		1: 4, 2: 4, 3: 4, 4: 4, 5: 4, 6: 4, 7: 3, 8: 3, 9: 3, 10: 1, 11: 4, 12: 3, 13: 0,
		14: 3, 15: 4, 16: 4, 17: 3, 18: 4, 19: 3, 20: 4, 21: 4, 22: 0, 23: 2, 24: 0, 25: 4, 26: 0,
		27: 4, 28: 1, 29: 1, 30: 2, 31: 2, 32: 2, 33: 1, 34: 2, 35: 1, 36: 2, 37: 1, 38: 1, 39: 3,
		40: 4, 41: 1, 42: 1, 43: 2, 44: 1, 45: 2, 46: 0, 47: 0, 48: 1, 49: 2, 50: 3, 51: 3, 52: 3,
		53: 2, 54: 0, 55: 0, 56: 0, 57: 2, 58: 2, 59: 2, 60: 2, 61: 1, 62: 2, 63: 1, 64: 3, 65: 3,
		66: 0, 67: 0, 68: 0, 69: 3, 70: 3, 71: 3, 72: 3, 73: 3, 74: 3, 75: 3, 76: 2, 77: 3, 78: 2,
		79: 1, 80: 4, 81: 1, 82: 1, 83: 1, 84: 3, 85: 2, 86: 4, 87: 3, 88: 2, 89: 2, 90: 1, 91: 0,
		92: 0, 93: 0, 94: 4, 95: 3, 96: 1, 97: 3, 98: 3,
	}
	for roomNumber, fileNumber := range expectedFileNumbers {
		assert.Equal(t, fileNumber, index.Rooms[roomNumber].FileNumber)
		assert.Zero(t, index.Rooms[roomNumber].FileOffset)
	}

	// Check directory of scripts
	var expectedScriptOffsets = map[scumm.RoomNumber][]uint32{
		39: {0x245b7},
		16: {0x05f56},
		88: {
			0x004b6, 0x005c5, 0x007d6, 0x0089b, 0x009a3, 0x00ab9, 0x00afb, 0x00dbd, 0x00f33,
			0x00f7f, 0x0177e, 0x01bce, 0x01c0c, 0x01c98, 0x01d60, 0x01ef4, 0x02052, 0x022cc,
			0x0247c, 0x0248d, 0x02531,
		},
		83: {0x074fd, 0x07755},
		81: {0x062d8},
		25: {
			0x14086, 0x14291, 0x142ce, 0x147e1, 0x149e5, 0x14f9a, 0x151ae, 0x15216, 0x1526d,
			0x15ec5, 0x15fb8,
		},
		78: {
			0x06ab3, 0x06ace, 0x06af7, 0x06b91, 0x06d7a, 0x06dbe, 0x06f41, 0x0744c, 0x075cd,
			0x07634, 0x07671, 0x077f9, 0x0786d, 0x07978, 0x07a14,
		},
		61: {0x0855e, 0x0903d, 0x095e0, 0x097a7, 0x09b3c},
		20: {
			0x0d99c, 0x0db56, 0x0f830, 0x0fb1d, 0x0fbb2, 0x0fbc3, 0x10351, 0x103e6, 0x104cf,
			0x105c7,
		},
		33: {0x138e7, 0x1399a, 0x139f1, 0x13a0b, 0x13a8d, 0x13ab0, 0x13b17, 0x13b1e},
		37: {0x0efba, 0x0f10e},
		41: {0x07f77, 0x08045},
		69: {0x07762},
		85: {0x08434},
		63: {0x05885},
		21: {0x0c557},
		10: {
			0x0da2c, 0x0e87b, 0x0ec30, 0x0ed71, 0x0f21c, 0x0f227, 0x0f23b, 0x0f275, 0x0f29e,
			0x0f3ae, 0x0f3f9, 0x0f419, 0x0f477, 0x0f48e, 0x0f62d, 0x0f634, 0x0f8f3, 0x0f94e,
			0x0f994, 0x0f9be, 0x0f9c5, 0x0fd09, 0x0fe68, 0x0fee6, 0x0ff52, 0x0ff87, 0x0ff8e,
			0x0ff95, 0x0ff9c, 0x0ffa3, 0x0ffc9, 0x10011, 0x10081, 0x100f3, 0x1010d, 0x101b9,
			0x101db, 0x101f4, 0x1033b,
		},
		15: {0x0ec84},
		31: {0x09aed},
		97: {0x02acd},
		89: {0x09880, 0x0a42d, 0x0ab92, 0x0abf7, 0x0add3},
		28: {0x19098},
		77: {0x0b15f, 0x0b16f, 0x0b1ae, 0x0b1bd},
		19: {0x0bf95},
		90: {0x05d8c, 0x05e02, 0x0602f, 0x06251, 0x0626c, 0x062b0},
		43: {0x06fe5, 0x08e46}, 38: {0x0734f, 0x078ec},
		7:  {0x06953},
		65: {0x138f4, 0x1393a, 0x1397d, 0x139c8},
		79: {0x087b1},
		59: {0x10d7b},
		30: {0x0aa91},
		82: {0x084dc},
		14: {0x06293, 0x064f7},
		23: {0x08012},
		70: {0x0e3e4, 0x0e556},
		2:  {0x09563, 0x095d0, 0x09600, 0x09817, 0x09897, 0x09926, 0x09939, 0x099c3, 0x09d1a},
		57: {0x06553},
	}
	for i, scriptOffset := range expectedScriptOffsets {
		assert.Equal(t, scriptOffset, index.Rooms[i].ScriptOffsets)
	}

	// Check directory of sounds
	var expectedSoundOffsets = map[scumm.RoomNumber][]uint32{
		2:  {0x09d45},
		8:  {0x065c1, 0x0663a, 0x06548},
		11: {0x0dabc, 0x0dbad, 0x0dc31, 0x0dcc0},
		14: {0x06542, 0x06667, 0x066ee},
		15: {0x0ee37, 0x0ef61}, 38: {0x07f26},
		19: {0x0e393, 0x0e4b8, 0x0c0fb},
		25: {0x16013},
		28: {0x19956},
		29: {0x09b81},
		30: {0x0adfb},
		31: {0x09fb9},
		34: {0x13883},
		35: {0x0fd35},
		37: {0x109a8, 0x10b28, 0x10bb7, 0x10c2c},
		41: {0x0821e, 0x0829f, 0x08318, 0x0839c, 0x0842b},
		43: {0x08f48},
		45: {0x07399, 0x074a2, 0x07593, 0x0769c},
		48: {0x070c7, 0x07130, 0x071ee},
		51: {0x08313, 0x09b4e, 0x0841c},
		53: {0x0feaf, 0x0fd6a, 0x0fe28, 0x0f802, 0x0f89d, 0x0fbb5, 0x0fccb, 0x0ff48, 0x0ffd7},
		58: {0x182b8, 0x1833d, 0x184bd, 0x17420},
		60: {0x065de, 0x0666d},
		64: {0x0b277, 0x0b306},
		70: {0x0e5de},
		71: {0x04468, 0x04507},
		75: {0x07a29, 0x07c6b, 0x07d3d, 0x07db2, 0x07e51, 0x07ee0},
		77: {0x0b43f}, 10: {0x19054, 0x18e89, 0x18f54, 0x10d74},
		78: {0x07a3e},
		81: {0x071b6},
		83: {0x0776d, 0x07b65, 0x07c6e, 0x0a1c2},
		85: {0x084b7},
		88: {0x025b9, 0x02648, 0x026d7, 0x0275c},
		94: {
			0x0759b, 0x06001, 0x090c7, 0x0c77d, 0x0cf4c, 0x0d609, 0x0dcc6, 0x0e383, 0x0ea40,
			0x0f8a4, 0x004c7, 0x1099a, 0x11db5, 0x14921, 0x16a2c, 0x179fc, 0x185d8, 0x19eba,
			0x1bc05,
		},
		95: {0x02f3d},
		96: {0x02cd2},
		97: {0x02c27},
		98: {0x02cc6},
	}
	for i, soundOffset := range expectedSoundOffsets {
		assert.Equal(t, soundOffset, index.Rooms[i].SoundOffsets)
	}

	// Check directory of costumes
	var expectedCostumeOffsets = map[scumm.RoomNumber][]uint32{
		2:  {0x0a7dd, 0x0aada, 0x0d2b7},
		7:  {0x06f2d},
		10: {0x190e3, 0x1ae26, 0x1b00a},
		11: {0x0dd39, 0x0dee2, 0x0e41c},
		14: {0x06918, 0x08153},
		15: {0x0f086, 0x0ffa9, 0x1255a, 0x1231c},
		16: {0x06317, 0x07784},
		17: {0x04126, 0x04631},
		18: {0x08401},
		19: {0x0e5e2, 0x0feef, 0x0f7dc, 0x0ef3c, 0x10adb, 0x1316a, 0x13b82},
		20: {0x107b0, 0x19eac, 0x15d95, 0x144be, 0x166fe},
		25: {0x160a2, 0x197dc, 0x1cf88, 0x201b2},
		28: {0x1ac6d, 0x1d465, 0x203f1},
		29: {0x0a819, 0x0bb01, 0x0cad5},
		30: {0x0ae70},
		31: {0x0a07c, 0x0c4b1},
		32: {0x06b74},
		33: {0x13b25},
		35: {0x0fdae, 0x115a8, 0x12fd4, 0x134c1},
		36: {0x07b82, 0x0b6d9},
		37: {0x10cbb, 0x141b4, 0x15762},
		38: {0x08c3d, 0x0b769, 0x0d380},
		41: {0x084a4},
		42: {0x08757, 0x099a7, 0x09e51, 0x0a3bc, 0x0a9c2},
		43: {0x08fd7},
		44: {0x0527e},
		45: {0x08816, 0x0ce0b, 0x077c1, 0x0d9b7},
		48: {0x07275, 0x08786},
		49: {0x05f55},
		51: {0x09bd5, 0x0bdaf},
		53: {0x100fc, 0x16a63, 0x18c1f, 0x12984, 0x14781, 0x11af4, 0x1239f},
		57: {0x072fb},
		58: {0x18588},
		59: {0x14c57, 0x1a576, 0x1b276, 0x1fdc4, 0x18bfe, 0x19954, 0x1fa37},
		60: {0x067c4},
		61: {0x09dc4},
		64: {0x0b395},
		69: {0x07804},
		70: {0x10b86},
		72: {0x0b0c6, 0x0ec6a, 0x065eb},
		73: {0x04728},
		74: {0x04fd8},
		75: {0x07f7b},
		77: {0x0cfef, 0x0d8c3},
		78: {0x08255},
		80: {0x08911},
		83: {0x0d669, 0x0c32b, 0x0e2a6, 0x0f45e, 0x0f791},
		85: {0x0a1d7, 0x0a34a, 0x0bffe, 0x0d28e, 0x0e085},
		87: {0x044c2},
		88: {0x027eb},
	}
	for i, costumeOffset := range expectedCostumeOffsets {
		assert.Equal(t, costumeOffset, index.Rooms[i].CostumeOffsets)
	}

	// TODO: Check directory of objects
	// Still not sure about the meaning of this index. Let's ignore it for now.
}

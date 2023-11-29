# SCUMM v4 Technical Specifications

Learning SCUMM internals is not easy. There are several different versions that drags tons of
breaking changes. And some websites documents some formats assuming they are universal. 

This document is an attempt to write a technical specification of SCUMM v4 as detailed as possible.
This is the version used by The Secret of Monkey Island and Loom VGA and CD version (EGA, floppy
disk version uses SCUMM v3).

## Resource files

In this version, we might expect to find the following resource files:

- `000.LFL`, the index file. 
- `9xx.LFL`, the charset files. Where _xx_ is the charset number. 
- `DISKnn.LEC`, the data file of the n-th floppy disk of the game. E.g., `DISK01.LEC`. 

### Index file

The index file in SCUMM v4 is not XOR-encoded. It is comprised by a sequence of blocks with the following structure:

| Offset | Size | Format      | Description |
| ------ | ---- | ----------- | ----------- |
| 0x00   | 4    | uint32 (LE) | Block size  |
| 0x04   | 2    | string      | Block type  |
| 0x06   | *    |             | Block data  |

The following blocks are present in the index file:

- `RN`, the room names table. It is a sequence of pairs of room number and names, terminated with
  the null character.

  | Offset | Size | Format         | Description      |
  | ------ | ---- | -------------- | ---------------- |
  | 0x00   | 1    | uint8          | Room number      |
  | 0x01   | 9    | string         | Room name        |
  | ...    | ...  | ...            | ...              |
  | ???    | 1    | Null char 0x00 | Termination mark |

- `0R`, the directory of rooms. It is a sequence of pairs of file number and file offset that
  indicates the data resource file where the i-th room can be found, and its offset. The offset is
  typically zero, as it indicates the offset in the data section of the `LF` block data where `RO`
  block can be found. As SCUMM always encode the `LF` blocks by putting the `RO` sub-block in the
  first place, the offset is always zero. However, in the sake of consistency it would be good idea
  to honor this value and use it to locate the `RO` block in the `LF` data, instead of assuming it
  is the first element. See Data Files section below for more details.

  | Offset | Size | Format      | Description             |
  | ------ | ---- | ----------- | ----------------------- |
  | 0x00   | 1    | uint8       | File number             |
  | 0x01   | 4    | uint32 (LE) | Offset in LF block data |

- `0S`, the directory of global scripts. It is a sequence of pairs of room number and file offset
  that indicates where the global script can be found in the `LF` block data of that room. See Data
  Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF block) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF block data   |

- `0N`, the directory of sounds. It is a sequence of pairs of room number and file offset that
  indicates where the sounds can be found in the `LF` block data of that room. See Data
  Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF block) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF block data   |

- `0C`, the directory of costumes. It is a sequence of pairs of room number and file offset that
  indicates where the costumes (actually animations) can be found in the `LF` block data of that
  room. See Data Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF block) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF block data   |

- `0O`, the directory of objects. It is a sequence of bytes describing a game object. The first 3
  bytes encode the object class, and the next byte encodes the owner (in its low nibble) and the
  state (in the high nibble). 

  | Offset | Size | Format      | Description                                |
  | ------ | ---- | ----------- | ------------------------------------------ |
  | 0x00   | 3    | uint24 (LE) | Object class                               |
  | 0x03   | 1    | byte        | Owner (low nibble) and state (high nibble) |

Special considerations:

- Rooms, global scripts, sounds, costumes and objects are identified by an integer number. This is
  what the index file is basically describing. The virtual machine use to refer to those resources
  by that ID. Thus, it is a good idea to read the index such as every resource is indexed by its ID.
- It is weird to find the directory of objects there that is not describing where the object is
  located in the resource data file, but just indicating the class, owner and state of those
  objects. It seems this directory was put there just by convenience. 
- The directory of rooms have a fixed size of 100 elements, no matter if the game uses less than
  that. The remaining entries in the directory are filled with zeroes, indicating a room in a disk 0
  and offset 0.
- The directories of global scripts, sounds and costumes have a fixed size of 200, no matter if the
  game uses less than that. The remaning entries of the directory are filled with zeroes, indicating
  room 0 and offset 0.
- The directory of objects have a fixed size of 1000 elements, no matter if the game uses less than
  that. 

### Charset file format

These are the LFL files whose number is above 900. They have the following structure

| Offset | Format      | Description            |
| ------ | ----------- | ---------------------- |
| 0x00   | uint32 (LE) | Charset data size - 11 |
| 0x04   | raw data    | Charset data           |

The actual charset data size is the result from adding 11 to the charset data size field. I still
don't know the reason.

In practice, the LFL file length is 15 bytes more than the charset data size indicated in the
header. That is the 11 bytes needed to adjust the size plus 4 bytes from the size field itself.

An example of charset file decoding can be found in the source file
[engines/scumm/resource_v4.cpp][1] of ScummVM.

TBD: meaning of the charset data interpreted from https://wiki.scummvm.org/index.php?title=SCUMM/Technical_Reference/Charset_resources#SCUMM_V4. 

### Data files

#### Encoding

The data files are encoded with an XOR operation with the value 0x69. This means every single byte
from the file have to be XORed with value 0x69 to obtain the actual data. An example of this can be
found in the source file [engines/scumm/resource.cpp][2] of ScummVM.

Apart from that, the data format consists in chunks of blocks, each one with the following header:

| Offset | Format          | Description                  |
| ------ | --------------- | ---------------------------- |
| 0x00   | uint32 (LE)     | Block size (header included) |
| 0x04   | string(2 bytes) | Block name                   |

One particularity introduced in v4 is that the blocks in the data files form a tree of resources.
This means that some blocks are embed into other blocks.

#### Structure

The outer tree structure is fixed, having the following blocks:

- `LE`, the root file container.
  - One single `FO` block, the room offset container. This is where room blocks are listed.
  - Repeated `LF` blocks, containing room, scripts, costumes, sounds, etc. This is the equivalent of
    a LFL file in previous versions of SCUMM.

In other words, the data file can be seen as an archive of LFL files that includes an index to list
the offsets of every file in the archive. The `LF` block has a similar structure of a LFL file in
SCUMM v3.

#### LE blocks

LE blocks are pure containers. They do not have any other data other than the initial `FO` and the
LFL files encoded in multiple `LF` blocks.

#### FO blocks

The FO block is an index of LF blocks included in the LE archive. It has the following structure:

| Offset | Format      | Description                          |
| ------ | ----------- | ------------------------------------ |
| 0x00   | uint8       | Number of LE blocks                  |
| +0x00  | uint8       | LF ID (or room ID)                   |
| +0x01  | uint32 (LE) | Absolute file offset of the LF block |



[1]: https://github.com/scummvm/scummvm/blob/master/engines/scumm/resource_v4.cpp#L175
[2]: https://github.com/scummvm/scummvm/blob/master/engines/scumm/resource.cpp#L105
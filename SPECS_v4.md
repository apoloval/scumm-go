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

The index file in SCUMM v4 is not XOR-encoded. It is comprised by a sequence of chunks with the following structure:

| Offset | Size | Format      | Description |
| ------ | ---- | ----------- | ----------- |
| 0x00   | 4    | uint32 (LE) | Chunk size  |
| 0x04   | 2    | string      | Chunk type  |
| 0x06   | *    |             | Chunk data  |

The following chunks are present in the index file:

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
  typically zero, as it indicates the offset in the data section of the `LF` chunk data where `RO`
  chunk can be found. As SCUMM always encode the `LF` chunks by putting the `RO` sub-chunk in the
  first place, the offset is always zero. However, in the sake of consistency it would be good idea
  to honor this value and use it to locate the `RO` chunk in the `LF` data, instead of assuming it
  is the first element. See Data Files section below for more details.

  | Offset | Size | Format      | Description             |
  | ------ | ---- | ----------- | ----------------------- |
  | 0x00   | 1    | uint8       | File number             |
  | 0x01   | 4    | uint32 (LE) | Offset in LF chunk data |

- `0S`, the directory of global scripts. It is a sequence of pairs of room number and file offset
  that indicates where the global script can be found in the `LF` chunk data of that room. See Data
  Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF chunk) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF chunk data   |

- `0N`, the directory of sounds. It is a sequence of pairs of room number and file offset that
  indicates where the sounds can be found in the `LF` chunk data of that room. See Data
  Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF chunk) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF chunk data   |

- `0C`, the directory of costumes. It is a sequence of pairs of room number and file offset that
  indicates where the costumes (actually animations) can be found in the `LF` chunk data of that
  room. See Data Files section below for more details.

  | Offset | Size | Format      | Description               |
  | ------ | ---- | ----------- | ------------------------- |
  | 0x00   | 1    | uint8       | Room number (or LF chunk) |
  | 0x01   | 4    | uint32 (LE) | Offset in LF chunk data   |

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

| Offset | Size | Format      | Description            |
| ------ | ---- | ----------- | ---------------------- |
| 0      | 4    | uint32 (LE) | Charset data size - 11 |
| 4      | 2    | uint16 (LE) | Magic number 0x0363    |
| 6      | 15   | bytes       | Color map              |
| 21     | 1    | byte        | Bits per pixel         |
| 22     | 1    | byte        | Font height            |
| 23     | 1    | uint16 (LE) | Number of characters   |
| 24     | 1024 | uint32 (LE) | Character offset       |

The actual charset data size is the result from adding 11 to the charset data size field. I still
don't know the reason. In practice, the LFL file length is 15 bytes more than the charset data size
indicated in the header. That is the 11 bytes needed to adjust the size plus 4 bytes from the size
field itself.

It is assumed the second element is a magic number. Other sources such as the ScummVM documentation
says this is just an unused field. At least, all the charset files in The Secret of Monkey Island
and Loom VGA version have the same value there. So it can be used as a magic number.

The character offset is respect the end of the color map. In SCUMM v4 format, this means the byte 21
of the file. For each character, a record can be found at offset+21 with the form:

| Offset | Size            | Format | Description        |
| ------ | --------------- | ------ | ------------------ |
| 0      | 1               | byte   | Width in pixels    |
| 1      | 1               | byte   | Height in pixels   |
| 3      | 1               | byte   | X offset in pixels |
| 4      | 1               | byte   | Y offset in pixels |
| 5      | W * H * BPP / 8 | bytes  | Glypth data        |

The record indicates the size in pixels of the gryph, the offset when it is rendered on the screen
(if any), and the data. 

The data is encoded in left-to-right, top-to-bottom order with top bits of each byte as the first
pixel. For example:

```
For 1 bit per pixel:
Bit position:  7      0 7      0 ...
Words of data: 01234567 89ABCDEF

For 2 bits per pixel:
Bit position:  7      0 7      0 ...
Words of data: 00112233 44556677

For 4 bits per pixel:
Bit position:  7      0 7      1 ...
Words of data: 00001111 22223333

For 8 bits per pixel:
Bit position:  7      0 7      1 ...
Words of data: 00000000 11111111
```

Special considerations:

- It looks like some character data contains one extra byte after the glyph bitstream. One could
  expect that character records are written one after other, so the offset of one character is the
  address of the next byte after the last byte of the previous character glyph data. However, for
  some random characters this is not true, having an extra byte at the end of glyph data. I could
  not find any pattern for this. Perhaps a memory alignment clause used in the code that generated
  the file.

### Data files

#### Encoding

The data files are encoded with an XOR operation with the value 0x69. This means every single byte
from the file have to be XORed with value 0x69 to obtain the actual data. An example of this can be
found in the source file [engines/scumm/resource.cpp][2] of ScummVM.

Apart from that, the data format consists in chunks of data, each one with the following header:

| Offset | Format          | Description                  |
| ------ | --------------- | ---------------------------- |
| 0x00   | uint32 (LE)     | Chunk size (header included) |
| 0x04   | string(2 bytes) | Chunk name                   |

One particularity introduced in v4 is that the chunks in the data files form a tree of resources.
This means that some chunks are embed into other chunks.

#### Structure

The outer tree structure is fixed, having the following chunks:

- `LE`, the root file container.
  - One single `FO` chunk, the room offset container. This is where room chunks are listed.
  - Repeated `LF` chunks, containing room, scripts, costumes, sounds, etc. This is the equivalent of
    a LFL file in previous versions of SCUMM.

In other words, the data file can be seen as an archive of LFL files that includes an index to list
the offsets of every file in the archive. The `LF` chunk has a similar structure of a LFL file in
SCUMM v3.

#### LE chunks

LE chunks are pure containers. They do not have any other data other than the initial `FO` and the
LFL files encoded in multiple `LF` chunks.

#### FO chunks

The FO chunk is an index of LF chunks included in the LE archive. It has the following structure:

| Offset | Format      | Description                          |
| ------ | ----------- | ------------------------------------ |
| 0x00   | uint8       | Number of LE chunks                  |
| +0x00  | uint8       | LF ID (or room ID)                   |
| +0x01  | uint32 (LE) | Absolute file offset of the LF chunk |

#### LF chunks

The LF is a container chunk. It essentially represents a LFL file in previos versions of SCUMM. Now,
instead of having a bunch of LFL files in a floppy disk, they are gathered together in a data
resource file and written in multiple LF chunks in that file.

The chunk starts with a uint16 value that describes the room ID it contains. Please take care of
this before processing any sub-chunk, as it is a uncommon that chunks have any content before their
children. After that, the following sub-chunks are defined in order:

- One `RO` container chunk, the room chunk. It describes the room elements, such as palette,
  bitmaps, objects, local scripts, etc. 
- Zero or more `SC` chunks, the global script chunks. They describe scripts that are not directly
  associated to the room.
- Zero or more `SO` container chunks, the sound chunks. They contain sound resources. 
- Zero or more `CO` chunks, the costume chunks. They contain animations.


#### RO chunks

The RO is a container chunk. It describes the elements of one room.

The following sub-chunks can be found in this order:

- One `HD` chunk, containing a header that describes the basic properties of the room such as its
  size and the number of objects.
- One `CC` chunk, containing the color cycle of the room. Its meaning is still unknown to me. 
- One `SP` chunk, whose contents are still unknown to me.
- One `BX` chunk, containing the boxes and boundaries of the room scene.
- One `PA` chunk, containing the color palette used by the room.
- One `SA` chunk, whose contents are still unknown to me.
- One `BM` chunk, containing the bitmap of the room background.
- Zero or more `OI` chunks, each one describing one room object. Expect as many as objects declared in the `HD` chunk.
- One `NL` chunk, whose contents are still unknown to me. Suspect a sort of script.
- One `SL` chunk, whose contents are still unknown to me. Suspect a sort of script.
- Zero or more `OC` chunks, each one describing the script of one room object. Expect as many as
  objects declared in the `HD` chunk.
- One `EX` chunk, containing the room exit script that will be executed when the player gets out
  from the room.
- One `EN` chunk, containing the room enter script that will be executed when the player gets into
  the room.
- One `LC` chunk, containing a descriptor that tells how many local script chunks can be found after
  this chunk.
- Zero or more `LS` chunks, containg the scripts local to the room.

#### PA chunks

The chunk body starts with a header describing the palette:

| Offset | Size | Format | Description                           |
| ------ | ---- | ------ | ------------------------------------- |
| 0      | 2    | uint16 | Palette size in bytes (typically 768) |


After that, for each possible byte value, the RGB components are described:

| Offset | Size | Format | Description               |
| ------ | ---- | ------ | ------------------------- |
| 0      | 1    | byte   | R (red) color component   |
| 1      | 1    | byte   | G (green) color component |
| 2      | 1    | byte   | B (blue) color component  |

Typically, the size of a PA chunk must be 776 bytes. That's 768 bytes from the RGB codes, plus 2
bytes from the header word, plus 6 bytes from the chunk header.

#### LC chunks

The LC chunk is a local script count descriptor. It is used to describe how many local scripts the
room has. And hence, how many LS chunks can be found after it in the RO chunk content.

The chunk body has the following structure:

| Offset | Size | Format | Description               |
| ------ | ---- | ------ | ------------------------- |
| 0      | 1    | byte   | Number of local scripts   |
| 1      | 1    | byte   | Padding byte, typically 0 |


#### LS chunks

LS chunks describe a local script. The start with a byte indicating the ID of the script, which is
always greater or equal to 200.

| Offset | Size | Format | Description       |
| ------ | ---- | ------ | ----------------- |
| 0      | 1    | byte   | Script ID (>=200) |
| 1      | n    |        | Bytecode          |

The last byte of the bytecode is typically `$A0`, one of the opcodes for `stopObjectCode`. The
instruction used to terminate an script.

#### SC chunks

The SC chunks are pure bytecode. They do not have any header of any sort. The whole chunk body is
full of bytecode.

The last byte is typically `$A0`, one of the opcodes for `stopObjectCode`. The instruction used to
terminate an script.

## Virtual Machine

### Bootscript

Once initialization is completed, the virtual machine executes the global script with ID 1, which is
known as the Bootscript. This script will initialize the game logic, loading the room that shows the
game intro and so.

[1]: https://github.com/scummvm/scummvm/blob/master/engines/scumm/resource_v4.cpp#L175
[2]: https://github.com/scummvm/scummvm/blob/master/engines/scumm/resource.cpp#L105
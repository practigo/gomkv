package gomkv

import (
	"fmt"
	"io"
	"os"
)

// Range is the byte-address range of the memory [Start, End).
type Range struct {
	Start int64
	End   int64 // not inclusive
}

func (r Range) Size() int64 {
	return r.End - r.Start
}

// Element ...
// [ID...Size...Data)
// |            |
// At           Offset
// [--   Size    -->)
type Element struct {
	ID     uint32 // predefined hex
	At     int64  // offset to file start
	Size   int64  // total size including "header"
	Offset int64  // data start relative to At
	Level  uint   // a.k.a., depth

	Children Elements
}

// DataRange ...
func (e *Element) DataRange() Range {
	return Range{e.At + e.Offset, e.At + e.Size}
}

// Elements ...
type Elements []*Element

// Iter loops thru the element-tree in a depth-first manner.
func (es Elements) Iter(f func(*Element) error) error {
	for _, e := range es {
		if err := f(e); err != nil {
			return err
		}
		if len(e.Children) > 0 {
			if err := e.Children.Iter(f); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetVarN ...
/*
bits, big-endian
1xxx xxxx                                                                              - value 0 to  2^7-2
01xx xxxx  xxxx xxxx                                                                   - value 0 to 2^14-2
001x xxxx  xxxx xxxx  xxxx xxxx                                                        - value 0 to 2^21-2
0001 xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx                                             - value 0 to 2^28-2
0000 1xxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx                                  - value 0 to 2^35-2
0000 01xx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx                       - value 0 to 2^42-2
0000 001x  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx            - value 0 to 2^49-2
0000 0001  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx - value 0 to 2^56-2
*/
func GetVarN(b byte) int64 {
	for i := 7; i >= 0; i-- {
		if b&(0x1<<i) > 0 {
			return int64(8 - i)
		}
	}
	// should never happen
	return 0
}

// ParseVarBits parses the variable bits for
// class ID (unmasked, as spec providing raw ID)
// or data size (masked, convert to number for calculation).
func ParseVarBits(bs []byte, masked bool) uint64 {
	// fmt.Printf("%08b\n", bs)
	n := len(bs)
	if n < 1 || n > 8 {
		return 0
	}

	acc := uint64(0)
	for i := 0; i < n; i++ {
		if i == 0 && masked {
			acc += ((0x1 << (8 - n)) - 1) & uint64(bs[0]) // only the fisrt byte has mask
		} else {
			acc = acc<<8 + uint64(bs[i])
		}
	}

	return acc
}

// ReadByte ...
func ReadByte(r io.ReaderAt, offset int64) (byte, error) {
	bs := make([]byte, 1)
	if _, err := r.ReadAt(bs, offset); err != nil {
		return 0, nil
	}
	return bs[0], nil
}

// ReadVarBits ...
func ReadVarBits(m io.ReaderAt, offset int64, masked bool) (n int64, v uint64, err error) {
	first, err := ReadByte(m, offset)
	if err != nil {
		return
	}

	n = GetVarN(first)
	// fmt.Printf("%08b %d\n", first, n)
	bs := make([]byte, n)
	if _, err = m.ReadAt(bs, offset); err != nil {
		return
	}

	v = ParseVarBits(bs, masked)

	return
}

// ReadElement ...
func ReadElement(m io.ReaderAt, offset int64) (e *Element, err error) {
	nC, id, err := ReadVarBits(m, offset, false)
	if err != nil {
		err = fmt.Errorf("read class ID: %w", err)
		return
	}

	nS, size, err := ReadVarBits(m, offset+nC, true)
	if err != nil {
		err = fmt.Errorf("read data size: %w", err)
		return
	}
	return &Element{
		ID:     uint32(id),
		At:     offset,
		Offset: nC + nS,
		Size:   int64(size) + nC + nS,
	}, nil
}

func getElements(m io.ReaderAt, r Range, level uint, nested map[uint32]bool) (es Elements, err error) {
	var (
		e   *Element
		cur = r.Start
	)

	for {
		e, err = ReadElement(m, cur)
		if err != nil {
			return
		}
		e.Level = level
		if isNested, ok := nested[e.ID]; ok && isNested {
			e.Children, err = getElements(m, e.DataRange(), level+1, nested)
		}
		es = append(es, e)
		// update
		cur += e.Size
		if cur >= r.End {
			break
		}
	}

	return
}

// EBMLReader reads the elements from an EBML file.
type EBMLReader interface {
	// GetElements return the element-tree recursively w.r.t. class-id map nested.
	GetElements() (Elements, error)

	// ReadData reads the data bytes pointed by the element.
	ReadData(e *Element) ([]byte, error)
}

type mkvReader struct {
	r      io.ReaderAt
	s      int64
	nested map[uint32]bool
}

func (er *mkvReader) GetElements() (Elements, error) {
	return getElements(er.r, Range{0, er.s}, 0, er.nested)
}

func (er *mkvReader) ReadData(e *Element) (bs []byte, err error) {
	r := e.DataRange()
	s := r.Size()
	if s > 0 {
		bs = make([]byte, s)
		_, err = er.r.ReadAt(bs, r.Start)
	}
	return
}

// DefaultNested ...
func DefaultNested() map[uint32]bool {
	return map[uint32]bool{
		ElementEMBL:    true,
		ElementSegment: true,
		ElementInfo:    true,
		ElementTracks:  true,
		ElementTags:    true,
		TrackEntry:     true,
		Tag:            true,
	}
}

// Open opens a MKV file and return the reader.
func Open(path string) (r EBMLReader, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &mkvReader{
		r:      file,
		s:      info.Size(),
		nested: DefaultNested(),
	}, nil
}

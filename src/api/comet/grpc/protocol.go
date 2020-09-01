package grpc

import (
	"bdim/src/pkg/bytes"
	"bdim/src/pkg/encoding/binary"
)

const (
	// MaxBodySize max package body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_opSize        = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize // + _opSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	//_opOffset     = _verOffset
	// _seqOffset    = _opOffset + _opSize
	// _heartOffset  = _seqOffset
)

// WriteTo write a package to bytes writer.
func (p *Package) WriteTo(b *bytes.Writer) {
	var (
		packLen = _rawHeaderSize + int32(len(p.Body))
		buf     = b.Peek(_rawHeaderSize)
	)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	if p.Body != nil {
		b.Write(p.Body)
	}
}

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package models

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type CursorT struct {
	Size    uint32
	Page    uint32
	Limit   uint32
	Reverse bool
	IdxType uint16
	LastKey []byte
	IdxFrom []byte
	IdxLast []byte
	QueryID []byte
}

func (t *CursorT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	lastKeyOffset := flatbuffers.UOffsetT(0)
	if t.LastKey != nil {
		lastKeyOffset = builder.CreateByteString(t.LastKey)
	}
	idxFromOffset := flatbuffers.UOffsetT(0)
	if t.IdxFrom != nil {
		idxFromOffset = builder.CreateByteString(t.IdxFrom)
	}
	idxLastOffset := flatbuffers.UOffsetT(0)
	if t.IdxLast != nil {
		idxLastOffset = builder.CreateByteString(t.IdxLast)
	}
	queryIDOffset := flatbuffers.UOffsetT(0)
	if t.QueryID != nil {
		queryIDOffset = builder.CreateByteString(t.QueryID)
	}
	CursorStart(builder)
	CursorAddSize(builder, t.Size)
	CursorAddPage(builder, t.Page)
	CursorAddLimit(builder, t.Limit)
	CursorAddReverse(builder, t.Reverse)
	CursorAddIdxType(builder, t.IdxType)
	CursorAddLastKey(builder, lastKeyOffset)
	CursorAddIdxFrom(builder, idxFromOffset)
	CursorAddIdxLast(builder, idxLastOffset)
	CursorAddQueryID(builder, queryIDOffset)
	return CursorEnd(builder)
}

func (rcv *Cursor) UnPackTo(t *CursorT) {
	t.Size = rcv.Size()
	t.Page = rcv.Page()
	t.Limit = rcv.Limit()
	t.Reverse = rcv.Reverse()
	t.IdxType = rcv.IdxType()
	t.LastKey = rcv.LastKeyBytes()
	t.IdxFrom = rcv.IdxFromBytes()
	t.IdxLast = rcv.IdxLastBytes()
	t.QueryID = rcv.QueryIDBytes()
}

func (rcv *Cursor) UnPack() *CursorT {
	if rcv == nil {
		return nil
	}
	t := &CursorT{}
	rcv.UnPackTo(t)
	return t
}

type Cursor struct {
	_tab flatbuffers.Table
}

func GetRootAsCursor(buf []byte, offset flatbuffers.UOffsetT) *Cursor {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Cursor{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Cursor) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Cursor) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Cursor) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Cursor) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *Cursor) Page() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Cursor) MutatePage(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *Cursor) Limit() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Cursor) MutateLimit(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *Cursor) Reverse() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Cursor) MutateReverse(n bool) bool {
	return rcv._tab.MutateBoolSlot(10, n)
}

func (rcv *Cursor) IdxType() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Cursor) MutateIdxType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *Cursor) LastKey(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Cursor) LastKeyLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Cursor) LastKeyBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Cursor) MutateLastKey(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *Cursor) IdxFrom(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Cursor) IdxFromLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Cursor) IdxFromBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Cursor) MutateIdxFrom(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *Cursor) IdxLast(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Cursor) IdxLastLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Cursor) IdxLastBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Cursor) MutateIdxLast(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *Cursor) QueryID(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Cursor) QueryIDLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Cursor) QueryIDBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Cursor) MutateQueryID(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func CursorStart(builder *flatbuffers.Builder) {
	builder.StartObject(9)
}
func CursorAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func CursorAddPage(builder *flatbuffers.Builder, page uint32) {
	builder.PrependUint32Slot(1, page, 0)
}
func CursorAddLimit(builder *flatbuffers.Builder, limit uint32) {
	builder.PrependUint32Slot(2, limit, 0)
}
func CursorAddReverse(builder *flatbuffers.Builder, reverse bool) {
	builder.PrependBoolSlot(3, reverse, false)
}
func CursorAddIdxType(builder *flatbuffers.Builder, idxType uint16) {
	builder.PrependUint16Slot(4, idxType, 0)
}
func CursorAddLastKey(builder *flatbuffers.Builder, lastKey flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(lastKey), 0)
}
func CursorStartLastKeyVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CursorAddIdxFrom(builder *flatbuffers.Builder, idxFrom flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(idxFrom), 0)
}
func CursorStartIdxFromVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CursorAddIdxLast(builder *flatbuffers.Builder, idxLast flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(idxLast), 0)
}
func CursorStartIdxLastVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CursorAddQueryID(builder *flatbuffers.Builder, queryID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(queryID), 0)
}
func CursorStartQueryIDVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CursorEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

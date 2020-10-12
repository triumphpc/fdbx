// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package models

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ValueT struct {
	Blob bool
	GZip bool
	Size uint32
	Hash uint64
	Data []byte
}

func (t *ValueT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	DataOffset := flatbuffers.UOffsetT(0)
	if t.Data != nil {
		DataOffset = builder.CreateByteString(t.Data)
	}
	ValueStart(builder)
	ValueAddBlob(builder, t.Blob)
	ValueAddGZip(builder, t.GZip)
	ValueAddSize(builder, t.Size)
	ValueAddHash(builder, t.Hash)
	ValueAddData(builder, DataOffset)
	return ValueEnd(builder)
}

func (rcv *Value) UnPackTo(t *ValueT) {
	t.Blob = rcv.Blob()
	t.GZip = rcv.GZip()
	t.Size = rcv.Size()
	t.Hash = rcv.Hash()
	t.Data = rcv.DataBytes()
}

func (rcv *Value) UnPack() *ValueT {
	if rcv == nil {
		return nil
	}
	t := &ValueT{}
	rcv.UnPackTo(t)
	return t
}

type Value struct {
	_tab flatbuffers.Table
}

func GetRootAsValue(buf []byte, offset flatbuffers.UOffsetT) *Value {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Value{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Value) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Value) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Value) Blob() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Value) MutateBlob(n bool) bool {
	return rcv._tab.MutateBoolSlot(4, n)
}

func (rcv *Value) GZip() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Value) MutateGZip(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func (rcv *Value) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Value) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *Value) Hash() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Value) MutateHash(n uint64) bool {
	return rcv._tab.MutateUint64Slot(10, n)
}

func (rcv *Value) Data(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *Value) DataLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Value) DataBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Value) MutateData(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func ValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func ValueAddBlob(builder *flatbuffers.Builder, Blob bool) {
	builder.PrependBoolSlot(0, Blob, false)
}
func ValueAddGZip(builder *flatbuffers.Builder, GZip bool) {
	builder.PrependBoolSlot(1, GZip, false)
}
func ValueAddSize(builder *flatbuffers.Builder, Size uint32) {
	builder.PrependUint32Slot(2, Size, 0)
}
func ValueAddHash(builder *flatbuffers.Builder, Hash uint64) {
	builder.PrependUint64Slot(3, Hash, 0)
}
func ValueAddData(builder *flatbuffers.Builder, Data flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(Data), 0)
}
func ValueStartDataVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func ValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

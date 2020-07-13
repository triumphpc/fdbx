// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package models

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RowStateT struct {
	XMin uint64
	XMax uint64
	CMin uint32
	CMax uint32
}

func (t *RowStateT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil {
		return 0
	}
	return CreateRowState(builder, t.XMin, t.XMax, t.CMin, t.CMax)
}
func (rcv *RowState) UnPackTo(t *RowStateT) {
	t.XMin = rcv.XMin()
	t.XMax = rcv.XMax()
	t.CMin = rcv.CMin()
	t.CMax = rcv.CMax()
}

func (rcv *RowState) UnPack() *RowStateT {
	if rcv == nil {
		return nil
	}
	t := &RowStateT{}
	rcv.UnPackTo(t)
	return t
}

type RowState struct {
	_tab flatbuffers.Struct
}

func (rcv *RowState) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RowState) Table() flatbuffers.Table {
	return rcv._tab.Table
}

func (rcv *RowState) XMin() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *RowState) MutateXMin(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *RowState) XMax() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(8))
}
func (rcv *RowState) MutateXMax(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(8), n)
}

func (rcv *RowState) CMin() uint32 {
	return rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(16))
}
func (rcv *RowState) MutateCMin(n uint32) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(16), n)
}

func (rcv *RowState) CMax() uint32 {
	return rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(20))
}
func (rcv *RowState) MutateCMax(n uint32) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(20), n)
}

func CreateRowState(builder *flatbuffers.Builder, XMin uint64, XMax uint64, CMin uint32, CMax uint32) flatbuffers.UOffsetT {
	builder.Prep(8, 24)
	builder.PrependUint32(CMax)
	builder.PrependUint32(CMin)
	builder.PrependUint64(XMax)
	builder.PrependUint64(XMin)
	return builder.Offset()
}
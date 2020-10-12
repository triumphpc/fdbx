package db_test

import (
	"encoding/binary"
	"testing"

	"github.com/shestakovda/errx"
	"github.com/shestakovda/fdbx/v2"
	"github.com/shestakovda/fdbx/v2/db"
	"github.com/stretchr/testify/suite"
)

const TestDB byte = 0x12

// TestInterface - внешние тесты библиотеки
func TestInterface(t *testing.T) {
	suite.Run(t, new(InterfaceSuite))
}

type InterfaceSuite struct {
	suite.Suite
}

func (s *InterfaceSuite) TestConnection() {
	if _, err := db.ConnectV610(0xFF); s.Error(err) {
		s.True(errx.Is(err, db.ErrConnect))
	}

	cn, err := db.ConnectV610(
		TestDB,
		db.ClusterFile(""),
	)
	s.Require().NoError(err)
	s.Require().NoError(cn.Clear())
	s.Equal(TestDB, cn.DB())

	var buf [8]byte
	var val fdbx.Value
	var pair fdbx.Pair
	var lgt db.ListGetter
	const num uint64 = 123
	const add uint64 = 321
	binary.LittleEndian.PutUint64(buf[:], num)

	key1 := fdbx.Key("key1")
	key2 := fdbx.Key("key2")
	key3 := fdbx.Key("key3")

	s.Require().NoError(cn.Write(func(w db.Writer) error {
		w.Upsert(fdbx.NewPair(key1, fdbx.Value("val1")))
		w.Upsert(fdbx.NewPair(key2, fdbx.Value(buf[:])))
		w.Versioned(key3)
		return nil
	}))

	s.Require().NoError(cn.Read(func(r db.Reader) error {
		if pair, err = r.Data(key1); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Equal("val1", val.String())
			}
		}
		if pair, err = r.Data(key2); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Equal(num, binary.LittleEndian.Uint64(val))
			}
		}
		if pair, err = r.Data(key3); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Len(val, 10)
			}
		}
		return nil
	}))

	s.Require().NoError(cn.Write(func(w db.Writer) error {
		w.Upsert(fdbx.NewPair(key1, fdbx.Value("val2")))
		w.Increment(key2, add)
		w.Delete(key3)
		return nil
	}))

	s.Require().NoError(cn.Read(func(r db.Reader) error {
		if pair, err = r.Data(key1); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Equal("val2", val.String())
			}
		}
		if pair, err = r.Data(key2); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Equal(num+add, binary.LittleEndian.Uint64(val))
			}
		}
		if pair, err = r.Data(key3); s.NoError(err) {
			if val, err = pair.Value(); s.NoError(err) {
				s.Empty(val)
			}
		}
		if lgt, err = r.List(nil, nil, 0, false); s.NoError(err) {
			s.Len(lgt(), 2)
		}
		return nil
	}))
}

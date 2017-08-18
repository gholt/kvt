package kvt_test

import (
	"testing"

	"github.com/gholt/kvt"
)

func TestValueTimestampUnmarshalJunk(t *testing.T) {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`ksjdflksjdf;lsdkjf;`))
	if err == nil {
		t.Fatal(err)
	}
}

func TestValueTimestampUnmarshalJunk2(t *testing.T) {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`[1,2,3]`))
	if err == nil || err.Error() != "expected [value,timestamp] from: [1,2,3]" {
		t.Fatal(err)
	}
}

func TestValueTimestampUnmarshalJunk3(t *testing.T) {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`[1,2]`))
	if err == nil || err.Error() != "invalid value from: [1,2]" {
		t.Fatal(err)
	}
}

func TestValueTimestampUnmarshalJunk4(t *testing.T) {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`["one","two"]`))
	if err == nil || err.Error() != `invalid timestamp from: ["one","two"]` {
		t.Fatal(err)
	}
}

func TestValueTimestampUnmarshalJunk5(t *testing.T) {
	vt := &kvt.ValueTimestamp{}
	err := vt.UnmarshalJSON([]byte(`["one",2.1]`))
	if err == nil || err.Error() != `invalid timestamp from: ["one",2.1]` {
		t.Fatal(err)
	}
}

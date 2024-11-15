package tcp

import "testing"

func TestEncode(t *testing.T) {
	b := Join(TlvInt8(0, 1),
		TlvInt16(1, -1),
		TlvInt32(2, 2),
		TlvInt64(3, 0),
		TlvString(4, "Hello World"),
	)
	tlvs := GetAll(b)
	if len(tlvs) != 5 {
		t.Logf("expected of number of elems is 4, actual = %v", len(tlvs))
		t.FailNow()
	}
	if v, err := GetTlv(4, b); err != nil || v.GetString() != "Hello World" {
		t.Logf("expected message is Hello World, actual = %v", v.GetString())
		t.FailNow()
	}

	if v, err := GetTlv(1, b); err != nil || v.GetInt16() != -1 {
		t.Logf("expected message is -1, actual = %v", v.GetInt16())
		t.FailNow()
	}

	if v, err := GetTlv(2, b); err != nil || v.GetInt32() != 2 {
		t.Logf("expected message is 2, actual = %v", v.GetInt32())
		t.FailNow()
	}

	if v, err := GetTlv(3, b); err != nil || v.GetInt64() != 0 {
		t.Logf("expected message is 0, actual = %v", v.GetInt64())
		t.FailNow()
	}
}

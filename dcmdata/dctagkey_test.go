package dcmdata

import (
	"strings"
	"testing"
)

func TestSetDcmTagKey(t *testing.T) {
	var v DcmTagKey
	v.Set(0x0001, 0x0001)
	if v.group != 0x0001 {
		t.Error("excepted 0x0001, got ", v.group)
	}
	if v.element != 0x0001 {
		t.Error("excepted 0x0001, got ", v.element)
	}

}

func TestNewDcmTagKey(t *testing.T) {
	var v = NewDcmTagKey()
	if v.group != 0xffff {
		t.Error("excepted 0xffff, got ", v.group)
	}
	if v.element != 0xffff {
		t.Error("excepted 0xffff, got ", v.element)
	}
}

func TestHasValidGroup(t *testing.T) {
	cases := []struct {
		in   DcmTagKey
		want bool
	}{
		{*NewDcmTagKey(), false},
		{DcmTagKey{0x0001, 0x0001}, false},
		{DcmTagKey{0x0003, 0x0001}, false},
		{DcmTagKey{0x0005, 0x0001}, false},
		{DcmTagKey{0x0007, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0001}, true},
	}
	for _, c := range cases {
		got := c.in.HasValidGroup()
		if got != c.want {
			t.Errorf("%v HasValidGroup(), want %v got %v", c.in, c.want, got)
		}
	}
}

func BenchmarkHasValidGroup(b *testing.B) {
	for i := 0x0000; i <= 0xFFFF; i++ {
		var v DcmTagKey
		v.group = uint16(i)
		switch v.group {
		case 1, 3, 5, 7, 0xFFFF:
			if v.HasValidGroup() != false {
				b.Error(i, v, "excepted false, got ", v.HasValidGroup())
			}
		default:
			if v.HasValidGroup() != true {
				b.Error(i, v, "excepted true, got ", v.HasValidGroup())
			}
		}
	}
}

func TestIsGroupLength(t *testing.T) {
	cases := []struct {
		in   DcmTagKey
		want bool
	}{
		{*NewDcmTagKey(), false},
		{DcmTagKey{0x0001, 0x0001}, false},
		{DcmTagKey{0x0001, 0x0000}, false},
		{DcmTagKey{0x0003, 0x0001}, false},
		{DcmTagKey{0x0003, 0x0000}, false},
		{DcmTagKey{0x0005, 0x0001}, false},
		{DcmTagKey{0x0005, 0x0000}, false},
		{DcmTagKey{0x0007, 0x0001}, false},
		{DcmTagKey{0x0007, 0x0000}, false},
		{DcmTagKey{0x0009, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0000}, true},
	}
	for _, c := range cases {
		got := c.in.IsGroupLength()
		if got != c.want {
			t.Errorf("%v IsGroupLength(), want %v got %v", c.in, c.want, got)
		}
	}

}

func BenchmarkIsGroupLength(b *testing.B) {
	for i := 0x0000; i <= 0xFFFF; i++ {
		for j := 0x0000; j <= 0xFFFF; j++ {
			var v DcmTagKey
			v.group = uint16(i)
			v.element = uint16(j)
			if j == 0x0000 {
				switch v.group {
				case 1, 3, 5, 7, 0xFFFF:
					if v.IsGroupLength() != false {
						b.Error(i, v, "excepted false, got ", v.IsGroupLength())
					}
				default:
					if v.IsGroupLength() != true {
						b.Error(i, v, "excepted true, got ", v.IsGroupLength())
					}
				}
			} else {
				switch v.group {
				case 1, 3, 5, 7, 0xFFFF:
					if v.IsGroupLength() != false {
						b.Error(i, v, "excepted false, got ", v.IsGroupLength())
					}
				default:
					if v.IsGroupLength() != false {
						b.Error(i, v, "excepted false, got ", v.IsGroupLength())
					}
				}

			}

		}
	}
}

func TestIsPrivate(t *testing.T) {
	cases := []struct {
		in   DcmTagKey
		want bool
	}{
		{*NewDcmTagKey(), false},
		{DcmTagKey{0x0001, 0x0001}, false},
		{DcmTagKey{0x0002, 0x0001}, false},
		{DcmTagKey{0x0003, 0x0001}, false},
		{DcmTagKey{0x0004, 0x0001}, false},
		{DcmTagKey{0x0005, 0x0001}, false},
		{DcmTagKey{0x0006, 0x0001}, false},
		{DcmTagKey{0x0007, 0x0001}, false},
		{DcmTagKey{0x0008, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0001}, true},
	}
	for _, c := range cases {
		got := c.in.IsPrivate()
		if got != c.want {
			t.Errorf("%v IsPrivate(), want %v got %v", c.in, c.want, got)
		}
	}

}

func BenchmarkIsPrivate(b *testing.B) {
	for i := 0x0000; i <= 0xFFFF; i++ {
		var v DcmTagKey
		v.group = uint16(i)
		switch v.group {
		case 1, 3, 5, 7, 0xFFFF:
			if v.IsPrivate() != false {
				b.Error(i, v, "excepted false, got ", v.IsPrivate())
			}
		default:
			if (v.group & 1) != 0 {
				if v.IsPrivate() != true {
					b.Error(i, v, "excepted true, got ", v.IsPrivate())
				}

			} else {
				if v.IsPrivate() != false {
					b.Error(i, v, "excepted false, got ", v.IsPrivate())
				}

			}
		}
	}
}

func TestIsPrivateReservation(t *testing.T) {
	cases := []struct {
		in   DcmTagKey
		want bool
	}{
		{*NewDcmTagKey(), false},
		{DcmTagKey{0x0001, 0x0001}, false},
		{DcmTagKey{0x0002, 0x0001}, false},
		{DcmTagKey{0x0003, 0x0001}, false},
		{DcmTagKey{0x0004, 0x0001}, false},
		{DcmTagKey{0x0005, 0x0001}, false},
		{DcmTagKey{0x0006, 0x0001}, false},
		{DcmTagKey{0x0007, 0x0001}, false},
		{DcmTagKey{0x0008, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0011}, true},
	}
	for _, c := range cases {
		got := c.in.IsPrivateReservation()
		if got != c.want {
			t.Errorf("%v IsPrivateReservation(), want %v got %v", c.in, c.want, got)
		}
	}

}

func BenchmarkIsPrivateReservation(b *testing.B) {
	for i := 0x0000; i <= 0xFFFF; i++ {
		for j := 0x0000; j <= 0xFFFF; j++ {
			v := DcmTagKey{uint16(i), uint16(j)}
			if v.element >= 0x0010 && v.element <= 0x00FF {
				switch v.group {
				case 1, 3, 5, 7, 0xFFFF:
					if v.IsPrivateReservation() != false {
						b.Error(i, v, "excepted false, got ", v.IsPrivateReservation())
					}
				default:
					if (v.group & 1) != 0 {
						if v.IsPrivateReservation() != true {
							b.Error(i, v, "excepted true, got ", v.IsPrivateReservation())
						}

					} else {
						if v.IsPrivateReservation() != false {
							b.Error(i, v, "excepted false, got ", v.IsPrivateReservation())
						}

					}
				}
			} else {
				switch v.group {
				case 1, 3, 5, 7, 0xFFFF:
					if v.IsPrivateReservation() != false {
						b.Error(i, v, "excepted false, got ", v.IsPrivateReservation())
					}
				default:
					if (v.group & 1) != 0 {
						if v.IsPrivateReservation() != false {
							b.Error(i, v, "excepted false, got ", v.IsPrivateReservation())
						}

					} else {
						if v.IsPrivateReservation() != false {
							b.Error(i, v, "excepted false, got ", v.IsPrivateReservation())
						}

					}
				}

			}

		}
	}

}

func TestHash(t *testing.T) {
	v := DcmTagKey{0x0002, 0x0002}
	if v.Hash() != 131074 {
		t.Error("excepted 131074, got ", v.Hash())
	}
}

func TestToString(t *testing.T) {
	v := NewDcmTagKey()
	if v.ToString() != "(????,????)" {
		t.Error("excepted (????,????), got ", v.ToString())
	}
	v.group = 0x001F
	v.element = 0x002F
	if strings.ToUpper(v.ToString()) != "(001F,002F)" {
		t.Error("excepted (001F,002F), got ", v.ToString())
	}
}

func TestIsSignableTag(t *testing.T) {
	cases := []struct {
		in   DcmTagKey
		want bool
	}{
		{*NewDcmTagKey(), true},
		{DcmTagKey{0x0001, 0x0000}, false},
		{DcmTagKey{0x0001, 0x0001}, false},
		{DcmTagKey{0x0002, 0x0001}, false},
		{DcmTagKey{0x0003, 0x0001}, false},
		{DcmTagKey{0x0004, 0x0001}, false},
		{DcmTagKey{0x0005, 0x0001}, false},
		{DcmTagKey{0x0006, 0x0001}, false},
		{DcmTagKey{0x0007, 0x0001}, false},
		{DcmTagKey{0x0008, 0x0001}, false},
		{DcmTagKey{0x0009, 0x0001}, true},
		{DcmTagKey{0x0009, 0x0011}, true},
		{DcmTagKey{0xFFFA, 0x0011}, false},
		{DcmTagKey{0x4ffe, 0x0011}, true},
		{DcmTagKey{0x4ffe, 0x0001}, false},
		{DcmTagKey{0xfffc, 0xfffc}, false},
		{DcmTagKey{0xfffc, 0x0001}, true},
		{DcmTagKey{0xFFFe, 0xe00d}, false},
		{DcmTagKey{0xFFFe, 0xe0dd}, false},
		{DcmTagKey{0xFFFe, 0x0001}, true},
	}
	for _, c := range cases {
		got := c.in.IsSignableTag()
		if got != c.want {
			t.Errorf("%v IsSignableTag(), want %v got %v", c.in, c.want, got)
		}
	}
}

func BenchmarkIsSignableTag(b *testing.B) {
	for i := 0x0000; i <= 0xFFFF; i++ {
		for j := 0x0000; j <= 0xFFFF; j++ {
			v := DcmTagKey{uint16(i), uint16(j)}

			if v.element == 0x0000 {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if (v.group == 0x0008) && (v.element == 0x0001) {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if v.group < 0x0008 {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if v.group == 0xFFFA {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if (v.group == 0x4FFE) && (v.element == 0x0001) {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if (v.group == 0xFFFC) && (v.element == 0xFFFC) {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else if (v.group == 0xFFFE) && ((v.element == 0xE00D) || (v.element == 0xE0DD)) {
				if v.IsSignableTag() != false {
					b.Error(v, "excepted false, got ", v.IsSignableTag())
				}
			} else {
				if v.IsSignableTag() != true {
					b.Error(v, "excepted true, got ", v.IsSignableTag())
				}
			}
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		base DcmTagKey
		in   DcmTagKey
		want bool
	}{
		{DcmTagKey{0x0010, 0x001F}, DcmTagKey{0x0010, 0x001F}, true},
		{DcmTagKey{0xFFFF, 0x001F}, DcmTagKey{0xFFFF, 0x001F}, true},
		{DcmTagKey{0xFFFF, 0x001F}, DcmTagKey{0xFFFF, 0x001F}, true},
		{DcmTagKey{0xFFFF, 0x001F}, DcmTagKey{0xFFFF, 0x001E}, false},
	}
	for _, c := range cases {
		got := c.base.Equal(c.in)
		if got != c.want {
			t.Errorf("%s Equal(%s)== %v, want %v ", c.base.ToString(), c.in.ToString(), got, c.want)
		}
	}
}

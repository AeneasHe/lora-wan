package semtech

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDatR(t *testing.T) {
	Convey("Given an empty DatR", t, func() {
		var d DatR

		Convey("Then MarshalJSON returns '0'", func() {
			b, err := d.MarshalJSON()
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, "0")
		})

		Convey("Given LoRa=SF7BW125", func() {
			d.LoRa = "SF7BW125"
			Convey("Then MarshalJSON returns '\"SF7BW125\"'", func() {
				b, err := d.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(b), ShouldEqual, `"SF7BW125"`)
			})
		})

		Convey("Given FSK=1234", func() {
			d.FSK = 1234
			Convey("Then MarshalJSON returns '1234'", func() {
				b, err := d.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(b), ShouldEqual, "1234")
			})
		})

		Convey("Given the string '1234'", func() {
			s := "1234"
			Convey("Then UnmarshalJSON returns FSK=1234", func() {
				err := d.UnmarshalJSON([]byte(s))
				So(err, ShouldBeNil)
				So(d.FSK, ShouldEqual, 1234)
			})
		})

		Convey("Given the string '\"SF7BW125\"'", func() {
			s := `"SF7BW125"`
			Convey("Then UnmarshalJSON returns LoRa=SF7BW125", func() {
				err := d.UnmarshalJSON([]byte(s))
				So(err, ShouldBeNil)
				So(d.LoRa, ShouldEqual, "SF7BW125")
			})
		})
	})
}

func TestCompactTime(t *testing.T) {
	Convey("Given the date 'Mon Jan 2 15:04:05 -0700 MST 2006'", t, func() {
		tStr := "Mon Jan 2 15:04:05 -0700 MST 2006"
		ts, err := time.Parse(tStr, tStr)
		So(err, ShouldBeNil)

		Convey("MarshalJSON returns '\"2006-01-02T22:04:05Z\"'", func() {

			b, err := CompactTime(ts).MarshalJSON()
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, `"2006-01-02T22:04:05Z"`)
		})

		Convey("Given the JSON value of the date (\"2006-01-02T22:04:05Z\")", func() {
			s := `"2006-01-02T22:04:05Z"`
			Convey("UnmarshalJSON returns the correct date", func() {
				var ct CompactTime
				err := ct.UnmarshalJSON([]byte(s))
				So(err, ShouldBeNil)
				So(time.Time(ct).Equal(ts), ShouldBeTrue)
			})
		})
	})
}

func TestGetPacketType(t *testing.T) {
	Convey("Given an empty slice []byte{}", t, func() {
		var b []byte

		Convey("Then GetPacketType returns an error (length)", func() {
			_, err := GetPacketType(b)
			So(err, ShouldResemble, errors.New("lorawan/semtech: at least 4 bytes of data are expected"))
		})

		Convey("Given the slice []byte{2, 1, 3, 4}", func() {
			b = []byte{2, 1, 3, 4}
			Convey("Then GetPacketType returns an error (protocol version)", func() {
				_, err := GetPacketType(b)
				So(err, ShouldResemble, errors.New("lorawan/semtech: unknown protocol version"))
			})
		})

		Convey("Given the slice []byte{1, 1, 3, 4}", func() {
			b = []byte{1, 1, 3, 4}
			Convey("Then GetPacketType returns PullACK", func() {
				t, err := GetPacketType(b)
				So(err, ShouldBeNil)
				So(t, ShouldEqual, PullACK)
			})
		})
	})
}

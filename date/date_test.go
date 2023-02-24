package date

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	tm, _ := time.Parse("2006/01/02 15:04:05.000000", "2023/02/25 16:07:22.123789")
	data := []struct {
		input  time.Time
		format string
		expect string
	}{
		{
			input:  tm,
			format: "m-dd/yyyy hh:M:ss",
			expect: "2-25/2023 16:7:22",
		},
		{
			input:  tm,
			format: "mm-dd/yyyy hh:MM:ss",
			expect: "02-25/2023 16:07:22",
		},
		{
			input:  tm,
			format: "mm-dd/yyyy hh:MM:ss.fff",
			expect: "02-25/2023 16:07:22.123",
		},
		{
			input:  tm,
			format: "mm-dd/yyyy hh:MM:ss.ffffff",
			expect: "02-25/2023 16:07:22.123789",
		},
	}
	for _, d := range data {
		if res := Get(d.input, d.format); res != d.expect {
			t.Errorf("%s != %s\n", res, d.expect)
		}
	}
}
func TestGetOther(t *testing.T) {
	now := time.Now()
	if now.Format("2006/01/02 15:04:05") != Get(now, "yyyy/mm/dd HH:MM:SS") {
		t.Errorf("Get Now error")
	}

	if GetString("2006/01/02 15:04:05", "2006/01/02 15:04:05", "yyyy/mm/dd HH:MM:SS") != "2006/01/02 15:04:05" {
		t.Error("GetString error")
	}

	if GetMicro(now.UnixMicro(), "yyyymmdd HH:MM:SS") != GetMilli(now.UnixMilli(), "yyyymmdd HH:MM:SS") {
		t.Error("GetMicro != GetMilli")
	}

}

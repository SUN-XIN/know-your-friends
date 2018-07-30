package helper

import (
	"testing"
	"time"
)

func TestIsInPeriod(t *testing.T) {
	var start, end int64

	// 2018-07-27 23h00 UTC
	start = 1532732400
	// 2018-07-28 06h00 UTC
	end = 1532757600
	if !IsInPeriod(start, end) {
		t.Errorf("Test1 Expect in period")
	}

	// 2018-07-28 01h00 UTC
	start = 1532739600
	// 2018-07-28 07h00 UTC
	end = 1532761200
	if !IsInPeriod(start, end) {
		t.Errorf("Test2 Expect in period")
	}

	// 2018-07-28 09h00 UTC
	start = 1532768400
	// 2018-07-28 10h00 UTC
	end = 1532772000
	if IsInPeriod(start, end) {
		t.Errorf("Test3 Expect not in period")
	}
}

func TestGetLastDays(t *testing.T) {
	// 2018-07-01 19h10
	startDay := time.Unix(1530465000, 0)
	days := GetLastDays(startDay)

	if len(days) != DEFAULT_ROLLING_DAYS {
		t.Errorf("Expect %d days, but get %d days", DEFAULT_ROLLING_DAYS-1, len(days))
		return
	}

	if days[0] != 1530403200 {
		t.Errorf("Expect %d for 2018-07-01 00h00 UTC, but get %d", 1530403200, days[0])
	}

	if days[1] != 1530316800 {
		t.Errorf("Expect %d for 2018-06-30 00h00 UTC, but get %d", 1530316800, days[1])
	}

	if days[2] != 1530230400 {
		t.Errorf("Expect %d for 2018-06-29 00h00 UTC, but get %d", 1530230400, days[2])
	}

	if days[3] != 1530144000 {
		t.Errorf("Expect %d for 2018-06-28 00h00 UTC, but get %d", 1530144000, days[3])
	}

	if days[4] != 1530057600 {
		t.Errorf("Expect %d for 2018-06-27 00h00 UTC, but get %d", 1530057600, days[4])
	}

	if days[5] != 1529971200 {
		t.Errorf("Expect %d for 2018-06-26 00h00 UTC, but get %d", 1529971200, days[5])
	}

	if days[6] != 1529884800 {
		t.Errorf("Expect %d for 2018-06-25 00h00 UTC, but get %d", 1529884800, days[6])
	}
}

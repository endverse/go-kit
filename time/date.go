package time

import (
	"time"
	gotime "time"
)

func Today() Time {
	return Time{gotime.Now()}
}

func Yesterday() Time {
	return Time{gotime.Now().AddDate(0, 0, -1)}
}

func Tomorrow() Time {
	return Time{gotime.Now().AddDate(0, 0, 1)}
}

func AnchorTime(years int, months int, days int) Time {
	return Time{gotime.Now().In(Location()).AddDate(years, months, days)}
}

// ********************* The days of the week. ********************* //
func Monday() Time {
	offset := int(gotime.Monday - gotime.Now().Weekday())
	if offset > 0 {
		offset = -6
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Tuesday() Time {
	offset := int(gotime.Tuesday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Wednesday() Time {
	offset := int(gotime.Wednesday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Thursday() Time {
	offset := int(gotime.Thursday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Friday() Time {
	offset := int(gotime.Friday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Saturday() Time {
	offset := int(gotime.Saturday - gotime.Now().Weekday())
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func Sunday() Time {
	offset := int(gotime.Sunday - gotime.Now().Weekday())
	if offset < 0 {
		offset = 7 - int(gotime.Now().Weekday())
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func LastMonday() Time {
	return Time{Monday().AddDate(0, 0, -7)}
}

func LastSunday() Time {
	offset := int(gotime.Sunday - gotime.Now().Weekday())
	if offset == 0 {
		offset = 7 - int(gotime.Now().Weekday())
	}
	return Time{gotime.Now().AddDate(0, 0, offset)}
}

func NowInLocation(location *time.Location) Time {
	return Time{gotime.Now().In(location)}
}

func Now() Time {
	return Time{gotime.Now().In(CST())}
}

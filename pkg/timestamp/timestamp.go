package timestamp

import "kms/internal/enums"

// FormatValue formatValue
type FormatValue = enums.Enum
type goFormatValue = enums.Enum

// all format value
var (
	FormatSet = enums.NewEnumSet(nil)
	YearFull  = FormatSet.Reg("yyyy")
	YearShort = FormatSet.Reg("yy")

	MonthFull       = FormatSet.Reg("MMMM")
	MonthAbbr       = FormatSet.Reg("MMM")
	MonthZeroPadded = FormatSet.Reg("MM")
	MonthShort      = FormatSet.Reg("M")

	DayOfYearZeroPadded  = FormatSet.Reg("dddd")
	DayOfMonthZeroPadded = FormatSet.Reg("dd")
	DayOfMonthShort      = FormatSet.Reg("d")

	DayOfWeekFullName = FormatSet.Reg("DDDD")
	DayOfWeekAbbr     = FormatSet.Reg("DDD")

	TwentyFourHourZeroPadded = FormatSet.Reg("HH")
	TwelveHourZeroPadded     = FormatSet.Reg("hh")
	TwelveHour               = FormatSet.Reg("h")

	AMPMUpper = FormatSet.Reg("A")
	AMPMLower = FormatSet.Reg("a")

	MinuteZeroPadded = FormatSet.Reg("mm")
	Minute           = FormatSet.Reg("m")

	SecondZeroPadded = FormatSet.Reg("ss")
	Second           = FormatSet.Reg("s")
	MicroSecond      = FormatSet.Reg("S")

	TimezoneFullName     = FormatSet.Reg("ZZZ")
	TimezoneWithColon    = FormatSet.Reg("zz")
	TimezoneWithoutColon = FormatSet.Reg("Z")

	// Go format value
	goFormatSet = enums.NewEnumSet(nil)
	GoLongMonth = goFormatSet.Reg("January")
	GoMonth     = goFormatSet.Reg("Jan")
	GoNumMonth  = goFormatSet.Reg("1")
	GoZeroMonth = goFormatSet.Reg("01")

	GoLongWeekDay = goFormatSet.Reg("Monday")
	GoWeekDay     = goFormatSet.Reg("Mon")

	GoDay         = goFormatSet.Reg("2")
	GoZeroDay     = goFormatSet.Reg("02")
	GoZeroYearDay = goFormatSet.Reg("002")

	GoHour       = goFormatSet.Reg("15")
	GoHour12     = goFormatSet.Reg("3")
	GoZeroHour12 = goFormatSet.Reg("03")

	GoMinute     = goFormatSet.Reg("4")
	GoZeroMinute = goFormatSet.Reg("04")

	GoSecond     = goFormatSet.Reg("5")
	GoZeroSecond = goFormatSet.Reg("05")

	GoMicrosecond = goFormatSet.Reg("000000")

	GoLongYear = goFormatSet.Reg("2006")
	GoYear     = goFormatSet.Reg("06")

	GoPM = goFormatSet.Reg("PM")
	Gopm = goFormatSet.Reg("pm")

	GoTZ = goFormatSet.Reg("MST")

	GoISO8601TZ      = goFormatSet.Reg("Z0700")
	GoISO8601ColonTZ = goFormatSet.Reg("Z07:00")
)

var (
	formatValueMap = map[FormatValue]goFormatValue{
		YearFull:                 GoLongYear,
		YearShort:                GoYear,
		MonthFull:                GoLongMonth,
		MonthAbbr:                GoMonth,
		MonthZeroPadded:          GoZeroMonth,
		MonthShort:               GoNumMonth,
		DayOfYearZeroPadded:      GoZeroYearDay,
		DayOfMonthZeroPadded:     GoZeroDay,
		DayOfMonthShort:          GoDay,
		DayOfWeekFullName:        GoLongWeekDay,
		DayOfWeekAbbr:            GoWeekDay,
		TwentyFourHourZeroPadded: GoHour,
		TwelveHourZeroPadded:     GoZeroHour12,
		TwelveHour:               GoHour12,
		AMPMUpper:                GoPM,
		AMPMLower:                Gopm,
		MinuteZeroPadded:         GoZeroMinute,
		Minute:                   GoMinute,
		SecondZeroPadded:         GoZeroSecond,
		Second:                   GoSecond,
		MicroSecond:              GoMicrosecond,
		TimezoneFullName:         GoTZ,
		TimezoneWithColon:        GoISO8601ColonTZ,
		TimezoneWithoutColon:     GoISO8601TZ,
	}
)

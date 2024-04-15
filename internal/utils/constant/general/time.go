package general

import "time"

const (
	FullTimeFormat        string = "2006-01-02 15:04:05"
	DisplayDateTimeFormat string = "02 Jan 2006 15:04:05"
	DateFormat            string = "2006-01-02"
)

const (
	NumJan = iota + 1
	NumFeb
	NumMar
	NumApr
	NumMay
	NumJune
	NumJuly
	NumAug
	NumSep
	NumOct
	NumNov
	NumDec
)

const (
	RomanJan  string = "I"
	RomanFeb  string = "II"
	RomanMar  string = "III"
	RomanApr  string = "IV"
	RomanMay  string = "V"
	RomanJune string = "VI"
	RomanJuly string = "VII"
	RomanAug  string = "VIII"
	RomanSep  string = "IX"
	RomanOct  string = "X"
	RomanNov  string = "XI"
	RomanDec  string = "XII"
)

const (
	AuthCookies string = "auth"
)

const (
	ENVProduction string = "production"
)

const (
	UpdatedBySystem int = 0
)

const (
	Time1Min = 1 * time.Minute
	Time5Min = 5 * time.Minute
	Time1Day = 24 * time.Hour
)

const (
	TimeLocationWIB string = "Asia/Jakarta"
)

const (
	FilterTimeStartUTC string = "07:00:01"
	FilterTimeEndUTC   string = "06:59:59"
)

var (
	MonthMap = map[string]string{
		"January":   "januari",
		"February":  "februari",
		"March":     "maret",
		"April":     "april",
		"May":       "mei",
		"June":      "juni",
		"July":      "juli",
		"August":    "agustus",
		"September": "september",
		"October":   "oktober",
		"November":  "november",
		"December":  "desember",
	}
)

func ConvertMonthtoRoman(month int) string {
	switch month {
	case NumJan:
		return RomanJan
	case NumFeb:
		return RomanFeb
	case NumMar:
		return RomanMar
	case NumApr:
		return RomanApr
	case NumMay:
		return RomanMay
	case NumJune:
		return RomanJune
	case NumJuly:
		return RomanJuly
	case NumAug:
		return RomanAug
	case NumSep:
		return RomanSep
	case NumOct:
		return RomanOct
	case NumNov:
		return RomanNov
	case NumDec:
		return RomanDec
	}

	return ""
}

func ConvertMonthtoString(month int) string {
	switch month {
	case NumJan:
		return "January"
	case NumFeb:
		return "February"
	case NumMar:
		return "March"
	case NumApr:
		return "April"
	case NumMay:
		return "May"
	case NumJune:
		return "June"
	case NumJuly:
		return "July"
	case NumAug:
		return "August"
	case NumSep:
		return "September"
	case NumOct:
		return "October"
	case NumNov:
		return "November"
	case NumDec:
		return "December"
	}

	return ""
}

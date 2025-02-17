package pkgep

import "time"

const (
	DateStringLayout  = "2006-01-02"
	StdDateTimeLayout = "2006-01-02 15:04:05"
)

// Convert date string format of "YYYY-MM-DD" to datetime.
func Atodt(dateString string) (time.Time, error) {

	d, err := time.Parse(DateStringLayout, dateString)
	if err != nil {
		return time.Time{}, err
	}

	// handle BE. and AC.
	if d.Year() > time.Now().Year() {
		d = time.Date(d.Year()-543, d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	}

	return d, nil
}

func Atotime(dateString string) *time.Time {

	if dateString == "" {
		return nil
	}

	d, err := time.Parse(DateStringLayout, dateString)
	if err != nil {
		return nil
	}

	// handle BE. and AC.
	if d.Year() > time.Now().Year() {
		d = time.Date(d.Year()-543, d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	}

	return &d
}

func FormatDateTimeSTD(dt time.Time) string {
	return dt.Format(StdDateTimeLayout)
}

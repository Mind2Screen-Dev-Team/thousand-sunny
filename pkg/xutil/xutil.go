package xutil

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	DATE_LAYOUT = "2006-01-02"
)

// default page: 1
// default size: 25
func CalcOffsetLimit(page, size, defaultPage, defaultSize int) (offset, limit int) {
	if defaultPage <= 0 {
		defaultPage = 1
	}

	if defaultSize <= 0 {
		defaultSize = 25
	}

	if page <= 0 {
		page = defaultPage // default to first page
	}

	if size <= 0 {
		size = defaultSize // default page size
	}

	offset = (page - 1) * size
	limit = size
	return
}

// toSnakeCase converts CamelCase or PascalCase to snake_case.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// GetSnakeCaseKey flattens the URL path into a snake_case key.
func GetSnakeCaseKey(r *http.Request) string {
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, s := range segments {
		segments[i] = toSnakeCase(s)
	}
	return strings.Join(segments, "_")
}

func GetSnakeCaseKeyURL(u *url.URL) string {
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	for i, s := range segments {
		segments[i] = toSnakeCase(s)
	}
	return strings.Join(segments, "_")
}

func PgTimestamptz(t time.Time, v bool) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: v}
}

func PgDate(t time.Time, v bool) pgtype.Date {
	return pgtype.Date{Time: t, Valid: v}
}

func MapIntoSlice[T comparable](m map[T]struct{}) []T {
	i := 0
	s := make([]T, len(m))
	for v := range m {
		s[i] = v
		i++
	}
	return s
}

func ConvertToStartEndTime(startDateStr, endDateStr string, loc *time.Location) (time.Time, time.Time, error) {
	var (
		startDate, err1 = time.ParseInLocation(DATE_LAYOUT, startDateStr, loc)
		endDate, err2   = time.ParseInLocation(DATE_LAYOUT, endDateStr, loc)
	)
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing dates: %v %v", err1, err2)
	}

	if startDate.IsZero() {
		return time.Time{}, time.Time{}, errors.New("start date is zero value")
	}

	if !startDate.Before(endDate) && !startDate.Equal(endDate) {
		return time.Time{}, time.Time{}, errors.New("start date must be before end date")
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)

	return startDate, endDate, nil
}

func PositiveInt64(input int64, def int64) int64 {
	if input <= 0 {
		return def
	}
	return input
}

func CountDaysBetween(start, end time.Time) int {
	// Normalize to midnight to compare date parts only
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	if !start.Before(end) {
		return 0
	}

	// Difference in days, add 1 to include end date
	return int(end.Sub(start).Hours()/24) + 1
}

// FormatFloatDecimal formats a float64 value to 'precision' decimal places
// and parses it back to float64.
func FormatFloatDecimal(val float64, precision int) float64 {
	val = SanitizeFloat(val)
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func SanitizeFloat(f float64) float64 {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return 0
	}
	return f
}

type MonthRange struct {
	StartDate time.Time
	EndDate   time.Time
}

func GenerateMonthlyRanges(startDateStr string, loc *time.Location, max int) ([]MonthRange, error) {
	startDate, err := time.ParseInLocation(DATE_LAYOUT, startDateStr, loc)
	if err != nil {
		return nil, err
	}

	var (
		ranges       = make([]MonthRange, 0)
		start        = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, loc)
		now          = time.Now().In(loc)
		currentMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	)

	// Loop from start month up to and including current month
	for !start.After(currentMonth) {
		if max > 0 && len(ranges) >= max {
			break
		}

		nextMonth := start.AddDate(0, 1, 0)
		end := nextMonth.Add(-time.Nanosecond)

		ranges = append(ranges, MonthRange{StartDate: start, EndDate: end})

		start = nextMonth
	}

	return ranges, nil
}

func NumericToFloat64Safe(n pgtype.Numeric) float64 {
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0.00
	}
	return f.Float64
}

func Float64ToPgNumericSafe(val float64) pgtype.Numeric {
	var num pgtype.Numeric

	// Guard against NaN and Infinity — those can't be represented safely
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return numericZero()
	}

	// Convert float64 to string with 2 decimals precision, e.g. "123.45"
	str := fmt.Sprintf("%.2f", val)

	// ScanScientific requires pointer receiver, scan from string
	err := (&num).ScanScientific(str)
	if err != nil || !num.Valid || num.Int == nil {
		// If error or invalid numeric, return zero numeric
		return numericZero()
	}

	return num
}

// Creates a valid pgtype.Numeric with value 0.00
func numericZero() pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(0), // The integer part (0)
		Exp:   -2,            // The exponent: -2 means two decimal digits (0.00)
		Valid: true,          // Must mark it as valid
	}
}

// ChunkSlice splits a slice of any type into chunks of the given size.
func ChunkSlice[T any](input []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(input); i += chunkSize {
		end := min(i + chunkSize, len(input))
		chunks = append(chunks, input[i:end])
	}
	return chunks
}

package csv

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func makeHeader(r interface{}, delim, sdelim string) []byte {
	st := reflect.TypeOf(r)
	var fields []string
	for x := 0; x < st.NumField(); x++ {
		sf := st.Field(x)
		fn := strings.Split(sf.Tag.Get("csv"), ",")[0]
		if fn == "" {
			fn = sf.Name
		}
		if fn != "-" {
			fields = append(fields, `"`+fn+`"`)
		}
	}
	return []byte(strings.Join(fields, delim))
}

func MarshalCSV(a any, delim, sdelim string) []string {
	var rs []string
	for i := 0; i < reflect.ValueOf(a).Len(); i++ {
		r := reflect.ValueOf(a).Index(i).Interface()
		if i == 0 {
			rs = append(rs, string(makeHeader(r, delim, sdelim)))
		}
		rs = append(rs, makeRow(r, delim, sdelim))
	}
	return rs
}

// MakeRow :: make csv row from record
func makeRow(r interface{}, delim, sdelim string) string {
	// st := reflect.TypeOf(r)
	var fields []string
	re := strings.NewReplacer(";", " ", "\n", " ", "\r", "")
	v := reflect.ValueOf(r)
	st := reflect.TypeOf(r)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i).Interface()
		n := st.Field(i).Tag.Get("csv")
		formt := st.Field(i).Tag.Get("time")
		if formt == "" {
			formt = "2006-01-02 15:04:05"
		}
		if n != "-" {
			switch v := f.(type) {
			case int:
				fields = append(fields, fmt.Sprint(v))
			case string:
				fields = append(fields, sdelim+re.Replace(v)+sdelim)
			case sql.NullString:
				fields = append(fields, sdelim+re.Replace(v.String)+sdelim)
			case float32:
				fields = append(fields, strings.Replace(fmt.Sprintf("%.2f", v), ".", ",", 1))
			case float64:
				fields = append(fields, strings.Replace(fmt.Sprintf("%.2f", v), ".", ",", 1))
			case sql.NullInt32:
				if !v.Valid {
					fields = append(fields, "")
					continue
				}
				fields = append(fields, fmt.Sprint(v.Int32))
			case sql.NullInt64:
				if !v.Valid {
					fields = append(fields, "")
					continue
				}
				fields = append(fields, fmt.Sprint(v.Int64))
			case sql.NullFloat64:
				if !v.Valid {
					fields = append(fields, "")
					continue
				}
				fields = append(fields, strings.Replace(fmt.Sprintf("%.2f", v.Float64), ".", ",", 1))
			case time.Time:
				if v.IsZero() {
					fields = append(fields, "")
					continue
				}
				fields = append(fields, sdelim+v.Format(formt)+sdelim)
			case sql.NullTime:
				if !v.Valid || v.Time.IsZero() {
					fields = append(fields, "")
					continue
				}
				fields = append(fields, sdelim+v.Time.Format(formt)+sdelim)
			case bool:
				if v {
					fields = append(fields, sdelim+"да"+sdelim)
				} else {
					fields = append(fields, sdelim+"нет"+sdelim)
				}
			default:
				fields = append(fields, fmt.Sprintf("%v", v))
			}
		}
	}

	return strings.Join(fields, delim)
}

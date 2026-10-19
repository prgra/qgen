package vase

import (
	"database/sql"
	"testing"
	"time"

	"github.com/prgra/qgen/config"
)

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{name: "international prefix", phone: "+7 (999) 123-45-67", want: "79991234567"},
		{name: "domestic prefix", phone: "8 (999) 123-45-67", want: "79991234567"},
		{name: "without prefix", phone: "999 123-45-67", want: "79991234567"},
		{name: "invalid country code", phone: "+1 999 123-4567", want: ""},
		{name: "invalid length", phone: "999123", want: ""},
		{name: "invalid characters", phone: "tel9991234567", want: ""},
		{name: "empty", phone: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePhone(tt.phone); got != tt.want {
				t.Errorf("normalizePhone(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}

// TestCalcPassport проверяет разбор паспортных номеров разных форматов.
func TestCalcPassport(t *testing.T) {
	issueDate := sql.NullTime{Time: time.Date(2015, 2, 10, 0, 0, 0, 0, time.UTC), Valid: true}
	tests := []struct {
		name      string
		num       string
		issue     sql.NullTime
		dataType  int
		series    string
		docNumber string
		unstruct  string
	}{
		{name: "10 цифр без пробелов, есть дата выдачи", num: "4509123456", issue: issueDate,
			dataType: 0, series: "4509", docNumber: "123456", unstruct: ""},
		{name: "серия и номер с пробелами, есть дата выдачи", num: "40 11 654321", issue: issueDate,
			dataType: 0, series: "4011", docNumber: "654321", unstruct: ""},
		{name: "10 цифр без даты выдачи", num: "4509123456",
			dataType: 1, series: "", docNumber: "", unstruct: "4509123456"},
		{name: "с пробелами без даты выдачи", num: "40 11 654321",
			dataType: 1, series: "", docNumber: "", unstruct: "40 11 654321"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AbonIdentRow{DocNumber: tt.num, PassportUnstruct: tt.num, PassportIssueDate: tt.issue}
			a.Calc(config.Config{})
			if a.PassportDataType != tt.dataType || a.DocSeries != tt.series ||
				a.DocNumber != tt.docNumber || a.PassportUnstruct != tt.unstruct {
				t.Errorf("Calc(%q): type=%d series=%q number=%q unstruct=%q, want type=%d series=%q number=%q unstruct=%q",
					tt.num, a.PassportDataType, a.DocSeries, a.DocNumber, a.PassportUnstruct,
					tt.dataType, tt.series, tt.docNumber, tt.unstruct)
			}
		})
	}
}

// TestSanitize проверяет удаление переводов строк из полей выборки.
func TestSanitize(t *testing.T) {
	a := AbonIdentRow{
		FirstName: "Иван\r\nПетрович",
		LastName:  "Иванов\nвторая строка",
		WifiSSID:  "SSID\rсо второй строки",
		Street:    sql.NullString{String: "Полевая\r\n5", Valid: true},
		WifiMAC:   sql.NullString{String: "0A-0B-0C-0D-0E-0F", Valid: true},
	}
	a.Calc(config.Config{})
	if a.FirstName != "Иван Петрович" || a.LastName != "Иванов вторая строка" ||
		a.WifiSSID != "SSID со второй строки" || a.Street.String != "Полевая 5" {
		t.Errorf("переводы строк не удалены: %q, %q, %q, %q",
			a.FirstName, a.LastName, a.WifiSSID, a.Street.String)
	}
}

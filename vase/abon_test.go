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

// TestCalcContractDates проверяет даты договора и интервала актуальности.
func TestCalcContractDates(t *testing.T) {
	start := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(2019, 3, 4, 0, 0, 0, 0, time.UTC)

	active := AbonIdentRow{Status: 0, ContractDate: start, ContractEndDate: last, ActualTo: last}
	active.Calc(config.Config{})
	if !active.ContractEndDate.IsZero() || !active.ActualTo.IsZero() || !active.ActualFrom.Equal(start) {
		t.Errorf("активный: end=%v actual_from=%v actual_to=%v", active.ContractEndDate, active.ActualFrom, active.ActualTo)
	}

	closed := AbonIdentRow{Status: 1, ContractDate: start, ContractEndDate: last}
	closed.Calc(config.Config{})
	if !closed.ContractEndDate.Equal(last) || !closed.ActualTo.Equal(last) || !closed.ActualFrom.Equal(start) {
		t.Errorf("отключен: end=%v actual_from=%v actual_to=%v", closed.ContractEndDate, closed.ActualFrom, closed.ActualTo)
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

// TestSplitCompany проверяет строки юрлица и его оконечного пользователя.
func TestSplitCompany(t *testing.T) {
	a := AbonIdentRow{
		ID: "500", Login: "romashka", Dogovor: "Д-500",
		ContactPhones: "8 (999) 123-45-67",
		FIORaw:        sql.NullString{String: `ООО "Ромашка"`, Valid: true},
		DocNumber:     "4509 123456", PassportIssueDate: sql.NullTime{Time: time.Date(2015, 2, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		City:      sql.NullString{String: "Москва", Valid: true},
		Street:    sql.NullString{String: "Полевая", Valid: true},
		House:     sql.NullString{String: "5", Valid: true},
		WifiMAC:   sql.NullString{String: "0a:0b:0c:0d:0e:0f", Valid: true},
		CompanyID: 7,
		CompName:  sql.NullString{String: `ООО «Ромашка»`, Valid: true},
		CompINN:   sql.NullString{String: "77 0123 4567", Valid: true},
		CompRepr:  sql.NullString{String: "Петров Пётр Петрович", Valid: true},
	}
	a.Calc(config.Config{Country: "Россия"})
	ul, eu := a.splitCompany()

	if ul.ID != "500" || ul.AbonentType != 1 || ul.ULName != "ООО Ромашка" || ul.ULINN != "7701234567" ||
		ul.ULContact != "Петров Пётр Петрович" || ul.ULContactPhones != "79991234567" {
		t.Errorf("ЮЛ: %+v", ul)
	}
	if ul.LastName != "" || ul.FIOUnstructured != "" || ul.DocNumber != "" || ul.DocSeries != "" {
		t.Errorf("у ЮЛ заполнены ФИО/паспорт: %q %q %q %q", ul.LastName, ul.FIOUnstructured, ul.DocSeries, ul.DocNumber)
	}
	if ul.BillAddrStreet != "Полевая" || ul.InstAddrStreet != "Полевая" {
		t.Errorf("адреса ЮЛ: bill=%q inst=%q", ul.BillAddrStreet, ul.InstAddrStreet)
	}

	if eu.ID != "500_1" || eu.AbonentType != 2 || eu.Login != "romashka" || eu.Dogovor != "Д-500" {
		t.Errorf("оконечный пользователь: id=%q type=%d login=%q dogovor=%q", eu.ID, eu.AbonentType, eu.Login, eu.Dogovor)
	}
	if eu.LastName != "Петров" || eu.FirstName != "Пётр" || eu.MiddleName != "Петрович" || eu.DocSeries != "4509" || eu.DocNumber != "123456" {
		t.Errorf("ФИО/паспорт оконечного пользователя: %q %q %q %q %q", eu.LastName, eu.FirstName, eu.MiddleName, eu.DocSeries, eu.DocNumber)
	}
	if eu.ULName != "" || eu.InstAddrStreet != "" || eu.RegAddrStreet != "Полевая" || eu.BillAddrStreet != "Полевая" || eu.WifiMAC.String != "" {
		t.Errorf("поля оконечного пользователя: ul=%q inst=%q reg=%q bill=%q mac=%q",
			eu.ULName, eu.InstAddrStreet, eu.RegAddrStreet, eu.BillAddrStreet, eu.WifiMAC.String)
	}
}

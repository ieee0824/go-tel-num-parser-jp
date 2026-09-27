package tnp

import "testing"

func TestIsTelNumber(t *testing.T) {
	SetIgnoreTypes()
	t.Cleanup(func() { SetIgnoreTypes() })

	tests := []struct {
		name  string
		input string
		want  bool
		kind  TelType
	}{
		{"fixed line", "03-5321-1111", true, FixedLinePhone},
		{"parentheses", "03(5321)1111", true, FixedLinePhone},
		{"mobile", "090-1234-5678", true, MobilePhone},
		{"new mobile prefix", "060-1234-5678", true, MobilePhone},
		{"data mobile 11 digits", "020-123-45678", true, M2M},
		{"data mobile 14 digits", "0200-12345-67890", true, M2M},
		{"data mobile 14 digits without hyphens", "02001234567890", true, M2M},
		{"paging", "020-412-34567", true, PocketBell},
		{"IP phone", "050-1234-5678", true, IPPhone},
		{"incoming charge without hyphens", "0120123456", true, IncomingCharge},
		{"0800 incoming charge", "0800-123-4567", true, IncomingCharge},
		{"unified number without hyphens", "0570123456", true, UnifiedNumber},
		{"information charge", "0990-123-456", true, InformationCharge},
		{"FMC", "0600-123-4567", true, FMC},
		{"embedded number", "call 03-5321-1111", false, -1},
		{"extra digit", "01201234567", false, -1},
		{"0800 with 10 digits", "0800-123-456", false, -1},
		{"reserved message service", "0170123456", false, -1},
		{"reserved mass call service", "0180123456", false, -1},
		{"0600 is not mobile", "0600-123-456", false, -1},
		{"invalid number", "004-0031", false, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, kind := IsTelNumber(tt.input)
			if got != tt.want || kind != tt.kind {
				t.Errorf("IsTelNumber(%q) = (%v, %v), want (%v, %v)", tt.input, got, kind, tt.want, tt.kind)
			}
		})
	}
}

func TestCropTelNumber(t *testing.T) {
	SetIgnoreTypes()
	t.Cleanup(func() { SetIgnoreTypes() })

	tests := []struct {
		input string
		want  string
	}{
		{"Call 0120123456 or 03-5321-1111", "0120123456"},
		{"Call 03(5321)1111", "03-5321-1111"},
		{"01201234567 then 050-1234-5678", "050-1234-5678"},
	}
	for _, tt := range tests {
		got, err := CropTelNumber(tt.input)
		if err != nil || got != tt.want {
			t.Errorf("CropTelNumber(%q) = (%q, %v), want %q", tt.input, got, err, tt.want)
		}
	}
	if got, err := CropTelNumber("01201234567"); err == nil {
		t.Errorf("CropTelNumber returned %q for a number with an extra digit", got)
	}
}

func TestSetIgnoreTypes(t *testing.T) {
	SetIgnoreTypes(IncomingCharge)
	t.Cleanup(func() { SetIgnoreTypes() })
	if got, _ := IsTelNumber("0120123456"); got {
		t.Fatal("incoming charge number was not ignored")
	}
	SetIgnoreTypes(MobilePhone)
	if got, kind := IsTelNumber("0120123456"); !got || kind != IncomingCharge {
		t.Fatalf("SetIgnoreTypes did not replace the previous selection: (%v, %v)", got, kind)
	}
	SetIgnoreTypes()
	if got, kind := IsTelNumber("090-1234-5678"); !got || kind != MobilePhone {
		t.Errorf("SetIgnoreTypes did not clear the selection: (%v, %v)", got, kind)
	}
}

func TestTelTypeStringOutsideRange(t *testing.T) {
	if got := TelType(-2).String(); got != "not tel number" {
		t.Errorf("TelType(-2).String() = %q", got)
	}
}

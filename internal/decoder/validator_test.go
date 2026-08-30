package decoder

import "testing"

func TestDecodeLegacy12Samples(t *testing.T) {
	tests := []struct {
		name           string
		serial         string
		wantDevice     string
		wantModel      string
		wantModelYear  int
		wantProduction int
		wantWeek       int
	}{
		{
			name:           "MacBook Pro M1",
			serial:         "C02FN15CQ05Q",
			wantDevice:     "MacBook Pro (13-inch, 2020, Apple M1)",
			wantModel:      "MacBookPro17,1",
			wantModelYear:  2020,
			wantProduction: 2021,
			wantWeek:       19,
		},
		{
			name:           "iPod touch 4 DCP9",
			serial:         "C3WF2E79DCP9",
			wantDevice:     "iPod touch (4th generation)",
			wantModel:      "",
			wantModelYear:  0,
			wantProduction: 2011,
			wantWeek:       2,
		},
		{
			name:           "iPod touch 4 DCP7",
			serial:         "C3TDXN3VDCP7",
			wantDevice:     "iPod touch (4th generation)",
			wantModel:      "",
			wantModelYear:  0,
			wantProduction: 2010,
			wantWeek:       52,
		},
		{
			name:           "iPod touch 5 F4K1",
			serial:         "CCOKG09WF4K1",
			wantDevice:     "iPod touch (5th generation)",
			wantModel:      "",
			wantModelYear:  0,
			wantProduction: 2013,
			wantWeek:       13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Decode(tt.serial)

			if !result.Valid {
				t.Fatalf(
					"expected serial to be valid, got invalid: %s",
					result.Reason,
				)
			}

			if result.Product == nil {
				t.Fatalf("expected product to be recognized")
			}

			if result.Product.Name != tt.wantDevice {
				t.Errorf(
					"device = %q, want %q",
					result.Product.Name,
					tt.wantDevice,
				)
			}

			if result.Year != tt.wantProduction {
				t.Errorf(
					"production year = %d, want %d",
					result.Year,
					tt.wantProduction,
				)
			}

			if result.Week != tt.wantWeek {
				t.Errorf(
					"week = %d, want %d",
					result.Week,
					tt.wantWeek,
				)
			}

			if tt.wantModel == "" {
				if result.Model != nil {
					t.Errorf(
						"model = %q, want no model metadata",
						result.Model.Model,
					)
				}

				return
			}

			if result.Model == nil {
				t.Fatalf("expected model metadata")
			}

			if result.Model.Model != tt.wantModel {
				t.Errorf(
					"model = %q, want %q",
					result.Model.Model,
					tt.wantModel,
				)
			}

			if result.Model.ModelYear != tt.wantModelYear {
				t.Errorf(
					"model year = %d, want %d",
					result.Model.ModelYear,
					tt.wantModelYear,
				)
			}
		})
	}
}

func TestDecodeRandomizedSerial(t *testing.T) {
	result := Decode("V43G295RXW")

	if result.Valid {
		t.Fatalf("expected randomized serial to be invalid")
	}

	if result.Format != FormatRandomized {
		t.Fatalf(
			"format = %q, want %q",
			result.Format,
			FormatRandomized,
		)
	}
}

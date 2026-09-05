package domain_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/idoceb00/lorren/internal/domain"
)

func TestNewDailyLog(t *testing.T) {
	fixedDate := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   domain.NewDailyLogInput
		want    *domain.DailyLog
		wantErr bool
	}{
		{
			name: "valid input with typical sleep hours",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				Training:   true,
				SleepHours: 7.5,
				Breakfast:  "oats",
			},
			want: &domain.DailyLog{
				Date:       fixedDate,
				Training:   true,
				SleepHours: 7.5,
				Breakfast:  "oats",
			},
			wantErr: false,
		},
		{
			name: "sleep hours at lower bound",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: 0,
			},
			want: &domain.DailyLog{
				Date:       fixedDate,
				SleepHours: 0,
			},
			wantErr: false,
		},
		{
			name: "sleep hours at upper bound",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: 24,
			},
			want: &domain.DailyLog{
				Date:       fixedDate,
				SleepHours: 24,
			},
			wantErr: false,
		},
		{
			name: "negative sleep hours returns error",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: -1,
			},
			wantErr: true,
		},
		{
			name: "sleep hours above upper bound returns error",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: 25,
			},
			wantErr: true,
		},
		{
			name: "sleep hours round up to two decimals",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: 7.126,
			},
			want: &domain.DailyLog{
				Date:       fixedDate,
				SleepHours: 7.13,
			},
			wantErr: false,
		},
		{
			name: "sleep hours round down to two decimals",
			input: domain.NewDailyLogInput{
				Date:       fixedDate,
				SleepHours: 7.124,
			},
			want: &domain.DailyLog{
				Date:       fixedDate,
				SleepHours: 7.12,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := domain.NewDailyLog(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewDailyLog() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewDailyLog() succeeded unexpectedly")
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewDailyLog() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

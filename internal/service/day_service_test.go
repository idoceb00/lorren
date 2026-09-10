package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/idoceb00/lorren/internal/domain"
	"github.com/idoceb00/lorren/internal/service"
)

// fakeInterviewer is a test double for domain.DailyInterviewer. askFunc lets
// each test case control what AskDailyLog returns. receivedExisting
// records what RecordDay actually passed in, so tests can assert on it.
type fakeInterviewer struct {
	askFunc          func(existing *domain.DailyLog) (*domain.DailyLog, error)
	called           bool
	receivedExisting *domain.DailyLog
}

func (f *fakeInterviewer) AskDailyLog(existing *domain.DailyLog) (*domain.DailyLog, error) {
	f.called = true
	f.receivedExisting = existing
	return f.askFunc(existing)
}

// fakeRepository is a test double for domain.DailyRepository.
type fakeRepository struct {
	findFunc   func(date time.Time) (*domain.DailyLog, error)
	saveFunc   func(log *domain.DailyLog) (string, error)
	saveCalled bool
}

func (f *fakeRepository) FindByDate(date time.Time) (*domain.DailyLog, error) {
	return f.findFunc(date)
}

func (f *fakeRepository) SaveDailyLog(log *domain.DailyLog) (string, error) {
	f.saveCalled = true
	return f.saveFunc(log)
}

func TestRecordDay(t *testing.T) {
	fixedDate := time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)
	existingLog := &domain.DailyLog{Date: fixedDate, Breakfast: "oats"}
	newLog := &domain.DailyLog{Date: fixedDate, Breakfast: "eggs"}

	const savedPath = "/vault/Diario/2026-09-08.md"

	tests := []struct {
		name string

		interviewer *fakeInterviewer
		repo        *fakeRepository
		date        time.Time
		wantErr     bool

		// wantExisting is what we expect RecordDay to have passed into
		// AskDailyLog: nil when there was nothing to find, the found
		// log otherwise. Only checked if the interviewer was called.
		wantExisting *domain.DailyLog
		wantSaved    bool
	}{
		{
			name: "no existing log: interviewer receives nil, result is saved",
			interviewer: &fakeInterviewer{
				askFunc: func(existing *domain.DailyLog) (*domain.DailyLog, error) {
					return newLog, nil
				},
			},
			repo: &fakeRepository{
				findFunc: func(date time.Time) (*domain.DailyLog, error) {
					return nil, domain.ErrNotFound
				},
				saveFunc: func(log *domain.DailyLog) (string, error) { return savedPath, nil },
			},
			date:         fixedDate,
			wantErr:      false,
			wantExisting: nil,
			wantSaved:    true,
		},
		{
			name: "existing log found: interviewer receives it, result is saved",
			interviewer: &fakeInterviewer{
				askFunc: func(existing *domain.DailyLog) (*domain.DailyLog, error) {
					return newLog, nil
				},
			},
			repo: &fakeRepository{
				findFunc: func(date time.Time) (*domain.DailyLog, error) {
					return existingLog, nil
				},
				saveFunc: func(log *domain.DailyLog) (string, error) { return savedPath, nil },
			},
			date:         fixedDate,
			wantErr:      false,
			wantExisting: existingLog,
			wantSaved:    true,
		},
		{
			name: "FindByDate returns a real error: aborts before asking or saving",
			interviewer: &fakeInterviewer{
				askFunc: func(existing *domain.DailyLog) (*domain.DailyLog, error) {
					t.Fatal("AskDailyLog should not be called")
					return nil, nil
				},
			},
			repo: &fakeRepository{
				findFunc: func(date time.Time) (*domain.DailyLog, error) {
					return nil, fmt.Errorf("disk error")
				},
				saveFunc: func(log *domain.DailyLog) (string, error) {
					t.Fatal("SaveDailyLog should not be called")
					return "", nil
				},
			},
			date:    fixedDate,
			wantErr: true,
		},
		{
			name: "interviewer returns an error: does not save",
			interviewer: &fakeInterviewer{
				askFunc: func(existing *domain.DailyLog) (*domain.DailyLog, error) {
					return nil, fmt.Errorf("form cancelled")
				},
			},
			repo: &fakeRepository{
				findFunc: func(date time.Time) (*domain.DailyLog, error) {
					return nil, domain.ErrNotFound
				},
				saveFunc: func(log *domain.DailyLog) (string, error) {
					t.Fatal("SaveDailyLog should not be called")
					return "", nil
				},
			},
			date:    fixedDate,
			wantErr: true,
		},
		{
			name: "SaveDailyLog fails: error propagates",
			interviewer: &fakeInterviewer{
				askFunc: func(existing *domain.DailyLog) (*domain.DailyLog, error) {
					return newLog, nil
				},
			},
			repo: &fakeRepository{
				findFunc: func(date time.Time) (*domain.DailyLog, error) {
					return nil, domain.ErrNotFound
				},
				saveFunc: func(log *domain.DailyLog) (string, error) {
					return "", fmt.Errorf("permission denied")
				},
			},
			date:    fixedDate,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotErr := service.RecordDay(tt.interviewer, tt.repo, tt.date)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RecordDay() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RecordDay() succeeded unexpectedly")
			}

			if tt.interviewer.called && tt.wantExisting != tt.interviewer.receivedExisting {
				t.Errorf("AskDailyLog received existing = %v, want %v", tt.interviewer.receivedExisting, tt.wantExisting)
			}
			if tt.wantSaved && !tt.repo.saveCalled {
				t.Error("expected SaveDailyLog to be called, but it wasn't")
			}
			if tt.wantSaved && gotPath != savedPath {
				t.Errorf("RecordDay() path = %q, want %q", gotPath, savedPath)
			}
		})
	}
}

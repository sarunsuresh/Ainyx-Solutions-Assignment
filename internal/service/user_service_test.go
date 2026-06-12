package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  now.AddDate(-25, -1, 0), // 25 years ago, last month
			want: 25,
		},
		{
			name: "birthday is today",
			dob:  now.AddDate(-30, 0, 0),
			want: 30,
		},
		{
			name: "birthday not yet this year",
			dob:  now.AddDate(-25, 1, 0), // 25 years ago but next month
			want: 24,
		},
		{
			name: "newborn",
			dob:  now,
			want: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CalculateAge(tc.dob); got != tc.want {
				t.Errorf("CalculateAge() = %d, want %d", got, tc.want)
			}
		})
	}
}

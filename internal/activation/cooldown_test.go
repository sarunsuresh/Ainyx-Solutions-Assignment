package activation

import (
    "testing"
    "time"
)

func TestCooldownFor(t *testing.T) {
    cases := []struct{
        resendCount int
        want        time.Duration
    }{
        {0, 1 * time.Minute},
        {1, 3 * time.Minute},
        {4, 60 * time.Minute},
        {7, 24 * time.Hour},
        {99, 24 * time.Hour},  // beyond last — stays at 24h
    }
    for _, tc := range cases {
        got := CooldownFor(tc.resendCount)
        if got != tc.want {
            t.Errorf("CooldownFor(%d) = %v, want %v", tc.resendCount, got, tc.want)
        }
    }
}

func TestCanResend(t *testing.T) {
    now := time.Now()

    // just sent — cannot resend (count=0, cooldown=1min)
    if CanResend(0, now) {
        t.Error("should not be able to resend immediately")
    }

    // sent 2 minutes ago, count=0 (1min cooldown) — can resend
    if !CanResend(0, now.Add(-2*time.Minute)) {
        t.Error("should be able to resend after 1 min cooldown")
    }

    // sent 4 mins ago, count=1 (3min cooldown) — can resend
    if !CanResend(1, now.Add(-4*time.Minute)) {
        t.Error("should be able to resend after 3 min cooldown")
    }
}
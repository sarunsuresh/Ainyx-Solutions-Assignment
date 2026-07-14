package activation

import "time"


var cooldownSequence = []time.Duration{
    1 * time.Minute,
    3 * time.Minute,
    5 * time.Minute,
    10 * time.Minute,
    60 * time.Minute,
    120 * time.Minute,
    180 * time.Minute,
    24 * time.Hour,
}


func CooldownFor(resendCount int) time.Duration {
    if resendCount >= len(cooldownSequence) {
        return cooldownSequence[len(cooldownSequence)-1]
    }
    return cooldownSequence[resendCount]
}


func CanResend(resendCount int, lastResendAt time.Time) bool {
    if lastResendAt.IsZero() {
        return true  // never sent before
    }
    required := CooldownFor(resendCount)
    return time.Since(lastResendAt) >= required
}


func TimeUntilResend(resendCount int, lastResendAt time.Time) time.Duration {
    if CanResend(resendCount, lastResendAt) {
        return 0
    }
    required := CooldownFor(resendCount)
    elapsed := time.Since(lastResendAt)
    return required - elapsed
}
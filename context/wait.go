package context

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func PollImmediate(interval, timeout time.Duration, condition wait.ConditionFunc) error {
	if timeout == 0 {
		return wait.PollImmediateInfinite(interval, condition)
	} else {
		return wait.PollImmediate(interval, timeout, condition)
	}
}

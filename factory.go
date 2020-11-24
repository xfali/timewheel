// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package timewheel

import (
	"fmt"
	"time"
)

const (
	DefaultMaxDuration = 24 * time.Hour
	DefaultMinDuration = 20 * time.Millisecond
	MaxTimerBufferSize = 1000
	MaxPurgeBufferSize = 1000
)

type AsyncOpt func(tw *HieraTimeWheel)

func New(opts ...AsyncOpt) *HieraTimeWheel {
	tw := &HieraTimeWheel{
		maxExpire:  DefaultMaxDuration,
		hieraTimes: parseHieraTimes(DefaultMaxDuration, DefaultMinDuration),
		addSize:    MaxTimerBufferSize,
		rmSize:     MaxPurgeBufferSize,
	}
	for _, opt := range opts {
		opt(tw)
	}
	return tw
}

func parseHieraTimes(max, min time.Duration) []time.Duration {
	hiera := []time.Duration{time.Hour, time.Minute, time.Second}
	if max < min {
		panic("max equal min")
	} else if max == min {
		return []time.Duration{max}
	}

	m, n := 0, len(hiera)
	for ; max < hiera[m]; m++ {
	}
	for ; min > hiera[n-1]; n-- {
	}
	hiera = hiera[m:n]
	if min < hiera[len(hiera)-1] {
		if hiera[len(hiera)-1] % min != 0 {
			panic(fmt.Sprintf("timewheel tick is invalid: %d, min tick: %d", hiera[len(hiera)-1], min))
		}
		hiera = append(hiera, min)
	}
	return hiera
}

func AsyncOptSetDuration(max, min time.Duration) AsyncOpt {
	return func(tw *HieraTimeWheel) {
		tw.maxExpire = max
		tw.hieraTimes = parseHieraTimes(max, min)
	}
}

func AsyncOptSetMinDuration(time time.Duration) AsyncOpt {
	return func(tw *HieraTimeWheel) {
		tw.maxExpire = time
	}
}

func AsyncOptSetMaxTimerBufferSize(size int) AsyncOpt {
	return func(tw *HieraTimeWheel) {
		tw.addSize = size
	}
}

func AsyncOptSetMaxPurgeBufferSize(size int) AsyncOpt {
	return func(tw *HieraTimeWheel) {
		tw.rmSize = size
	}
}

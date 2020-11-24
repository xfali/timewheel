/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12
 * @time 14:00
 * @version V1.0
 * Description:
 */

package timewheel

import (
	"errors"
	"github.com/xfali/goutils/atomic"
	"time"
)

//Hierarchical Timing Wheels
type SyncHieraTimeWheel struct {
	timeWheels []TimeWheel
	hieraTimes []time.Duration
	stop       atomic.AtomicBool
	maxExpire  time.Duration
}

//创建一个通用的时间轮，分层数据格式为：时间由大到小排列，如hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 20*time.Millisecond }
func NewSyncHiera(duration time.Duration, hieraTimes []time.Duration) TimeWheel {
	if len(hieraTimes) < 2 || duration < hieraTimes[0] {
		return nil
	}

	tw := &SyncHieraTimeWheel{}
	deep := len(hieraTimes)

	tw.timeWheels = make([]TimeWheel, deep)

	secondTick := false

	time := duration / hieraTimes[0]
	tw.maxExpire = time * hieraTimes[0]
	if time > 0 {
		secondTick = true
		wheel := NewSyncOne(hieraTimes[0], time*hieraTimes[0])
		tw.timeWheels[0] = wheel
	}

	for j := 1; j < deep; j++ {
		i := j
		time = (duration % hieraTimes[i-1]) / hieraTimes[i]
		if secondTick {
			wheel := NewSyncOne(hieraTimes[i], hieraTimes[i-1])
			wheel.Add(func() {
				tw.timeWheels[i-1].Tick(hieraTimes[i-1])
			}, hieraTimes[i-1], true)
			tw.timeWheels[i] = wheel
		} else {
			if time > 0 {
				secondTick = true
				wheel := NewSyncOne(hieraTimes[i], time*hieraTimes[i])
				tw.timeWheels[i] = wheel
			}
		}
	}

	tw.hieraTimes = hieraTimes
	tw.stop = atomic.AtomicBool(1)
	return tw
}

func (htw *SyncHieraTimeWheel) Start() {
	htw.stop = 0
}

func (htw *SyncHieraTimeWheel) Stop() {
	htw.stop.Set()
}

func (htw *SyncHieraTimeWheel) Tick(duration time.Duration) {
	if htw.stop.IsSet() {
		return
	}
	htw.timeWheels[len(htw.timeWheels)-1].Tick(duration)
}

func (htw *SyncHieraTimeWheel) Add(callback OnTimeout, expire time.Duration, repeat bool) (Timer, error) {
	if expire < htw.hieraTimes[len(htw.hieraTimes)-1] {
		return nil, errors.New("expire time is too small")
	}

	if expire > htw.maxExpire {
		return nil, errors.New("expire time is too large")
	}

	absoluteTime := htw.absoluteTime(expire)
	if repeat {
		return htw.addTime(0, func() {
			callback()
			htw.Add(callback, expire, repeat)
		}, absoluteTime, repeat)
	} else {
		return htw.addTime(0, callback, absoluteTime, repeat)
	}
}

func (htw *SyncHieraTimeWheel) RollTime() time.Duration {
	return 0
}

func (htw *SyncHieraTimeWheel) parse(expire time.Duration) int {
	deep := 0
	nextTime := expire / htw.hieraTimes[deep]
	if nextTime > 0 {
		return deep
	}
	deep++
	for deep < len(htw.hieraTimes) {
		nextTime = expire % htw.hieraTimes[deep-1] / htw.hieraTimes[deep]
		if nextTime > 0 {
			return deep
		}
		deep++
	}
	return deep
}

func (htw *SyncHieraTimeWheel) absoluteTime(expire time.Duration) time.Duration {
	deep := htw.parse(expire)
	deep++
	for deep < len(htw.hieraTimes) {
		expire += htw.timeWheels[deep].RollTime()
		deep++
	}
	return expire
}

func (htw *SyncHieraTimeWheel) addTime(deep int, callback OnTimeout, expire time.Duration, repeat bool) (Timer, error) {
	var nextTime time.Duration
	if deep == 0 {
		nextTime = expire / htw.hieraTimes[deep]
	} else {
		nextTime = expire % htw.hieraTimes[deep-1] / htw.hieraTimes[deep]
	}

	if deep == len(htw.hieraTimes)-1 {
		if nextTime > 0 {
			return htw.timeWheels[deep].Add(func() {
				callback()
			}, nextTime*htw.hieraTimes[deep], false)
		} else {
			callback()
			return nil, nil
		}
	} else {
		if nextTime > 0 {
			return htw.timeWheels[deep].Add(func() {
				htw.addTime(deep+1, callback, expire, repeat)
			}, nextTime*htw.hieraTimes[deep], false)
		} else {
			return htw.addTime(deep+1, callback, expire, repeat)
		}
	}
}

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
	"fmt"
	"github.com/xfali/goutils/atomic"
	"time"
)

type HieraTimer struct {
	TimerData
	tw       *HieraTimeWheel
	timer    Timer
	pastTime time.Duration
	rmFlag   atomic.AtomicBool
}

//Hierarchical Timing Wheels
type HieraTimeWheel struct {
	timeWheels []TimeWheel
	hieraTimes []time.Duration
	stop       chan bool
	addChan    chan *HieraTimer
	rmChan     chan *HieraTimer
	maxExpire  time.Duration

	addSize, rmSize int
}

//创建一个通用的时间轮，分层数据格式为：时间由大到小排列，如hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 20*time.Millisecond }
func NewAsyncHiera(duration time.Duration, hieraTimes []time.Duration, addMax int, rmMax int) *HieraTimeWheel {
	tw := &HieraTimeWheel{
		maxExpire:  duration,
		hieraTimes: hieraTimes,
		addSize:    addMax,
		rmSize:     rmMax,
	}
	return tw
}

func (htw *HieraTimeWheel) init() {
	duration := htw.maxExpire
	// 这里曾经限制为至少两层，目前改为允许1层时间轮
	if len(htw.hieraTimes) < 1 || duration < htw.hieraTimes[0] {
		panic("hieraTimes less than 1 or max duration less than hieraTimes[0]")
	}

	deep := len(htw.hieraTimes)
	htw.timeWheels = make([]TimeWheel, deep)
	htw.stop = make(chan bool)
	htw.addChan = make(chan *HieraTimer, htw.addSize)
	htw.rmChan = make(chan *HieraTimer, htw.rmSize)

	secondTick := false

	time := duration / htw.hieraTimes[0]
	htw.maxExpire = time * htw.hieraTimes[0]

	if time > 0 {
		secondTick = true
		wheel := NewSyncOne(htw.hieraTimes[0], time*htw.hieraTimes[0])
		htw.timeWheels[0] = wheel
	}

	for j := 1; j < deep; j++ {
		i := j
		time = (duration % htw.hieraTimes[i-1]) / htw.hieraTimes[i]
		if secondTick {
			wheel := NewSyncOne(htw.hieraTimes[i], htw.hieraTimes[i-1])
			wheel.Add(func() {
				htw.timeWheels[i-1].Tick(htw.hieraTimes[i-1])
			}, htw.hieraTimes[i-1], true)
			htw.timeWheels[i] = wheel
		} else {
			if time > 0 {
				secondTick = true
				wheel := NewSyncOne(htw.hieraTimes[i], time*htw.hieraTimes[i])
				htw.timeWheels[i] = wheel
			}
		}
	}
}

func (htw *HieraTimeWheel) Start() {
	htw.init()

	go func() {
		now := time.Now()
		cur := now
		tickTime := htw.hieraTimes[len(htw.hieraTimes)-1]
		for {
			//FIXME: 增加timer和tick跳动必须二选一，否则增加的timer会计时不准确。
			//但是当大量同时注册timer时，有可能造成间隔了多个tick才开始回调
			select {
			case <-htw.stop:
				return
			case timer, ok := <-htw.addChan:
				if ok {
					htw.add2Slot(timer)
				}
			default:
				passTime := time.Since(now)
				if passTime < tickTime {
					time.Sleep(tickTime - passTime)
				}
				cur = time.Now()
				htw.Tick(tickTime)
				now = cur
			}
			select {
			case <-htw.stop:
				return
			case rmCh, ok := <-htw.rmChan:
				if ok {
					htw.removeTimer(rmCh)
				}
			default:
			}
		}
	}()
}

func (htw *HieraTimeWheel) Stop() {
	close(htw.stop)
}

func (htw *HieraTimeWheel) Tick(duration time.Duration) {
	htw.timeWheels[len(htw.timeWheels)-1].Tick(duration)
}

func (htw *HieraTimeWheel) Add(callback OnTimeout, expire time.Duration, repeat bool) (Timer, error) {
	if expire < htw.hieraTimes[len(htw.hieraTimes)-1] {
		return nil, fmt.Errorf("expire time is too small: expire %d, min: %d", expire, htw.hieraTimes[len(htw.hieraTimes)-1])
	}
	if expire > htw.maxExpire {
		return nil, fmt.Errorf("expire time is too large: expire %d, max: %d", expire, htw.maxExpire)
	}

	aTimer := &HieraTimer{
		TimerData: TimerData{callback, expire, repeat},
		tw:        htw,
		rmFlag:    0,
	}
	htw.addChan <- aTimer
	return aTimer, nil
}

func (htw *HieraTimeWheel) add2Slot(timer *HieraTimer) {
	absoluteTime := htw.absoluteTime(timer.Expire)
	if timer.rmFlag.IsSet() {
		return
	}
	if timer.Repeat {
		callback := timer.Callback
		timer.Callback = func() {
			callback()
			timer.Callback = callback
			htw.add2Slot(timer)
		}
		htw.addTime(0, timer, absoluteTime)
	} else {
		htw.addTime(0, timer, absoluteTime)
	}
}

func (htw *HieraTimeWheel) removeTimer(timer *HieraTimer) {
	if timer.timer != nil {
		timer.timer.Cancel()
	}
}

func (htw *HieraTimeWheel) RollTime() time.Duration {
	var time time.Duration = 0
	for i := 0; i < len(htw.timeWheels); i++ {
		time += htw.timeWheels[i].RollTime()
	}
	return time
}

func (htw *HieraTimeWheel) parse(expire time.Duration) int {
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

func (htw *HieraTimeWheel) absoluteTime(expire time.Duration) time.Duration {
	deep := htw.parse(expire)
	deep++
	for deep < len(htw.hieraTimes) {
		expire += htw.timeWheels[deep].RollTime()
		deep++
	}
	return expire
}

func (htw *HieraTimeWheel) addTime(deep int, timer *HieraTimer, expire time.Duration) {
	var nextTime time.Duration
	if deep == 0 {
		nextTime = expire / htw.hieraTimes[deep]
	} else {
		nextTime = expire % htw.hieraTimes[deep-1] / htw.hieraTimes[deep]
	}

	if deep == len(htw.hieraTimes)-1 {
		if nextTime > 0 {
			tmpTimer, _ := htw.timeWheels[deep].Add(func() {
				timer.Callback()
			}, nextTime*htw.hieraTimes[deep], false)
			timer.timer = tmpTimer
		} else {
			timer.Callback()
		}
	} else {
		if nextTime > 0 {
			tmpTimer, _ := htw.timeWheels[deep].Add(func() {
				htw.addTime(deep+1, timer, expire)
			}, nextTime*htw.hieraTimes[deep], false)
			timer.timer = tmpTimer
		} else {
			htw.addTime(deep+1, timer, expire)
		}
	}
}

func (aTimer *HieraTimer) Cancel() {
	aTimer.rmFlag.Set()
	aTimer.tw.rmChan <- aTimer
}

func (aTimer *HieraTimer) PastTime() time.Duration {
	//NOTICE:异步时间轮的Tick与Add在同一个select，所以需要+1
	return aTimer.pastTime
}

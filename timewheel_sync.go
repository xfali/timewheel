/**
 * Copyright (C) 2018-2020, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12
 * @time 9:38
 * @version V1.0
 * Description:
 */

package timewheel

import (
	"container/list"
	"errors"
	"github.com/xfali/goutils/atomic"
	"time"
)

type SyncTimer struct {
	TimerData
	tw       *TimeWheelsync
	slot     int
	initSlot int
	rmFlag   atomic.AtomicBool
}

type TimeWheelsync struct {
	slots    []*list.List
	tickTime time.Duration
	index    int
	stop     atomic.AtomicBool
	lastTime time.Duration
}

//创建一个单层的同步时间轮
//tickTime：一个tick的时间
//duration：最长的过期时间，不能小于tickTime
func NewSyncOne(tickTime time.Duration, duration time.Duration) *TimeWheelsync {
	if tickTime > duration {
		return nil
	}

	tw := &TimeWheelsync{
		slots:    make([]*list.List, duration/tickTime),
		tickTime: tickTime,
		index:    0,
		stop:     0,
	}

	for i := 0; i < len(tw.slots); i++ {
		tw.slots[i] = list.New()
	}
	return tw
}

func (tw *TimeWheelsync) Start() {
	tw.stop = 0
}

func (tw *TimeWheelsync) Stop() {
	tw.stop.Set()
}

func (tw *TimeWheelsync) add2Slot(timer *SyncTimer) {
	var index int
	length := len(tw.slots)
	if timer.Expire < 0 {
		duration := tw.tickTime * time.Duration(length)
		index = int(duration + timer.Expire/tw.tickTime)
		index = (index + tw.index) % length
	} else {
		index = int(timer.Expire / tw.tickTime)
		if index == 0 {
			index = (tw.index + 1) % length
		} else {
			index = (index + tw.index) % length
		}
	}

	timer.initSlot = tw.index
	timer.slot = index
	l := tw.slots[index]
	if !timer.rmFlag.IsSet() {
		l.PushBack(timer)
	}
}

func (tw *TimeWheelsync) removeTimer(timer *SyncTimer) {
	l := tw.slots[timer.slot]
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == timer || e.Value.(*SyncTimer).rmFlag.IsSet() {
			l.Remove(e)
			return
		}
	}
}

func (tw *TimeWheelsync) Tick(duration time.Duration) {
	if tw.stop.IsSet() {
		return
	}
	tw.lastTime += duration
	if tw.lastTime >= tw.tickTime {
		tw.index = (tw.index + 1) % len(tw.slots)
		l := tw.slots[tw.index]
		var n *list.Element
		for e := l.Front(); e != nil; e = n {
			n = e.Next()
			timer := e.Value.(*SyncTimer)
			timer.Callback()
			l.Remove(e)
			if timer.Repeat {
				tw.add2Slot(timer)
			}
		}

		tw.lastTime = 0
	}
}

func (tw *TimeWheelsync) Add(callback OnTimeout, expire time.Duration, repeat bool) (Timer, error) {
	if expire > tw.tickTime*time.Duration(len(tw.slots)) {
		return nil, errors.New("expireTime out of range")
	}

	aTimer := &SyncTimer{
		TimerData: TimerData{callback, expire, repeat},
		rmFlag:    0,
		tw:        tw,
	}
	tw.add2Slot(aTimer)

	return aTimer, nil
}

func (tw *TimeWheelsync) RollTime() time.Duration {
	return time.Duration(tw.index) * tw.tickTime
}

func (aTimer *SyncTimer) Cancel() {
	aTimer.tw.removeTimer(aTimer)
}

func (aTimer *SyncTimer) PastTime() time.Duration {
	return time.Duration(aTimer.tw.index-aTimer.initSlot) * aTimer.tw.tickTime
}

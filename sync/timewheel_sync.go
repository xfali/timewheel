/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 9:38
 * @version V1.0
 * Description:
 */

package sync

import (
    "time"
    "container/list"
    "errors"
    "timewheel/utils"
    "timewheel"
)

type SyncTimer struct {
    timer  *timewheel.Timer
    slot   int
    rmFlag utils.AtomicBool
}

type TimeWheelsync struct {
    slots    [] *list.List
    tickTime time.Duration
    index    int
    stop     utils.AtomicBool
    lastTime time.Duration
}

func New(tickTime time.Duration, duration time.Duration) *TimeWheelsync {
    if tickTime > duration {
        return nil
    }

    tw := &TimeWheelsync{
        slots:    make([] *list.List, duration/tickTime),
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
    if timer.timer.Time < 0 {
        duration := tw.tickTime * time.Duration(length)
        index = int(duration + timer.timer.Time / tw.tickTime)
        index = (index + tw.index) % length
    } else {
        index = int(timer.timer.Time / tw.tickTime)
        if index == 0 {
            index = (tw.index + 1) % length
        } else {
            index = (index + tw.index) % length
        }
    }

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
            timer.timer.Callback(timer.timer.Data)
            l.Remove(e)
        }

        tw.lastTime = 0
    }
}

func (tw *TimeWheelsync) Add(timer *timewheel.Timer) (timewheel.CancelFunc, error) {
    if timer.Time > tw.tickTime*time.Duration(len(tw.slots)) {
        return nil, errors.New("expireTime out of range")
    }

    aTimer := &SyncTimer{
        timer: timer,
        rmFlag : 0,
    }
    tw.add2Slot(aTimer)

    return func() { tw.Cancel(aTimer) }, nil
}

func (atw *TimeWheelsync) Cancel(atimer *SyncTimer) {
    atw.removeTimer(atimer)
}

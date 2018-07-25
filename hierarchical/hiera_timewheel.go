/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 14:00
 * @version V1.0
 * Description: 
 */

package hierarchical

import (
    "time"
    "github.com/xfali/timewheel"
    "github.com/xfali/timewheel/sync"
    "github.com/xfali/goutils/atomic"
    "fmt"
)

//Hierarchical Timing Wheels
type SyncHieraTimeWheel struct {
    timeWheels [4] timewheel.TimeWheel
    tickTime   time.Duration
    stop     atomic.AtomicBool
}

func NewSyncHieraTimeWheel(tickTime time.Duration, duration time.Duration) *SyncHieraTimeWheel {
    tw := &SyncHieraTimeWheel{}
    secondTick := false
    hour := duration / time.Hour
    if hour > 0 {
        secondTick = true
        wheel := sync.New(time.Hour, hour*time.Hour)
        tw.timeWheels[0] = wheel
    }

    minute := (duration % time.Hour) / time.Minute
    if hour > 0 {
        wheel := sync.New(time.Minute, time.Hour)
        wheel.Add(func() {
            tw.timeWheels[0].Tick(time.Hour)
        }, time.Hour, true)
        tw.timeWheels[1] = wheel
    } else {
        if minute > 0 {
            secondTick = true
            wheel := sync.New(time.Minute, minute*time.Minute)
            tw.timeWheels[1] = wheel
        }
    }

    second := (duration % time.Minute) / time.Second
    if secondTick {
        wheel := sync.New(time.Second, time.Minute)
        wheel.Add(func() {
            fmt.Printf("Minute tick\n")
            tw.timeWheels[1].Tick(time.Minute)
        }, time.Minute, true)
        tw.timeWheels[2] = wheel
    } else {
        if second > 0 {
            secondTick = true
            wheel := sync.New(time.Second, second*time.Second)
            tw.timeWheels[2] = wheel
        }
    }

    millisecond := (duration % time.Second) / time.Millisecond
    if secondTick {
        wheel := sync.New(tickTime, time.Second)
        wheel.Add(func() {
            fmt.Printf("Second tick\n")
            tw.timeWheels[2].Tick(time.Second)
        }, time.Second, true)
        tw.timeWheels[3] = wheel
    } else {
        if millisecond > 0 {
            wheel := sync.New(tickTime, millisecond*time.Millisecond)
            tw.timeWheels[3] = wheel
        }
    }
    tw.tickTime = tickTime
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
    //fmt.Println(duration / time.Millisecond)
    htw.timeWheels[3].Tick(duration)
}

func (htw *SyncHieraTimeWheel) Add(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    return htw.addHour(callback, expire, repeat)
}

func (htw *SyncHieraTimeWheel)addHour(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    hour := expire / time.Hour
    if hour > 0 {
        fmt.Println("addHour")
        return htw.timeWheels[0].Add(func() {
            htw.addMinute(callback, expire, false)
        }, hour*time.Hour, repeat)
    } else {
        return htw.addMinute(callback, expire, repeat)
    }
}

func (htw *SyncHieraTimeWheel)addMinute(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    minute := expire % time.Hour / time.Minute
    if minute > 0 {
        fmt.Println("addMinute")
        return htw.timeWheels[1].Add(func() {
            htw.addSecond(callback, expire, false)
        }, minute*time.Minute, repeat)
    } else {
        return htw.addSecond(callback, expire, repeat)
    }
}

func (htw *SyncHieraTimeWheel)addSecond(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    second := expire % time.Minute / time.Second
    if second > 0 {
        fmt.Println("addSecond")
        return htw.timeWheels[2].Add(func() {
            htw.addMilliSecond(callback, expire, false)
        }, second*time.Second, repeat)
    } else {
        return htw.addMilliSecond(callback, expire, repeat)
    }
}

type undoTimer bool
func (undo *undoTimer) Cancel() {
}
func (undo *undoTimer) PastTime() (time.Duration) {
    return 0
}
var undo = new(undoTimer)

func (htw *SyncHieraTimeWheel)addMilliSecond(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    millisecond := expire % time.Second / time.Millisecond
    if millisecond > 0 {
        fmt.Println("addMilliSecond")
        return htw.timeWheels[3].Add(callback, millisecond*time.Millisecond, repeat)
    } else {
        callback()
        return undo, nil
    }
}


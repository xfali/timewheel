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
    "errors"
)

//Hierarchical Timing Wheels
type HieraTimeWheel struct {
    timeWheels [] timewheel.TimeWheel
    hieraTimes   []time.Duration
    stop     atomic.AtomicBool
}

//创建一个通用的时间轮，分层数据格式为：时间由大到小排列，如hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 20*time.Millisecond }
func NewHieraTimeWheel(duration time.Duration, hieraTimes []time.Duration) *HieraTimeWheel {
    if len(hieraTimes) < 2 {
        return nil
    }

    tw := &HieraTimeWheel{}
    deep := len(hieraTimes)
    tw.timeWheels = make([]timewheel.TimeWheel, deep)

    secondTick := false

    time := duration / hieraTimes[0]
    if time > 0 {
        secondTick = true
        wheel := sync.New(hieraTimes[0], time*hieraTimes[0])
        tw.timeWheels[0] = wheel
    }

    for j:=1; j<deep; j++ {
        i := j
        time = (duration % hieraTimes[i-1]) / hieraTimes[i]
        if secondTick {
            wheel := sync.New(hieraTimes[i], hieraTimes[i-1])
            wheel.Add(func() {
                tw.timeWheels[i-1].Tick(hieraTimes[i-1])
            }, hieraTimes[i-1], true)
            tw.timeWheels[i] = wheel
        } else {
            if time > 0 {
                secondTick = true
                wheel := sync.New(hieraTimes[i], time*hieraTimes[i])
                tw.timeWheels[i] = wheel
            }
        }
    }

    tw.hieraTimes = hieraTimes
    tw.stop = atomic.AtomicBool(1)
    return tw
}

func (htw *HieraTimeWheel) Start() {
    htw.stop = 0
}

func (htw *HieraTimeWheel) Stop() {
    htw.stop.Set()
}

func (htw *HieraTimeWheel) Tick(duration time.Duration) {
    if htw.stop.IsSet() {
        return
    }
    //fmt.Println(duration / time.Millisecond)
    htw.timeWheels[len(htw.timeWheels)-1].Tick(duration)
}

func (htw *HieraTimeWheel) Add(callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    if expire < htw.hieraTimes[len(htw.hieraTimes)-1] {
        return nil, errors.New("expire time is too small")
    }

    return htw.addTime(0, callback, expire, repeat)
}

func (htw *HieraTimeWheel)addTime(deep int, callback timewheel.OnTimeout, expire time.Duration, repeat bool) (timewheel.Timer, error) {
    var nextTime time.Duration
    if deep == 0 {
        nextTime = expire / htw.hieraTimes[deep]
    } else {
        nextTime = expire % htw.hieraTimes[deep-1] / htw.hieraTimes[deep]
    }

    if deep == len(htw.hieraTimes)-1 {
        fmt.Println("finally: ", deep)
        if nextTime > 0 {
            now := time.Now()
            return htw.timeWheels[deep].Add(func() {
                callback()
                if repeat {
                    fmt.Println("repeat in callback", expire, " time ", nextTime, " deep ", deep, " time ", time.Since(now))
                    htw.addTime(0, callback, expire, repeat)
                }
            }, nextTime*htw.hieraTimes[deep], false)
        } else {
            callback()
            if repeat {
                fmt.Println("repeat in 0")
                htw.addTime(0, callback, expire, repeat)
            }
            return undo, nil
        }
    } else {
        if nextTime > 0 {
            fmt.Println("addTime: ", deep)
            now := time.Now()
            return htw.timeWheels[deep].Add(func() {
                fmt.Println("addTime ", time.Since(now))
                htw.addTime(deep + 1, callback, expire, repeat)
            }, nextTime*htw.hieraTimes[deep], false)
        } else {
            return htw.addTime(deep + 1, callback, expire, repeat)
        }
    }
}


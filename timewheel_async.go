/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 9:38
 * @version V1.0
 * Description:
 */

package timewheel

import (
    "time"
    "container/list"
    "errors"
    "github.com/xfali/goutils/atomic"
)

type ASyncTimer struct {
    TimerData
    tw *TimeWheelAsync
    slot     int
    initSlot int
    rmFlag   atomic.AtomicBool
}

type TimeWheelAsync struct {
    slots    [] *list.List
    tickTime time.Duration
    stop chan bool
    addChan  chan *ASyncTimer
    rmChan   chan *ASyncTimer
    index    int
}

func NewAsyncOne(tickTime time.Duration, duration time.Duration, addMax int, rmMax int) *TimeWheelAsync {
    if tickTime > duration {
        return nil
    }

    tw := &TimeWheelAsync{
        slots:    make([] *list.List, duration/tickTime),
        tickTime: tickTime,
        stop:     make(chan bool),
        addChan:  make(chan *ASyncTimer, addMax),
        rmChan:   make(chan *ASyncTimer, rmMax),
        index:    0,
    }

    for i := 0; i < len(tw.slots); i++ {
        tw.slots[i] = list.New()
    }
    return tw
}

func (tw *TimeWheelAsync) Start() {
    go func() {
        now := time.Now()
        cur := now
        for {
            //FIXME: 增加timer和tick跳动必须二选一，否则增加的timer会计时不准确。
            //但是当大量同时注册timer时，有可能造成间隔了多个tick才开始回调
            select {
            case <-tw.stop:
                return
            case timer, ok := <-tw.addChan:
                if ok {
                    tw.add2Slot(timer)
                }
            default:
                passTime := time.Since(now)
                if passTime < tw.tickTime {
                    time.Sleep(tw.tickTime - passTime)
                }
                cur = time.Now()
                tw.Tick(tw.tickTime)
                now = cur
            }
            select {
            case <-tw.stop:
                return
            case rmCh, ok := <-tw.rmChan:
                if ok {
                    tw.removeTimer(rmCh)
                }
            default:
            }
        }
    }()
}

func (tw *TimeWheelAsync) Stop() {
    close(tw.stop)
}

func (tw *TimeWheelAsync) add2Slot(timer *ASyncTimer) {
    var index int
    length := len(tw.slots)
    if timer.Expire < 0 {
        duration := tw.tickTime * time.Duration(length)
        index = int(duration + timer.Expire / tw.tickTime)
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

func (tw *TimeWheelAsync) removeTimer(timer *ASyncTimer) {
    l := tw.slots[timer.slot]
    for e := l.Front(); e != nil; e = e.Next() {
        if e.Value == timer || e.Value.(*ASyncTimer).rmFlag.IsSet() {
            l.Remove(e)
            return
        }
    }
}

func (tw *TimeWheelAsync) Tick(duration time.Duration) {
    tw.index = (tw.index + 1) % len(tw.slots)
    l := tw.slots[tw.index]
    var n *list.Element
    for e := l.Front(); e != nil; e = n {
        n = e.Next()
        timer := e.Value.(*ASyncTimer)
        timer.Callback()
        l.Remove(e)
        if timer.Repeat {
            tw.add2Slot(timer)
        }
    }
}

func (tw *TimeWheelAsync) Add(callback OnTimeout, expire time.Duration, repeat bool) (Timer, error)  {
    if expire > tw.tickTime * time.Duration(len(tw.slots)) {
        return nil, errors.New("expireTime out of range")
    }
    aTimer := &ASyncTimer{
        TimerData: TimerData{callback, expire, repeat},
        tw : tw,
        rmFlag : 0,
    }
    tw.addChan <- aTimer
    return aTimer, nil
}

func (atw *TimeWheelAsync) Cancel(atimer *ASyncTimer)  {
    atimer.rmFlag.Set()
    atw.rmChan <- atimer
}

func (atw *TimeWheelAsync) RollTime() (time.Duration) {
    return time.Duration(atw.index) * atw.tickTime
}

func (aTimer *ASyncTimer) Cancel() {
    aTimer.rmFlag.Set()
    aTimer.tw.rmChan <- aTimer
}

func (aTimer *ASyncTimer) PastTime() (time.Duration) {
    //NOTICE:异步时间轮的Tick与Add在同一个select，所以需要+1
    return time.Duration(aTimer.tw.index - aTimer.initSlot + 1) * aTimer.tw.tickTime
}

//func (tw *TimeWheelAsync) Remove(timer timewheel.TimerCancel) {
//   timer.Cancel(tw)
//}

//func (atimer *ASyncTimer) Cancel(tw timewheel.TimeWheel) (error) {
//    atimer.rmFlag.Set()
//    atw ,ok := tw.(*TimeWheelAsync)
//    if !ok {
//        return errors.New("Error type. Expect type: TimeWheelAsync")
//    }
//    atw.rmChan <- atimer
//    return nil
//}
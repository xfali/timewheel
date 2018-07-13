/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 9:38
 * @version V1.0
 * Description:
 */

package async

import (
    "time"
    "container/list"
    "errors"
    "timewheel/utils"
    "timewheel"
)

const (
    AddChanSize    = 10
    RemoveChanSize = 10
)

type ASyncTimer struct {
    timer *timewheel.Timer
    slot     int
    rmFlag   utils.AtomicBool
}

type TimeWheelAsync struct {
    slots    [] *list.List
    tickTime time.Duration
    stop chan bool
    addChan  chan *ASyncTimer
    rmChan   chan *ASyncTimer
    index    int
}

func New(tickTime time.Duration, duration time.Duration) *TimeWheelAsync {
    if tickTime > duration {
        return nil
    }

    tw := &TimeWheelAsync{
        slots:    make([] *list.List, duration/tickTime),
        tickTime: tickTime,
        stop:     make(chan bool),
        addChan:  make(chan *ASyncTimer, AddChanSize),
        rmChan:   make(chan *ASyncTimer, RemoveChanSize),
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
            //增加timer和tick跳动必须二选一，否则增加的timer会计时不准确
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
                tw.Tick(passTime)
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
        timer.timer.Callback(timer.timer.Data)
        l.Remove(e)
    }
}

func (tw *TimeWheelAsync) Add(timer *timewheel.Timer) (timewheel.CancelFunc, error)  {
    if timer.Time > tw.tickTime * time.Duration(len(tw.slots)) {
        return nil, errors.New("expireTime out of range")
    }
    aTimer := &ASyncTimer{
        timer: timer,
        rmFlag : 0,
    }
    tw.addChan <- aTimer
    return func(){ tw.Cancel(aTimer) }, nil
}

func (atw *TimeWheelAsync) Cancel(atimer *ASyncTimer)  {
    atimer.rmFlag.Set()
    atw.rmChan <- atimer
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
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
    "fmt"
    "sync/atomic"
)

const (
    AddChanSize    = 10
    RemoveChanSize = 10
)

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) == 1 }
func (b *atomicBool) set() { atomic.StoreInt32((*int32)(b), 1) }

type OnTimeout func(interface{})()

type Timer struct {
    callback OnTimeout
    time     time.Duration
    data interface{}
    slot int
    rmFlag atomicBool
}

type TimeWheel struct {
    slots    [] *list.List
    tickTime time.Duration
    stop chan bool
    addChan  chan *Timer
    rmChan   chan *Timer
    index    int
}

func NewTimeWheel(tickTime time.Duration, duration time.Duration) *TimeWheel {
    tw := &TimeWheel{
        slots:    make([] *list.List, duration/tickTime),
        tickTime: tickTime,
        stop:     make(chan bool),
        addChan:  make(chan *Timer, AddChanSize),
        rmChan:   make(chan *Timer, RemoveChanSize),
        index:    0,
    }
    for i := 0; i < len(tw.slots); i++ {
        tw.slots[i] = list.New()
    }
    return tw
}

func (tw *TimeWheel) Start() {
    go func() {
        now := time.Now()
        cur := now
        for {
            select {
            case <-tw.stop:
                return
            case timer, ok := <-tw.addChan:
                if ok {
                    tw.add2Slot(timer)
                }
            case rmCh, ok := <-tw.rmChan:
                if ok {
                    tw.removeTimer(rmCh)
                }
            default:
                passTime := time.Now().Sub(now)
                if passTime < tw.tickTime {
                    time.Sleep(tw.tickTime - passTime)
                }
                cur = time.Now()
                tw.tick(cur.Sub(now))
                now = cur
            }
        }
    }()
}

func (tw *TimeWheel) Stop() {
    close(tw.stop)
}

func (tw *TimeWheel) add2Slot(timer *Timer) {
    index := int(timer.time / tw.tickTime) + tw.index
    timer.slot = index
    fmt.Printf("slot: %d\n", index)
    l := tw.slots[index]
    if !timer.rmFlag.isSet() {
        l.PushBack(timer)
    }
}

func (tw *TimeWheel) removeTimer(timer *Timer) {
    l := tw.slots[timer.slot]
    for e := l.Front(); e != nil; e = e.Next() {
        if e.Value == timer || e.Value.(*Timer).rmFlag.isSet() {
            l.Remove(e)
            return
        }
    }
}

func (tw *TimeWheel) tick(duration time.Duration) {
    tw.index = (tw.index + 1) % len(tw.slots)
    l := tw.slots[tw.index]
    var n *list.Element
    for e := l.Front(); e != nil; e = n {
        n = e.Next()
        timer := e.Value.(*Timer)
        timer.callback(timer.data)
        l.Remove(e)
    }
}

func (tw *TimeWheel) Add(callback OnTimeout, expireTime time.Duration, data interface{}) (*Timer, error)  {
    if expireTime > tw.tickTime * time.Duration(len(tw.slots)) {
        return nil, errors.New("expireTime out of range")
    }
    timer := &Timer{
        callback:callback,
        time: expireTime,
        data: data,
        rmFlag : 0,
    }
    tw.addChan <- timer
    return timer, nil
}

func (tw *TimeWheel) Remove(timer *Timer) {
    timer.rmFlag.set()
    tw.rmChan <- timer
}

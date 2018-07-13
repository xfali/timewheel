/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 14:00
 * @version V1.0
 * Description: 
 */

package timewheel

import (
    "time"
)

//Hierarchical Timing Wheels
type HieraTimeWheel struct {
    timeWheels [4] *TimeWheel
    tickTime time.Duration
    stop      chan bool
}
//
//func NewHieraTimeWheel(tickTime time.Duration, duration time.Duration) *HieraTimeWheel {
//    tw := &HieraTimeWheel{}
//    hour := duration / time.Hour
//    if hour > 0 {
//        wheel := NewTimeWheel(time.Hour, hour*time.Hour)
//        tw.timeWheels[0] = wheel
//    }
//
//    minute := (duration % time.Hour) / time.Minute
//    if hour > 0 {
//        wheel := NewTimeWheel(time.Minute, time.Hour)
//        tw.timeWheels[1] = wheel
//    } else {
//        if minute > 0 {
//            wheel := NewTimeWheel(time.Minute, minute*time.Minute)
//            tw.timeWheels[1] = wheel
//        }
//    }
//
//    second := (duration % time.Minute) / time.Second
//    if minute > 0 {
//        wheel := NewTimeWheel(time.Second, time.Minute)
//        tw.timeWheels[2] = wheel
//    } else {
//        if second > 0 {
//            wheel := NewTimeWheel(time.Second, second*time.Second)
//            tw.timeWheels[2] = wheel
//        }
//    }
//
//    millisecond := (duration % time.Second) / time.Millisecond
//    if second > 0 {
//        wheel := NewTimeWheel(tickTime, time.Millisecond)
//        tw.timeWheels[3] = wheel
//    } else {
//        if millisecond > 0 {
//            wheel := NewTimeWheel(tickTime, millisecond*time.Millisecond)
//            tw.timeWheels[3] = wheel
//        }
//    }
//    tw.tickTime = tickTime
//    return tw
//}
//
//func (htw *HieraTimeWheel) Start() {
//    go func() {
//        now := time.Now()
//        cur := now
//        for {
//            select {
//            case <-htw.stop:
//                return
//            default:
//                passTime := time.Now().Sub(now)
//                if passTime < htw.tickTime {
//                    time.Sleep(htw.tickTime - passTime)
//                }
//                cur = time.Now()
//                htw.tick(cur.Sub(now))
//                now = cur
//            }
//        }
//    }()
//}
//
//func (htw *HieraTimeWheel) Stop() {
//    close(htw.stop)
//}
//
//func (htw *HieraTimeWheel) tick(duration time.Duration) {
//
//}
//
//func (htw *HieraTimeWheel) Add(callback OnTimeout, expireTime time.Duration, data interface{}) (*Timer, error)  {
//    if expireTime > time.Hour {
//        return htw.timeWheels[0].Add(func(i interface{}) {
//            htw.timeWheels[1].Add(func(i interface{}) {
//                htw.timeWheels[2].Add(func(i interface{}) {
//                    htw.timeWheels[3].Add(callback, expireTime % time.Second, data)
//                }, expireTime / time.Second, data)
//            }, expireTime / time.Minute, data)
//        }, expireTime / time.Hour, data)
//    }
//
//    return nil, nil
//}
//
//func (htw *HieraTimeWheel) Remove(timer *Timer) {
//
//}

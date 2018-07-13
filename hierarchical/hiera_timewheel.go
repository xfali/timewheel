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
    "timewheel"
    "timewheel/sync"
)

//Hierarchical Timing Wheels
type HieraTimeWheel struct {
    timeWheels [4] timewheel.TimeWheel
    tickTime time.Duration
    stop      chan bool
}

func NewHieraTimeWheel(tickTime time.Duration, duration time.Duration) *HieraTimeWheel {
   tw := &HieraTimeWheel{}
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
       wheel.Add(timewheel.NewTimer(func(data interface{}){
            tw.timeWheels[0].Tick(time.Hour)
       }, -1, nil))
       tw.timeWheels[1] = wheel
   } else {
       if minute > 0 {
           secondTick = true
           wheel := sync.New(time.Minute, minute*time.Minute)
           tw.timeWheels[1] = wheel
       }
   }

   second := (duration % time.Minute) / time.Second
   if minute > 0 {
       wheel := sync.New(time.Second, time.Minute)
       wheel.Add(timewheel.NewTimer(func(data interface{}){
           tw.timeWheels[1].Tick(time.Minute)
       }, -1, nil))
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
       wheel.Add(timewheel.NewTimer(func(data interface{}){
           tw.timeWheels[2].Tick(time.Second)
       }, -1, nil))
       tw.timeWheels[3] = wheel
   } else {
       if millisecond > 0 {
           wheel := sync.New(tickTime, millisecond*time.Millisecond)
           tw.timeWheels[3] = wheel
       }
   }
   tw.tickTime = tickTime
   return tw
}

func (htw *HieraTimeWheel) Start() {
   go func() {
       now := time.Now()
       cur := now
       for {
           select {
           case <-htw.stop:
               return
           default:
               passTime := time.Since(now)
               if passTime < htw.tickTime {
                   time.Sleep(htw.tickTime - passTime)
               }
               cur = time.Now()
               htw.Tick(htw.tickTime)
               now = cur
           }
       }
   }()
}

func (htw *HieraTimeWheel) Stop() {
   close(htw.stop)
}

func (htw *HieraTimeWheel) Tick(duration time.Duration) {
    htw.timeWheels[3].Tick(duration)
}

func (htw *HieraTimeWheel) Add(timer *timewheel.Timer) (timewheel.CancelFunc, error)  {
   if expireTime > time.Hour {
       //return htw.timeWheels[0].Add(func(i interface{}) {
       //    htw.timeWheels[1].Add(func(i interface{}) {
       //        htw.timeWheels[2].Add(func(i interface{}) {
       //            htw.timeWheels[3].Add(callback, expireTime % time.Second, data)
       //        }, expireTime / time.Second, data)
       //    }, expireTime / time.Minute, data)
       //}, expireTime / time.Hour, data)
   }

   return nil, nil
}

func (htw *HieraTimeWheel) Remove(timer *Timer) {

}

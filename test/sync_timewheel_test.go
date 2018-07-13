/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 15:16
 * @version V1.0
 * Description: 
 */

package test

import (
    "testing"
    "time"
    "fmt"
    "timewheel"
    "timewheel/sync"
)

func TestSyncTimeWheel(t *testing.T) {
    tw := sync.New(100*time.Millisecond, time.Minute)
    tw.Start()

    f := func(data interface{}) {
        fmt.Printf("timeout %d ms\n", time.Since(data.(time.Time))/time.Millisecond, )
    }

    tw.Add(timewheel.NewTimer(f, 0*time.Second, time.Now()))
    tw.Add(timewheel.NewTimer(f, 1*time.Second, time.Now()))
    cancel, _ := tw.Add(timewheel.NewTimer(f, 2*time.Second, time.Now()))
    cancel()
    tw.Add(timewheel.NewTimer(f, 3*time.Second, time.Now()))
    tw.Add(timewheel.NewTimer(f, 4*time.Second, time.Now()))
    tw.Add(timewheel.NewTimer(f, 1*time.Hour, time.Now()))

    cur := time.Now()
    timeout := time.After(10*time.Second)
    for {
        select {
        case <- timeout:
            fmt.Println("close")
            tw.Stop()
            time.Sleep(time.Second)
            return
        default:

        }
        time.Sleep(10*time.Millisecond)
        tick := time.Now()
        tw.Tick(tick.Sub(cur))
        cur = tick
    }

}

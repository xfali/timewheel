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
    "timewheel"
    "time"
    "fmt"
)

func TestTimeWheel(t *testing.T) {
    tw := timewheel.NewTimeWheel(100*time.Millisecond, time.Minute)
    tw.Start()

    now := time.Now()
    f := func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
    }

    tw.Add(f, 1*time.Second, "test1")
    timer, _ := tw.Add(f, 2*time.Second, "test2")
    tw.Remove(timer)
    tw.Add(f, 3*time.Second, "test3")
    time.Sleep(time.Second)
    tw.Add(f, 4*time.Second, "test4")
    tw.Add(f, 1*time.Hour, "test4")

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

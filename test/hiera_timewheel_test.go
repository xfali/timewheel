/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/24 
 * @time 18:26
 * @version V1.0
 * Description: 
 */

package test

import (
    "testing"
    "fmt"
    "time"
    "github.com/xfali/timewheel/hierarchical"
)

func TestSyncHieraTimeWheel1(t *testing.T) {
    tw := hierarchical.NewSyncHieraTimeWheel(100*time.Millisecond, 2*time.Hour)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test1\n", time.Since(now)/time.Millisecond)
    }, 0*time.Second, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test2\n", time.Since(now)/time.Millisecond)
    }, 1*time.Second, false)
    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)
    timer.Cancel()
    tw.Add(func() {
        fmt.Printf("timeout %d ms test4\n", time.Since(now)/time.Millisecond)
    }, 3*time.Second, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test5\n", time.Since(now)/time.Millisecond)
    }, 4*time.Second, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test6\n", time.Since(now)/time.Millisecond)
    }, 1*time.Hour, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test7\n", time.Since(now)/time.Millisecond)
    }, -1, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test8\n", time.Since(now)/time.Millisecond)
    }, -110*time.Millisecond, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test9\n", time.Since(now)/time.Millisecond)
    }, -2*time.Second, false)

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

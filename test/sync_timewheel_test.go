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
    "github.com/xfali/timewheel"
    "github.com/xfali/timewheel/sync"
)

func TestSyncTimeWheel(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
    tw.Start()

    now := time.Now()
    f := func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
    }

    tw.Add(timewheel.NewTimer(f, 0*time.Second, "test0"))
    tw.Add(timewheel.NewTimer(f, 1*time.Second, "test1"))
    cancel, _ := tw.Add(timewheel.NewTimer(f, 2*time.Second, "test2"))
    cancel()
    tw.Add(timewheel.NewTimer(f, 3*time.Second, "test3"))
    tw.Add(timewheel.NewTimer(f, 4*time.Second, "test4"))
    tw.Add(timewheel.NewTimer(f, 1*time.Hour, "test5"))

    tw.Add(timewheel.NewTimer(f, -1, "test6"))

    tw.Add(timewheel.NewTimer(f, -110*time.Millisecond, "test7"))

    tw.Add(timewheel.NewTimer(f, -2*time.Second, "test8"))

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

func TestSyncTimeWheel2(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 5*time.Second)
    tw.Start()

    now := time.Now()
    f := func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
    }

    tw.Add(timewheel.NewTimer(func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
        tw.Add(timewheel.NewTimer(f , 3*time.Second, "test1 in test0"))
    }, 3*time.Second, "test0"))


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

func TestSyncTimeWheel3(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
    tw.Start()

    type mydata struct {
        str string
        time time.Time
    }

    f := func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(data.(mydata).time)/time.Millisecond, data.(mydata).str)
    }

    for i:=1; i<=50; i++ {
        tw.Add(timewheel.NewTimer(f, time.Duration(i*100)*time.Millisecond, mydata{fmt.Sprintf("test%d", i), time.Now()}))
    }

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

func TestSyncTimeWheel4(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
    tw.Start()

    now := time.Now()
    f := func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
    }

    cancel, _ := tw.Add(timewheel.NewTimer(f, 2*time.Second, "Should be cancel"))
    tw.Add(timewheel.NewTimer(func(data interface{}) {
        fmt.Printf("timeout %d ms %s\n", time.Since(now)/time.Millisecond, data)
        cancel()
    }, 1*time.Second, "test0"))

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
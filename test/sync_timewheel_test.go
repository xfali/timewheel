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
    "github.com/xfali/timewheel/sync"
)

func TestSyncTimeWheel(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
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

func TestSyncTimeWheel2(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 5*time.Second)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test1\n", time.Since(now)/time.Millisecond)
        tw.Add(func() {
            fmt.Printf("timeout %d ms test1 in test0\n", time.Since(now)/time.Millisecond)
        } , 3*time.Second, false)
    }, 3*time.Second, false)


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

    for i:=1; i<=50; i++ {
        data := mydata{fmt.Sprintf("test%d", i), time.Now()}
        tw.Add(func() {
            fmt.Printf("timeout %d ms %s\n", time.Since(data.time)/time.Millisecond, data.str)
        }, time.Duration(i*100)*time.Millisecond, false)
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

    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms Should be cancel\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
        timer.Cancel()
    }, 1*time.Second, false)

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

func TestSyncTimeWheel5(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
    }, 100*time.Millisecond, true)

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

func TestSyncTimeWheel6(t *testing.T) {
    tw := sync.New(100*time.Millisecond, 8*time.Second)
    tw.Start()

    now := time.Now()

    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
    }, 3*time.Second, false)

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
            fmt.Println(timer.PastTime())
        }
        time.Sleep(100*time.Millisecond)
        tick := time.Now()
        tw.Tick(tick.Sub(cur))
        cur = tick
    }
}
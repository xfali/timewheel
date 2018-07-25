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
)

func TestAsyncTimeWheel(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
    }, 0, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test1\n", time.Since(now)/time.Millisecond)
    }, 1*time.Second, false)
    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test2\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)
    timer.Cancel()
    tw.Add(func() {
        fmt.Printf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
    }, 3*time.Second, false)
    time.Sleep(time.Second)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test4 + 1s\n", time.Since(now)/time.Millisecond)
    }, 4*time.Second, false)
    tw.Add(func() {
        fmt.Printf("timeout %d ms test5 + 1s\n", time.Since(now)/time.Millisecond)
    }, 1*time.Hour, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test6 + 1s\n", time.Since(now)/time.Millisecond)
    }, -1, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test7 + 1s\n", time.Since(now)/time.Millisecond)
    }, -110*time.Millisecond, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test8 + 1s\n", time.Since(now)/time.Millisecond)
    }, -2*time.Second, false)

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

func TestAsyncTimeWheel2(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
        tw.Add(func() {
            fmt.Printf("timeout %d ms test1 in test0\n", time.Since(now)/time.Millisecond)
        } , 3*time.Second, false)
    }, 3*time.Second, false)

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

func TestAsyncTimeWheel3(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
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

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

func TestAsyncTimeWheel4(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
    tw.Start()

    now := time.Now()

    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms Should be cancel\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)

    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
        timer.Cancel()
    }, 1*time.Second, false)

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

func TestAsyncTimeWheel5(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
    }, 100*time.Millisecond, true)

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }

}

func TestASyncTimeWheel6(t *testing.T) {
    tw := timewheel.NewAsyncOne(100*time.Millisecond, 8*time.Second, 10, 10)
    tw.Start()

    now := time.Now()

    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test0\n", time.Since(now)/time.Millisecond)
    }, 3*time.Second, false)

    for i:=0;i<3;i++ {
        time.Sleep(time.Second)
        fmt.Println(timer.PastTime())
    }

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

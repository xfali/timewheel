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
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
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

    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

func TestSyncHieraTimeWheel2(t *testing.T) {
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
    tw.Start()

    now := time.Now()


    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)
    timer.Cancel()


    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

func TestSyncHieraTimeWheel3(t *testing.T) {
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
    tw.Start()

    now := time.Now()


    timer, _ := tw.Add(func() {
        fmt.Printf("timeout %d ms test1\n", time.Since(now)/time.Millisecond)
    }, 2*time.Second, false)

    tw.Add(func() {
        fmt.Printf("timeout: %d test0 cancel test1\n", time.Since(now)/time.Millisecond)
        timer.Cancel()
    }, time.Second, false)


    for {
        select {
        case <-time.After(5*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

func TestSyncHieraTimeWheel4(t *testing.T) {
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout: %d \n", time.Since(now)/time.Millisecond)
    }, time.Second + 500*time.Millisecond, true)


    for {
        select {
        case <-time.After(10*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

func TestSyncHieraTimeWheel5(t *testing.T) {
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
    tw.Start()

    now := time.Now()

    tw.Add(func() {
        fmt.Printf("timeout: %d test0 cancel test1\n", time.Since(now)/time.Millisecond)
        tw.Add(func() {
            fmt.Printf("timeout %d ms test1\n", time.Since(now)/time.Millisecond)
        }, 2*time.Second, false)
    }, time.Second, false)


    for {
        select {
        case <-time.After(5*time.Second):
            tw.Stop()
            time.Sleep(time.Second)
            return
        }
    }
}

func TestAsyncTimeWheel6(t *testing.T) {
    hieraTimes := []time.Duration{ time.Hour, time.Minute, time.Second, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
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

func TestAsyncTimeWheel7(t *testing.T) {
    hieraTimes := []time.Duration{ time.Minute, 100*time.Millisecond }
    tw := hierarchical.NewHieraTimeWheel(2*time.Hour, hieraTimes, 10, 10)
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
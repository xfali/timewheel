// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"fmt"
	"github.com/xfali/timewheel"
	"testing"
	"time"
)

func TestNewErr(t *testing.T) {
	timewheel.New(timewheel.AsyncOptSetDuration(time.Second, 3*time.Millisecond))
}

func TestNew(t *testing.T) {
	tw := timewheel.New(timewheel.AsyncOptSetMaxDuration(time.Second))
	test0(tw, t)
	tw = timewheel.New(timewheel.AsyncOptSetDuration(time.Second, time.Second))
	test0(tw, t)
	tw = timewheel.New(timewheel.AsyncOptSetDuration(time.Minute, time.Second))
	test0(tw, t)
	tw = timewheel.New(timewheel.AsyncOptSetMinDuration(25 * time.Millisecond))
	test0(tw, t)

	t.Run("repeat 1 Hiera", func(t *testing.T) {
		tw = timewheel.New(timewheel.AsyncOptSetDuration(time.Second, time.Second))
		test0_0(tw, t)
	})
}

func TestNew2(t *testing.T) {
	t.Run("repeat 2 Hiera", func(t *testing.T) {
		tw := timewheel.New(timewheel.AsyncOptSetDuration(10*time.Second, timewheel.DefaultMinDuration))
		test0_0(tw, t)
	})
}

func TestAsyncTimer1(t *testing.T) {
	tw := timewheel.New()
	test1(tw, t)
}

func TestAsyncTimer2(t *testing.T) {
	tw := timewheel.New()
	test2(tw, t)
}

func TestAsyncTimer3(t *testing.T) {
	tw := timewheel.New()
	test3(tw, t)
}

func TestAsyncTimer4(t *testing.T) {
	tw := timewheel.New()
	test4(tw, t)
}

func TestAsyncTimer5(t *testing.T) {
	tw := timewheel.New()
	test5(tw, t)
}

func TestAsyncTimer6(t *testing.T) {
	tw := timewheel.New()
	test6(tw, t)
}

func TestAsyncTimer7(t *testing.T) {
	tw := timewheel.New(timewheel.AsyncOptSetDuration(time.Second, timewheel.DefaultMinDuration))
	test7(tw, t)
}

func test0(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	now := time.Now()

	_, err := tw.Add(func() {
		t.Logf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
	}, time.Second, false)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-time.After(2 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test0_0(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	now := time.Now()

	_, err := tw.Add(func() {
		t.Logf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
		t.Log(tw.RollTime())
	}, 1000*time.Millisecond, true)
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test1(tw timewheel.TimeWheel, t *testing.T) {
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
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test2(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	now := time.Now()

	timer, _ := tw.Add(func() {
		fmt.Printf("timeout %d ms test3\n", time.Since(now)/time.Millisecond)
	}, 500*time.Millisecond, false)
	timer.Cancel()

	for {
		select {
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test3(tw timewheel.TimeWheel, t *testing.T) {
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
		case <-time.After(5 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test4(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	now := time.Now()

	tw.Add(func() {
		fmt.Printf("timeout: %d \n", time.Since(now)/time.Millisecond)
	}, time.Second+500*time.Millisecond, true)

	for {
		select {
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test5(tw timewheel.TimeWheel, t *testing.T) {
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
		case <-time.After(5 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

func test6(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	type mydata struct {
		str  string
		time time.Time
	}

	for i := 1; i <= 50; i++ {
		data := mydata{fmt.Sprintf("test%d", i), time.Now()}
		tw.Add(func() {
			fmt.Printf("timeout %d ms %s\n", time.Since(data.time)/time.Millisecond, data.str)
		}, time.Duration(i*100)*time.Millisecond, false)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}

}

func test7(tw timewheel.TimeWheel, t *testing.T) {
	tw.Start()

	now := time.Now()
	tw.Add(func() {
		fmt.Println("timeout no repeat: ", time.Since(now))
	}, time.Second, false)

	tw.Add(func() {
		fmt.Println("timeout repeat: ", time.Since(now))
	}, time.Second, true)

	for {
		select {
		case <-time.After(10 * time.Second):
			tw.Stop()
			time.Sleep(time.Second)
			return
		}
	}
}

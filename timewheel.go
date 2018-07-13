/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/13 
 * @time 8:54
 * @version V1.0
 * Description: 
 */

package timewheel

import (
    "time"
)

type OnTimeout func(interface{}) ()

type Timer struct {
    Callback OnTimeout
    Time     time.Duration
    Data     interface{}
}

type CancelFunc func()

type TimeWheel interface {
    Start()

    Stop()

    Tick(duration time.Duration)

    Add(*Timer) (CancelFunc, error)
}

func NewTimer(Callback OnTimeout, Time time.Duration, Data interface{}) (*Timer) {
    return &Timer{
        Callback: Callback,
        Time:     Time,
        Data:     Data,
    }
}

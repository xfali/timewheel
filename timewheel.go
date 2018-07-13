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
    //到期时回调
    Callback OnTimeout
    //过期时间，如果为0，则在下一个tick回调
    Time     time.Duration
    //回调时回传的数据
    Data     interface{}
}

type CancelFunc func()

type TimeWheel interface {
    //开启时间轮
    Start()

    //关闭时间轮
    Stop()

    //每一个时钟跳动的时间片操作
    Tick(duration time.Duration)

    //*Timer:增加一个计时器
    //CancelFunc 取消方法
    //正常为nil，其他返回具体错误
    Add(*Timer) (CancelFunc, error)
}

func NewTimer(Callback OnTimeout, Time time.Duration, Data interface{}) (*Timer) {
    return &Timer{
        Callback: Callback,
        Time:     Time,
        Data:     Data,
    }
}

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

type OnTimeout func() ()

type CancelFunc func()

type Timer struct {
    Callback OnTimeout
    Expire time.Duration
    Repeat bool
}

type TimeWheel interface {
    //开启时间轮
    Start()

    //关闭时间轮
    Stop()

    //每一个时钟跳动的时间片操作
    Tick(time.Duration)

    //参数：callback: 超时回调
    //参数：expire: 超时时间
    //参数：repeat: 是否重复
    //返回：CancelFunc: 取消方法
    //返回：err: 正常为nil，其他返回具体错误
    Add(callback OnTimeout, expire time.Duration, repeat bool) (CancelFunc, error)
}


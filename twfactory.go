/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/25 
 * @time 13:50
 * @version V1.0
 * Description: 
 */

package timewheel

import (
    "time"
    "github.com/xfali/timewheel/async"
    "github.com/xfali/timewheel/sync"
    "github.com/xfali/timewheel/hierarchical"
)

func NewAsyncOne(tickTime time.Duration, duration time.Duration, addMax int, rmMax int) (TimeWheel) {
    return async.New(tickTime, duration, addMax, rmMax)
}

func NewSyncOne(tickTime time.Duration, duration time.Duration) (TimeWheel) {
    return sync.New(tickTime, duration)
}

func NewAsyncHiera(duration time.Duration, hieraTimes []time.Duration, addMax int, rmMax int) (TimeWheel) {
    return hierarchical.NewHieraTimeWheel(duration , hieraTimes, addMax, rmMax)
}

func NewSyncHiera(duration time.Duration, hieraTimes []time.Duration) TimeWheel {
    return hierarchical.NewSyncHieraTimeWheel(duration, hieraTimes)
}

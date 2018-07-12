/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/12 
 * @time 14:00
 * @version V1.0
 * Description: 
 */

package timewheel

//
////Hierarchical Timing Wheels
//type HieraTimeWheel struct {
//    TimeWheel [4] *TimeWheel
//}
//
//func NewHieraTimeWheel(tickTime time.Duration, duration time.Duration) {
//    tw := &HieraTimeWheel{}
//    hour := duration / time.Hour
//    if hour > 0 {
//        wheel := TimeWheel{make([] list.List, hour)}
//        tw.TimeWheel[0] = &wheel
//    }
//
//    minute := (duration % time.Hour) / time.Minute
//    if hour > 0 {
//        wheel := TimeWheel{make([] list.List, 60)}
//        tw.TimeWheel[1] = &wheel
//    } else {
//        if minute > 0 {
//            wheel := TimeWheel{make([] list.List, minute)}
//            tw.TimeWheel[1] = &wheel
//        }
//    }
//
//    second := (duration % time.Minute) / time.Second
//    if minute > 0 {
//        wheel := TimeWheel{make([] list.List, 60)}
//        tw.TimeWheel[2] = &wheel
//    } else {
//        if second > 0 {
//            wheel := TimeWheel{make([] list.List, second)}
//            tw.TimeWheel[2] = &wheel
//        }
//    }
//
//    millisecond := (duration % time.Second) / time.Millisecond
//    if second > 0 {
//        wheel := TimeWheel{make([] list.List, 1000)}
//        tw.TimeWheel[3] = &wheel
//    } else {
//        if millisecond > 0 {
//            wheel := TimeWheel{make([] list.List, millisecond)}
//            tw.TimeWheel[3] = &wheel
//        }
//    }
//}

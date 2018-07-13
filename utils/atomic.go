/**
 * Copyright (C) 2018, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @date 2018/7/13 
 * @time 8:57
 * @version V1.0
 * Description: 
 */

package utils

import "sync/atomic"

type AtomicBool int32

func (b *AtomicBool) IsSet() bool { return atomic.LoadInt32((*int32)(b)) == 1 }
func (b *AtomicBool) Set() { atomic.StoreInt32((*int32)(b), 1) }

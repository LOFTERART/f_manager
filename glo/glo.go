/**
 * @Author: LOFTER
 * @Description:
 * @File:  glo
 * @Date: 2021/2/5 4:05 下午
 */
package glo

import (
	"upload/cache"
	"upload/chunk"
)

var (
	IPAddress string
	FD        chunk.FileDealer
	ROP       cache.RedisOp
)

/**
 @author: yuefei
 @date: 2021/3/2
 @note:
**/

package util

import (
	"fmt"
	"testing"
)

func TestGetStandardTimeStampStr(t *testing.T) {
	stampStr := TransIntToCstTime(1614332847383)
	fmt.Printf(stampStr)
}

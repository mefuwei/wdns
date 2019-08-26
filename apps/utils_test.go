/*
 Auther : F.W
 Create time  2018/12/30
*/
package apps

import (
	"fmt"
	"testing"
)

func TestVerifyIP(t *testing.T) {
	ok := VerifyIP("fe80::bcb4:e5ff:fef1:261f")
	fmt.Println(ok)
}

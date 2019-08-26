/*
 Auther : F.W
 Create time  2018/12/30
*/
package apps

import "net"

func VerifyIP(s string) bool {

	if net.ParseIP(s).To16() == nil {
		return false
	}
	return true

}

// 域名解析记录冲突验证
func VerifyDomainConflict(s *JsonSerializer, list []*JsonSerializer) bool {

	parseType := s.ParseType

	for _, v := range list {

		if *s == *v {
			return false
		}

	}

	if parseType == "A" || parseType == "AAAA" {
		for _, v := range list {
			if v.ParseType == "CNAME" {
				return false
			}
		}
		return true

	} else if parseType == "CNAME" {
		for _, v := range list {
			if v.ParseType == "CNAME" {
				return false
			}
		}
		return true
	} else {
		return true
	}

}

/*
 Auther : F.W
 Create time  2018/12/30
*/
package apps

import (
	"net"
)

func VerifyIP(s string) bool {

	if net.ParseIP(s).To16() == nil {
		return false
	}
	return true

}

// 域名解析记录冲突验证
//func VerifyRecordRules(s *JsonSerializer, list []*JsonSerializer) (e error) {
//
//	for _, v := range list {
//
//		if *s == *v {
//			return errors.New("记录已经存在")
//		}
//
//	}
//
//	for _, one := range list {
//		if s.Area == one.Area {
//			switch s.Rtype {
//			case "A":
//				if one.Rtype == "CNAME" {
//					return errors.New("A记录与CNAME记录冲突")
//				}
//			case "AAAA":
//
//				if one.Rtype == "CNAME" {
//					return errors.New("AAAA记录与CNAME记录冲突")
//				}
//			case "CNAME":
//
//				if one.Rtype == "A" {
//					return errors.New("CNAME记录与A记录冲突")
//				}
//				if one.Rtype == "AAAA" {
//					return errors.New("CNAME记录与AAAA记录冲突")
//				}
//			}
//
//		}
//
//	}
//	return nil
//}

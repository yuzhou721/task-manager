package utils

import (
	"crypto/sha1"
	"fmt"
	"log"
	"sort"
	"strings"
)

//Sha 取得sha1加密字符串
func Sha(data []string) (r string) {
	sort.Strings(data)
	t := strings.Join(data, "")
	b := sha1.Sum([]byte(t))
	log.Printf("pubtoken:%x", string(b[:]))
	return fmt.Sprintf("%x", string(b[:]))
}

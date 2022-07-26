package payment

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

func (p *Payment) Sign(paramsMap map[string]interface{}) string {
	var paramsArr []string
	for k, v := range paramsMap {
		if k == "other_settle_params" {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		switch k {
		case "app_id", "thirdparty_id", "sign":
		default:
			paramsArr = append(paramsArr, value)
		}
	}

	paramsArr = append(paramsArr, p.Salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))

}

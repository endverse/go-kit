package alert

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	ahash "github.com/endverse/go-kit/hash"
)

func generateHMAC(params map[string]string, sk string) string {
	var (
		result string
		tmp    []string
	)

	for k := range params {
		tmp = append(tmp, k)
	}
	sort.Strings(tmp)

	for _, v := range tmp {
		result = fmt.Sprintf("%s%s=%s&", result, v, params[v])
	}

	result = fmt.Sprintf("%ssk=%s", result, sk)

	return strings.ToUpper(ahash.MD5(result))
}

func generateUserIDs(users []BossHiUser) string {
	ids := make([]string, 0)
	for _, user := range users {
		ids = append(ids, strconv.FormatInt(user.UserId, 10))
	}

	return strings.Join(ids, ",")
}

func generatePath(host, path string) string {
	return host + path
}

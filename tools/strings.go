package tools

import "encoding/json"

func StringCompress(strs []string) (str string, err error) {
	bs, err := json.Marshal(strs)
	return string(bs), err
}

func StringDecompress(str string) (strs []string, err error) {
	if str == "" {
		return make([]string, 0), nil
	}

	bs := []byte(str)
	err = json.Unmarshal(bs, &strs)
	return strs, err
}

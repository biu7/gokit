package json

import "github.com/bytedance/sonic"

func Unmarshal[T any](data []byte) (T, error) {
	var ret T
	err := sonic.Unmarshal(data, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func Marshal[T any](ret *T) ([]byte, error) {
	return sonic.Marshal(ret)
}

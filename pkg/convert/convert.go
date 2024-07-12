package convert

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

// func (s StrTo) UInt32() (uint32, error) {
// 	v, err := strconv.Atoi(s.String())
// 	return uint32(v), err
// }

// chatgpt给出的UInt32优化后的代码
func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.ParseUint(string(s), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}
func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

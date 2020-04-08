package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindEncoding(t *testing.T) {
	require.Equal(t, GSM7BIT, FindEncoding("abc30hb3bk2lopzSD=2-^"))
	require.Equal(t, UCS2, FindEncoding("Trần Lập và ban nhạc Bức tường huyền thoại"))
	require.Equal(t, UCS2, FindEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ"))
}

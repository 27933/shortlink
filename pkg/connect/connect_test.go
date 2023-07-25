package connect

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	convey.Convey("基础用例", t, func() {
		url := "https://www.liwenzhou.com/posts/Go/unit-test-5/"
		got := Get(url)
		// 断言 -> 判断是否满足条件
		convey.So(got, convey.ShouldEqual, true) // 断言
		// convey.ShouldBeTrue(got)
	})
	convey.Convey("url错误的示例", t, func() {
		url := "https://www.27933.com/"
		got := Get(url)
		// 断言 -> 判断是否满足条件
		convey.So(got, convey.ShouldEqual, false) // 断言
		// convey.ShouldBeTrue(got)
	})
}

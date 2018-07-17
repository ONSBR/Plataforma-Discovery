package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldGetDependecyDomain(t *testing.T) {
	Convey("should get all processes that depends on some entity", t, func() {
		instances, err := GetProcessesWithDependsOn("ec498841-59e5-47fd-8075-136d79155705", []string{"conta", "operacao"})
		So(err, ShouldBeNil)
		So(len(instances), ShouldBeGreaterThan, 0)
	})
}

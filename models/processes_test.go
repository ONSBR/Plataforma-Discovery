package models

import (
	"testing"

	"github.com/ONSBR/Plataforma-Discovery/helpers"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldMountDataSourceChain(t *testing.T) {
	Convey("should mount datasource graph from dependency domain", t, func() {
		instances, err := helpers.GetProcessesWithDependsOn("ec498841-59e5-47fd-8075-136d79155705", []string{"conta", "operacao"})
		So(err, ShouldBeNil)
		chain := NewDataSourceChain(instances)
		So(len(chain), ShouldBeGreaterThan, 0)
	})
}

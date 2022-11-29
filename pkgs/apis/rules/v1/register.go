package v1

import (
	"gitee.com/cpds/cpds-analyzer/pkgs/rules"
	"github.com/emicklei/go-restful"
)

func AddToContainer(container *restful.Container, r *rules.Rules) {
	webservice := new(restful.WebService)
	webservice.Path("/rules")

	handler := newRulesHandler()

	webservice.Route(webservice.GET("/").
		To(handler.GetRules))

	container.Add(webservice)
}

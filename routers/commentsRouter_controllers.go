package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["mepmgr/controllers:ConfigController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:ConfigController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           "/",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:ConfigController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:ConfigController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           "/mepGroup",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/mepGroup",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/mepGroup/:mepGroupName",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           "/mepGroup/:mepGroupName",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "AddMep",
			Router:           "/mepGroup/:mepGroupName/mep",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepGroupController"],
		beego.ControllerComments{
			Method:           "DeleteMep",
			Router:           "/mepGroup/:mepGroupName/mep/:mepName",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepProcessLogController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepProcessLogController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "/process/log",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "grouptopology/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"],
		beego.ControllerComments{
			Method:           "GetDetail",
			Router:           "grouptopology/detail",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"],
		beego.ControllerComments{
			Method:           "GetByLocation",
			Router:           "locationtopology/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:MepTopologyController"],
		beego.ControllerComments{
			Method:           "GetDetailByLocation",
			Router:           "locationtopology/detail",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "Create",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           "/:mepId",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           "/:mepId",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "UpdateStatus",
			Router:           "/:mepId/status",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:NodeController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:NodeController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "/detail/:mepId",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:PerferenceController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:PerferenceController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/perference",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:PerferenceController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:PerferenceController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           "/perference",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["mepmgr/controllers:CertController"] = append(beego.GlobalControllerRouter["mepmgr/controllers:CertController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           "/:certId",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}

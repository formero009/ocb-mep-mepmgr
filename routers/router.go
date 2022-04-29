// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"github.com/astaxie/beego"
	"mepmgr/controllers"
	"mepmgr/manager"
)

func init() {
	//route config
	ns := beego.NewNamespace("/mepm/v1/mepmgr",
		beego.NSNamespace("/meps",
			beego.NSInclude(
				&controllers.NodeController{},
			),
			beego.NSInclude(
				&controllers.MepGroupController{MepGroupMg: manager.NewDefaultMepGroupManager()},
			),
			beego.NSInclude(
				&controllers.MepTopologyController{MepGTopologyMg: manager.NewDefaultMepTopologyManager()},
			),
			beego.NSInclude(
				&controllers.MepProcessLogController{MepProcessLogManager: manager.NewDefaultMepProcessLogManager()},
			),
			beego.NSInclude(
				&controllers.PerferenceController{PerferenceMg: manager.NewDefaultPreferenceManager()},
			),
		),
		beego.NSNamespace("/config", beego.NSInclude(
			&controllers.ConfigController{},
		),
		),
		//v2.4.9 cert manager
		beego.NSNamespace("/cert", beego.NSInclude(
			&controllers.CertController{CertMg: manager.NewDefaultCertManager()},
		),
		),
	)
	beego.AddNamespace(ns)
	beego.SetStaticPath("/swagger", "swagger")
}

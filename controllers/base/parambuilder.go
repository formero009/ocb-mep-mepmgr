package base

import (
	"fmt"
	"mepmgr/common"
	"mepmgr/util/cameltrans"
	"strconv"
)

//const value for pageSearch
const (
	defaultPageNo   = 1
	defaultPageSize = 20
)

type ParamBuilderController struct {
	ResultHandlerController
}

func (c *ParamBuilderController) BuildQueryParam() *common.QueryParam {
	no, size := c.buildPageParam()
	return &common.QueryParam{
		PageNo:   no,
		PageSize: size,
		Order:    cameltrans.Camel2Snake(c.Input().Get("order")),
	}
}

func (c *ParamBuilderController) buildPageParam() (no int64, size int64) {
	pageNo := c.Input().Get("page_no")
	pageSize := c.Input().Get("page_size")
	if pageNo == "" {
		pageNo = strconv.Itoa(defaultPageNo)
	}

	if pageSize == "" {
		pageSize = strconv.Itoa(defaultPageSize)
	}

	no, err := strconv.ParseInt(pageNo, 10, 64)
	// pageNo must bigger than zero.
	if err != nil || no < 1 {
		c.AbortBadRequest("Invalid pageNo in query.")
	}
	// pageSize must bigger than zero.
	size, err = strconv.ParseInt(pageSize, 10, 64)
	if err != nil || size < 1 {
		c.AbortBadRequest("Invalid pageSize in query.")
	}
	return
}

func (c *ParamBuilderController) GetIDFromURL() int64 {
	return c.GetIntParamFromURL(":id")
}

func (c *ParamBuilderController) GetStrIDFromURL() string {
	paramStr := c.Ctx.Input.Param(":id")
	if len(paramStr) == 0 {
		c.AbortBadRequest(fmt.Sprintf("Invalid %s in URL", ":id"))
	}
	return paramStr
}
func (c *ParamBuilderController) GetStrUUIDFromURL(uuid string) string {
	paramStr := c.Ctx.Input.Param(":" + uuid)
	if len(paramStr) == 0 {
		c.AbortBadRequest(fmt.Sprintf("Invalid %s in URL", ":"+uuid))
	}
	return paramStr
}
func (c *ParamBuilderController) GetIntParamFromURL(param string) int64 {
	paramStr := c.Ctx.Input.Param(param)
	if len(paramStr) == 0 {
		c.AbortBadRequest(fmt.Sprintf("Invalid %s in URL", param))
	}

	paramInt, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil || paramInt < 0 {
		c.AbortBadRequest(fmt.Sprintf("Invalid %s in URL", param))
	}

	return paramInt
}

func (c *ParamBuilderController) GetStrParamFromURL(param string) string {
	paramStr := c.Ctx.Input.Param(param)
	if len(paramStr) == 0 {
		c.AbortBadRequest(fmt.Sprintf("Invalid %s in URL", param))
	}
	return paramStr
}

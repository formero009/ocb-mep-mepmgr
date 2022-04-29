/**
 @author: yuefei
 @date: 2021/2/4
 @note:
**/

package usermd

type UserInfo struct {
	UserId         string   `json:"userId"`
	Username       string   `json:"username"`
	Name           string   `json:"name"`
	Phone          string   `json:"phone"`
	Email          string   `json:"email"`
	Brand          string   `json:"brand"`
	Roles          []string `json:"roles"`
	Authorizations []string `json:"authorizations"`
}

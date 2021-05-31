package models

type PuUser struct {
	Username      string `json:"pu_username"`
	Password      string `json:"pu_password"`
	Cookies       string `json:"pu_cookies"`
	GenPermission string `json:"gen_permission"`
	RemindJoin    string `json:"remind_join"`
	RemindSignIn  string `json:"remind_signin"`
	RemindSignOut string `json:"remind_signout"`
	VipPermission string `json:"vip_permission"`
	AutoJoin      string `json:"auto_join"`
	AutoSignIn    string `json:"auto_signin"`
	AutoSignOut   string `json:"auto_signout"`
}

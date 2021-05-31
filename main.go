package main

import (
	"./http"
	"./models"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

func main() {
	for {
		go getPuUsersTask()
		time.Sleep(time.Second * 60 * 5)
	}
}

func getPuUsersTask() {
	allUsers := gjson.Get(http.Get("https://xmwluo.com/wx/AppCenter/PuAssistant/UserApi/GetBindPuUsers.php"), "bind_pu_users")
	allUsers.ForEach(func(key, value gjson.Result) bool {
		var user models.PuUser
		err := json.Unmarshal([]byte(value.String()), &user)
		if err == nil {
			//fmt.Println(user.Username)
			go puUserTask(user)
			time.Sleep(time.Second)
		} else {
			fmt.Println("解析用户json失败：" + err.Error())
		}
		return true // keep iterating
	})
}

func puUserTask(user models.PuUser) {
	for page := 1; page <= 5; page++ {
		eventListData := http.Post("http://centos.api.xmwluo.com/PuApi/GetEventList.php", "cookies="+user.Cookies+"&page="+strconv.Itoa(page))
		//if page == 3 {
		//	fmt.Println(eventListData)
		//}
		if gjson.Get(eventListData, "code").Int() == 1 {
			eventList := gjson.Get(eventListData, "data")
			eventList.ForEach(func(key, event gjson.Result) bool {
				fmt.Println(event.Get("title").String())
				go exEventTask(user, event)
				time.Sleep(time.Second)
				return true
			})
		} else {
			//获取不到活动列表数据说明cookies无效，则启动一个线程登录账号，登录后端会自动将登录成功后的cookies保存到数据库
			go login(user.Username, user.Password)
			break
		}
		time.Sleep(time.Second)
	}
}

func exEventTask(user models.PuUser, event gjson.Result) {
	eventId := event.Get("id").String()
	eventTitle := event.Get("title").String()
	eventDetailData := http.Post("http://centos.api.xmwluo.com/PuApi/GetEventDetail.php", "cookies="+user.Cookies+"&eid="+eventId)
	if gjson.Get(eventDetailData, "result").String() == "success" {
		switch gjson.Get(eventDetailData, "joinStatus").String() {
		case "dot1":
			//活动发布
		case "dot2":
			//活动开始报名
			if user.VipPermission == "1" && user.AutoJoin == "1" {
				//go act(user,eventId,eventTitle,"join")
			}
		case "dot3":
			//已报名该活动
			//fmt.Println(user.Username,eventTitle,"已报名该活动")
			switch gjson.Get(eventDetailData, "signInStatus").String() {
			case "dot1":
				//不可签到
			case "dot2":
				//可签到
				if user.VipPermission == "1" && user.AutoSignIn == "1" {
					go act(user, eventId, eventTitle, "signin")
				}
			case "dot3":
				//已签到 检查签退状态
				//fmt.Println(user.Username,eventTitle,"已签到，"+gjson.Get(eventDetailData,"isNeedSignOut").String())
				switch gjson.Get(eventDetailData, "isNeedSignOut").String() {
				case "需要签退":
					//fmt.Println(user.Username,eventTitle,"签退状态:"+gjson.Get(eventDetailData,"signOutStatus").String())
					switch gjson.Get(eventDetailData, "signOutStatus").String() {
					case "dot1":
						//不可签退
					case "dot2":
						//可签退
						//fmt.Println(user.Username,eventTitle,"可签退")
						if user.VipPermission == "1" && user.AutoSignOut == "1" {
							go act(user, eventId, eventTitle, "signout")
						}
					case "dot3":
						//已签退
					}
				case "无需签退":
					//无需签退
					go addEventLog(user.Username, eventTitle, "无需签退")
				}
			}
		}
	}
}

func addEventLog(username string, eventName string, eventInfo string) {
	http.Post("https://xmwluo.com/wx/AppCenter/PuAssistant/UserApi/AddPuEventLog.php", "pu_username="+username+"&event_name="+eventName+"&event_info="+eventInfo)
}

func login(username string, password string) string {
	return http.Post("https://xmwluo.com/wx/AppCenter/PuAssistant/UserApi/UpdatePuUserCookies.php", "username="+username+"&password="+password)
}

func act(user models.PuUser, eid string, eventTitle, act string) {
	actResult := http.Post("http://centos.api.xmwluo.com/PuApi/EventAct.php", "cookies="+user.Cookies+"&eid="+eid+"&act="+act)

	if act == "join" {
		if gjson.Get(actResult, "info").String() == "报名成功" {
			fmt.Println(user.Username, eventTitle, act, gjson.Get(actResult, "info").String())
			go addEventLog(user.Username, eventTitle, gjson.Get(actResult, "info").String())
		}
	} else {
		fmt.Println(user.Username, eventTitle, act, gjson.Get(actResult, "info").String())
		go addEventLog(user.Username, eventTitle, gjson.Get(actResult, "info").String())
	}
}

package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"strings"
	"time"
)

var ADMIN_LOGINS = []string{"admin"}

type User struct {
	BaseModel
	Login       string `sql:"size:255;not null"`
	Password    string `sql:"size:255;not null"`
	Email       string `sql:"size:255"`
	Avatar      string `sql:"size:255"`
	GitHub      string
	Twitter     string
	HomePage    string
	Tagline     string
	Description string
	Location    string
	Topics      []Topic
	Replies     []Reply
}

func (u User) BeforeCreate() error {
	u.Login = strings.ToLower(u.Login)
	return nil
}

func (u User) GavatarURL(size string) string {
	emailMD5 := u.EncodePassword(u.Email)
	return fmt.Sprintf("https://gravatar.com/avatar/%v?s=%v", emailMD5, size)
}

func (u User) NotifyChannelId() string {
	return fmt.Sprintf("notify/%v", u.Id)
}

func (u User) SameAs(obj interface{}) bool {
	return obj.(User).Id == u.Id
}

func (u User) IsAdmin() bool {
	for _, str := range ADMIN_LOGINS {
		if u.Login == str {
			return true
		}
	}
	return false
}

func (u User) UnReadNotificationsCount() (count int) {
	db.Model(&Notification{}).Where("`user_id` = ? and `read` = 0", u.Id).Count(&count)
	return
}

func PushNotifyInfoToUser(userId int32, note Notification, isNew bool) {
	u := User{}
	u.Id = userId
	unreadCount := u.UnReadNotificationsCount()

	// Set read, update client unread_count
	if !isNew {
		go PushMessage(u.NotifyChannelId(), &NotifyInfo{UnreadCount: unreadCount, IsNew: false})
		return
	}

	actor := User{}
	if note.Id > 0 {
		db.First(&actor, note.ActorId)
	}
	info := NotifyInfo{
		UnreadCount: unreadCount,
		IsNew:       true,
		Title:       note.NotifyableTitle(),
		Avatar:      actor.GavatarURL("256x256"),
		Path:        note.NotifyableURL(),
	}

	revel.INFO.Println("[Push] Notify:", info)

	go PushMessage(u.NotifyChannelId(), &info)
}

func (u User) EncodePassword(raw string) (md5Digest string) {
	data := []byte(raw)
	result := md5.Sum(data)
	md5Digest = hex.EncodeToString(result[:])
	return
}

func (u User) Signup(login string, email string, password string, passwordConfirm string) (user User, v revel.Validation) {
	u.Login = strings.ToLower(strings.Trim(login, " "))
	u.Email = strings.ToLower(strings.Trim(email, " "))

	v.Required(email).Key("Email").Message("不能为空")
	v.MinSize(login, 5).Key("用户名").Message("最少要 5 个字符")
	v.MinSize(password, 6).Key("密码").Message("最少要 6 个字符")
	v.Email(email).Key("Email").Message("格式不正确")

	if password != passwordConfirm {
		v.Error("密码与确认密码不一致")
	}

	var existCount int
	db.Model(&User{}).Where("login = ?", login).Count(&existCount)
	fmt.Println("login same as: ", login, " have ", existCount)
	if existCount > 0 {
		v.Error("帐号已经被注册")
	}

	if v.HasErrors() {
		return u, v
	}

	u.Password = u.EncodePassword(password)

	err := db.Create(&u).Error
	if err != nil {
		v.Error(fmt.Sprintf("服务器异常, %v", err))
	}
	fmt.Println("created user: ", u)
	return u, v
}

func (u User) Signin(login string, password string) (user User, v revel.Validation) {
	login = strings.Trim(login, " ")

	if len(password) == 0 {
		v.Error("还未输入密码")
	}

	db.First(&user, "login = ? and password = ?", strings.ToLower(login), u.EncodePassword(password))
	fmt.Println("first user:", user)
	if user.Id == 0 {
		v.Error("帐号密码不正确")
	}
	return user, v
}

func UpdateUserProfile(u User) (user User, v revel.Validation) {
	v.Email(u.Email).Key("Email").Message("格式不正确")
	if v.HasErrors() {
		return u, v
	}
	willUpdateUser := User{
		Email:       u.Email,
		Location:    u.Location,
		Description: u.Description,
		GitHub:      u.GitHub,
		Twitter:     u.Twitter,
		Tagline:     u.Tagline,
	}
	err := db.First(&u, u.Id).Updates(willUpdateUser).Error
	if err != nil {
		v.Error(err.Error())
	}
	return u, v
}

func (u User) UpdatePassword(oldPassword, newPassword, confirmPassword string) (v revel.Validation) {
	user := User{}
	v.Required(oldPassword).Key("旧密码").Message("不能为空")
	db.First(&user, "id = ? and password = ?", u.Id, u.EncodePassword(oldPassword))
	if user.NewRecord() {
		v.Error("旧密码不正确")
	}
	v.MinSize(newPassword, 6).Key("新密码").Message("最少要 6 个子符")
	if newPassword != confirmPassword {
		v.Error("新密码与确认新密码输入的内容不一致")
	}
	if v.HasErrors() {
		return v
	}

	err := db.Model(u).Update("password", u.EncodePassword(newPassword)).Error
	if err != nil {
		v.Error(err.Error())
	}
	return v
}

func FindUserByLogin(login string) (u User, err error) {
	err = db.Where("login = ?", strings.ToLower(login)).First(&u).Error
	return
}

func UsersCountCached() (count int) {
	if err := cache.Get("users/total", &count); err != nil {
		if err = db.Model(User{}).Count(&count).Error; err == nil {
			go cache.Set("users/total", count, 30*time.Minute)
		}
	}

	return
}

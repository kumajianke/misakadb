package miusers

import (
	"encoding/json"
	"fmt"
	"misakadb/clilog"
	"misakadb/safe"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserJSON struct {
	Password string `json:"password"`
	Role     string `json:"role"`
	Remote   bool   `json:"remote"`
}

type UserManager struct {
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

var (
	UserFile = "./profiles/user.dat"
)

func (u *UserManager) InitUser() {

	clilog.Info("正在初始化数据表中...")
	random_password_no_salt := uuid.New().String()
	empty, err_empty := safe.EncryptByte([]byte("{}"))
	if err_empty != nil {
		clilog.Error("加密数据失败！")
		os.Exit(0)
	}
	os.WriteFile(UserFile, empty, 0600)

	userJsonMap := u.AddUser("root", random_password_no_salt)
	u.SaveUserFile(userJsonMap)

	clilog.Success(
		fmt.Sprintf(
			"初始化用户信息表完毕, 这是root的密码 %s。不包含前后符号，请在第一时间进行修改，命令：misaka-tools chpwd root。",
			random_password_no_salt,
		),
	)
}

func (u *UserManager) SaveUserFile(userMap map[string]UserJSON) {

	jsonData, err := json.Marshal(userMap)

	// 加密用户信息
	cipherData, err := safe.EncryptByte(jsonData)
	if err != nil {
		clilog.Error("加密数据失败！")
		os.Exit(0)
	}

	err = os.WriteFile(UserFile, cipherData, 0600)
	if err != nil {
		clilog.Error("无法写入文件，请检查根目录中是否存在profiles文件夹。")
		panic("system exit.")
	}

}

func (u *UserManager) LoadUserFile() map[string]UserJSON {
	cryData, err := os.ReadFile(UserFile)
	if err != nil {
		clilog.Error("无法读取原始的用户数据，请检查对应文件是否存在!")
		os.Exit(0)
	}

	plain, err := safe.DecryptByte(cryData)
	if err != nil {
		clilog.Error("无法解密的用户数据!")
		os.Exit(0)
	}

	userJson := make(map[string]UserJSON)

	err = json.Unmarshal(plain, &userJson)
	if err != nil {
		clilog.Error("无法解密的用户数据!")
		os.Exit(0)
	}
	return userJson
}

func (u *UserManager) AddUser(username string, password string) map[string]UserJSON {
	password_hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}
	userJSON := UserJSON{
		Password: string(password_hash[:]),
		Role:     "default",
		Remote:   true,
	}

	var user map[string]UserJSON
	user = u.LoadUserFile()
	if _, ok := user[username]; ok {
		clilog.Error(fmt.Sprintf("用户 %s 已存在", username))
		os.Exit(0)
	}

	user[username] = userJSON
	u.SaveUserFile(user)
	return user
}

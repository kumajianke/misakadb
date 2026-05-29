package miusers

import (
	"encoding/json"
	"errors"
	"fmt"
	"misakadb/safe"
	"os"
	"strings"

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

func (u *UserManager) InitUser() (string, error) {
	random_password_no_salt := uuid.New().String()
	empty, err := safe.EncryptByte([]byte("{}"))
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll("./profiles", 0700); err != nil {
		return "", err
	}

	err = os.WriteFile(UserFile, empty, 0600)
	if err != nil {
		return "", errors.New("无法写入用户文件，请检查 profiles 目录权限")
	}

	if _, err := u.AddUserWithRole("root", random_password_no_salt, "root"); err != nil {
		return "", err
	}

	return random_password_no_salt, nil
}

func (u *UserManager) SaveUserFile(userMap map[string]UserJSON) error {
	jsonData, err := json.Marshal(userMap)
	if err != nil {
		return err
	}

	// 加密用户信息
	cipherData, err := safe.EncryptByte(jsonData)
	if err != nil {
		return err
	}

	err = os.WriteFile(UserFile, cipherData, 0600)
	if err != nil {
		return errors.New("无法写入用户文件，请检查 profiles 目录权限")
	}

	return nil
}

func (u *UserManager) LoadUserFile() (map[string]UserJSON, error) {
	cryData, err := os.ReadFile(UserFile)
	if err != nil {
		return nil, errors.New("无法读取用户数据，请检查 user.dat 是否存在")
	}

	plain, err := safe.DecryptByte(cryData)
	if err != nil {
		return nil, errors.New("无法解密用户数据，请检查 AES 密钥是否匹配")
	}

	userJson := make(map[string]UserJSON)

	err = json.Unmarshal(plain, &userJson)
	if err != nil {
		return nil, errors.New("用户数据格式错误")
	}

	return userJson, nil
}

func (u *UserManager) AddUser(username string, password string) (map[string]UserJSON, error) {
	return u.AddUserWithRole(username, password, "default")
}

func (u *UserManager) AddUserWithRole(username string, password string, role string) (map[string]UserJSON, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	role = strings.TrimSpace(role)

	if username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if password == "" {
		return nil, errors.New("密码不能为空")
	}
	if role == "" {
		return nil, errors.New("角色不能为空")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	userJSON := UserJSON{
		Password: string(passwordHash),
		Role:     role,
		Remote:   true,
	}

	user, err := u.LoadUserFile()
	if err != nil {
		return nil, err
	}
	if _, ok := user[username]; ok {
		return nil, fmt.Errorf("用户 %s 已存在", username)
	}

	user[username] = userJSON
	if err := u.SaveUserFile(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManager) VerifyPassword(username string, password string) error {
	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	user, ok := userMap[username]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.New("密码错误")
	}

	return nil
}

func (u *UserManager) VerifyRole(username string, role string) error {
	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	user, ok := userMap[username]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	if user.Role != strings.TrimSpace(role) {
		return fmt.Errorf("用户 %s 没有 %s 权限", username, role)
	}

	return nil
}

func (u *UserManager) ChangePassword(username string, password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("新密码不能为空")
	}

	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	user, ok := userMap[username]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(passwordHash)
	userMap[username] = user
	return u.SaveUserFile(userMap)
}

func (u *UserManager) ChangeRole(username string, role string) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return errors.New("角色不能为空")
	}

	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	user, ok := userMap[username]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	user.Role = role
	userMap[username] = user
	return u.SaveUserFile(userMap)
}

func (u *UserManager) SetRemote(username string, remote bool) error {
	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	user, ok := userMap[username]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	user.Remote = remote
	userMap[username] = user
	return u.SaveUserFile(userMap)
}

func (u *UserManager) RemoveUser(username string) error {
	if strings.TrimSpace(username) == "root" {
		return errors.New("root 用户不允许删除")
	}

	userMap, err := u.LoadUserFile()
	if err != nil {
		return err
	}

	if _, ok := userMap[username]; !ok {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	delete(userMap, username)
	return u.SaveUserFile(userMap)
}

package logic

import (
	"errors"
	"web_app/dao/mysql"
	"web_app/middleware"
	"web_app/models"
	"web_app/utils"
)

func RegisterUser(req *models.RegisterRequest) error {

	// 检查用户名是否已存在
	if exists, err := mysql.CheckUserExists(req.Username); err != nil {
		return errors.New("数据库错误: " + err.Error())
	} else if exists {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败: " + err.Error())
	}

	// 生成雪花算法ID
	userID := utils.GenerateID()

	// 插入用户到数据库
	//调用dao层插入用户
	err = mysql.InsertUser(userID, req.Username, req.Email, hashedPassword)
	if err != nil {
		return errors.New("插入用户失败: " + err.Error())
	}
	return nil
}

func Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 1. 查询用户
	user, err := mysql.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 验证密码
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("密码错误")
	}

	// 3. 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成token失败: " + err.Error())
	}

	return &models.LoginResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	}, nil
}

func GetUserInfo(userID int64) (*models.User, error) {
	user, err := mysql.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	return user, nil
}

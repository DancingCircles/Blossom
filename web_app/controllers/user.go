// Package controllers 处理HTTP请求的控制器
package controllers

import (
	"net/http"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserController 用户控制器
type UserController struct{}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息"
// @Success 200 {object} models.Response
// @Router /api/v1/register [post]
func (uc *UserController) Register(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数验证失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误: "+err.Error()))
		return
	}

	// 2. 调用逻辑层进行注册处理
	if err := logic.RegisterUser(&req); err != nil {
		zap.L().Error("用户注册失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeServerError, err.Error()))
		return
	}

	// 3. 返回响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{
		"message": "注册成功",
	}))
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Router /api/v1/login [post]
func (uc *UserController) Login(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(models.CodeInvalidParams, "参数错误"))
		return
	}

	// 2. 调用逻辑层校验用户并生成token
	loginResp, err := logic.Login(&req)
	if err != nil {
		zap.L().Error("用户登录失败", zap.Error(err))
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, err.Error()))
		return
	}

	// 3. 返回登录成功响应
	c.JSON(http.StatusOK, models.NewSuccessResponse(loginResp))
}

// GetUserInfo 获取用户信息
// @Summary 获取当前用户信息
// @Description 根据token获取用户信息
// @Tags 用户
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.User}
// @Router /api/v1/user/info [get]
func (uc *UserController) GetUserInfo(c *gin.Context) {
	// 1. 从context中获取当前用户ID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.NewErrorResponse(models.CodeUnauthorized, "用户未登录"))
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "用户ID格式解析错误"))
		return
	}

	// 2. 调用逻辑层获取用户信息
	user, err := logic.GetUserInfo(userID)
	if err != nil {
		zap.L().Error("获取用户信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(models.CodeServerError, "获取用户信息失败"))
		return
	}

	// 3. 返回用户信息
	c.JSON(http.StatusOK, models.NewSuccessResponse(user))
}

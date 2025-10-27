// API 调用封装 - 简单易懂的版本

// API 基础地址（根据你的后端配置修改）
const API_BASE_URL = 'http://localhost:8082/api/v1';

// 通用的 API 请求函数
async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    
    // 默认配置
    const config = {
        headers: {
            'Content-Type': 'application/json',
        },
        ...options
    };

    // 如果有 token，添加到请求头
    const token = localStorage.getItem('token');
    if (token) {
        config.headers['Authorization'] = `Bearer ${token}`;
    }

    try {
        const response = await fetch(url, config);
        const data = await response.json();

        if (!response.ok) {
            // 创建错误对象，包含HTTP状态码和响应数据
            const error = new Error(data.message || '请求失败');
            error.status = response.status;
            error.code = data.code;
            error.response = data;
            throw error;
        }

        return data;
    } catch (error) {
        console.error('API 请求错误:', error);
        throw error;
    }
}

// ========== 用户相关 API ==========

// 用户注册
async function register(username, email, password) {
    return apiRequest('/register', {
        method: 'POST',
        body: JSON.stringify({ username, email, password })
    });
}

// 用户登录
async function login(username, password) {
    const response = await apiRequest('/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
    });
    
    // 登录成功后保存 token（后端返回格式：{code, message, data}）
    if (response.data && response.data.token) {
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('username', response.data.username);
        localStorage.setItem('user_id', response.data.user_id);
    }
    
    return response.data;
}

// 用户登出
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('user_id');
    window.location.href = 'index.html';
}

// 检查是否登录
function isLoggedIn() {
    return !!localStorage.getItem('token');
}

// 获取当前用户名
function getCurrentUsername() {
    return localStorage.getItem('username') || '游客';
}

// ========== 话题相关 API ==========

// 获取话题列表
async function getTopics(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    return apiRequest(`/topics?${queryString}`);
}

// 获取单个话题详情
async function getTopic(id) {
    return apiRequest(`/topics/${id}`);
}

// 创建话题
async function createTopic(title, content, category) {
    return apiRequest('/topics', {
        method: 'POST',
        body: JSON.stringify({ title, content, category })
    });
}

// 点赞话题
async function likeTopic(topicId) {
    return apiRequest(`/topics/${topicId}/vote?type=like`, {
        method: 'POST'
    });
}

// 点踩话题
async function dislikeTopic(topicId) {
    return apiRequest(`/topics/${topicId}/vote?type=dislike`, {
        method: 'POST'
    });
}

// ========== 搜索相关 API ==========

// 搜索话题
async function searchTopics(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    return apiRequest(`/search?${queryString}`);
}

// 获取搜索建议
async function getSearchSuggestions(prefix) {
    return apiRequest(`/search/suggest?prefix=${encodeURIComponent(prefix)}`);
}

// 获取分类热门话题
async function getHotTopicsByCategory(category, size = 10) {
    return apiRequest(`/search/hot?category=${encodeURIComponent(category)}&size=${size}`);
}

// 获取分类统计
async function getCategoryStats() {
    return apiRequest('/search/stats');
}

// ========== 评论相关 API ==========

// 获取话题评论列表
async function getComments(topicId, params = {}) {
    const queryString = new URLSearchParams(params).toString();
    return apiRequest(`/topics/${topicId}/comments?${queryString}`);
}

// 创建评论
async function createComment(topicId, content, parentId = null) {
    const body = { content };
    if (parentId) {
        body.parent_id = parentId;
    }
    return apiRequest(`/topics/${topicId}/comments`, {
        method: 'POST',
        body: JSON.stringify(body)
    });
}

// 删除评论
async function deleteComment(commentId) {
    return apiRequest(`/comments/${commentId}`, {
        method: 'DELETE'
    });
}

// ========== 工具函数 ==========

// 显示提示消息
function showMessage(message, type = 'info') {
    // 创建提示框
    const msgBox = document.createElement('div');
    msgBox.className = `message-box message-${type}`;
    msgBox.textContent = message;
    
    // 添加样式
    msgBox.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 16px 24px;
        background: ${type === 'error' ? '#f44336' : type === 'success' ? '#4CAF50' : '#2196F3'};
        color: white;
        border-radius: 4px;
        font-weight: 600;
        z-index: 10000;
        animation: slideIn 0.3s ease;
    `;
    
    document.body.appendChild(msgBox);
    
    // 3秒后自动移除
    setTimeout(() => {
        msgBox.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => msgBox.remove(), 300);
    }, 3000);
}

// 添加动画样式
if (!document.querySelector('#message-animations')) {
    const style = document.createElement('style');
    style.id = 'message-animations';
    style.textContent = `
        @keyframes slideIn {
            from {
                transform: translateX(400px);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
        @keyframes slideOut {
            from {
                transform: translateX(0);
                opacity: 1;
            }
            to {
                transform: translateX(400px);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);
}


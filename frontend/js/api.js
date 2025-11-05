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
    // 创建提示框容器
    const msgBox = document.createElement('div');
    msgBox.className = `message-box message-${type}`;
    
    // 根据类型设置图标
    const icons = {
        'success': '✓',
        'error': '✗',
        'info': 'ℹ'
    };
    
    const colors = {
        'success': '#2e7d32',
        'error': '#b71c1c',
        'info': '#1565c0'
    };
    
    msgBox.innerHTML = `
        <div class="message-icon">${icons[type] || 'ℹ'}</div>
        <div class="message-text">${message}</div>
    `;
    
    // 添加样式
    msgBox.style.cssText = `
        position: fixed;
        top: 100px;
        left: 50%;
        transform: translateX(-50%) translateY(-100px);
        padding: 20px 40px;
        background: ${colors[type]};
        color: #f5ebe0;
        border: 5px solid #3d0000;
        box-shadow: 8px 8px 0 rgba(61, 0, 0, 0.3);
        font-weight: 900;
        font-size: 18px;
        text-transform: uppercase;
        letter-spacing: 1px;
        z-index: 10000;
        animation: messageSlideDown 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
        display: flex;
        align-items: center;
        gap: 15px;
        font-family: 'Impact', 'Arial Black', sans-serif;
        min-width: 300px;
        justify-content: center;
    `;
    
    document.body.appendChild(msgBox);
    
    // 3秒后自动移除
    setTimeout(() => {
        msgBox.style.animation = 'messageSlideUp 0.3s ease forwards';
        setTimeout(() => msgBox.remove(), 300);
    }, 3000);
}

// 自定义确认对话框
function showConfirm(message, onConfirm, onCancel) {
    // 创建遮罩层
    const overlay = document.createElement('div');
    overlay.className = 'confirm-overlay';
    overlay.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.7);
        z-index: 9999;
        display: flex;
        align-items: center;
        justify-content: center;
        animation: fadeIn 0.2s ease;
    `;
    
    // 创建对话框
    const dialog = document.createElement('div');
    dialog.className = 'confirm-dialog';
    dialog.innerHTML = `
        <div class="confirm-header">确认操作</div>
        <div class="confirm-message">${message}</div>
        <div class="confirm-buttons">
            <button class="confirm-btn confirm-yes">确定</button>
            <button class="confirm-btn confirm-no">取消</button>
        </div>
    `;
    
    dialog.style.cssText = `
        background: #f5ebe0;
        border: 6px solid #3d0000;
        box-shadow: 12px 12px 0 rgba(61, 0, 0, 0.3);
        padding: 40px;
        min-width: 400px;
        animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    `;
    
    overlay.appendChild(dialog);
    document.body.appendChild(overlay);
    
    // 按钮点击事件
    const yesBtn = dialog.querySelector('.confirm-yes');
    const noBtn = dialog.querySelector('.confirm-no');
    
    yesBtn.onclick = () => {
        overlay.style.animation = 'fadeOut 0.2s ease';
        setTimeout(() => {
            overlay.remove();
            if (onConfirm) onConfirm();
        }, 200);
    };
    
    noBtn.onclick = () => {
        overlay.style.animation = 'fadeOut 0.2s ease';
        setTimeout(() => {
            overlay.remove();
            if (onCancel) onCancel();
        }, 200);
    };
    
    // 点击遮罩关闭
    overlay.onclick = (e) => {
        if (e.target === overlay) {
            noBtn.click();
        }
    };
}

// 添加动画和样式
if (!document.querySelector('#message-animations')) {
    const style = document.createElement('style');
    style.id = 'message-animations';
    style.textContent = `
        @keyframes messageSlideDown {
            from {
                transform: translateX(-50%) translateY(-100px);
                opacity: 0;
            }
            to {
                transform: translateX(-50%) translateY(0);
                opacity: 1;
            }
        }
        @keyframes messageSlideUp {
            from {
                transform: translateX(-50%) translateY(0);
                opacity: 1;
            }
            to {
                transform: translateX(-50%) translateY(-100px);
                opacity: 0;
            }
        }
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        @keyframes fadeOut {
            from { opacity: 1; }
            to { opacity: 0; }
        }
        @keyframes scaleIn {
            from {
                transform: scale(0.8);
                opacity: 0;
            }
            to {
                transform: scale(1);
                opacity: 1;
            }
        }
        
        .message-icon {
            font-size: 28px;
            font-weight: 900;
        }
        
        .message-text {
            flex: 1;
        }
        
        .confirm-header {
            font-family: 'Impact', 'Arial Black', sans-serif;
            font-size: 32px;
            font-weight: 900;
            text-transform: uppercase;
            color: #3d0000;
            margin-bottom: 24px;
            text-align: center;
            letter-spacing: -1px;
        }
        
        .confirm-message {
            font-size: 18px;
            color: #3d0000;
            margin-bottom: 32px;
            line-height: 1.6;
            text-align: center;
            font-weight: 600;
        }
        
        .confirm-buttons {
            display: flex;
            gap: 20px;
            justify-content: center;
        }
        
        .confirm-btn {
            padding: 14px 36px;
            font-size: 16px;
            font-weight: 900;
            text-transform: uppercase;
            border: 4px solid #3d0000;
            cursor: pointer;
            transition: all 0.2s ease;
            font-family: 'Impact', 'Arial Black', sans-serif;
            letter-spacing: 1px;
            box-shadow: 4px 4px 0 rgba(61, 0, 0, 0.2);
        }
        
        .confirm-yes {
            background: #b71c1c;
            color: #f5ebe0;
        }
        
        .confirm-yes:hover {
            background: #8b0000;
            transform: translate(-2px, -2px);
            box-shadow: 6px 6px 0 rgba(61, 0, 0, 0.3);
        }
        
        .confirm-no {
            background: #f5ebe0;
            color: #3d0000;
        }
        
        .confirm-no:hover {
            background: #fff;
            transform: translate(-2px, -2px);
            box-shadow: 6px 6px 0 rgba(61, 0, 0, 0.3);
        }
    `;
    document.head.appendChild(style);
}


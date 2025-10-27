// 话题详情页面逻辑

let currentTopicId = null;
let currentPage = 1;
let totalPages = 1;

document.addEventListener('DOMContentLoaded', () => {
    // 从URL获取话题ID
    const urlParams = new URLSearchParams(window.location.search);
    currentTopicId = urlParams.get('id');

    if (!currentTopicId) {
        showMessage('话题不存在', 'error');
        setTimeout(() => {
            window.location.href = 'index.html';
        }, 1500);
        return;
    }

    initDetailPage();
});

async function initDetailPage() {
    // 加载话题详情
    await loadTopicDetail();
    
    // 加载评论列表
    await loadComments(1);
    
    // 初始化评论表单
    initCommentForm();
    
    // 初始化加载更多按钮
    initLoadMore();
}

// 加载话题详情
async function loadTopicDetail() {
    const container = document.getElementById('topic-detail');
    
    try {
        const response = await getTopic(currentTopicId);
        const topic = response.data;
        
        // 渲染话题详情
        container.innerHTML = renderTopicDetail(topic);
        
        // 绑定投票按钮事件
        bindVoteButtons(topic);
        
    } catch (error) {
        // 如果是404错误，显示友好提示并跳转
        if (error.status === 404) {
            container.innerHTML = `
                <div class="error-message" style="text-align: center; padding: 60px 20px;">
                    <div style="font-size: 64px; margin-bottom: 20px;">😕</div>
                    <h2 style="margin-bottom: 10px;">话题不存在</h2>
                    <p style="color: #666; margin-bottom: 20px;">该话题可能已被删除或不存在</p>
                    <p style="color: #999; font-size: 14px;">3秒后自动返回首页...</p>
                </div>
            `;
            showMessage('话题不存在，即将返回首页', 'error');
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 3000);
        } else {
            // 其他错误显示重新加载按钮
            container.innerHTML = `
                <div class="error-message">
                    <p>加载失败：${error.message}</p>
                    <button onclick="window.location.reload()" class="btn-secondary">重新加载</button>
                </div>
            `;
        }
    }
}

// 渲染话题详情
function renderTopicDetail(topic) {
    const tagClasses = {
        'tech': '',
        'design': 'tag-design',
        'discuss': 'tag-discuss',
        'share': 'tag-share',
        'product': 'tag-product'
    };

    const tagNames = {
        'tech': '技术',
        'design': '设计',
        'discuss': '讨论',
        'share': '分享',
        'product': '产品'
    };

    return `
        <div class="topic-header">
            <div class="user-info">
                <div class="user-avatar">👤</div>
                <div>
                    <h4 class="username">${topic.username || '匿名用户'}</h4>
                    <span class="post-time">${formatTime(topic.created_at)}</span>
                </div>
            </div>
            <span class="tag ${tagClasses[topic.category]}">${tagNames[topic.category]}</span>
        </div>
        <h2 class="topic-title">${escapeHtml(topic.title)}</h2>
        <div class="topic-content">${escapeHtml(topic.content).replace(/\n/g, '<br>')}</div>
        <div class="topic-stats">
            <button class="stat-item vote-btn" data-type="like" data-topic-id="${topic.id}">
                👍 ${topic.like_count || 0}
            </button>
            <button class="stat-item vote-btn" data-type="dislike" data-topic-id="${topic.id}">
                👎 ${topic.dislike_count || 0}
            </button>
            <span class="stat-item">💬 ${topic.comment_count || 0}</span>
            <span class="stat-item">👁️ ${topic.view_count || 0}</span>
        </div>
    `;
}

// 绑定投票按钮
function bindVoteButtons(topic) {
    const voteButtons = document.querySelectorAll('.vote-btn');
    
    voteButtons.forEach(btn => {
        btn.addEventListener('click', async (e) => {
            if (!isLoggedIn()) {
                showMessage('请先登录', 'error');
                setTimeout(() => {
                    window.location.href = 'login.html';
                }, 1000);
                return;
            }

            const type = btn.dataset.type;
            const topicId = btn.dataset.topicId;

            try {
                if (type === 'like') {
                    await likeTopic(topicId);
                } else {
                    await dislikeTopic(topicId);
                }
                
                showMessage('投票成功', 'success');
                
                // 重新加载话题详情
                await loadTopicDetail();
                
            } catch (error) {
                showMessage(error.message || '投票失败', 'error');
            }
        });
    });
}

// 加载评论列表
async function loadComments(page = 1) {
    const container = document.getElementById('comments-list');
    
    try {
        const response = await getComments(currentTopicId, { page, page_size: 20 });
        const data = response.data;
        
        // 更新评论数量
        document.getElementById('comment-count').textContent = data.total || 0;
        
        // 如果是第一页，清空容器
        if (page === 1) {
            container.innerHTML = '';
        }
        
        // 渲染评论
        if (data.comments && data.comments.length > 0) {
            data.comments.forEach(comment => {
                container.appendChild(renderComment(comment));
            });
        } else if (page === 1) {
            container.innerHTML = '<div class="empty-state">还没有评论，快来抢沙发吧！</div>';
        }
        
        // 更新分页信息
        currentPage = page;
        totalPages = data.total_pages || 1;
        
        // 显示/隐藏"加载更多"按钮
        const loadMoreContainer = document.getElementById('load-more-container');
        if (data.has_more) {
            loadMoreContainer.style.display = 'block';
        } else {
            loadMoreContainer.style.display = 'none';
        }
        
    } catch (error) {
        // 如果是404错误（话题不存在），不显示评论错误
        // 因为话题详情已经会处理并跳转
        if (error.status === 404) {
            if (page === 1) {
                container.innerHTML = '<div class="empty-state">话题不存在，无法加载评论</div>';
            }
            return;
        }
        
        if (page === 1) {
            container.innerHTML = `<div class="error-message">加载评论失败：${error.message}</div>`;
        } else {
            showMessage('加载失败：' + error.message, 'error');
        }
    }
}

// 渲染单个评论
function renderComment(comment) {
    const div = document.createElement('div');
    div.className = 'comment-item';
    div.dataset.commentId = comment.id;
    
    const currentUserId = localStorage.getItem('user_id');
    const isAuthor = currentUserId && String(comment.user_id) === String(currentUserId);
    
    div.innerHTML = `
        <div class="comment-avatar">👤</div>
        <div class="comment-body">
            <div class="comment-header">
                <span class="comment-username">${escapeHtml(comment.username || '匿名用户')}</span>
                <span class="comment-time">${formatTime(comment.created_at)}</span>
            </div>
            <div class="comment-content">${escapeHtml(comment.content)}</div>
            ${isAuthor ? `
                <div class="comment-actions">
                    <button class="btn-link delete-comment-btn" data-comment-id="${comment.id}">删除</button>
                </div>
            ` : ''}
        </div>
    `;
    
    // 绑定删除按钮事件
    if (isAuthor) {
        const deleteBtn = div.querySelector('.delete-comment-btn');
        deleteBtn.addEventListener('click', () => handleDeleteComment(comment.id));
    }
    
    return div;
}

// 初始化评论表单
function initCommentForm() {
    const formContainer = document.getElementById('comment-form-container');
    const form = document.getElementById('comment-form');
    const input = document.getElementById('comment-input');
    const charCount = document.getElementById('comment-char-count');
    
    // 检查登录状态
    if (isLoggedIn()) {
        formContainer.style.display = 'block';
    } else {
        formContainer.innerHTML = `
            <div class="login-prompt">
                <p>登录后才能发表评论</p>
                <a href="login.html" class="btn-primary btn-small">立即登录</a>
            </div>
        `;
        formContainer.style.display = 'block';
        return;
    }
    
    // 字符计数
    input.addEventListener('input', () => {
        charCount.textContent = input.value.length;
    });
    
    // 提交评论
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const content = input.value.trim();
        
        if (!content) {
            showMessage('请输入评论内容', 'error');
            return;
        }
        
        if (content.length < 1) {
            showMessage('评论内容至少1个字符', 'error');
            return;
        }
        
        try {
            await createComment(currentTopicId, content);
            showMessage('评论发表成功', 'success');
            
            // 清空输入框
            input.value = '';
            charCount.textContent = '0';
            
            // 重新加载评论列表
            await loadComments(1);
            
            // 重新加载话题详情（更新评论数）
            await loadTopicDetail();
            
        } catch (error) {
            showMessage(error.message || '发表评论失败', 'error');
        }
    });
}

// 初始化"加载更多"按钮
function initLoadMore() {
    const btn = document.getElementById('load-more-btn');
    
    btn.addEventListener('click', async () => {
        btn.disabled = true;
        btn.textContent = '加载中...';
        
        await loadComments(currentPage + 1);
        
        btn.disabled = false;
        btn.textContent = '加载更多评论';
    });
}

// 处理删除评论
async function handleDeleteComment(commentId) {
    if (!confirm('确定要删除这条评论吗？')) {
        return;
    }
    
    try {
        await deleteComment(commentId);
        showMessage('评论已删除', 'success');
        
        // 从DOM中移除评论
        const commentElement = document.querySelector(`[data-comment-id="${commentId}"]`);
        if (commentElement) {
            commentElement.remove();
        }
        
        // 重新加载话题详情（更新评论数）
        await loadTopicDetail();
        
        // 更新评论计数
        const commentCount = document.getElementById('comment-count');
        const currentCount = parseInt(commentCount.textContent) || 0;
        commentCount.textContent = Math.max(0, currentCount - 1);
        
    } catch (error) {
        showMessage(error.message || '删除失败', 'error');
    }
}

// 格式化时间
function formatTime(timeStr) {
    const time = new Date(timeStr);
    const now = new Date();
    const diff = now - time;
    
    const minute = 60 * 1000;
    const hour = 60 * minute;
    const day = 24 * hour;
    
    if (diff < minute) {
        return '刚刚';
    } else if (diff < hour) {
        return Math.floor(diff / minute) + '分钟前';
    } else if (diff < day) {
        return Math.floor(diff / hour) + '小时前';
    } else if (diff < 7 * day) {
        return Math.floor(diff / day) + '天前';
    } else {
        return time.toLocaleDateString('zh-CN');
    }
}

// HTML转义
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}


// 评论树形结构渲染系统

// 构建评论树
function buildCommentTree(comments) {
    const commentMap = {};
    const rootComments = [];
    
    // 第一遍：建立ID映射
    comments.forEach(comment => {
        comment.replies = [];
        commentMap[comment.id] = comment;
    });
    
    // 第二遍：建立父子关系
    comments.forEach(comment => {
        if (comment.parent_id && commentMap[comment.parent_id]) {
            // 添加到父评论的replies中
            commentMap[comment.parent_id].replies.push(comment);
        } else {
            // 顶级评论
            rootComments.push(comment);
        }
    });
    
    return rootComments;
}

// 渲染评论树（递归）- 抖音风格
function renderCommentTree(comment, level = 0) {
    const wrapper = document.createElement('div');
    wrapper.className = 'comment-wrapper';
    wrapper.dataset.level = level;
    
    const div = document.createElement('div');
    div.className = 'comment-item';
    if (level > 0) {
        div.classList.add('reply-comment');
        div.style.marginLeft = `${Math.min(level * 40, 120)}px`;
    }
    div.dataset.commentId = comment.id;
    
    const currentUserId = localStorage.getItem('user_id');
    const isAuthor = currentUserId && String(comment.user_id) === String(currentUserId);
    const isLoggedInUser = isLoggedIn();
    
    const replyCount = comment.replies ? comment.replies.length : 0;
    
    div.innerHTML = `
        <div class="comment-avatar">👤</div>
        <div class="comment-body">
            <div class="comment-header">
                <div class="comment-user">
                    <span class="comment-username">${escapeHtml(comment.username || '匿名用户')}</span>
                    ${level > 0 ? '<span class="reply-badge">回复</span>' : ''}
                </div>
                <span class="comment-time">${formatTime(comment.created_at)}</span>
            </div>
            <div class="comment-content">${escapeHtml(comment.content)}</div>
            <div class="comment-actions">
                ${isLoggedInUser ? `
                    <button class="btn-link reply-comment-btn">💬 回复</button>
                ` : ''}
                ${replyCount > 0 ? `
                    <button class="btn-link toggle-replies-btn">
                        <span class="toggle-text">展开</span> ${replyCount} 条回复
                    </button>
                ` : ''}
                ${isAuthor ? `
                    <button class="btn-link delete-comment-btn">🗑️ 删除</button>
                ` : ''}
            </div>
            <div class="reply-form-container" style="display: none;">
                <textarea 
                    class="reply-textarea" 
                    placeholder="写下你的回复..." 
                    rows="3"
                    maxlength="1000"></textarea>
                <div class="reply-form-actions">
                    <span class="reply-char-count"><span class="char-number">0</span>/1000</span>
                    <button class="btn-cancel-reply">取消</button>
                    <button class="btn-submit-reply">发表回复</button>
                </div>
            </div>
        </div>
    `;
    
    wrapper.appendChild(div);
    
    // 创建回复容器（抖音风格：默认折叠）
    if (replyCount > 0) {
        const repliesContainer = document.createElement('div');
        repliesContainer.className = 'replies-container';
        repliesContainer.style.display = 'none'; // 默认折叠
        
        // 递归渲染子回复
        comment.replies.forEach(reply => {
            repliesContainer.appendChild(renderCommentTree(reply, level + 1));
        });
        
        wrapper.appendChild(repliesContainer);
    }
    
    // 绑定事件
    bindCommentEvents(wrapper, comment);
    
    return wrapper;
}

// 绑定评论事件
function bindCommentEvents(wrapper, comment) {
    const isLoggedInUser = isLoggedIn();
    const currentUserId = localStorage.getItem('user_id');
    const isAuthor = currentUserId && String(comment.user_id) === String(currentUserId);
    
    // 回复按钮
    if (isLoggedInUser) {
        const replyBtn = wrapper.querySelector('.reply-comment-btn');
        if (replyBtn) {
            replyBtn.addEventListener('click', () => {
                toggleReplyForm(wrapper, comment);
            });
        }
    }
    
    // 展开/折叠回复（抖音风格）
    const toggleBtn = wrapper.querySelector('.toggle-replies-btn');
    if (toggleBtn) {
        toggleBtn.addEventListener('click', () => {
            const repliesContainer = wrapper.querySelector('.replies-container');
            const toggleText = toggleBtn.querySelector('.toggle-text');
            const replyCount = comment.replies ? comment.replies.length : 0;
            
            if (repliesContainer.style.display === 'none') {
                // 展开
                repliesContainer.style.display = 'block';
                toggleText.textContent = '收起';
                toggleBtn.innerHTML = `<span class="toggle-text">收起</span> ${replyCount} 条回复`;
            } else {
                // 收起
                repliesContainer.style.display = 'none';
                toggleText.textContent = '展开';
                toggleBtn.innerHTML = `<span class="toggle-text">展开</span> ${replyCount} 条回复`;
            }
        });
    }
    
    // 删除按钮
    if (isAuthor) {
        const deleteBtn = wrapper.querySelector('.delete-comment-btn');
        if (deleteBtn) {
            deleteBtn.addEventListener('click', async () => {
                if (confirm('确定要删除这条评论吗？')) {
                    try {
                        await deleteComment(comment.id);
                        showMessage('删除成功', 'success');
                        await loadComments(1);
                        await loadTopicDetail();
                    } catch (error) {
                        showMessage(error.message || '删除失败', 'error');
                    }
                }
            });
        }
    }
}

// 切换回复表单
function toggleReplyForm(wrapper, comment) {
    const replyContainer = wrapper.querySelector('.reply-form-container');
    const textarea = wrapper.querySelector('.reply-textarea');
    const charNumber = wrapper.querySelector('.char-number');
    const cancelBtn = wrapper.querySelector('.btn-cancel-reply');
    const submitBtn = wrapper.querySelector('.btn-submit-reply');
    
    // 隐藏所有其他回复框
    document.querySelectorAll('.reply-form-container').forEach(form => {
        if (form !== replyContainer) {
            form.style.display = 'none';
        }
    });
    
    // 切换显示
    const isVisible = replyContainer.style.display !== 'none';
    replyContainer.style.display = isVisible ? 'none' : 'block';
    
    if (!isVisible) {
        textarea.focus();
        textarea.placeholder = `回复 @${comment.username}...`;
    }
    
    // 字符计数
    textarea.oninput = () => {
        charNumber.textContent = textarea.value.length;
    };
    
    // 取消
    cancelBtn.onclick = () => {
        replyContainer.style.display = 'none';
        textarea.value = '';
        charNumber.textContent = '0';
    };
    
    // 提交
    submitBtn.onclick = async () => {
        const content = textarea.value.trim();
        
        if (!content) {
            showMessage('请输入回复内容', 'error');
            return;
        }
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '发表中...';
            
            await createComment(currentTopicId, content, comment.id);
            showMessage('回复成功', 'success');
            
            textarea.value = '';
            charNumber.textContent = '0';
            replyContainer.style.display = 'none';
            
            // 重新加载评论
            await loadComments(1);
            await loadTopicDetail();
            
        } catch (error) {
            showMessage(error.message || '回复失败', 'error');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = '发表回复';
        }
    };
}


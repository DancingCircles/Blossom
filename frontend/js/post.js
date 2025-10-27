// 发帖页面逻辑

document.addEventListener('DOMContentLoaded', () => {
    // 检查登录状态
    if (!isLoggedIn()) {
        showMessage('请先登录', 'error');
        setTimeout(() => {
            window.location.href = 'login.html';
        }, 1500);
        return;
    }

    initPostPage();
});

function initPostPage() {
    const form = document.getElementById('post-form');
    const titleInput = document.getElementById('post-title');
    const contentInput = document.getElementById('post-content');
    const categoryInputs = document.querySelectorAll('input[name="category"]');
    
    // 预览元素
    const previewUsername = document.getElementById('preview-username');
    const previewTitle = document.getElementById('preview-title');
    const previewContent = document.getElementById('preview-content');
    const previewTag = document.getElementById('preview-tag');

    // 设置当前用户名
    previewUsername.textContent = getCurrentUsername();

    // 字符计数
    const titleCount = document.getElementById('title-count');
    const contentCount = document.getElementById('content-count');

    titleInput.addEventListener('input', () => {
        titleCount.textContent = titleInput.value.length;
        updatePreview();
    });

    contentInput.addEventListener('input', () => {
        contentCount.textContent = contentInput.value.length;
        updatePreview();
    });

    // 分类选择
    categoryInputs.forEach(input => {
        input.addEventListener('change', updatePreview);
    });

    // 更新预览
    function updatePreview() {
        const title = titleInput.value.trim() || '话题标题会显示在这里';
        const content = contentInput.value.trim() || '话题内容会显示在这里...';
        const category = document.querySelector('input[name="category"]:checked').value;

        previewTitle.textContent = title;
        previewContent.textContent = content.length > 200 
            ? content.substring(0, 200) + '...' 
            : content;

        // 更新标签
        const tagNames = {
            'tech': '技术',
            'design': '设计',
            'discuss': '讨论',
            'share': '分享',
            'product': '产品'
        };
        
        const tagClasses = {
            'tech': '',
            'design': 'tag-design',
            'discuss': 'tag-discuss',
            'share': 'tag-share',
            'product': 'tag-product'
        };

        previewTag.textContent = tagNames[category];
        previewTag.className = 'tag ' + tagClasses[category];
    }

    // 表单提交
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const title = titleInput.value.trim();
        const content = contentInput.value.trim();
        const category = document.querySelector('input[name="category"]:checked').value;

        // 验证
        if (!title) {
            showMessage('请输入话题标题', 'error');
            return;
        }

        if (title.length < 5) {
            showMessage('标题至少5个字符', 'error');
            return;
        }

        if (!content) {
            showMessage('请输入话题内容', 'error');
            return;
        }

        if (content.length < 10) {
            showMessage('内容至少10个字符', 'error');
            return;
        }

        try {
            // 调用创建话题 API
            const result = await createTopic(title, content, category);
            
            showMessage('话题发布成功！', 'success');
            
            // 1.5秒后跳转到首页
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1500);
            
        } catch (error) {
            showMessage(error.message || '发布失败，请稍后重试', 'error');
        }
    });
}


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
    
    // 字符计数
    const titleCount = document.getElementById('title-count');
    const contentCount = document.getElementById('content-count');

    titleInput.addEventListener('input', () => {
        titleCount.textContent = titleInput.value.length;
    });

    contentInput.addEventListener('input', () => {
        contentCount.textContent = contentInput.value.length;
    });

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

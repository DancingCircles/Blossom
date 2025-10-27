// 登录注册页面逻辑

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', () => {
    // 如果已登录，跳转到首页
    if (isLoggedIn()) {
        window.location.href = 'index.html';
        return;
    }

    initAuthPage();
});

function initAuthPage() {
    // 获取元素
    const tabs = document.querySelectorAll('.auth-tab');
    const forms = document.querySelectorAll('.auth-form');
    const switchLinks = document.querySelectorAll('.switch-tab');
    
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');

    // 标签切换
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            switchTab(tab.dataset.tab);
        });
    });

    // 文字链接切换
    switchLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            switchTab(link.dataset.tab);
        });
    });

    // 切换标签函数
    function switchTab(tabName) {
        // 更新标签状态
        tabs.forEach(t => {
            if (t.dataset.tab === tabName) {
                t.classList.add('active');
            } else {
                t.classList.remove('active');
            }
        });

        // 更新表单显示
        forms.forEach(f => {
            if (f.id === `${tabName}-form`) {
                f.classList.add('active');
            } else {
                f.classList.remove('active');
            }
        });
    }

    // 登录表单提交
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const username = document.getElementById('login-username').value.trim();
        const password = document.getElementById('login-password').value;

        // 简单验证
        if (!username || !password) {
            showMessage('请填写完整信息', 'error');
            return;
        }

        try {
            // 调用登录 API
            const result = await login(username, password);
            
            showMessage('登录成功！', 'success');
            
            // 1秒后跳转到首页
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1000);
            
        } catch (error) {
            showMessage(error.message || '登录失败，请检查用户名和密码', 'error');
        }
    });

    // 注册表单提交
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const username = document.getElementById('register-username').value.trim();
        const email = document.getElementById('register-email').value.trim();
        const password = document.getElementById('register-password').value;
        const confirm = document.getElementById('register-confirm').value;

        // 验证
        if (!username || !email || !password || !confirm) {
            showMessage('请填写完整信息', 'error');
            return;
        }

        if (username.length < 4 || username.length > 20) {
            showMessage('用户名长度应为4-20个字符', 'error');
            return;
        }

        if (password.length < 6) {
            showMessage('密码长度至少6个字符', 'error');
            return;
        }

        if (password !== confirm) {
            showMessage('两次输入的密码不一致', 'error');
            return;
        }

        try {
            // 调用注册 API
            await register(username, email, password);
            
            showMessage('注册成功！请登录', 'success');
            
            // 切换到登录表单
            setTimeout(() => {
                switchTab('login');
                // 清空注册表单
                registerForm.reset();
            }, 1000);
            
        } catch (error) {
            showMessage(error.message || '注册失败，请稍后重试', 'error');
        }
    });
}


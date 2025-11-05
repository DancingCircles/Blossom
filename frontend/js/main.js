// Blossom 论坛 - 主 JavaScript 文件

// ========== 全局分页状态 ==========
let currentPage = 1;
let totalPages = 1;
let currentSortType = 'hot';
let currentSearchKeyword = ''; // 当前搜索关键词
const pageSize = 8; // 每页8条记录

// ========== 页面加载动画 ==========
document.addEventListener('DOMContentLoaded', () => {
    initAnimations();
    initInteractions();
    initParallax();
    initUserMenu();
    initSorting(); // 这个会自动调用 loadTopics 和 bindVoteButtons
    initPagination(); // 初始化分页功能
});

// 初始化入场动画
function initAnimations() {
    // Logo 入场动画
    const logo = document.querySelector('.bubble-logo');
    if (logo) {
        logo.style.opacity = '0';
        logo.style.transform = 'scale(0.9)';
        setTimeout(() => {
            logo.style.transition = 'all 1s cubic-bezier(0.34, 1.56, 0.64, 1)';
            logo.style.opacity = '1';
            logo.style.transform = 'scale(1)';
        }, 200);
    }

    // 话题卡片滚动入场动画
    const observerOptions = {
        threshold: 0.15,
        rootMargin: '0px 0px -80px 0px'
    };

    const cardObserver = new IntersectionObserver((entries) => {
        entries.forEach((entry, index) => {
            if (entry.isIntersecting) {
                setTimeout(() => {
                    entry.target.style.opacity = '1';
                    entry.target.style.transform = 'translateY(0) rotate(0)';
                }, index * 120);
                cardObserver.unobserve(entry.target);
            }
        });
    }, observerOptions);

    const topicCards = document.querySelectorAll('.topic-card');
    topicCards.forEach(card => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(50px) rotate(-2deg)';
        card.style.transition = 'all 0.8s cubic-bezier(0.34, 1.56, 0.64, 1)';
        cardObserver.observe(card);
    });
}

// ========== 交互功能 ==========
function initInteractions() {
    // 话题搜索框交互
    const topicSearchInput = document.getElementById('topic-search-input');
    const topicSearchBtn = document.getElementById('topic-search-btn');
    
    if (topicSearchInput) {
        // 回车搜索
        topicSearchInput.addEventListener('keypress', async (e) => {
            if (e.key === 'Enter') {
                const query = topicSearchInput.value.trim();
                if (query) {
                    console.log('回车搜索:', query);
                    await performSearch(query);
                } else {
                    showMessage('请输入搜索关键词', 'error');
                }
            }
        });
    }

    // 搜索按钮点击事件
    if (topicSearchBtn) {
        topicSearchBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            e.stopPropagation();
            const query = topicSearchInput ? topicSearchInput.value.trim() : '';
            console.log('点击搜索按钮:', query);
            if (query) {
                await performSearch(query);
            } else {
                showMessage('请输入搜索关键词', 'error');
            }
        });
    }

    // 话题卡片点击
    document.querySelectorAll('.topic-card').forEach(card => {
        card.addEventListener('click', handleTopicCardClick);
        
        // 添加悬停音效反馈（可选）
        card.addEventListener('mouseenter', () => {
            card.style.transition = 'all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)';
        });
    });

}

// 话题卡片点击处理
function handleTopicCardClick(e) {
    // 如果点击的是投票按钮，不触发卡片跳转
    if (e.target.closest('.vote-btn')) {
        return;
    }
    
    const card = e.currentTarget;
    const topicId = card.dataset.topicId;
    
    if (!topicId) {
        console.warn('话题ID不存在');
        return;
    }
    
    // 添加点击波纹效果
    createRipple(e, card);
    
    // 跳转到话题详情页
    setTimeout(() => {
        window.location.href = `detail.html?id=${topicId}`;
    }, 200);
}

// 创建点击波纹效果
function createRipple(e, element) {
    const ripple = document.createElement('div');
    ripple.style.position = 'absolute';
    ripple.style.borderRadius = '50%';
    ripple.style.background = 'rgba(91, 159, 237, 0.4)';
    ripple.style.width = '30px';
    ripple.style.height = '30px';
    ripple.style.pointerEvents = 'none';
    ripple.style.zIndex = '10';
    
    const rect = element.getBoundingClientRect();
    ripple.style.left = (e.clientX - rect.left - 15) + 'px';
    ripple.style.top = (e.clientY - rect.top - 15) + 'px';
    
    element.style.position = 'relative';
    element.appendChild(ripple);
    
    ripple.animate([
        { transform: 'scale(1)', opacity: 1 },
        { transform: 'scale(25)', opacity: 0 }
    ], {
        duration: 800,
        easing: 'cubic-bezier(0.4, 0, 0.2, 1)'
    }).onfinish = () => ripple.remove();
}

// ========== 视差效果 ==========
function initParallax() {
    let rafId = null;
    let lastMouseX = 0;
    let lastMouseY = 0;
    
    document.addEventListener('mousemove', (e) => {
        lastMouseX = e.clientX;
        lastMouseY = e.clientY;
        
        if (!rafId) {
            rafId = requestAnimationFrame(updateParallax);
        }
    });
    
    function updateParallax() {
        const mouseX = lastMouseX / window.innerWidth - 0.5;
        const mouseY = lastMouseY / window.innerHeight - 0.5;
        
        // 星星视差效果
        document.querySelectorAll('.star-icon').forEach((star, index) => {
            const speed = (index + 1) * 15;
            star.style.transform = `translate(${mouseX * speed}px, ${mouseY * speed}px)`;
        });
        
        
        // 云朵视差
        const cloud = document.querySelector('.cloud-icon');
        if (cloud) {
            cloud.style.transform = `translate(${mouseX * 30}px, ${mouseY * 20}px)`;
        }
        
        rafId = null;
    }
}

// ========== 平滑滚动 ==========
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// ========== 滚动方向检测 & 导航栏自动隐藏 ==========
let lastScrollTop = 0;
let scrollTimeout = null;
const navbar = document.querySelector('.navbar');

window.addEventListener('scroll', () => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    
    // 清除之前的定时器
    if (scrollTimeout) {
        clearTimeout(scrollTimeout);
    }
    
    // 在顶部时始终显示导航栏
    if (scrollTop <= 10) {
        navbar?.classList.remove('hidden');
        navbar?.classList.remove('scrolled');
        lastScrollTop = scrollTop;
        return;
    }
    
    // 滚动超过 50px 时增加导航栏不透明度
    if (scrollTop > 50) {
        navbar?.classList.add('scrolled');
    } else {
        navbar?.classList.remove('scrolled');
    }
    
    // 滚动方向检测（至少滚动 5px 才触发）
    if (Math.abs(scrollTop - lastScrollTop) > 5) {
        if (scrollTop > lastScrollTop && scrollTop > 100) {
            // 向下滚动且超过100px：隐藏导航栏
            navbar?.classList.add('hidden');
            document.body.classList.add('scrolling-down');
            document.body.classList.remove('scrolling-up');
        } else {
            // 向上滚动：显示导航栏
            navbar?.classList.remove('hidden');
            document.body.classList.add('scrolling-up');
            document.body.classList.remove('scrolling-down');
        }
    }
    
    lastScrollTop = scrollTop;
}, { passive: true });

// ========== 键盘快捷键 ==========
document.addEventListener('keydown', (e) => {
    // Ctrl/Cmd + K: 聚焦搜索框
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        const searchInput = document.querySelector('.search-input');
        if (searchInput) {
            searchInput.focus();
        }
    }
    
    // ESC: 关闭菜单
    if (e.key === 'Escape') {
        const menuBtn = document.querySelector('.menu-btn');
        if (menuBtn) {
            // 触发菜单关闭逻辑
            menuBtn.click();
        }
    }
});

// ========== 工具函数 ==========

// 防抖函数
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 节流函数
function throttle(func, limit) {
    let inThrottle;
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

// 随机数生成
function random(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

// ========== 用户菜单 ==========
function initUserMenu() {
    const userMenuBtn = document.getElementById('user-menu-btn');
    const userDropdown = document.getElementById('user-dropdown');
    const usernameDisplay = document.getElementById('username-display');
    const notLoggedIn = document.getElementById('not-logged-in');
    const loggedIn = document.getElementById('logged-in');
    const userName = document.getElementById('user-name');
    const logoutBtn = document.getElementById('logout-btn');
    const postBtn = document.getElementById('post-btn');
    
    // 如果关键元素不存在，直接返回
    if (!userMenuBtn || !usernameDisplay) {
        return;
    }
    
    // 发帖按钮点击验证
    if (postBtn) {
        postBtn.addEventListener('click', (e) => {
            if (!isLoggedIn()) {
                e.preventDefault();
                showMessage('请先登录后再发帖', 'error');
                setTimeout(() => {
                    window.location.href = 'login.html';
                }, 1000);
            }
        });
    }

    // 更新用户UI状态
    function updateUserUI() {
        if (isLoggedIn()) {
            const username = getCurrentUsername();
            usernameDisplay.textContent = username;
            if (userName) userName.textContent = username;
            if (notLoggedIn) notLoggedIn.style.display = 'none';
            if (loggedIn) loggedIn.style.display = 'block';
        } else {
            usernameDisplay.textContent = '登录';
            if (notLoggedIn) notLoggedIn.style.display = 'block';
            if (loggedIn) loggedIn.style.display = 'none';
        }
    }
    
    // 检查登录状态并更新UI
    updateUserUI();

    // 切换下拉菜单
    userMenuBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        userDropdown.classList.toggle('active');
    });

    // 点击页面其他地方关闭菜单
    document.addEventListener('click', () => {
        userDropdown.classList.remove('active');
    });

    // 阻止下拉菜单内的点击事件冒泡
    userDropdown.addEventListener('click', (e) => {
        e.stopPropagation();
    });

    // 退出登录
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            showConfirm('确定要退出登录吗？', () => {
                logout();
            });
        });
    }
}


// ========== 排序功能 ==========
function initSorting() {
    const sortButtons = document.querySelectorAll('.sort-btn');

    sortButtons.forEach(button => {
        button.addEventListener('click', async () => {
            // 更新按钮状态
            sortButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            const sortType = button.dataset.sort;
            console.log('排序方式:', sortType);

            // 重置到第一页
            currentPage = 1;
            currentSortType = sortType;
            currentSearchKeyword = ''; // 清除搜索状态
            
            // 恢复标题为正常状态
            updateSectionTitle('热门话题');

            // 调用 API 重新获取话题列表
            await loadTopics(sortType, currentPage, pageSize);
        });
    });
    
    // 页面加载时默认加载热门话题
    currentSortType = 'hot';
    currentSearchKeyword = '';
    loadTopics('hot', 1, pageSize);
}

// 加载话题列表
async function loadTopics(sortType = 'hot', page = 1, pageSize = 8) {
    const container = document.querySelector('.topics-grid');
    
    // 显示加载状态
    if (page === 1) {
        container.innerHTML = '<div class="loading" style="grid-column: 1/-1; text-align: center; padding: 60px 20px; color: #999;">加载中...</div>';
    }
    
    try {
        const response = await getTopics({ sort: sortType, page, page_size: pageSize });
        const topics = response.data.topics;
        const total = response.data.total || 0;
        
        // 计算总页数
        totalPages = Math.ceil(total / pageSize);
        if (totalPages === 0) totalPages = 1;
        
        if (!topics || topics.length === 0) {
            container.innerHTML = '<div class="empty-state" style="grid-column: 1/-1; text-align: center; padding: 60px 20px; color: #999;">暂无话题</div>';
            updatePaginationUI();
            return;
        }
        
        // 更新标题为正常状态
        updateSectionTitle('热门话题');
        
        // 渲染话题列表 - 每次都清空
        renderTopics(topics, true);
        
        // 更新分页UI
        updatePaginationUI();
        
    } catch (error) {
        console.error('加载话题失败:', error);
        container.innerHTML = '<div class="error-message" style="grid-column: 1/-1; text-align: center; padding: 60px 20px; color: #f44336;">加载失败，请刷新重试</div>';
        showMessage('加载话题失败', 'error');
    }
}

// 执行搜索
async function performSearch(keyword) {
    const container = document.querySelector('.topics-grid');
    
    // 显示加载状态
    container.innerHTML = '<div class="loading" style="grid-column: 1/-1; text-align: center; padding: 60px 20px; color: #999;">搜索中...</div>';
    
    // 更新标题
    updateSectionTitle(`搜索结果: "${keyword}"`);
    
    // 清除排序按钮的选中状态
    document.querySelectorAll('.sort-btn').forEach(btn => btn.classList.remove('active'));
    
    // 重置分页状态
    currentPage = 1;
    currentSortType = 'search';
    currentSearchKeyword = keyword;
    
    try {
        const response = await searchTopics({
            keyword: keyword,
            page: currentPage,
            page_size: pageSize
        });
        
        const searchResult = response.data;
        const total = searchResult.total || 0;
        
        // 计算总页数
        totalPages = Math.ceil(total / pageSize);
        if (totalPages === 0) totalPages = 1;
        
        if (!searchResult || !searchResult.topics || searchResult.topics.length === 0) {
            container.innerHTML = `
                <div class="empty-state" style="grid-column: 1/-1; text-align: center; padding: 60px 20px;">
                    <div style="font-size: 64px; margin-bottom: 20px;">🤷‍♂️</div>
                    <div style="font-size: 20px; color: #3d0000; font-weight: 900; margin-bottom: 10px; font-family: Impact, sans-serif;">没找到匹配的话题</div>
                    <div style="font-size: 14px; color: #666; font-weight: 600;">换个关键词试试？</div>
                </div>
            `;
            updatePaginationUI();
            return;
        }
        
        // 显示搜索统计
        console.log(`找到 ${searchResult.total} 条结果，耗时 ${searchResult.took}ms`);
        showMessage(`找到 ${searchResult.total} 条结果`, 'success');
        
        // 渲染搜索结果
        renderTopics(searchResult.topics, true);
        
        // 更新分页UI
        updatePaginationUI();
        
    } catch (error) {
        console.error('搜索失败:', error);
        container.innerHTML = '<div class="error-message" style="grid-column: 1/-1; text-align: center; padding: 60px 20px; color: #f44336;">搜索失败，请稍后重试</div>';
        showMessage('搜索失败', 'error');
        updatePaginationUI();
    }
}

// 更新章节标题
function updateSectionTitle(title) {
    const titleElement = document.querySelector('.section-title .title-text');
    if (titleElement) {
        titleElement.textContent = title;
    }
}

// 渲染话题列表
function renderTopics(topics, clearFirst = true) {
    const container = document.querySelector('.topics-grid');
    
    if (clearFirst) {
        container.innerHTML = '';
    }
    
    topics.forEach(topic => {
        const card = createTopicCard(topic);
        container.appendChild(card);
    });
    
    // 重新初始化动画
    initAnimations();
    
    // 重新绑定交互事件
    document.querySelectorAll('.topic-card').forEach(card => {
        card.addEventListener('click', handleTopicCardClick);
    });
    
    // 重新绑定投票按钮
    bindVoteButtons();
}

// 绑定投票按钮事件
function bindVoteButtons() {
    const voteButtons = document.querySelectorAll('.vote-btn');
    
    voteButtons.forEach(button => {
        // 移除旧的事件监听器（如果有）
        button.replaceWith(button.cloneNode(true));
    });
    
    // 重新获取并绑定
    document.querySelectorAll('.vote-btn').forEach(button => {
        button.addEventListener('click', async (e) => {
            e.stopPropagation(); // 防止触发卡片点击
            
            if (!isLoggedIn()) {
                showMessage('请先登录', 'error');
                setTimeout(() => {
                    window.location.href = 'login.html';
                }, 1000);
                return;
            }

            const type = button.dataset.type;
            const topicId = button.dataset.topicId;

            try {
                if (type === 'like') {
                    await likeTopic(topicId);
                } else {
                    await dislikeTopic(topicId);
                }
                
                showMessage('投票成功', 'success');
                
                // 重新加载当前排序的话题列表
                const activeSort = document.querySelector('.sort-btn.active');
                const sortType = activeSort ? activeSort.dataset.sort : 'hot';
                await loadTopics(sortType);
                
            } catch (error) {
                showMessage(error.message || '投票失败', 'error');
            }
        });
    });
}

// 创建话题卡片
function createTopicCard(topic) {
    const article = document.createElement('article');
    article.className = 'topic-card grid-card';
    article.dataset.topicId = topic.id;
    
    // 分类标签映射
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
    
    const avatarEmojis = {
        'tech': '👤',
        'design': '🎨',
        'discuss': '📱',
        'share': '💡',
        'product': '🚀'
    };
    
    article.innerHTML = `
        <div class="card-content">
            <div class="topic-header">
                <div class="user-avatar">${avatarEmojis[topic.category] || '👤'}</div>
                <div class="user-info">
                    <h4 class="username">${escapeHtml(topic.username || '匿名用户')}</h4>
                    <span class="post-time">${formatTimeAgo(topic.created_at)}</span>
                </div>
            </div>
            <h3 class="topic-title">${escapeHtml(topic.title)}</h3>
            <p class="topic-excerpt">${escapeHtml(truncateText(topic.content, 120))}</p>
            <div class="topic-footer">
                <span class="tag ${tagClasses[topic.category]}">${tagNames[topic.category]}</span>
                <div class="topic-stats">
                    <button class="stat-item vote-btn" data-type="like" data-topic-id="${topic.id}">
                        👍 ${topic.like_count || 0}
                    </button>
                    <span class="stat-item">💬 ${topic.comment_count || 0}</span>
                    <span class="stat-item">👁️ ${topic.view_count || 0}</span>
                </div>
            </div>
        </div>
    `;
    
    return article;
}

// 格式化时间
function formatTimeAgo(timeStr) {
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

// 截断文本
function truncateText(text, maxLength) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// HTML转义
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ========== 分页功能 ==========
function initPagination() {
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    
    if (prevBtn) {
        prevBtn.addEventListener('click', async () => {
            if (currentPage > 1) {
                currentPage--;
                if (currentSortType === 'search') {
                    const response = await searchTopics({ keyword: currentSearchKeyword, page: currentPage, page_size: pageSize });
                    const container = document.querySelector('.topics-grid');
                    container.innerHTML = ''; // 清空
                    renderTopics(response.data.topics, true);
                    updatePaginationUI();
                } else {
                    await loadTopics(currentSortType, currentPage, pageSize);
                }
            }
        });
    }
    
    if (nextBtn) {
        nextBtn.addEventListener('click', async () => {
            if (currentPage < totalPages) {
                currentPage++;
                if (currentSortType === 'search') {
                    const response = await searchTopics({ keyword: currentSearchKeyword, page: currentPage, page_size: pageSize });
                    const container = document.querySelector('.topics-grid');
                    container.innerHTML = ''; // 清空
                    renderTopics(response.data.topics, true);
                    updatePaginationUI();
                } else {
                    await loadTopics(currentSortType, currentPage, pageSize);
                }
            }
        });
    }
}

function updatePaginationUI() {
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    const currentPageSpan = document.getElementById('current-page');
    const totalPagesSpan = document.getElementById('total-pages');
    
    // 更新页码显示
    if (currentPageSpan) {
        currentPageSpan.textContent = currentPage;
    }
    if (totalPagesSpan) {
        totalPagesSpan.textContent = totalPages;
    }
    
    // 更新按钮状态
    if (prevBtn) {
        prevBtn.disabled = currentPage <= 1;
    }
    if (nextBtn) {
        nextBtn.disabled = currentPage >= totalPages;
    }
}

// ========== 开发环境日志 ==========
console.log('%c🌸 Blossom Forum', 'font-size: 24px; font-weight: bold; color: #5B9FED;');
console.log('%c欢迎来到 Blossom 创意论坛 - 思想绽放的地方!', 'font-size: 14px; color: #666;');
console.log('%c按 Ctrl+K (Mac: Cmd+K) 可快速聚焦搜索框', 'font-size: 12px; color: #999;');



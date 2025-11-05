# 向日葵海报设计系统

> 复古风格的配色和字体方案，适用于海报、宣传页、品牌设计等场景

---

## 🎨 配色方案

### 主色调 - 红色系

```css
--primary-red: #b71c1c;    /* 主红色 */
--dark-red: #8b0000;       /* 深红色 */
--deep-red: #3d0000;       /* 深褐红色 */
```

| 颜色名称 | 十六进制 | RGB | 用途 |
|---------|---------|-----|------|
| 主红色 | `#b71c1c` | `rgb(183, 28, 28)` | 标题、横幅、装饰框、主要CTA按钮 |
| 深红色 | `#8b0000` | `rgb(139, 0, 0)` | 阴影、纹理、边框、悬停状态 |
| 深褐红色 | `#3d0000` | `rgb(61, 0, 0)` | 正文文字、小标题 |

### 背景色 - 米黄色系

```css
--bg-light: #f5ebe0;       /* 浅米色 */
--bg-medium: #ede0d4;      /* 中米色 */
```

| 颜色名称 | 十六进制 | RGB | 用途 |
|---------|---------|-----|------|
| 浅米色 | `#f5ebe0` | `rgb(245, 235, 224)` | 主背景、卡片背景、浅色区域 |
| 中米色 | `#ede0d4` | `rgb(237, 224, 212)` | 渐变背景、分隔区域 |

### 文字颜色

```css
--text-light: #f5ebe0;     /* 浅色文字 */
--text-dark: #3d0000;      /* 深色文字 */
```

### 装饰色

```css
--white-transparent: rgba(255, 255, 255, 0.3);  /* 半透明白色 */
```

---

## ✍️ 字体方案

### 标题字体

**Impact**
- 用途：所有大标题、重点文字、数字
- 特点：粗体、醒目、压缩字体
- 替代字体：Anton, Bebas Neue, Oswald

```css
font-family: 'Impact', sans-serif;
font-weight: 900;
letter-spacing: -5px; /* 紧凑效果 */
```

### 正文字体

**Arial**
- 用途：正文、描述文字、副标题
- 特点：清晰、易读、无衬线
- 替代字体：Helvetica, Roboto, 'Microsoft YaHei'

```css
font-family: 'Arial', sans-serif;
font-weight: 400 | 700;
```

---

## 📐 使用示例

### CSS 变量定义

```css
:root {
    /* 主色调 */
    --primary-red: #b71c1c;
    --dark-red: #8b0000;
    --deep-red: #3d0000;
    
    /* 背景色 */
    --bg-light: #f5ebe0;
    --bg-medium: #ede0d4;
    
    /* 文字颜色 */
    --text-light: #f5ebe0;
    --text-dark: #3d0000;
    
    /* 装饰色 */
    --white-transparent: rgba(255, 255, 255, 0.3);
}
```

### 渐变背景

```css
/* 温暖的米色渐变 */
background: linear-gradient(to bottom, var(--bg-light), var(--bg-medium));

/* 也可以使用固定颜色 */
background: linear-gradient(to bottom, #f5ebe0, #ede0d4);
```

### 标题样式

```css
.title {
    font-family: 'Impact', sans-serif;
    font-size: clamp(60px, 10vw, 140px);
    color: var(--primary-red);
    font-weight: 900;
    letter-spacing: -5px;
    text-transform: uppercase;
}
```

### 按钮样式

```css
.btn-primary {
    background: var(--primary-red);
    color: var(--text-light);
    font-family: 'Impact', sans-serif;
    padding: 15px 40px;
    border-radius: 50px;
    border: none;
    font-size: 24px;
    letter-spacing: 2px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.btn-primary:hover {
    background: var(--dark-red);
    transform: scale(1.05);
}
```

### 卡片样式

```css
.card {
    background: var(--bg-light);
    border: 3px solid var(--primary-red);
    border-radius: 20px;
    padding: 30px;
    color: var(--text-dark);
}
```

---

## 🎯 设计原则

1. **对比度高** - 深红色与米黄色形成强烈对比，确保可读性
2. **复古感** - 模仿20世纪中期的海报设计风格
3. **温暖色调** - 米黄色背景营造温暖、亲切的氛围
4. **醒目大胆** - 使用Impact字体和大尺寸文字吸引注意力

---

## 🖼️ 配色板预览

```
主红色系：
■ #b71c1c  ■ #8b0000  ■ #3d0000

背景米色系：
□ #f5ebe0  □ #ede0d4

文字色：
□ #f5ebe0 (浅)  ■ #3d0000 (深)
```

---

## 📦 快速导入

### Google Fonts 引入

```html
<!-- 在 HTML head 中引入 -->
<link href="https://fonts.googleapis.com/css2?family=Impact&display=swap" rel="stylesheet">
```

### CSS 导入

```css
@import url('https://fonts.googleapis.com/css2?family=Impact&family=Arial:wght@400;700&display=swap');
```

---

## 💡 适用场景

- ✅ 海报设计
- ✅ 宣传页面
- ✅ 品牌视觉
- ✅ 活动页面
- ✅ 复古风格网站
- ✅ 印刷品设计

---

## 🔗 相关资源

- [Adobe Color - 配色生成器](https://color.adobe.com)
- [Coolors - 配色方案](https://coolors.co)
- [Google Fonts](https://fonts.google.com)

---

> **创建日期：** 2025-11-04  
> **风格：** 复古海报风格  
> **灵感来源：** 向日葵主题海报


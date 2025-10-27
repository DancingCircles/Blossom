# PowerShell 脚本：自动修正 go.mod 版本号
# 使用方法：在提交前运行 .\scripts\fix-gomod.ps1

$goModFile = "web_app\go.mod"

Write-Host "🔍 检查 go.mod 版本..." -ForegroundColor Cyan

$content = Get-Content $goModFile -Raw

if ($content -match "go 1\.24|go 1\.22") {
    Write-Host "⚠️  检测到版本需要更新，修正为 1.23..." -ForegroundColor Yellow
    
    $content = $content -replace "go 1\.24\.\d+", "go 1.23"
    $content = $content -replace "go 1\.24", "go 1.23"
    $content = $content -replace "go 1\.22", "go 1.23"
    
    Set-Content -Path $goModFile -Value $content -NoNewline
    
    Write-Host "✅ go.mod 已修正为 1.23" -ForegroundColor Green
    
    # 运行 go mod tidy
    Set-Location web_app
    go mod tidy
    Set-Location ..
    
    Write-Host "✅ 依赖已更新" -ForegroundColor Green
} else {
    Write-Host "✅ go.mod 版本正确 (1.23)" -ForegroundColor Green
}

Write-Host ""
Write-Host "💡 提示：如果 go.mod 经常自动变化，建议：" -ForegroundColor Cyan
Write-Host "   1. 在编辑器设置中禁用 go.mod 的自动格式化" -ForegroundColor White
Write-Host "   2. 提交前运行此脚本" -ForegroundColor White


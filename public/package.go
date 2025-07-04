package public

import "embed"

//go:embed dist/*
var Public embed.FS

// 打包嵌入前端资源

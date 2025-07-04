package public

import "embed"

//go:embed dist/*
var Public embed.FS

// 暂时忽略，编译时需要将将所有前端资源打包到程序中

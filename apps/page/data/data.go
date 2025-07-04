package data

import "embed"

//go:embed menu-admin.json
var MenuAdmin embed.FS

//go:embed menu-user.json
var MenuUser embed.FS

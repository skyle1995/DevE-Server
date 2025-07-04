package page

import (
	"github.com/skyle1995/DevE-Server/apps/page/data"
	"github.com/skyle1995/DevE-Server/apps/page/model"
)

// Service 页面服务结构体
type Service struct{}

// NewService 创建页面服务实例
func NewService() *Service {
	return &Service{}
}

// GetMenu 根据用户角色获取菜单数据
func (s *Service) GetMenu(role int) (*model.MenuResponse, error) {
	var menuData []byte
	var err error

	// 根据角色返回不同的菜单数据
	if role == 0 { // 管理员
		menuData, err = data.MenuAdmin.ReadFile("menu-admin.json")
	} else { // 普通会员或VIP会员
		menuData, err = data.MenuUser.ReadFile("menu-user.json")
	}

	if err != nil {
		return nil, err
	}

	return &model.MenuResponse{
		Menu: menuData,
	}, nil
}

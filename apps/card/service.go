package card

import (
	"errors"

	"github.com/skyle1995/DevE-Server/apps/card/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/random"
	"gorm.io/gorm"
)

// Service 提供卡密相关的服务
type Service struct{}

// NewService 创建一个新的卡密服务实例
func NewService() *Service {
	return &Service{}
}

// GenerateCards 生成卡密
// @param req 生成卡密请求
// @param userID 当前用户ID
// @return 生成的卡密列表和错误信息
func (s *Service) GenerateCards(req model.GenerateCardRequest, userID int) ([]dbmodel.Card, error) {
	// 查询卡密类型
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", req.TypeID, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密类型不存在或无权限使用")
		}
		return nil, errors.New("查询卡密类型失败: " + result.Error.Error())
	}

	// 生成卡密
	cards := make([]dbmodel.Card, 0, req.Count)
	for i := 0; i < req.Count; i++ {
		// 生成卡号
		cardNo := ""
		if req.Prefix != "" {
			cardNo = req.Prefix + "-"
		}
		cardNo += random.GenerateRandomString(10)

		// 确定卡密长度，默认为16，最低限制为8
		keyLength := 16
		if req.KeyLength != nil && *req.KeyLength >= 8 {
			keyLength = *req.KeyLength
		} else if req.KeyLength != nil && *req.KeyLength < 8 {
			// 如果指定的长度小于8，则使用最低限制8
			keyLength = 8
		}

		// 生成卡密
		cardKey := random.GenerateRandomString(keyLength)

		// 创建卡密
		card := dbmodel.Card{
			CardNo:  cardNo,
			CardKey: cardKey,
			TypeID:  uint(req.TypeID),
			AppID:   cardType.AppID,
			Status:  0, // 0-未使用
			UserID:  userID,
		}

		// 设置最大换绑/解绑次数
		if req.MaxRebindCount != nil {
			card.MaxRebindCount = *req.MaxRebindCount
		} else {
			card.MaxRebindCount = cardType.DefaultMaxRebindCount
		}

		if req.MaxUnbindCount != nil {
			card.MaxUnbindCount = *req.MaxUnbindCount
		} else {
			card.MaxUnbindCount = cardType.DefaultMaxUnbindCount
		}

		// 添加到列表
		cards = append(cards, card)
	}

	// 批量保存卡密
	result = database.DB.Create(&cards)
	if result.Error != nil {
		return nil, errors.New("保存卡密失败: " + result.Error.Error())
	}

	return cards, nil
}

// GetCardList 获取卡密列表
// @param req 获取卡密列表请求
// @param userID 当前用户ID
// @return 卡密列表、总数和错误信息
func (s *Service) GetCardList(req model.GetCardListRequest, userID int) ([]dbmodel.Card, int64, error) {
	// 构建查询条件
	query := database.DB.Model(&dbmodel.Card{}).Where("user_id = ?", userID)

	// 应用筛选条件
	if req.AppID > 0 {
		query = query.Where("app_id = ?", req.AppID)
	}

	if req.TypeID > 0 {
		query = query.Where("type_id = ?", req.TypeID)
	}

	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	if req.CardNo != "" {
		query = query.Where("card_no LIKE ?", "%"+req.CardNo+"%")
	}

	// 获取总数
	var total int64
	result := query.Count(&total)
	if result.Error != nil {
		return nil, 0, errors.New("获取卡密总数失败: " + result.Error.Error())
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 排序
	query = query.Order("created_at DESC")

	// 查询卡密列表
	var cards []dbmodel.Card
	result = query.Find(&cards)
	if result.Error != nil {
		return nil, 0, errors.New("获取卡密列表失败: " + result.Error.Error())
	}

	return cards, total, nil
}

// UpdateCard 更新卡密
// @param req 更新卡密请求
// @param userID 当前用户ID
// @return 更新后的卡密和错误信息
func (s *Service) UpdateCard(req model.UpdateCardRequest, userID int) (*dbmodel.Card, error) {
	// 查询卡密
	var card dbmodel.Card
	result := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密不存在或无权限修改")
		}
		return nil, errors.New("查询卡密失败: " + result.Error.Error())
	}

	// 更新卡密信息
	if req.TypeID > 0 {
		// 查询卡密类型
		var cardType dbmodel.CardType
		result = database.DB.Where("id = ? AND user_id = ?", req.TypeID, userID).First(&cardType)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("卡密类型不存在或无权限使用")
			}
			return nil, errors.New("查询卡密类型失败: " + result.Error.Error())
		}

		card.TypeID = uint(req.TypeID)
		card.AppID = cardType.AppID
	}

	if req.Status >= 0 {
		card.Status = req.Status
	}

	if req.MaxRebindCount != nil {
		card.MaxRebindCount = *req.MaxRebindCount
	}

	if req.MaxUnbindCount != nil {
		card.MaxUnbindCount = *req.MaxUnbindCount
	}

	// 保存卡密
	result = database.DB.Save(&card)
	if result.Error != nil {
		return nil, errors.New("更新卡密失败: " + result.Error.Error())
	}

	return &card, nil
}

// DeleteCard 删除卡密
// @param req 删除卡密请求
// @param userID 当前用户ID
// @return 错误信息
func (s *Service) DeleteCard(req model.DeleteCardRequest, userID int) error {
	// 查询卡密
	var card dbmodel.Card
	result := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("卡密不存在或无权限删除")
		}
		return errors.New("查询卡密失败: " + result.Error.Error())
	}

	// 删除卡密
	result = database.DB.Delete(&card)
	if result.Error != nil {
		return errors.New("删除卡密失败: " + result.Error.Error())
	}

	return nil
}

// CreateCardType 创建卡密类型
// @param req 创建卡密类型请求
// @param userID 当前用户ID
// @return 创建的卡密类型和错误信息
func (s *Service) CreateCardType(req model.CreateCardTypeRequest, userID int) (*dbmodel.CardType, error) {
	// 检查应用是否存在
	var app dbmodel.App
	result := database.DB.Where("id = ? AND user_id = ?", req.AppID, userID).First(&app)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("应用不存在或无权限使用")
		}
		return nil, errors.New("查询应用失败: " + result.Error.Error())
	}

	// 创建卡密类型
	cardType := dbmodel.CardType{
		Name:                  req.Name,
		Duration:              req.Duration,
		TimeUnit:              req.TimeUnit,
		AppID:                 uint(req.AppID),
		Status:                req.Status,
		DefaultMaxRebindCount: req.DefaultMaxRebindCount,
		DefaultMaxUnbindCount: req.DefaultMaxUnbindCount,
		UserID:                userID,
	}

	// 保存卡密类型
	result = database.DB.Create(&cardType)
	if result.Error != nil {
		return nil, errors.New("创建卡密类型失败: " + result.Error.Error())
	}

	return &cardType, nil
}

// GetCardTypeList 获取卡密类型列表
// @param req 获取卡密类型列表请求
// @param userID 当前用户ID
// @return 卡密类型列表、总数和错误信息
func (s *Service) GetCardTypeList(req model.GetCardTypeListRequest, userID int) ([]dbmodel.CardType, int64, error) {
	// 构建查询条件
	query := database.DB.Model(&dbmodel.CardType{}).Where("user_id = ?", userID)

	// 应用筛选条件
	if req.AppID > 0 {
		query = query.Where("app_id = ?", req.AppID)
	}

	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	var total int64
	result := query.Count(&total)
	if result.Error != nil {
		return nil, 0, errors.New("获取卡密类型总数失败: " + result.Error.Error())
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 排序
	query = query.Order("created_at DESC")

	// 查询卡密类型列表
	var cardTypes []dbmodel.CardType
	result = query.Find(&cardTypes)
	if result.Error != nil {
		return nil, 0, errors.New("获取卡密类型列表失败: " + result.Error.Error())
	}

	return cardTypes, total, nil
}

// UpdateCardType 更新卡密类型
// @param req 更新卡密类型请求
// @param userID 当前用户ID
// @return 更新后的卡密类型和错误信息
func (s *Service) UpdateCardType(req model.UpdateCardTypeRequest, userID int) (*dbmodel.CardType, error) {
	// 查询卡密类型
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密类型不存在或无权限修改")
		}
		return nil, errors.New("查询卡密类型失败: " + result.Error.Error())
	}

	// 更新卡密类型信息
	if req.Name != "" {
		cardType.Name = req.Name
	}

	if req.Duration > 0 {
		cardType.Duration = req.Duration
	}

	if req.TimeUnit != "" {
		cardType.TimeUnit = req.TimeUnit
	}

	if req.AppID > 0 {
		// 检查应用是否存在
		var app dbmodel.App
		result = database.DB.Where("id = ? AND user_id = ?", req.AppID, userID).First(&app)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("应用不存在或无权限使用")
			}
			return nil, errors.New("查询应用失败: " + result.Error.Error())
		}

		cardType.AppID = uint(req.AppID)
	}

	if req.Status >= 0 {
		cardType.Status = req.Status
	}

	if req.DefaultMaxRebindCount >= 0 {
		cardType.DefaultMaxRebindCount = req.DefaultMaxRebindCount
	}

	if req.DefaultMaxUnbindCount >= 0 {
		cardType.DefaultMaxUnbindCount = req.DefaultMaxUnbindCount
	}

	// 保存卡密类型
	result = database.DB.Save(&cardType)
	if result.Error != nil {
		return nil, errors.New("更新卡密类型失败: " + result.Error.Error())
	}

	return &cardType, nil
}

// DeleteCardType 删除卡密类型
// @param req 删除卡密类型请求
// @param userID 当前用户ID
// @return 错误信息
func (s *Service) DeleteCardType(req model.DeleteCardTypeRequest, userID int) error {
	// 查询卡密类型
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", req.ID, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("卡密类型不存在或无权限删除")
		}
		return errors.New("查询卡密类型失败: " + result.Error.Error())
	}

	// 检查是否有关联的卡密
	var count int64
	result = database.DB.Model(&dbmodel.Card{}).Where("type_id = ?", req.ID).Count(&count)
	if result.Error != nil {
		return errors.New("检查关联卡密失败: " + result.Error.Error())
	}

	if count > 0 {
		return errors.New("该卡密类型下存在卡密，无法删除")
	}

	// 删除卡密类型
	result = database.DB.Delete(&cardType)
	if result.Error != nil {
		return errors.New("删除卡密类型失败: " + result.Error.Error())
	}

	return nil
}

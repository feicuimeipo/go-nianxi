package dao

import (
	"errors"
	"fmt"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/dto"
	"gitee.com/go-nianxi/go-admin/internal/admin/domain/vo"
	"gitee.com/go-nianxi/go-admin/internal/admin/model"
	"github.com/casbin/casbin/v2"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"strings"
)

type ApiDao struct {
	db             *gorm.DB
	casbinEnforcer *casbin.Enforcer
}

func NewApiDao(db *gorm.DB, casbinEnforcer *casbin.Enforcer) *ApiDao {
	return &ApiDao{
		db:             db,
		casbinEnforcer: casbinEnforcer,
	}
}

// 获取接口列表
func (a ApiDao) GetApis(req *vo.ApiListRequest) ([]*model.Api, int64, error) {
	var list []*model.Api
	db := a.db.Model(&model.Api{}).Order("created_at DESC")

	db.Where("application_id = ?", req.ApplicationId)

	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 根据接口ID获取接口列表
func (a ApiDao) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := a.db.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

// 获取接口树(按接口Category字段分类)
func (a ApiDao) GetApiTree(applicationId uint) ([]*dto.ApiTreeDto, error) {
	var apiList []*model.Api
	err := a.db.Where("application_id=?", applicationId).Order("category").Order("created_at").Find(&apiList).Error
	// 获取所有的分类
	var categoryList []string
	for _, api := range apiList {
		categoryList = append(categoryList, api.Category)
	}
	// 获取去重后的分类
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*dto.ApiTreeDto, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &dto.ApiTreeDto{
			ID:       -i,
			Desc:     category,
			Category: category,
			Children: nil,
		}
		for _, api := range apiList {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, api)
			}
		}
	}

	return apiTree, err
}

// 创建接口
func (a ApiDao) CreateApi(api *model.Api) error {
	
	err := a.db.Create(api).Error
	return err
}

// 更新接口
func (a ApiDao) UpdateApiById(apiId uint, api *model.Api) error {
	// 根据id获取接口信息
	var oldApi model.Api
	err := a.db.First(&oldApi, apiId).Error
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}
	api.ApplicationId = oldApi.ApplicationId

	err = a.db.Model(api).Where("id = ?", apiId).Updates(api).Error
	if err != nil {
		return err
	}
	// 更新了method和path就更新casbin中policy
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies := a.casbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// 接口在casbin的policy中存在才进行操作
		if len(policies) > 0 {
			// 先删除
			isRemoved, _ := a.casbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("更新权限接口失败")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			// 新增
			isAdded, _ := a.casbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("更新权限接口失败")
			}
			// 加载policy
			err := a.casbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("更新权限接口成功，权限接口策略加载失败")
			} else {
				return err
			}
		}
	}
	return err
}

// 批量删除接口
func (a ApiDao) BatchDeleteApiByIds(apiIds []uint) error {

	apis, err := a.GetApisById(apiIds)
	if err != nil {
		return errors.New("根据接口ID获取接口列表失败")
	}
	if len(apis) == 0 {
		return errors.New("根据接口ID未获取到接口列表")
	}

	err = a.db.Where("id IN (?)", apiIds).Unscoped().Delete(&model.Api{}).Error
	// 如果删除成功，删除casbin中policy
	if err == nil {
		for _, api := range apis {
			policies := a.casbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := a.casbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("删除权限接口失败")
				}
			}
		}
		// 重新加载策略
		err := a.casbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("删除权限接口成功，权限接口策略加载失败")
		} else {
			return err
		}
	}
	return err
}

// 根据接口路径和请求方式获取接口描述
func (a ApiDao) GetApiDescByPath(path string, method string) (string, error) {
	var api model.Api
	err := a.db.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}

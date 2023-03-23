package rules

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Operator interface {
	GetRules(filter, sortField, sortOrder string, pageNo, pageSize int) ([]Rule, error)

	CreateRule(rule *Rule) error

	UpdateRule(rule *Rule) error

	DeleteRuleByID(id int) error

	GetTotalPages(pageSize int) int
}

type operator struct {
	db *gorm.DB
}

func NewOperator(db *gorm.DB) Operator {
	return &operator{
		db: db.Session(&gorm.Session{}),
	}
}

func (o *operator) GetRules(filter, sortField, sortOrder string, pageNo, pageSize int) ([]Rule, error) {
	var query = o.db
	if filter != "" {
		query = query.Where("name LIKE ?", "%"+filter+"%")
	}

	query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))

	offset := (pageNo - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	var rules []Rule
	err := query.Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (o *operator) CreateRule(rule *Rule) error {
	rule.CreateTime = time.Now().Unix()
	rule.UpdateTime = time.Now().Unix()

	if err := o.db.Create(rule).Error; err != nil {
		return err
	}

	return nil
}

func (o *operator) UpdateRule(rule *Rule) error {
	result := o.db.Model(&rule).Updates(rule)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errors.New("nothing changed")
	}
	o.db.Model(&rule).Update("update_time", time.Now().Unix())

	return nil
}

func (o *operator) DeleteRuleByID(id int) error {
	result := o.db.Delete(&Rule{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (o *operator) GetTotalPages(pageSize int) int {
	var tableCount int64
	o.db.Count(&tableCount)

	pageCount := tableCount / int64(pageSize)
	if tableCount%int64(pageSize) != 0 && tableCount == 0 {
		pageCount++
	}
	return int(pageCount)
}

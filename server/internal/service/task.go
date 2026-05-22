package service

import (
	"context"
	"fmt"

	"pxx/internal/cache"
	"pxx/internal/model"
	"pxx/pkg/utils"

	"gorm.io/gorm"
)

type TaskService struct {
	db    *gorm.DB
	cache *cache.CacheManager
}

func NewTaskService(db *gorm.DB, cm *cache.CacheManager) *TaskService {
	return &TaskService{db: db, cache: cm}
}

func (s *TaskService) List(ctx context.Context, taskType string, page, pageSize int) ([]model.Task, int64, error) {
	offset, limit := utils.Paginate(page, pageSize)
	var total int64
	var items []model.Task
	query := s.db.WithContext(ctx).Model(&model.Task{}).Where("status = ?", model.TaskStatusOpen)
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	query.Count(&total).Order("created_at DESC").Offset(offset).Limit(limit).Preload("Publisher").Find(&items)
	return items, total, nil
}

func (s *TaskService) Create(ctx context.Context, task *model.Task) error {
	task.Status = model.TaskStatusOpen
	return s.db.WithContext(ctx).Create(task).Error
}

func (s *TaskService) Take(ctx context.Context, taskID, takerID uint) error {
	return s.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ? AND status = ?", taskID, model.TaskStatusOpen).
		Updates(map[string]interface{}{"status": model.TaskStatusTaken, "taker_id": takerID}).Error
}

func (s *TaskService) Done(ctx context.Context, taskID, userID uint) error {
	result := s.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ? AND status = ? AND (publisher_id = ? OR taker_id = ?)", taskID, model.TaskStatusTaken, userID, userID).
		Update("status", model.TaskStatusDone)
	if result.RowsAffected == 0 {
		return fmt.Errorf("操作失败")
	}
	return result.Error
}

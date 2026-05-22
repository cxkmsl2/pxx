package service

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"pxx/internal/model"
	"pxx/internal/ws"

	"gorm.io/gorm"
)

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) Send(ctx context.Context, fromID, toID uint, content, msgType string) (*model.Message, error) {
	msg := &model.Message{
		FromUserID: fromID, ToUserID: toID, Content: content, MsgType: msgType,
	}
	if err := s.db.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, err
	}
	pushData, _ := json.Marshal(map[string]interface{}{
		"type": "new_message", "from_id": fromID, "content": content,
	})
	ws.DefaultHub.SendToUser(toID, pushData)
	return msg, nil
}

func (s *MessageService) GetMessages(ctx context.Context, user1, user2 uint, limit int) ([]model.Message, error) {
	var msgs []model.Message
	err := s.db.WithContext(ctx).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", user1, user2, user2, user1).
		Order("created_at ASC").Limit(limit).Find(&msgs).Error
	return msgs, err
}

type ConversationResult struct {
	UserID   uint      `json:"user_id"`
	Nickname string    `json:"nickname"`
	LastMsg  string    `json:"last_msg"`
	LastTime time.Time `json:"last_time"`
	Unread   int       `json:"unread"`
}

func (s *MessageService) GetSystemMessages(ctx context.Context, userID uint) ([]model.Message, error) {
	var msgs []model.Message
	s.db.WithContext(ctx).Where("to_user_id = ? AND from_user_id = 0", userID).Order("created_at DESC").Limit(50).Find(&msgs)
	return msgs, nil
}

func (s *MessageService) MarkRead(ctx context.Context, userID, fromUserID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.Message{}).Where("to_user_id = ? AND from_user_id = ? AND read_at IS NULL", userID, fromUserID).Update("read_at", &now).Error
}




func (s *MessageService) GetConversations(ctx context.Context, userID uint) ([]ConversationResult, error) {
	var msgs []model.Message
	s.db.WithContext(ctx).
		Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").Limit(200).Find(&msgs)

	// Group by peer in memory
	type peerData struct {
		userID   uint
		nickname string
		lastMsg  string
		lastTime time.Time
		unread   int
	}
	peerMap := make(map[uint]*peerData)
	peerOrder := []uint{}

	for _, m := range msgs {
		var peerID uint
		if m.FromUserID == userID {
			peerID = m.ToUserID
		} else {
			peerID = m.FromUserID
		}
		if _, exists := peerMap[peerID]; !exists {
			peerMap[peerID] = &peerData{userID: peerID}
			peerOrder = append(peerOrder, peerID)
		}
		pd := peerMap[peerID]
		if m.CreatedAt.After(pd.lastTime) {
			pd.lastMsg = m.Content
			pd.lastTime = m.CreatedAt
		}
		if m.ToUserID == userID && m.ReadAt == nil {
			pd.unread++
		}
	}

	// Fetch nicknames
	result := make([]ConversationResult, 0, len(peerOrder))
	for _, pid := range peerOrder {
		pd := peerMap[pid]
		var user model.User
		s.db.WithContext(ctx).Select("nickname").First(&user, pid)
		result = append(result, ConversationResult{
			UserID: pd.userID, Nickname: user.Nickname,
			LastMsg: pd.lastMsg, LastTime: pd.lastTime, Unread: pd.unread,
		})
	}
	// Sort by lastTime descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastTime.After(result[j].LastTime)
	})
	return result, nil
}

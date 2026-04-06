package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/malcolmkzh/study-notifier/internal/modules/telegram/dto"
	telegrammodel "github.com/malcolmkzh/study-notifier/internal/modules/telegram/model"
	telegramrepository "github.com/malcolmkzh/study-notifier/internal/modules/telegram/repository"
	userrepository "github.com/malcolmkzh/study-notifier/internal/modules/user/repository"
	"github.com/malcolmkzh/study-notifier/internal/utilities/notification"
)

const (
	linkCodeLength = 6
	linkTTL        = 10 * time.Minute
)

const linkCodeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type Implementation struct {
	telegramLinkRepo telegramrepository.Utility
	userRepo         userrepository.Utility
	notification     notification.Utility
	now              func() time.Time
}

func NewService(
	telegramLinkRepo telegramrepository.Utility,
	userRepo userrepository.Utility,
	notificationUtility notification.Utility,
) (*Implementation, error) {
	if telegramLinkRepo == nil {
		return nil, errors.New("telegram link repository is required")
	}
	if userRepo == nil {
		return nil, errors.New("user repository is required")
	}
	if notificationUtility == nil {
		return nil, errors.New("notification utility is required")
	}

	return &Implementation{
		telegramLinkRepo: telegramLinkRepo,
		userRepo:         userRepo,
		notification:     notificationUtility,
		now:              time.Now,
	}, nil
}

func (s *Implementation) CreateLink(ctx context.Context, userID string) (*dto.CreateLinkResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if err := s.telegramLinkRepo.DeleteByUserID(ctx, userID); err != nil {
		return nil, err
	}

	code, err := randomString(linkCodeLength)
	if err != nil {
		return nil, err
	}

	link := telegrammodel.TelegramLink{
		UserID:    userID,
		Code:      code,
		ExpiresAt: s.now().UTC().Add(linkTTL),
	}

	if err := s.telegramLinkRepo.Create(ctx, &link); err != nil {
		return nil, err
	}

	return &dto.CreateLinkResponse{Code: code}, nil
}

func (s *Implementation) HandleMessage(ctx context.Context, chatID int64, text string) error {
	text = strings.TrimSpace(text)

	switch {
	case text == "/start":
		return s.notification.SendTelegramMessage(ctx, notification.SendTelegramMessageRequest{
			ChatID: strconv.FormatInt(chatID, 10),
			Text:   "Use /link <CODE> to connect your Telegram account.",
		})
	case strings.HasPrefix(text, "/link "):
		return s.handleLink(ctx, chatID, strings.TrimSpace(strings.TrimPrefix(text, "/link ")))
	default:
		return nil
	}
}

func (s *Implementation) handleLink(ctx context.Context, chatID int64, code string) error {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return s.notification.SendTelegramMessage(ctx, notification.SendTelegramMessageRequest{
			ChatID: strconv.FormatInt(chatID, 10),
			Text:   "Invalid code",
		})
	}

	link, err := s.telegramLinkRepo.GetByCode(ctx, code)
	if err != nil {
		return err
	}
	if link == nil {
		return s.notification.SendTelegramMessage(ctx, notification.SendTelegramMessageRequest{
			ChatID: strconv.FormatInt(chatID, 10),
			Text:   "Invalid code",
		})
	}

	if s.now().UTC().After(link.ExpiresAt.UTC()) {
		_ = s.telegramLinkRepo.DeleteByCode(ctx, code)
		return s.notification.SendTelegramMessage(ctx, notification.SendTelegramMessageRequest{
			ChatID: strconv.FormatInt(chatID, 10),
			Text:   "Code expired",
		})
	}

	if err := s.userRepo.UpdateTelegramChatID(ctx, link.UserID, strconv.FormatInt(chatID, 10)); err != nil {
		return err
	}

	if err := s.telegramLinkRepo.DeleteByCode(ctx, code); err != nil {
		return err
	}

	return s.notification.SendTelegramMessage(ctx, notification.SendTelegramMessageRequest{
		ChatID: strconv.FormatInt(chatID, 10),
		Text:   "Successfully linked your account!",
	})
}

func randomString(length int) (string, error) {
	var builder strings.Builder
	builder.Grow(length)

	limit := big.NewInt(int64(len(linkCodeCharset)))
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", fmt.Errorf("generate link code: %w", err)
		}

		builder.WriteByte(linkCodeCharset[index.Int64()])
	}

	return builder.String(), nil
}

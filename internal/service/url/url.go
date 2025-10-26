package url

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
)

type URLRepo interface {
	SaveURL(ctx context.Context, url, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
}

type Service struct {
	log  *slog.Logger
	repo URLRepo
}

func NewService(log *slog.Logger, repo URLRepo) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, url string, length *int) (string, error) {

	if length == nil {
		length = new(int)
		*length = 8
	}
	alias := newRandomString(*length)

	s.log.Info("alias creation attempt", slog.String("url", url))

	err := s.repo.SaveURL(ctx, url, alias)
	if err != nil {
		s.log.Error(
			"alias creation attempt failed",
			slog.String("url", url),
			slog.String("err", err.Error()),
		)

		return "", fmt.Errorf("failed to create an alias")
	}

	s.log.Info("alias creation attempt successful", slog.String("url", url), slog.String("alias", alias))

	return alias, nil
}

func (s *Service) Get(ctx context.Context, alias string) (string, error) {

	s.log.Info("url retrieval attempt", slog.String("alias", alias))

	url, err := s.repo.GetURL(ctx, alias)
	if err != nil {
		s.log.Error(
			"alias creation attempt failed",
			slog.String("alias", alias),
			slog.String("err", err.Error()),
		)
		return "", fmt.Errorf("failed to get an url by the alias")
	}

	s.log.Info("url retrieval attempt successful", slog.String("alias", alias))

	return url, nil
}

func newRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}

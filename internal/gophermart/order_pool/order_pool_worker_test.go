package orderpool

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func TestPositiveOrderPoolWorkers(t *testing.T) {
	srv := NewAccrualServer(t)
	srv.SetSuccessRequest(true)

	require.NoError(t, srv.Run(), "Не удалось запустить тестовый сервер")
	defer srv.Close()

	pool := NewTestOrderPool(t, srv.GetAddress())

	pool.Get().Run()
	defer pool.Get().Close()

	pool.Get().AddJob(models.Order{
		ID:         123,
		UserID:     uuid.New(),
		Status:     "NEW",
		Accrual:    nil,
		UploadedAt: time.Now(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bufLogger := pool.GetBufLogger()
	for {
		select {
		case <-ctx.Done():
			require.FailNow(t, "Не удалось дождаться выполнения задачи")
		default:
		}

		if bufLogger.Len() < 1 {
			time.Sleep(time.Second)
			continue
		}

		require.Contains(t, bufLogger.String(), "success updated order in pool", "Отсутствует лог об успешном обновлении метрики")
		break
	}
}

func TestNegativeOrderPoolWorkers(t *testing.T) {
	srv := NewAccrualServer(t)
	srv.SetSuccessRequest(false)

	require.NoError(t, srv.Run(), "Не удалось запустить тестовый сервер")
	defer srv.Close()

	pool := NewTestOrderPool(t, srv.GetAddress())

	pool.Get().Run()
	defer pool.Get().Close()

	pool.Get().AddJob(models.Order{
		ID:         123,
		UserID:     uuid.New(),
		Status:     "NEW",
		Accrual:    nil,
		UploadedAt: time.Now(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bufLogger := pool.GetBufLogger()
	for {
		select {
		case <-ctx.Done():
			require.FailNow(t, "Не удалось дождаться выполнения задачи")
		default:
		}

		if bufLogger.Len() < 1 {
			time.Sleep(time.Second)
			continue
		}

		require.Contains(t, bufLogger.String(), "error in order pool", "Отсутствует лог об ошибке при обновлении метрики")
		break
	}
}

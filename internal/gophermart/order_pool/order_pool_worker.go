package orderpool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
)

func (pool *OrderPool) worker() {
	defer pool.wg.Done()

	for job := range pool.jobs {
		pool.workOrder(job)
	}
}

func (pool *OrderPool) workOrder(order models.Order) {
	for order.Status != "PROCESSED" && order.Status != "INVALID" {

		result, err := pool.getAccrualSystemResult(order.ID)
		if err != nil {
			pool.results <- workerResult{
				err:   err,
				order: order,
			}
			return
		}

		if result.Status != order.Status {
			order.Status = result.Status
			order.Accrual = &result.Accrual

			if err = pool.pg.Order().Edit(context.Background(), &order); err != nil {
				pool.results <- workerResult{
					err:   err,
					order: order,
				}
				return
			}
		}

		time.Sleep(time.Second * 2)

	}

	pool.results <- workerResult{
		err:   nil,
		order: order,
	}
}

type accrualSystemResult struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (pool *OrderPool) getAccrualSystemResult(orderID int64) (*accrualSystemResult, error) {
	resp, err := pool.client.R().
		SetResult(&accrualSystemResult{}).
		Get(fmt.Sprintf("/api/orders/%d", orderID))

	if err != nil {
		return nil, fmt.Errorf("error request to accrual: %w", err)
	} else if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("invalid http status (not ok)")
	}

	result, ok := resp.Result().(*accrualSystemResult)
	if !ok {
		return nil, errors.New("invalid response")
	}

	result.Status = strings.ReplaceAll(result.Status, "REGISTERED", "PROCESSING")
	return result, nil
}

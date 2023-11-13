package orderpool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	funcError := func(err error, order models.Order) {
		order.Status = "INVALID"
		if e := pool.pg.Order().Edit(context.Background(), &order); e != nil {
			err = e
		}

		pool.results <- workerResult{
			err:   err,
			order: order,
		}
	}

	for order.Status != "PROCESSED" && order.Status != "INVALID" {
		result, err := pool.getAccrualSystemResult(order.ID)
		if err != nil {
			funcError(err, order)
			return
		}

		if result.Status != "REGISTERED" && result.Status != order.Status {
			order.Status = result.Status
			order.Accrual = &result.Accrual

			if err = pool.pg.Order().Edit(context.Background(), &order); err != nil {
				funcError(err, order)
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
		return nil, err
	} else if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("invalid http status (not ok)")
	}

	result, ok := resp.Result().(*accrualSystemResult)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return result, nil
}

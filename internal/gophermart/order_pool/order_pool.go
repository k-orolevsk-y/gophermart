package orderpool

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/k-orolevsk-y/gophermart/internal/gophermart/config"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/models"
	"github.com/k-orolevsk-y/gophermart/internal/gophermart/repository"
)

type OrderPool struct {
	jobs    chan models.Order
	results chan workerResult

	client *resty.Client
	logger *zap.Logger
	pg     *repository.Pg

	wg     sync.WaitGroup
	closed bool
}

type workerResult struct {
	err   error
	order models.Order
}

func NewOrderPool(logger *zap.Logger, pg *repository.Pg) *OrderPool {
	return &OrderPool{
		client: resty.New().SetBaseURL(config.Config.AccrualSystemAddress),
		logger: logger,
		pg:     pg,
	}
}

func (pool *OrderPool) Run() {
	pool.jobs = make(chan models.Order)
	pool.results = make(chan workerResult)

	for w := 1; w <= config.Config.WorkersLimit; w++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	go pool.workerResults()
}

func (pool *OrderPool) AddJob(order models.Order) {
	pool.jobs <- order
}

func (pool *OrderPool) Close() {
	close(pool.jobs)

	pool.wg.Wait()
	pool.closed = true

	close(pool.results)
}

func (pool *OrderPool) workerResults() {
	for {
		select {
		case result := <-pool.results:
			if result.err != nil {
				pool.logger.Error("Error in order pool", zap.Any("order", result.order), zap.Error(result.err))
			} else {
				pool.logger.Debug("Success updated order in pool", zap.Any("order", result.order))
			}
		default:
			if pool.closed {
				return
			}

			time.Sleep(time.Second)
		}
	}
}

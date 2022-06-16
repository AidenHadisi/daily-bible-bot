package cron

import (
	"time"

	"github.com/go-co-op/gocron"
)

type Cron interface {
	CreateJob(expression string, jobFun interface{}) error
	StartCrons() error
	StopCrons()
}

type SimpleCron struct {
	c *gocron.Scheduler
}

func NewSimpleCron() *SimpleCron {
	return &SimpleCron{
		c: gocron.NewScheduler(time.UTC),
	}
}

func (s *SimpleCron) CreateJob(expression string, jobFun interface{}) error {
	_, err := s.c.Cron(expression).Do(jobFun)
	if err != nil {
		return err
	}
	return nil
}

func (s *SimpleCron) StartCrons() error {
	s.c.StartAsync()
	return nil
}

func (s *SimpleCron) StopCrons() {
	s.c.Stop()
}

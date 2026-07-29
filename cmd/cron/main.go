package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/internal/infra/locks"
	"codedock.run/codedock-tunnel/internal/infra/postgres"
	"codedock.run/codedock-tunnel/internal/infra/redis"
	"codedock.run/codedock-tunnel/internal/repositories"
	"codedock.run/codedock-tunnel/internal/services"
	"codedock.run/codedock-tunnel/internal/workers"
)

var version = "dev"

func main() {
	cfg, err := config.LoadCron()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := postgres.Open(ctx, postgres.Config{DSN: cfg.Database.URL, MaxOpenConns: cfg.Database.MaxConns, MaxIdleConns: cfg.Database.MaxConns, ConnMaxLifetime: cfg.Database.MaxLifetime})
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := redis.Open(ctx, redis.Config{Host: cfg.Redis.Host, Port: cfg.Redis.Port, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()
	lease, err := locks.Acquire(ctx, redisClient.Raw(), "codedock-tunnel:cron", 10*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer lease.Release(context.Background())
	sessions, err := repositories.NewSessionRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	keys, err := repositories.NewAPIKeyRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	deviceLogins, err := repositories.NewDeviceLoginRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	organizations, err := repositories.NewOrganizationRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	billing, err := repositories.NewBillingRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	usageRepository, err := repositories.NewUsageRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	audit, err := repositories.NewAuditRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	usageService, err := services.NewUsageService(usageRepository)
	if err != nil {
		log.Fatal(err)
	}
	aggregation, err := services.NewUsageAggregationService(organizations, usageService)
	if err != nil {
		log.Fatal(err)
	}
	retention, err := services.NewRetentionService(organizations, billing, usageRepository, audit)
	if err != nil {
		log.Fatal(err)
	}
	cleanup, err := workers.NewCleanupJob(sessions, keys, deviceLogins)
	if err != nil {
		log.Fatal(err)
	}
	usageJob, err := workers.NewUsageJob(aggregation)
	if err != nil {
		log.Fatal(err)
	}
	retentionJob, err := workers.NewRetentionJob(retention)
	if err != nil {
		log.Fatal(err)
	}
	operations, err := redis.NewOperations(redisClient)
	if err != nil {
		log.Fatal(err)
	}
	operationalService, err := services.NewOperationsService(organizations, operations)
	if err != nil {
		log.Fatal(err)
	}
	reconciliation, err := workers.NewReconciliationJob(operationalService)
	if err != nil {
		log.Fatal(err)
	}
	runner, err := workers.NewRunner([]workers.Job{cleanup, usageJob, retentionJob, reconciliation}, time.Hour, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := runner.RunOnce(ctx); err != nil {
		log.Fatal(err)
	}
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/daffadon/learn-nats/consumer/utils"
	"github.com/daffadon/learn-nats/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() *pgxpool.Pool {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Failed to load ENV for database")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	return pool
}

func SaveToDB(ctx context.Context, pool *pgxpool.Pool, data []byte) {
	var iotReq types.IoTReq
	err := json.Unmarshal(data, &iotReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal data: %v\n", err)
		return
	}
	cond := utils.CheckCondition(iotReq.Temperature)
	query := "INSERT INTO \"iot\" (iot_id, temperature, condition) VALUES ($1, $2, $3)"
	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to acquire connection: %v\n", err)
		return
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, query, iotReq.IotID, iotReq.Temperature, cond)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute query: %v\n", err)
		return
	}
}

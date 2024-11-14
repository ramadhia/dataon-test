package usecase

import (
	"context"
)

type AlgoUsecase interface {
	AlgoTest1(ctx context.Context, args []string) (string, error)
	AlgoTest2(ctx context.Context, totalBelanja float32, uangDibayar float32) (string, error)
	AlgoTest3(ctx context.Context, arg string) (string, error)
	AlgoTest4(ctx context.Context, args AlgoTest4Request) (string, error)
}

type AlgoTest4Request struct {
	CutiBersama int
	CutiDurasi  int
	JoinDate    string
	CutiDate    string
}

package service

import "context"

func Init(ctx context.Context) error {
	if err := defaultOption.StartAutoRefresh(ctx); err != nil {
		return err
	}
	return nil
}

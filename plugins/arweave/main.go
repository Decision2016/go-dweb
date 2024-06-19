/**
  @author: decision
  @date: 2024/6/17
  @note:
**/

package main

import "context"

type ArweaveStorage struct {
}

func (s *ArweaveStorage) Initial(ctx context.Context) error {
	return nil
}

func (s *ArweaveStorage) Ping(ctx context.Context) error {
	return nil
}
func (s *ArweaveStorage) Exists(ctx context.Context, identity string) (bool,
	error) {
	return true, nil
}

func (s *ArweaveStorage) Upload(ctx context.Context, name string,
	source string) error {
	return nil
}

func (s *ArweaveStorage) Download(ctx context.Context, identity string,
	dst string) error {
	return nil
}

func (s *ArweaveStorage) Delete(ctx context.Context, identity string) error {
	return nil
}

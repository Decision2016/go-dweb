/**
  @author: decision
  @date: 2024/6/13
  @note:
**/

package plugins

import "context"

type FullIPFS struct {
}

func (i *FullIPFS) Initial(ctx context.Context) error {
	return nil
}

func (i *FullIPFS) Ping(ctx context.Context) error {
	return nil
}
func (i *FullIPFS) Exists(ctx context.Context, identity string) (bool, error) {
	return true, nil
}

func (i *FullIPFS) Upload(ctx context.Context, name string, source string) error {
	return nil
}

func (i *FullIPFS) Download(ctx context.Context, identity string, dst string) error {
	return nil
}

func (i *FullIPFS) Delete(ctx context.Context, identity string) error {
	return nil
}

package port

import "context"

type StorageClient interface {
	DownloadPDF(ctx context.Context, url string) ([]byte, error)
}

package objectStorage

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type ObjectStorage interface {
	UploadFile(ctx context.Context, prefix, objectName string) error
	DownloadFile(ctx context.Context, prefix, objectName string) error
	DownloadPictureFile(ctx context.Context, bucketName, prefix, objectName string) error
	SetMetadata(ctx context.Context, prefix, objectName string) error
	GenPublicLink(objectName string) string
}

//go:embed inventory-toll-file-upload.sa.json
var serviceAccountBytes []byte

type serviceAccount struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

func NewGCS(originalBucket, processedBucket string) (ObjectStorage, error) {
	var sa serviceAccount
	err := json.Unmarshal(serviceAccountBytes, &sa)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service account: %w", err)
	}
	if sa.ClientEmail == "" || sa.PrivateKey == "" {
		return nil, errors.New("service account not configured")
	}

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(serviceAccountBytes))
	if err != nil {
		return nil, err
	}
	return &gcs{
		client:                   client,
		originalBucket:           originalBucket,
		processedBucket:          processedBucket,
		serviceAccountEmail:      sa.ClientEmail,
		serviceAccountPrivateKey: []byte(sa.PrivateKey),
	}, nil
}

type gcs struct {
	client                   *storage.Client
	originalBucket           string
	processedBucket          string
	serviceAccountEmail      string
	serviceAccountPrivateKey []byte
}

func (g *gcs) UploadFile(ctx context.Context, fileName, objectName string) error {
	// Open local file.
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*80)
	defer cancel()
	// Upload an object with storage.Writer.
	wc := g.client.Bucket("test-inventory-toll-files").Object(objectName + "/" + fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	err = os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("os.Remove: %v", err)
	}
	return nil
}

func (g *gcs) DownloadFile(ctx context.Context, fileName, objectName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := g.client.Bucket("test-inventory-toll-files").Object(objectName + "/" + fileName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	return nil
}

func (g *gcs) DownloadPictureFile(ctx context.Context, bucketName, fileName, objectName string) error {
	f, err := os.Create(fileName + ".jpg")
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := g.client.Bucket(bucketName).Object(objectName + "/" + fileName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	return nil
}

func (g *gcs) SetMetadata(ctx context.Context, fileName, objectName string) error {
	// Open local file.
	ctx, cancel := context.WithTimeout(ctx, time.Second*80)
	defer cancel()

	// Upload an object with storage.Writer.
	o := g.client.Bucket("test-inventory-toll-files").Object(objectName + "/" + fileName)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}

	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})
	// Update the object CacheControl.
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		CacheControl: "public, max-age=3",
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %v", objectName+"/"+fileName, err)
	}

	return nil
}

const googleapisHost = "https://storage.googleapis.com/"

func (g *gcs) GenPublicLink(objectName string) string {
	return googleapisHost + fmt.Sprintf("%s/%s", g.originalBucket, objectName)
}

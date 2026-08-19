package biznetgio

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
)

type ObjectStorageService struct {
	client *Client
}

type ObjectStorageCreateRequest struct {
	ProductID        int64  `json:"product_id"`
	Cycle            string `json:"cycle"`
	Label            string `json:"label"`
	Quota            int64  `json:"quota,omitempty"`
	Promocode        string `json:"promocode,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type ObjectStorageQuotaUpgradeRequest struct {
	AddQuota         int64  `json:"add_quota,omitempty"`
	PayInvoiceWithCC string `json:"pay_invoice_with_cc,omitempty"`
}

type CredentialStatusRequest struct {
	Active bool `json:"active"`
}

type BucketCreateRequest struct {
	Name string `json:"name"`
	ACL  string `json:"acl,omitempty"`
}

type SetACLRequest struct {
	ACL string `json:"acl,omitempty"`
}

type GenerateObjectURLRequest struct {
	Expiry int64 `json:"expiry,omitempty"`
}

func (s *ObjectStorageService) Create(ctx context.Context, req ObjectStorageCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", "/object-storages", req, &out)
	return out, err
}

func (s *ObjectStorageService) AccountsList(ctx context.Context, status string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/object-storages/accounts"+statusQuery(status), nil, &out)
	return out, err
}

func (s *ObjectStorageService) AccountGet(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d", accountID), nil, &out)
	return out, err
}

func (s *ObjectStorageService) QuotaUpgrade(ctx context.Context, accountID int64, req ObjectStorageQuotaUpgradeRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/object-storages/accounts/%d", accountID), req, &out)
	return out, err
}

func (s *ObjectStorageService) Delete(ctx context.Context, accountID int64) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/object-storages/%d", accountID), nil, nil)
}

func (s *ObjectStorageService) ProductsList(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", "/object-storages/products", nil, &out)
	return out, err
}

func (s *ObjectStorageService) ProductGet(ctx context.Context, productID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/products/%d", productID), nil, &out)
	return out, err
}

func (s *ObjectStorageService) CredentialsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/credentials", accountID), nil, &out)
	return out, err
}

func (s *ObjectStorageService) CredentialCreate(ctx context.Context, accountID int64) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/object-storages/accounts/%d/credentials", accountID), nil, &out)
	return out, err
}

func (s *ObjectStorageService) CredentialUpdate(ctx context.Context, accountID int64, accessKey string, req CredentialStatusRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/object-storages/accounts/%d/credentials/%s", accountID, esc(accessKey)), req, &out)
	return out, err
}

func (s *ObjectStorageService) CredentialDelete(ctx context.Context, accountID int64, accessKey string) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/object-storages/accounts/%d/credentials/%s", accountID, esc(accessKey)), nil, nil)
}

func (s *ObjectStorageService) BucketsList(ctx context.Context, accountID int64) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/buckets", accountID), nil, &out)
	return out, err
}

func (s *ObjectStorageService) BucketCreate(ctx context.Context, accountID int64, req BucketCreateRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/object-storages/accounts/%d/buckets", accountID), req, &out)
	return out, err
}

func (s *ObjectStorageService) BucketSetACL(ctx context.Context, accountID int64, bucketName string, req SetACLRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s", accountID, esc(bucketName)), req, &out)
	return out, err
}

func (s *ObjectStorageService) BucketDelete(ctx context.Context, accountID int64, bucketName string) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s", accountID, esc(bucketName)), nil, nil)
}

func (s *ObjectStorageService) BucketUsage(ctx context.Context, accountID int64, bucketName string) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/usage", accountID, esc(bucketName)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) BucketInfo(ctx context.Context, accountID int64, bucketName, objectOrDirectory string) (map[string]any, error) {
	var out map[string]any
	path := fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/info", accountID, esc(bucketName))
	if objectOrDirectory != "" {
		path += "?object_or_directory=" + url.QueryEscape(objectOrDirectory)
	}
	err := s.client.doJSON(ctx, "GET", path, nil, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectsList(ctx context.Context, accountID int64, bucketName string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects", accountID, esc(bucketName)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectsListInDirectory(ctx context.Context, accountID int64, bucketName, directory string) ([]map[string]any, error) {
	var out []map[string]any
	err := s.client.doJSON(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects/%s/", accountID, esc(bucketName), esc(directory)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectCreateDirectory(ctx context.Context, accountID int64, bucketName, newDirectory string) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects/%s/", accountID, esc(bucketName), esc(newDirectory)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectDownload(ctx context.Context, accountID int64, bucketName, objectName string) ([]byte, error) {
	return s.client.raw(ctx, "GET", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects/%s", accountID, esc(bucketName), esc(objectName)), nil, "application/json")
}

func (s *ObjectStorageService) ObjectSetACL(ctx context.Context, accountID int64, bucketName, objectOrDirectory string, req SetACLRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects/%s", accountID, esc(bucketName), esc(objectOrDirectory)), req, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectDelete(ctx context.Context, accountID int64, bucketName, objectOrDirectory string) error {
	return s.client.doJSON(ctx, "DELETE", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/objects/%s", accountID, esc(bucketName), esc(objectOrDirectory)), nil, nil)
}

func (s *ObjectStorageService) ObjectUpload(ctx context.Context, accountID int64, bucketName, directory, filename string, file []byte) (map[string]any, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := fw.Write(file); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/upload", accountID, esc(bucketName))
	if directory != "" {
		path += "?directory=" + url.QueryEscape(directory)
	}
	var out map[string]any
	err = s.client.do(ctx, "POST", path, buf.Bytes(), w.FormDataContentType(), &out)
	return out, err
}

func (s *ObjectStorageService) ObjectGenerateURL(ctx context.Context, accountID int64, bucketName, objectName string, req GenerateObjectURLRequest) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/url/%s", accountID, esc(bucketName), esc(objectName)), req, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectCopy(ctx context.Context, accountID int64, bucketName, toBucketName, objectName string) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "POST", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/copy/%s/%s", accountID, esc(bucketName), esc(toBucketName), esc(objectName)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) ObjectMove(ctx context.Context, accountID int64, bucketName, toBucketName, objectName string) (map[string]any, error) {
	var out map[string]any
	err := s.client.doJSON(ctx, "PUT", fmt.Sprintf("/object-storages/accounts/%d/buckets/%s/move/%s/%s", accountID, esc(bucketName), esc(toBucketName), esc(objectName)), nil, &out)
	return out, err
}

func (s *ObjectStorageService) AccountID(v map[string]any) (int64, bool) {
	return mapInt64(v, "account_id")
}

func (s *ObjectStorageService) ProductID(v map[string]any) (int64, bool) {
	return mapInt64(v, "product_id")
}

func (s *ObjectStorageService) BucketName(v map[string]any) (string, bool) {
	if name, ok := mapString(v, "name"); ok {
		return name, true
	}
	return mapString(v, "bucket_name")
}

func (s *ObjectStorageService) AccessKey(v map[string]any) (string, bool) {
	return mapString(v, "access_key")
}

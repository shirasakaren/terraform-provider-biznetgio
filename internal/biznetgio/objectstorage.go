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

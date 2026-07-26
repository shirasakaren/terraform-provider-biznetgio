package biznetgio

import (
	"context"
	"fmt"
	"net/url"
)

type GPUService struct {
	client *Client
}

type GPUCreateRequest struct {
	ProductID         int64  `json:"product_id"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	ServiceName       string `json:"service_name,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
	Cycle             string `json:"cycle"`
}

type GPUOneTimeCreateRequest struct {
	ProductID         int64  `json:"product_id"`
	SelectOS          string `json:"select_os"`
	KeypairID         int64  `json:"keypair_id"`
	ServiceName       string `json:"service_name,omitempty"`
	SSHAndConsoleUser string `json:"ssh_and_console_user"`
	ConsolePassword   string `json:"console_password"`
	Promocode         string `json:"promocode,omitempty"`
	PayInvoiceWithCC  string `json:"pay_invoice_with_cc,omitempty"`
	AdditionalHours   int64  `json:"additional_hours,omitempty"`
}

type SetVMStatusRequest struct {
	Status string `json:"status"`
}

type GPURebuildRequest struct {
	SelectOS string `json:"select_os"`

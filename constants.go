package tremendous

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type DeliveryStatus string

const (
	DeliveryStatusPending DeliveryStatus = "PENDING"
	DeliveryStatusSuccess DeliveryStatus = "SUCCEEDED"
	DeliveryStatusFailed  DeliveryStatus = "FAILED"
)

type Scope string

const (
	ScopeDefault         = Scope("default")
	ScopeTeamManagement  = Scope("team_management")
	ScopeFraudPrevention = Scope("fraud_prevention")
)

const (
	LiveEndpoint    = "https://api.tremendous.com/api/v2"
	TestingEndpoint = "https://testflight.tremendous.com/api/v2"
)

type DeliveryMethod string

const (
	DeliveryMethodEmail = DeliveryMethod("EMAIL")
	DeliveryMethodLink  = DeliveryMethod("LINK")
	DeliveryMethodPhone = DeliveryMethod("PHONE")
)

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

type OrderStatus string

const (
	OrderStatusCart     OrderStatus = "CART"
	OrderStatusExecuted OrderStatus = "EXECUTED"
	OrderStatusFailed   OrderStatus = "FAILED"
)

type RoleType string

const (
	RoleTypeAdmin  RoleType = "ADMIN"
	RoleTypeMember RoleType = "MEMBER"
)

type FoundingSources struct {
	FundingSources []*FoundingSource `json:"funding_sources"`
}
type FoundingSource struct {
	Method string `json:"method"`
	Id     string `json:"id"`

	Type string `json:"type"`

	Meta map[string]interface{} `json:"meta"`
}

type FoundingSourceMeta struct {
	AvailableCents int `json:"available_cents"`
	PendingCents   int `json:"pending_cents"`
}

type OrdersList struct {
	Orders []*OrderResponse `json:"orders"`
}
type Orders struct {
	Id         string      `json:"id"`
	ExternalId string      `json:"external_id"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     OrderStatus `json:"status"`
	Channel    string      `json:"channel"`
	Payment    Payment     `json:"payment"`
	Reward     RewardOrder `json:"reward"`
}
type OrderResponse struct {
	Order struct {
		Id         string    `json:"id"`
		ExternalId string    `json:"external_id"`
		CampaignId string    `json:"campaign_id"`
		CreatedAt  time.Time `json:"created_at"`
		Channel    string    `json:"channel"`
		Status     string    `json:"status"`
		Payment    Payment   `json:"payment"`
		Rewards    []Reward  `json:"rewards"`
	} `json:"order"`
}

type Payment struct {
	FundingSourceId string  `json:"funding_source_id"`
	Subtotal        float64 `json:"subtotal"`
	Total           float64 `json:"total"`
	Fees            float64 `json:"fees"`
}

type Campaigns struct {
	Campaigns []*Campaign `json:"campaigns"`
}

type Campaign struct {
	Id          string   `json:"id"`
	Products    []string `json:"products"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
}

type AccessTokenRequest struct {
	ClientId     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	RedirectUri  string    `json:"redirect_uri,omitempty"`
	GrantType    GrantType `json:"grant_type"`
	Code         string    `json:"code,omitempty"`
	RefreshToken string    `json:"refresh_token"`
}

type Fields struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Id       string    `json:"id"`
	Label    string    `json:"label"`
	DataType string    `json:"data_type"`
	Data     FieldData `json:"data"`
	Required bool      `json:"required"`
	Scope    string    `json:"scope"`
}
type FieldData struct {
	Options []string `json:"options"`
}

type Rewards struct {
	Rewards []*Reward `json:"rewards"`
}

type RewardValue struct {
	Denomination float64 `json:"denomination"`
	CurrencyCode string  `json:"currency_code"`
}
type DeliveryMeta struct {
	SubjectLine string `json:"subject_line"`
	FromName    string `json:"from_name"`
}
type Delivery struct {
	Method DeliveryMethod `json:"method"`
	Status DeliveryStatus `json:"status"`
	Meta   DeliveryMeta   `json:"meta"`
	Link   string         `json:"link"`
}
type Recipient struct {
	Email            string            `json:"email"`
	Name             string            `json:"name"`
	Phone            string            `json:"phone"`
	RecipientAddress *RecipientAddress `json:"recipient_address"`
}
type RecipientAddress struct {
	FullName string `json:"full_name"`
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
}

type CustomField struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}
type RewardEvents struct {
	Type    string `json:"type"`
	DateUtc string `json:"date_utc"`
}
type RewardOrder struct {
	Id         string          `json:"id"`
	OrderId    string          `json:"order_id"`
	CreatedAt  string          `json:"created_at"`
	Products   []string        `json:"products"`
	Events     []*RewardEvents `json:"events,omitempty"`
	CampaignID string          `json:"campaign_id"`

	Value        RewardValue   `json:"value"`
	Delivery     Delivery      `json:"delivery"`
	Recipient    Recipient     `json:"recipient"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type Reward struct {
	Id         string          `json:"id"`
	OrderId    string          `json:"order_id"`
	CreatedAt  string          `json:"created_at"`
	Products   []*Product      `json:"products"`
	Events     []*RewardEvents `json:"events"`
	CampaignID string          `json:"campaign_id"`

	Value        RewardValue   `json:"value"`
	Delivery     Delivery      `json:"delivery"`
	Recipient    Recipient     `json:"recipient"`
	CustomFields []CustomField `json:"custom_fields"`
}
type Errors struct {
	Errors struct {
		Message string                 `json:"message"`
		Payload map[string]interface{} `json:"payload"`
	} `json:"errors"`
}

type Country struct {
	Abbr string `json:"abbr"`
}
type Image struct {
	Src string `json:"src"`
}
type Sku struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
type Product struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Countries []Country `json:"countries"`
	Images    []Image   `json:"images"`
	Skus      []Sku     `json:"skus"`
}

type Products struct {
	Products []*Product `json:"products"`
}

type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "PENDING"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusDeleted InvoiceStatus = "DELETED"
)

type Invoice struct {
	Id       string        `json:"id"`
	PoNumber string        `json:"po_number"`
	Amount   int           `json:"amount"`
	Status   InvoiceStatus `json:"status"`
}

type Invoices struct {
	Invoices []*Invoice `json:"invoices"`
}

type Org struct {
	ParentID string `json:"parent_id,omitempty"`
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Website  string `json:"website"`
}

type OrgAccessToken struct {
	AccessToken string `json:"access_token"`
}
type Organization struct {
	Organization Org `json:"organization"`
}
type Organizations struct {
	Organization []*Organization `json:"s"`
}

type User struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Role      RoleType `json:"role"`
	Status    string   `json:"status"`
	InviteUrl string   `json:"invite_url"`
}

type Member struct {
	Member User `json:"member"`
}
type Members struct {
	Members []User `json:"members"`
}

type Hook struct {
	Url string `json:"url"`

	Id string `json:"id,omitempty"`

	PrivateKey string `json:"private_key,omitempty"`
}

type Webhook struct {
	Webhook Hook `json:"webhook"`
}

type WebhookEvents struct {
	Events []string `json:"events"`
}

func formatResponse[T any](r *http.Response, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	p, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	println(string(p))

	var t T
	err = json.Unmarshal(p, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func joinURL(baseURL, p string) string {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "" // Handle error accordingly
	}

	parsedURL.Path = path.Join(strings.TrimRight(parsedURL.Path, "/"), strings.TrimLeft(p, "/"))
	return parsedURL.String()
}

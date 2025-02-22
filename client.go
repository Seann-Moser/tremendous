package tremendous

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiKey string

	clientID     string
	clientSecret string
	refreshToken string
	accessKey    string

	httpClient  *http.Client
	autoRefresh bool

	refresh  chan TokenResponse
	endpoint string
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		refresh:    make(chan TokenResponse, 10),
		endpoint:   LiveEndpoint,
	}
}
func (c *Client) Close() {
	close(c.refresh)
}

func (c *Client) OauthRefresh() <-chan TokenResponse {
	return c.refresh
}

func (c Client) SetEndpoint(endpoint string) Client {
	c.endpoint = endpoint
	return c
}

func (c Client) NewClientWithAPIKey(apiKey string) Client {
	c.apiKey = apiKey
	return c
}

func (c Client) NewClientWithOAuth(config OauthConfig, autoRefresh bool) Client {
	c.clientSecret = config.ClientSecret
	c.clientID = config.ClientId
	c.accessKey = config.AccessToken
	c.refreshToken = config.RefreshToken
	c.autoRefresh = autoRefresh
	return c
}

func (c Client) do(req *http.Request) (*http.Response, error) {
	key := ""
	if c.apiKey != "" {
		key = c.apiKey
	} else {
		key = c.accessKey
	}
	if key == "" {
		return nil, errors.New("no api key provided")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	println(key)
	return c.httpClient.Do(req)
}

func (c *Client) doRequest(ctx context.Context, method, p string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, joinURL(c.endpoint, p), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && c.refreshToken != "" && c.autoRefresh {
		token, err := c.SendOauthRequest(ctx, &AccessTokenRequest{
			ClientId:     c.clientID,
			ClientSecret: c.clientSecret,
			GrantType:    GrantTypeRefreshToken,
			RefreshToken: c.refreshToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		c.accessKey = token.AccessToken
		c.refreshToken = token.RefreshToken

		c.refresh <- *token
		resp, err := c.do(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted {
		return resp, nil
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("tremendous: unexpected status code: %d : %s", resp.StatusCode, string(b))
}

func (c *Client) CreateOrder(ctx context.Context, order *Orders) (*Orders, error) {
	return formatResponse[Orders](c.doRequest(ctx, http.MethodPost, "/orders", order))
}

func (c *Client) ListOrders(ctx context.Context) (*Orders, error) {
	return formatResponse[Orders](c.doRequest(ctx, http.MethodGet, "/orders", nil))
}

func (c *Client) RetrieveOrder(ctx context.Context, orderID string) (*Orders, error) {
	// /orders/{:id}
	return formatResponse[Orders](c.doRequest(ctx, http.MethodGet, "/orders/"+orderID, nil))
}

func (c *Client) ListRewards(ctx context.Context) (*Rewards, error) {
	return formatResponse[Rewards](c.doRequest(ctx, http.MethodGet, "/rewards", nil))
}

func (c *Client) RetrieveReward(ctx context.Context, rewardID string) (*Reward, error) {
	return formatResponse[Reward](c.doRequest(ctx, http.MethodGet, "/rewards/"+rewardID, nil))
}

func (c *Client) ApproveReward(ctx context.Context, rewardID string) (*Reward, error) {
	return formatResponse[Reward](c.doRequest(ctx, http.MethodPost, "/rewards/"+rewardID+"/approve", nil))
}

func (c *Client) ListCampaigns(ctx context.Context) (*Campaigns, error) {
	return formatResponse[Campaigns](c.doRequest(ctx, http.MethodGet, "/campaigns", nil))
}

func (c *Client) ListProducts(ctx context.Context) (*Products, error) {
	return formatResponse[Products](c.doRequest(ctx, http.MethodGet, "/products", nil))
}

func (c *Client) ListFundingSources(ctx context.Context) (*FoundingSources, error) {
	return formatResponse[FoundingSources](c.doRequest(ctx, http.MethodGet, "/funding_sources", nil))
}

func (c *Client) RetrieveFundingSource(ctx context.Context, foundingID string) (*FoundingSource, error) {
	return formatResponse[FoundingSource](c.doRequest(ctx, http.MethodGet, "/funding_sources/"+foundingID, nil))
}

func (c *Client) ListInvoices(ctx context.Context) (*Invoices, error) {
	return formatResponse[Invoices](c.doRequest(ctx, http.MethodGet, "/invoices", nil))
}

func (c *Client) RetrieveInvoice(ctx context.Context, invoiceId string) (*Invoice, error) {
	return formatResponse[Invoice](c.doRequest(ctx, http.MethodGet, "/invoices/"+invoiceId, nil))
}

func (c *Client) RetrieveInvoicePDF(ctx context.Context, invoiceId string) ([]byte, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/invoices/"+invoiceId+"/pdf", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve invoice PDF: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) DeleteInvoice(ctx context.Context, invoiceId string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "/invoices/"+invoiceId, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete invoice: status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrganization(ctx context.Context, org *Organization) (*Organization, error) {
	return formatResponse[Organization](c.doRequest(ctx, http.MethodPost, "/organizations", org))
}

func (c *Client) ListOrganizations(ctx context.Context) (*Organizations, error) {
	return formatResponse[Organizations](c.doRequest(ctx, http.MethodGet, "/organizations", nil))
}

func (c *Client) RetrieveOrganization(ctx context.Context, orgID string) (*Organization, error) {
	return formatResponse[Organization](c.doRequest(ctx, http.MethodGet, "/organizations/"+orgID, nil))
}

func (c *Client) CreateOrgAccessToken(ctx context.Context, orgID string) (*OrgAccessToken, error) {
	return formatResponse[OrgAccessToken](c.doRequest(ctx, http.MethodPost, "/organizations/"+orgID+"/access_token", nil))
}

func (c *Client) CreateMember(ctx context.Context, member *Member) (*Member, error) {
	return formatResponse[Member](c.doRequest(ctx, http.MethodPost, "/members", member))
}

func (c *Client) ListMembers(ctx context.Context) (*Members, error) {
	return formatResponse[Members](c.doRequest(ctx, http.MethodGet, "/members", nil))
}

func (c *Client) RetrieveMember(ctx context.Context, memberID string) (*Member, error) {
	return formatResponse[Member](c.doRequest(ctx, http.MethodGet, "/members/"+memberID, nil))
}

func (c *Client) ListFields(ctx context.Context) (*Fields, error) {
	return formatResponse[Fields](c.doRequest(ctx, http.MethodGet, "/fields", nil))
}

func (c *Client) ValidWebhook(r *http.Request, key string) (bool, error) {
	signatureHeader := r.Header.Get("Tremendous-Webhook-Signature")
	parts := bytes.SplitN([]byte(signatureHeader), []byte("="), 2)
	if len(parts) != 2 || string(parts[0]) != "sha256" {
		return false, fmt.Errorf("invalid algorithm")
	}

	hmacKey := []byte(key)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset body for future reads

	h := hmac.New(sha256.New, hmacKey)
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), parts[1]), nil
}

func (c *Client) CreateWebhook(ctx context.Context, url string) (*Webhook, error) {
	return formatResponse[Webhook](c.doRequest(ctx, http.MethodPost, "/webhooks", map[string]string{"url": url}))
}

func (c *Client) ShowWebhookEvents(ctx context.Context, webhookID string) (*WebhookEvents, error) {
	return formatResponse[WebhookEvents](c.doRequest(ctx, http.MethodGet, "/webhooks/"+webhookID+"/events", nil))
}

func (c *Client) SimulateWebhook(ctx context.Context, webhookID string, event string) error {
	type SimulatedEvent struct {
		Event string `json:"event"`
	}

	_, err := c.doRequest(ctx, http.MethodPost, "/webhooks/"+webhookID+"/simulate", SimulatedEvent{Event: event})
	return err
}

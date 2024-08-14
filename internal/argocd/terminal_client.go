package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"nhooyr.io/websocket"
)

type OperationType string

const OpStdin OperationType = "stdin"
const OpStdout OperationType = "stdout"
const OpResize OperationType = "resize"

type TerminalClientOptions struct {
	ProjectName  string `json:"projectName"`
	POD          string `json:"pod"`
	Container    string `json:"container"`
	ArgoCDServer string `json:"argoCDServer"`
	OptionsApp
	OptionsK8s
}

type TerminalClient struct {
	conn *websocket.Conn
}

type Operation struct {
	Operation OperationType `json:"operation"`
	Data      string        `json:"data,omitempty"`
	Rows      int           `json:"rows"`
	Cols      int           `json:"cols"`
}

func (tc *TerminalClient) Send(ctx context.Context, op Operation) error {
	bytes, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("marshal operation: %w", err)
	}
	return tc.conn.Write(ctx, websocket.MessageText, bytes)
}

func (tc *TerminalClient) Recv(ctx context.Context) (op Operation, err error) {
	msgType, msgReader, err := tc.conn.Reader(ctx)
	if err != nil {
		return Operation{}, fmt.Errorf("websocket read: %w", err)
	}

	if msgType != websocket.MessageText {
		return Operation{}, fmt.Errorf("unsupported message type: %v", msgType)
	}

	if err = json.NewDecoder(msgReader).Decode(&op); err != nil {
		return Operation{}, fmt.Errorf("unmarshal operation: %w", err)
	}
	return
}

func NewTerminalClient(ctx context.Context, u *url.URL, headers http.Header, client *http.Client) (*TerminalClient, error) {
	conn, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{HTTPHeader: headers, HTTPClient: client})
	if err != nil {
		return nil, fmt.Errorf("dial url: %s: %w", u.String(), err)
	}
	return &TerminalClient{
		conn: conn,
	}, nil
}

func BuildDefaultTerminalURL(options TerminalClientOptions) (*url.URL, error) {
	apiURL, err := url.Parse(options.ArgoCDServer + "/")
	if err != nil {
		return nil, err
	}

	apiURL = apiURL.JoinPath("/terminal")
	query := apiURL.Query()
	query.Add("pod", options.POD)
	query.Add("container", options.Container)
	query.Add("appName", options.AppName)
	query.Add("appNamespace", options.AppNamespace)
	query.Add("projectName", options.ProjectName)
	query.Add("namespace", options.Namespace)
	apiURL.RawQuery = query.Encode()
	apiURL.Scheme = "wss"
	return apiURL, nil
}

func BuildDefaultHeaders(token string) http.Header {
	return map[string][]string{
		"Content-Type": {"application/json"},
		"Cookie":       {fmt.Sprintf("argocd.token=%s", token)},
	}
}

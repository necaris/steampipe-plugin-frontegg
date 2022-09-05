package frontegg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableFronteggUser() *plugin.Table {
	return &plugin.Table{
		Name:        "frontegg_user",
		Description: "TODO",
		List: &plugin.ListConfig{
			Hydrate: listUser,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "User's unique UUID"},
			{Name: "sub", Type: proto.ColumnType_STRING, Description: "UUID"},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The user's email address"},
			{Name: "verified", Type: proto.ColumnType_BOOL, Description: "Has the user verified their email?"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "User's name"},
			{Name: "phone_number", Type: proto.ColumnType_STRING, Description: "Phone number, if entered"},
			{Name: "profile_picture_url", Type: proto.ColumnType_STRING, Description: "URL to a profile picture"},
			{Name: "provider", Type: proto.ColumnType_STRING, Description: "What provider"},
			{Name: "mfa_enrolled", Type: proto.ColumnType_BOOL, Description: "Have they enrolled MFA?"},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Arbitrary metadata"},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: "Tenant that is their default"},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time the user was created"},
			{Name: "is_locked", Type: proto.ColumnType_BOOL, Description: "Is the user locked?"},
		},
	}
}

type FronteggUsersResponse struct {
	Items []FronteggUser    `json:"items"`
	Links map[string]string `json:"_links"`
}

type FronteggUser struct {
	ID                string      `json:"id"`
	Sub               string      `json:"sub"`
	Email             string      `json:"email"`
	Verified          bool        `json:"verified"`
	Name              string      `json:"name"`
	PhoneNumber       string      `json:"phoneNumber"`
	ProfilePictureURL string      `json:"profilePictureUrl"`
	Provider          string      `json:"provider"`
	MfaEnrolled       bool        `json:"mfaEnrolled"`
	Metadata          interface{} `json:"metadata"`
	TenantID          string      `json:"tenantId"`
	CreatedAt         time.Time   `json:"createdAt"`
	IsLocked          bool        `json:"isLocked"`
	// TODO:
	//     "tenantIds": [
	//       "5b58adea-4699-4af1-be96-1e21d26693c3"
	//     ],
	//     "roles": [],
	//     "permissions": [],
	//     "tenants": [
	//       {
	//         "tenantId": "5b58adea-4699-4af1-be96-1e21d26693c3",
	//         "roles": [
	//           {
	//             "id": "0a4d54f7-05d0-4f7a-8f57-6a2c2474b2c1",
	//             "vendorId": "60adc590-2624-4f2e-9f75-ce4ed0deb974",
	//             "tenantId": null,
	//             "key": "",
	//             "name": "",
	//             "description": "",
	//             "isDefault": true,
	//             "firstUserRole": false,
	//             "createdAt": "2022-08-16T22:52:07.000Z",
	//             "updatedAt": "2022-09-01T22:38:37.000Z",
	//             "permissions": [
	//               "1800fdcc-9d26-468a-9b30-f8d983c73585",
	//             ],
	//             "level": 3
	//           },
	//         ]
	//       }
	//       ]
}

func listUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := NewClientFromEnv(ctx, d)
	if err != nil {
		return nil, err
	}
	baseUrl := "identity/resources/users/v2"
	nextUrl := baseUrl
	for nextUrl != "" {
		var result FronteggUsersResponse
		raw, err := client.RawRequest(http.MethodGet, nextUrl)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(raw.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		for _, u := range result.Items {
			logger.Warn("got one", "u", hclog.Fmt("%+v", u))
			d.StreamListItem(ctx, u)
		}
		if result.Links["next"] != "" {
			nextUrl = baseUrl + result.Links["next"]
		} else {
			nextUrl = ""
		}
	}
	return nil, nil
}

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := NewClientFromEnv(ctx, d)
	if err != nil {
		return nil, err
	}
	quals := d.KeyColumnQuals
	plugin.Logger(ctx).Warn("getUser", "quals", quals)
	id := quals["id"].GetValue()
	plugin.Logger(ctx).Warn("getUser", "id", id)
	result, err := client.RequestAsInterface(http.MethodGet, fmt.Sprintf("identity/resources/users/v1/%s/", id))
	if err != nil {
		return nil, err
	}
	return result, nil
}

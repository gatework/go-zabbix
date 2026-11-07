//go:build integration

package zabbix

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestIntegrationLoginWithUserGroupAccess(t *testing.T) {
	admin, ctx := loginClient(t)
	roles, err := admin.Role().Get(ctx, RoleGetParams{
		Output: []string{"roleid"}, Filter: map[string]any{"type": 1}, Limit: 1,
	})
	integrationCheck(t, err)
	if len(roles) != 1 || roles[0].RoleID == "" {
		t.Fatal("fixture has no user role")
	}

	for _, tc := range []struct {
		name      string
		debugMode int
	}{
		{name: "default GUI access"},
		{name: "debug enabled", debugMode: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := fmt.Sprintf("go-zabbix-auth-%d", time.Now().UnixNano())
			ids, err := admin.UserGroup().Create(ctx, UserGroupCreateParams{
				Name: name, GUIAccess: integrationValue(0), DebugMode: integrationValue(tc.debugMode),
				MFAStatus: integrationValue(0),
			})
			groupID := integrationID(t, ids, err)
			integrationDelete(t, admin.UserGroup().Delete, groupID)
			const password = "G0-zabbix-Session-Fixture!"
			ids, err = admin.User().Create(ctx, UserCreateParams{
				Username: name, Passwd: integrationValue(password), RoleID: &roles[0].RoleID,
				UserGroups: []UserGroupRef{{UsrGrpID: groupID}},
			})
			userID := integrationID(t, ids, err)
			integrationDelete(t, admin.User().Delete, userID)

			client, err := New(*zbxURL, WithTimeout(15*time.Second))
			integrationCheck(t, err)
			t.Cleanup(client.CloseIdleConnections)
			t.Cleanup(func() {
				cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := client.Logout(cleanup); err != nil {
					t.Errorf("fixture session cleanup: %v", err)
				}
			})

			session, err := client.LoginWithUserData(ctx, name, password)
			integrationCheck(t, err)
			if session == nil || session.UserID != userID || session.Username != name || session.SessionID == "" || !client.HasCredentials() {
				t.Fatal("fixture login did not return and install the expected user session")
			}
			if session.GUIAccess.String() != "0" || session.DebugMode.String() != strconv.Itoa(tc.debugMode) {
				t.Fatalf("group access = %s/%s, want 0/%d", session.GUIAccess, session.DebugMode, tc.debugMode)
			}
			integrationCheck(t, client.Logout(ctx))
			if client.HasCredentials() {
				t.Fatal("fixture logout retained session credentials")
			}
		})
	}
}

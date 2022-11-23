package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/google/go-cmp/cmp"

	"github.com/dtherhtun/service/business/data/store/user"
	"github.com/dtherhtun/service/business/data/tests"
	"github.com/dtherhtun/service/business/sys/auth"
	"github.com/dtherhtun/service/business/sys/database"
)

var dbc = tests.DBContainer{
	Image: "postgres:13-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, dbc)
	t.Cleanup(teardown)

	store := user.NewStore(log, db)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := context.Background()
			now := time.Date(2022, time.April, 9, 0, 0, 0, 0, time.UTC)

			nu := user.NewUser{
				Name:            "D Ther Htun",
				Email:           "dtherhtun.cw@gmail.com",
				Roles:           []string{auth.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

			usr, err := store.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "service project",
					Subject:   usr.ID,
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
					IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
				},
				Roles: []string{auth.RoleUser},
			}

			saved, err := store.QueryByID(ctx, claims, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould be able get back the same user. Diff:\n%s.", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould be able get back the same user.", tests.Success, testID)

			upd := user.UpdateUser{
				Name:  tests.StringPointer("Luo Sten"),
				Email: tests.StringPointer("dtherhtun@yahoo.com"),
			}

			claims = auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "service project",
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
					IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
				},
				Roles: []string{auth.RoleAdmin},
			}

			if err := store.Update(ctx, claims, usr.ID, upd, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able update user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able update user.", tests.Success, testID)

			saved, err = store.QueryByEmail(ctx, claims, *upd.Email)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by Email : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by Email.", tests.Success, testID)

			if saved.Name != *upd.Name {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Name)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Name)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Name.", tests.Success, testID)
			}

			if saved.Email != *upd.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Email)
				t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updates to Email.", tests.Success, testID)
			}

			if err := store.Delete(ctx, claims, usr.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = store.QueryByID(ctx, claims, usr.ID)
			if !errors.Is(err, database.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", tests.Success, testID)
		}
	}
}

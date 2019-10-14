package cache

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/log"
	"github.com/tinylib/msgp/msgp"
)

func userCacheKey(ID int64) string {
	return fmt.Sprintf("user-%d", ID)
}

// New return a new Cache storage
func New(client *memcache.Client, storage iface.Storage) iface.Storage {
	return &Cache{client, storage}
}

// Cache is the stuct of the Cache
type Cache struct {
	client  *memcache.Client
	storage iface.Storage
}

// Tx start a new transaction
func (c *Cache) Tx() (*sql.Tx, error) {
	return c.storage.Tx()
}

// AddUser add a new user to the cache storage
func (c *Cache) AddUser(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	return c.storage.AddUser(ctx, tx, name)
}

// DeleteUser remove an user from the cache storage
func (c *Cache) DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	_ = c.client.Delete(userCacheKey(userID))
	return c.storage.DeleteUser(ctx, tx, userID)
}

// FilterUsersID retrieve usersID from the cache storage
func (c *Cache) FilterUsersID(ctx context.Context, filter iface.FilterUsers) ([]int64, error) {
	return c.storage.FilterUsersID(ctx, filter)
}

// FetchUsers retrieve users from the cache storage
func (c *Cache) FetchUsers(ctx context.Context, IDs ...int64) ([]*entity.User, error) {
	keys := make([]string, 0, len(IDs))
	for _, ID := range IDs {
		keys = append(keys, fmt.Sprintf("user-%d", ID))
	}

	IDsToFetch := make([]int64, 0, len(IDs))
	hit := map[int64]bool{}
	musers := map[int64]*entity.User{}
	if items, err := c.client.GetMulti(keys); err != nil {
		log.Log(err)
	} else {
		for _, item := range items {
			var user entity.User
			if err := msgp.Decode(bytes.NewBuffer(item.Value), &user); err != nil {
				log.Log(err)
				continue
			}

			musers[user.ID] = &user
			hit[user.ID] = true
		}
		for _, ID := range IDs {
			if _, has := hit[ID]; !has {
				IDsToFetch = append(IDsToFetch, ID)
			}
		}
	}

	if len(IDsToFetch) != 0 {
		dbusers, err := c.storage.FetchUsers(ctx, IDsToFetch...)
		if err != nil {
			return nil, err
		}

		for _, user := range dbusers {
			var buf bytes.Buffer
			err := msgp.Encode(&buf, user)
			if err != nil {
				log.Log(err)
				continue
			}

			err = c.client.Set(&memcache.Item{
				Key:   fmt.Sprintf("user-%d", user.ID),
				Value: buf.Bytes(),
			})
			if err != nil {
				log.Log(err)
			}
			musers[user.ID] = user
		}
	}

	users := make([]*entity.User, 0, len(IDs))
	for _, ID := range IDs {
		users = append(users, musers[ID])
	}

	return users, nil
}

// AddEmail add a new email to the cache storage
func (c *Cache) AddEmail(ctx context.Context, tx *sql.Tx, userID int64, address string) (int64, error) {
	return c.storage.AddEmail(ctx, tx, userID, address)
}

// DeleteEmail remove an email from the cache storage
func (c *Cache) DeleteEmail(ctx context.Context, tx *sql.Tx, emailID int64) error {
	return c.storage.DeleteEmail(ctx, tx, emailID)
}

// DeleteEmailsByUserID remove emails from an userID
func (c *Cache) DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	return c.storage.DeleteEmailsByUserID(ctx, tx, userID)
}

// FilterEmails returns Emails for a given filter
func (c *Cache) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	return c.storage.FilterEmails(ctx, filter)
}

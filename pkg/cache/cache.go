package cache

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gogo/protobuf/proto"
	"github.com/rafaelsq/boiler/pkg/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/log"
)

func userCacheKey(ID int64) string {
	return fmt.Sprintf("user-%d", ID)
}

func userFilterCacheKey(filter iface.FilterUsers) string {
	return fmt.Sprintf("user-filter-%s|%d|%d", filter.Email, filter.Offset, filter.Limit)
}

func New(client *memcache.Client, storage iface.Storage) iface.Storage {
	return &Cache{client, storage}
}

type Cache struct {
	client  *memcache.Client
	storage iface.Storage
}

// begin transaction
func (c *Cache) Tx() (*sql.Tx, error) {
	return c.storage.Tx()
}

// user
func (c *Cache) AddUser(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	return c.storage.AddUser(ctx, tx, name)
}

func (c *Cache) DeleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	_ = c.client.Delete(userCacheKey(userID))
	return c.storage.DeleteUser(ctx, tx, userID)
}

func (c *Cache) FilterUsersID(ctx context.Context, filter iface.FilterUsers) ([]int64, error) {
	return c.storage.FilterUsersID(ctx, filter)
}

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
			if err := proto.Unmarshal(item.Value, &user); err != nil {
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
			b, err := proto.Marshal(user)
			if err != nil {
				log.Log(err)
				continue
			}
			c.client.Set(&memcache.Item{
				Key:   fmt.Sprintf("user-%d", user.ID),
				Value: b,
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

// email
func (c *Cache) AddEmail(ctx context.Context, tx *sql.Tx, userID int64, address string) (int64, error) {
	return c.storage.AddEmail(ctx, tx, userID, address)
}

func (c *Cache) DeleteEmail(ctx context.Context, tx *sql.Tx, emailID int64) error {
	return c.storage.DeleteEmail(ctx, tx, emailID)
}

func (c *Cache) DeleteEmailsByUserID(ctx context.Context, tx *sql.Tx, userID int64) error {
	return c.storage.DeleteEmailsByUserID(ctx, tx, userID)
}

func (c *Cache) FilterEmails(ctx context.Context, filter iface.FilterEmails) ([]*entity.Email, error) {
	return c.storage.FilterEmails(ctx, filter)
}

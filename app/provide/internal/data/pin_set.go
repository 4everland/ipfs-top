package data

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/4everland/ipfs-top/app/provide/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"
	"strings"
	"time"
)

const CidTable = "cids"

type PinSetRepo struct {
	data *gorm.DB
	log  *log.Helper

	filterPrefixKeys []string
}

type Cids struct {
	Cid       string          `gorm:"primary_key"`
	CreatedAt int64           `gorm:"autoCreateTime"`
	UpdatedAt int64           `gorm:"autoUpdateTime"`
	Pkey      JsonStringArray `gorm:"type:jsonb"`
}

func NewDataSetRepo(db *gorm.DB, conf *conf.Data, l log.Logger) *PinSetRepo {
	return &PinSetRepo{
		data: db,
		log:  log.NewHelper(l),

		filterPrefixKeys: conf.PinSet.FilterPrefix,
	}
}

func (d *PinSetRepo) AllKeys(ctx context.Context, endTime time.Time) (chan cid.Cid, chan error) {

	var dataCh = make(chan cid.Cid)
	var errCh = make(chan error)
	go func() {
		defer close(dataCh)
		defer close(errCh)
		var cids = make([]Cids, 0, 1000)
		var nextKey = ""
		for {
			tx := d.data.WithContext(ctx).Table(CidTable).Where("cid > ?", nextKey).Order("cid asc").Limit(1000).Find(&cids)
			if tx.Error != nil {
				errCh <- tx.Error
				return
			}
			for _, c := range cids {
				nextKey = c.Cid
				if endTime.Unix() <= c.UpdatedAt {
					continue
				}
				if len(c.Pkey) != 0 {
					for _, p := range d.filterPrefixKeys {
						if strings.HasPrefix(c.Pkey[0], p) {
							continue
						}
					}
					cc, err := cid.Decode(c.Cid)
					if err != nil {
						continue
					}
					dataCh <- cc
				}
				//dataCh <- c
			}
			if len(cids) < 1000 {
				return
			}

			cids = cids[:0]

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				continue
			}
		}

	}()

	return dataCh, errCh

}

// JsonStringArray represents a one-dimensional array of the PostgreSQL character types.
type JsonStringArray []string

// Scan implements the sql.Scanner interface.
func (a *JsonStringArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to StringArray", src)
}

func (a *JsonStringArray) scanBytes(src []byte) error {
	return json.Unmarshal(src, a)
}

// Value implements the driver.Valuer interface.
func (a JsonStringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(data), nil

}

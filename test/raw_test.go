package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
		rows,err := db.Query(`SElECT create_time FROM user WHERE id=? `,1)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var createTime time.Time
			err := rows.Scan(&createTime)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(createTime)
		}
}

package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
		rows,err := db.Query(`SElECT update_time FROM user WHERE id=? `,1)
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var updateTime time.Time
			err := rows.Scan(&updateTime)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(updateTime)
		}
}

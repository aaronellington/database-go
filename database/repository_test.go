package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/aaronellington/database-go/database"
	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"createdAt"`
	Email     string    `db:"email"`
}

func (e Order) TableName() string {
	return "user"
}

func (e Order) Joins() string {
	return ""
}

type OrderRepository struct {
	helper *database.Repository[Order]
}

func TestFoobar(t *testing.T) {
	connection := getTestConnection(t)

	orderRepo := OrderRepository{
		helper: database.NewRepository[Order](connection),
	}

	if _, err := orderRepo.helper.Insert(context.TODO(), Order{
		CreatedAt: time.Now(),
		Email:     "user@exmaple.com",
	}); err != nil {
		t.Fatal(err)
	}

	orders, err := orderRepo.helper.Find(
		context.Background(),
		database.Query{},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(orders) != 1 {
		t.Fatal("no orders found")
	}
}

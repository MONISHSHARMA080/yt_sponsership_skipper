package paymentbackendgo

type TemporaryFieldToVerifyPaymentLater struct {
	UserAccountID    int    `json:"user_account_id"`
	RecurringOrderID string `json:"recurring_order_id"`
	OnetimeOrderID   string `json:"onetime_order_id"`
}

// insert into employees
// (id, name, city, department, salary)
// values
//   (11, 'Diane', 'Berlin', 'hr', 70),
//   (21, 'Emma', 'London', 'it', 95),
//   (25, 'Frank', 'Berlin', 'it', 120),
//   (33, 'Alice', 'Berlin', 'sales', 100)
// on conflict do update set
//   city = excluded.city,
//   salary = excluded.salary;

//

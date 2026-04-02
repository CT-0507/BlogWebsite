```go
var OrderSaga = domain.SagaDefinition{
	Name: "order",
	Steps: []domain.Step{
		{
			Name:           "create_order",
			ActionType:     "CreateOrder",
			CompensateType: "CancelOrder",
			Next:           "reserve_inventory",
		},
		{
			Name:           "reserve_inventory",
			ActionType:     "ReserveInventory",
			CompensateType: "ReleaseInventory",
			Next:           "charge_payment",
		},
		{
			Name:           "charge_payment",
			ActionType:     "ChargePayment",
			CompensateType: "RefundPayment",
			Next:           "Complete",
		},
	},
}
```

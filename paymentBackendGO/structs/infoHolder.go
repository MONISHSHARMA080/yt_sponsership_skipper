
package structs

type InfoHolder struct{
	PriceForRecurring int64
	PriceForOneTime int64
	DecryptedKey string
	Name string
	Email string
	IsPaidUser bool
	// PlanType string
}
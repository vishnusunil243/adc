package payment_service

type CheckoutSessionRes struct {
	CheckoutSessionId string `json:"checkout_session_id"`
}

func NewCheckoutSessionRes(checkoutSessionId string) *CheckoutSessionRes {
	return &CheckoutSessionRes{
		CheckoutSessionId: checkoutSessionId,
	}
}

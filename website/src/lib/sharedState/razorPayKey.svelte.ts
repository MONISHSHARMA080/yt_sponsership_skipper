export const razorpayOrderId = $state<{ orderIdForRecurring: string | null, orderIdForOnetime: string | null, numberOfTimesKeyUsed: number, fetchingStatus: "fetching" | "error" | "success" }>({
  orderIdForOnetime: null, orderIdForRecurring: null, numberOfTimesKeyUsed: 0,
  fetchingStatus: "fetching"
})

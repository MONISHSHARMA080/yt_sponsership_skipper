
interface orderIdType { orderIdForRecurring: string | null, orderIdForOnetime: string | null, numberOfTimesKeyUsed: number, fetchingStatus: "fetching" | "error" | "success", areWeInAMiddleOfMultipleFetchCycle: boolean }

export const razorpayOrderId = $state<orderIdType>({
  orderIdForOnetime: null, orderIdForRecurring: null, numberOfTimesKeyUsed: 0,
  fetchingStatus: "fetching", areWeInAMiddleOfMultipleFetchCycle: false
})

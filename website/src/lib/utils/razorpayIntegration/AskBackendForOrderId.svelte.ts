import { razorpayOrderId } from "$lib/sharedState/razorPayKey.svelte";
import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { executeWithKeyRefresh } from "../ApiReqHelper/KeyRefreshhandler";

interface requestType {
  user_key: string
}

interface responseType {
  order_id_for_recurring: string,
  order_id_for_onetime: string,
  message: string,
  status_code: number,
}

export async function askBackendForOrderId(keyStateObj: keyStateObject): Promise<boolean> {
  console.log(`Starting askBackendForOrderId`);

  // Define function to create promise array
  const createPromiseArray = (userKey: string) => {
    const reqBody: requestType = { user_key: userKey };
    return [
      fetch('/api/makeAPayment', {
        headers: { 'Content-Type': 'application/json' },
        method: "POST",
        body: JSON.stringify(reqBody)
      })
    ];
  };

  // Define how to process individual promises
  const processPromise = async (resp1: Promise<Response>): Promise<responseType> => {
    try {
      const resp = await resp1;
      console.log(`Response received from makeAPayment`);

      if (!resp.ok) {
        console.error("The response is not ok from the backend:", resp);
        throw new Error("The response is not ok");
      }

      // Parse and validate JSON
      const responseInJson = await resp.json();
      console.log("Valid JSON received:", responseInJson);

      if (validateResponseInJSONTOBeMyType(responseInJson)) {
        return responseInJson;
      } else {
        throw new Error("Invalid response format");
      }
    } catch (error) {
      throw error;
    }
  };

  // Define success callback
  const onSuccess = (result: responseType) => {
    console.log(`Storing response in razorpay storage: ${result.order_id_for_recurring}, ${result.order_id_for_onetime}`);

    if (result.order_id_for_onetime) {
      razorpayOrderId.orderIdForOnetime = result.order_id_for_onetime;
    }
    if (result.order_id_for_recurring) {
      razorpayOrderId.orderIdForRecurring = result.order_id_for_recurring;
    }
  };

  // Execute the operation with key refresh capability
  const result = await executeWithKeyRefresh({
    keyStateObj,
    createPromiseArray,
    processPromise,
    onSuccess,
    validateResult: (res) => res.status_code === 200
  });

  console.log("Result of askBackendForOrderId:", result.success);
  return result.success;
}

// Helper validation function
function validateResponseInJSONTOBeMyType(data: any): boolean {
  try {
    if (typeof data !== 'object' || data === null) return false;

    if (typeof data.order_id_for_recurring !== 'string') return false;
    console.log("Validating order_id_for_recurring", data.plan_type);

    if (typeof data.message !== 'string') return false;
    console.log("Validating message", data.plan_type);

    if (typeof data.order_id_for_onetime !== 'string') return false;
    console.log("Validating order_id_for_onetime", data.plan_type);

    if (typeof data.status_code !== 'number') return false;
    console.log("Validation complete");

    return true;
  } catch (error) {
    console.log("Error during JSON validation:", error);
    return false;
  }
}

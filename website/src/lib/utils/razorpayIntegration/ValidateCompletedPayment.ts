import { z } from "zod";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";
import { didUserSelectOneTimePayment } from "$lib/sharedState/didUserSeletctOneTimePayment.svelte";
import { KeyUpdate } from "../updateKey";
import { keyUpdatedState } from "$lib/sharedState/updatedKeyReceived.svelte";
import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { executeWithKeyRefresh } from "../ApiReqHelper/KeyRefreshhandler";

const RazorpayPaymentSchema = z.object({
  razorpay_payment_id: z.string(),
  razorpay_order_id: z.string(),
  razorpay_signature: z.string()
});

const responseFormApiCall = z.object({
  status_code: z.number(),
  success: z.boolean(),
  message: z.string(),
  new_key: z.string().optional() // Add the new_key field that's returned from the server

});

type RazorpayPayment = z.infer<typeof RazorpayPaymentSchema>;
type ResponseFormApiCall = z.infer<typeof responseFormApiCall>;

interface ValidationResponseType {
  status_code: number,
  success: boolean,
  new_key?: string,
  message: string
}

interface RequestVerifyPaymentSignature {
  user_key: string;
  razorpay_payment_id: string;
  razorpay_order_id: string;
  razorpay_signature: string;
  email: string;
  did_user_selected_one_time_payment: boolean;
}

/** function validated the payment completed and will save the new key in the local storage and in the state and also will send it to the chrome extension */
export async function validateCompletedPayment(responseFromRazorpay: unknown, userKey: string, email: string, didUserSelectOneTimePaymentMethod: boolean): Promise<{ status_code: number, success: boolean, message: string, responseFromApiCall: null | ResponseFormApiCall }> {
  console.log(`\n\n++++ user selected one time payemnt, value change by me is  ${didUserSelectOneTimePayment.valueChangedByMe} and the value is ${didUserSelectOneTimePayment.didUserSelectOneTimePayment}+++++++\n\n`);
  console.log(`did the user selected one time payment (as a param) ->`, didUserSelectOneTimePayment);

  try {
    // Validate the incoming Razorpay response
    const validationResult = RazorpayPaymentSchema.safeParse(responseFromRazorpay);
    if (!validationResult.success) {
      console.error("Invalid Razorpay response format:", validationResult.error);
      return {
        status_code: 400,
        success: false,
        responseFromApiCall: null,
        message: "Invalid payment response format"
      };
    }

    const validatedPayment = validationResult.data;

    // Construct the request body according to the Go API's expected format
    const requestBody: RequestVerifyPaymentSignature = {
      user_key: userKey,
      razorpay_payment_id: validatedPayment.razorpay_payment_id,
      razorpay_order_id: validatedPayment.razorpay_order_id,
      razorpay_signature: validatedPayment.razorpay_signature,
      email: email,
      did_user_selected_one_time_payment: didUserSelectOneTimePaymentMethod
    };

    // Create an AsyncRequestQueue instance
    const asyncReqQueue = new AsyncRequestQueue<Response, ResponseFormApiCall>(10);

    // Prepare the request to validate the payment
    const promiseArray = [
      fetch('/api/validatePayment', {
        headers: { 'Content-Type': 'application/json' },
        method: "POST",
        body: JSON.stringify(requestBody)
      })
    ];


    // makign sure the key will be refetched if there is 426

    let res = await executeWithKeyRefresh<Response, ResponseFormApiCall>(keyFromChromeExtensionState, asyncReqQueue, processIndividualPromise, promiseArray)

    if (res.error !== null || res.success === false || res.result === null) {
      return {
        status_code: 500,
        success: false,
        responseFromApiCall: res.result,
        message: res.result?.message || "Error validating payment"
      }
    }

    // Process the request using the AsyncRequestQueue
    // const result = await asyncReqQueue.process(
    //   (promiseToProcess) => processIndividualPromise<ValidationResponseType>(promiseToProcess), promiseArray
    // );

    // Check if there was an error in the response
    if (res.error !== null) {
      console.log(`There is an error in the result: ${res.error}, and the message from the server is: ${res.result}`);
      return {
        status_code: 500,
        success: false,
        responseFromApiCall: res.result,
        message: res.result?.message || "Error validating payment"
      };
    }
    let newKey = res.result?.new_key
    console.log(`the new key returned is ${newKey}`);

    if (newKey === "" || newKey === null || newKey === undefined) {
      return {
        success: false,
        message: "new key is not there",
        responseFromApiCall: res.result,
        status_code: 200
      }
    }

    // let updateKeyClass = new KeyUpdate
    // let err= updateKeyClass.UpdateKey(newKey, true)
    //   console.log(`the error in saving the key to the storage is ${err}`);
    // if (err !== null){
    //   return {
    //     success: false,
    //     message: "there is a error in saving the key to the storage"
    //   }
    // }


    // now save the key to the global shared state as the watch class will watch it and write it to the storage automatically
    keyFromChromeExtensionState.isPaidUser = true
    console.log(`about to set the new that that we got form the backend`);

    keyFromChromeExtensionState.key = newKey

    keyUpdatedState.newKeyReceived = true// now send it to the chrome extension

    // here make a global shared state that has int and boolean to it, when we get the new key we will set the state to true and usign effect we will 
    // start to send the new key to the chrome extension

    // Return the successful response
    return {
      success: true,
      message: "success",
      responseFromApiCall: res.result,
      status_code: 200
    }
  } catch (error) {
    console.error("Error validating payment:", error);
    return {
      status_code: 500,
      responseFromApiCall: null,
      success: false,
      message: error instanceof Error ? error.message : "Unknown error occurred"
    };
  }
}

async function processIndividualPromise<T>(resp1: Promise<Response>): Promise<T> {
  try {
    const resp = await resp1;

    if (!resp.ok) {
      console.error("The response is not ok from the backend:", resp);
      throw new Error("The response is not ok");
    }

    // Parse the JSON response
    const responseData = await resp.json();

    // Validate the response using Zod

    const validationResult = responseFormApiCall.safeParse(responseData);
    console.log(`validating the json schema's result (in validating the payment) is ${validationResult.success
      } and the result is ${JSON.stringify(validationResult)}`);

    if (!validationResult.success) {
      console.error("Invalid response format:", validationResult.error);
      throw new Error("Invalid response format from server");
    }

    // Return the validated response
    return responseData as T;

  } catch (error) {
    // Re-throw to be caught by AsyncRequestQueue's error handler
    throw error;
  }
}

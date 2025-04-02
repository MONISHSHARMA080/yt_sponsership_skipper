import { z } from "zod";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";
import { didUserSelectOneTimePayment } from "$lib/sharedState/didUserSeletctOneTimePayment.svelte";
import { KeyUpdate } from "../updateKey";
import { keyUpdatedState } from "$lib/sharedState/updatedKeyReceived.svelte";
import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";

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
export async function validateCompletedPayment(
  responseFromRazorpay: unknown,
  userKey: string,
  email: string,
  didUserSelectOneTimePaymentMethod: boolean
) {
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

    // Process the request using the AsyncRequestQueue
    const result = await asyncReqQueue.process(promiseArray,
      (promiseToProcess) => processIndividualPromise<ValidationResponseType>(promiseToProcess)
    );

    // Check if there was an error in the response
    if (result[0].error !== null) {
      console.log(`There is an error in the result: ${result[0].error}, and the message from the server is: ${result[0].result?.message}`);
      return {
        status_code: 500,
        success: false,
        message: result[0].result?.message || "Error validating payment"
      };
    }
    let newKey = result[0].result?.new_key
    console.log(`the new key returned is ${newKey}`);

    if (newKey === "" || newKey === null || newKey === undefined) {
      return {
        success: false,
        message: "new key is not there"
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
    return result[0].result;

  } catch (error) {
    console.error("Error validating payment:", error);
    return {
      status_code: 500,
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
    const validationSchema = z.object({
      status_code: z.number(),
      new_key: z.string().optional(),
      success: z.boolean(),
      message: z.string(),
    });

    const validationResult = validationSchema.safeParse(responseData);
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

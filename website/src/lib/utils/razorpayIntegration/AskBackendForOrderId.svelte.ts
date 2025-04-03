import { razorpayOrderId } from "$lib/sharedState/razorPayKey.svelte";
import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { executeWithKeyRefresh } from "../ApiReqHelper/KeyRefreshhandler";
import { getNewKeyAndSaveItToTheGlobalState } from "../ApiReqHelper/updateTheKeyOnItsDeperaction.svelte";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";

interface requestType {
  user_key: string
}
interface responseType {
  order_id_for_recurring: string,
  order_id_for_onetime: string,
  message: string,
  status_code: number,
  // plan_type:"onetime"|"recurringpayment"
}

/**
 * this func is designed to run when the state changes(derived/effect), run the function and ask the bakend for the order Id and store it in the global state 
*/
export async function askBackendForOrderId(keyStateObj: keyStateObject,): Promise<boolean> {
  console.log(`1111--`);



  try {
    // razorpayOrderId.fetchingStatus = "fetching"
    if (!keyStateObj.isValidatedThroughBackend || keyStateObj.key === "" || keyStateObj.key === null) {
      console.log(`the key is not validated or is empty or null we are returning ->`, keyStateObj.key);

      return false
    }
    // if we are validated 
    let asyncReqQueue = new AsyncRequestQueue<Response, responseType>(10)

    // let reqBody: requestType = { user_key: keyStateObj.key }
    let promiseArray = factoryPromiseArrayForMakeAPayment(keyStateObj.key)
    let res = await executeWithKeyRefresh<Response, responseType>(keyStateObj, asyncReqQueue, processIndividualPromise, promiseArray)
    console.log(` the rsult from the executeWithKeyRefresh is --- ${JSON.stringify(res)}`);
    console.log(`the result.error is ${res.error}`);

    if (res.error !== null || res.success === false || res.result === null) {
      return false
    }
    let result = res.result
    // usually it will be as if it was not it will blow/thorw in the async request queue
    if (result.status_code !== 200) { return false }

    let oneTime = res.result.order_id_for_onetime
    let recurring = res.result.order_id_for_recurring
    if (oneTime === "" || recurring === "") {
      console.log(`the order id is empty (ont of them or both)  onetime and recurring is ->${oneTime} -- ${recurring}`);
      return false
    }

    razorpayOrderId.orderIdForOnetime = oneTime
    razorpayOrderId.orderIdForOnetime = recurring

    return true


    // first get/add a unifying return type  there such that(use generics) I get the razorpayOrderId here
    // still have to return true or false

  } catch (error) {
    console.log(` error in the executeWithKeyRefresh.ts is  ->`, error);
    return false
  }









  // let result = await asyncReqQueue.process((promiseToProcess) => processIndividualPromise<responseType>(promiseToProcess), promiseArray)
  //
  // // for (let index = 0; index < result.length; index++) {
  // let res = result[0]
  // if (res.error !== null || res.result === null) {
  //   console.log(`there is a error in the result array at ${0} and is ->`, res.error, "\n and the message form the server is ->", res.result?.message);
  //   // razorpayOrderId.fetchingStatus = "error"
  //   if (res.ifErrorThenOriginalPromiseResponse !== null) {
  //
  //     let response = await res.ifErrorThenOriginalPromiseResponse
  //     console.log(`the respoonse from awaiting the result in error -+ is ${JSON.stringify(response)}`);
  //     console.log(`the error response.status is ${response.status}`);
  //     if (response.status === 426) {
  //       // the key expired get a new one
  //       let res = await getNewKeyAndSaveItToTheGlobalState(keyStateObj.key)
  //       if (res !== null) {
  //         console.log(`there is a error in gettign the new key can't re run it error ->${res}`);
  //         return false
  //       } else {
  //         console.log(`the new key is here and we are re runnign the makeAPayment again `);
  //         // some how re run the function again 
  //         // abstract the above into seperate functions  and re run them again
  //         let newKey = keyFromChromeExtensionState.key
  //         if (newKey === "" || newKey === null) {
  //           console.log(`the new key is null or "" `);
  //           return false
  //         }
  //         let promiseArray = factoryPromiseArrayForMakeAPayment(newKey)
  //         let result = await asyncReqQueue.process((promiseToProcess) => processIndividualPromise<responseType>(promiseToProcess), promiseArray)
  //         let res = result[0]
  //         if (res.error !== null || res.result === null) {
  //           console.log(`there is a error in the result array after refetching the key at 0 and is ->`, res.error, "\n and the message form the server is ->", res.result?.message);
  //           return false
  //         }
  //         if (res.result.status_code !== 200) {
  //           console.error(`the response after gettign the new key and re running the fetch again is not 200 `);
  //           return false
  //         } else {
  //           console.log(`resppnse after refetching the key and re runnign the fetch is  200 and`);
  //
  //           return true
  //         }
  //       }
  //     }
  //   }
  //   return false
  // }
  //
  // console.log(`about to store the response in razorpay storage and it is -> ${res.result?.order_id_for_recurring} and the one time is ${res.result?.order_id_for_onetime}`);
  // if (res.result.status_code !== 200) {
  //
  //   // razorpayOrderId.fetchingStatus = "error"
  //   return false
  // }
  // //  if (res.result?.order_id_for_onetime ){
  // //       razorpayOrderId.orderIdForOnetime = res.result.order_id_for_onetime
  // //  }else if (res.result?.order_id_for_recurring ){
  // //       razorpayOrderId.orderIdForRecurring = res.result.order_id_for_recurring
  // //  }
  // // razorpayOrderId.fetchingStatus = "success"
  // if (res.result?.order_id_for_onetime) {
  //   razorpayOrderId.orderIdForOnetime = res.result.order_id_for_onetime;
  // }
  // if (res.result?.order_id_for_recurring) {
  //   razorpayOrderId.orderIdForRecurring = res.result.order_id_for_recurring;
  // }
  // // }
  // console.log("\n\n\n\n\n result array is ->", result, "\n\n\n\n\n and the razor pay state is 1 time ", razorpayOrderId.orderIdForOnetime, " and recurr ", razorpayOrderId.orderIdForRecurring);
  //
  // return true

}

// throw insde the catch cause I want my error to be caught by the async req queue 
async function processIndividualPromise<T>(resp1: Promise<Response>): Promise<T> {
  try {
    let resp = await resp1;

    console.log(`response come form the makeAPayment \n\n`);

    if (!resp.ok) {
      console.error("the response is not ok form the backend and we are throwing a error in the async request queue process indivudual func->", resp);
      throw new Error("the response is not ok bruh");
    }

    // Parse the JSON first
    const responseInJson = await resp.json();

    console.log("the valid JSON is ->", responseInJson);
    // Then validate it
    if (validateResponseInJSONTOBeMyType(responseInJson)) {
      // JSON is valid   
      let respReceived: T = responseInJson;
      return respReceived;
    } else {
      throw new Error("the json is not even valid or of the same type in the response");
    }
  } catch (error) {
    // This will catch any errors in the above code, including the parsing and validation errors
    throw error; // Re-throw to be caught by AsyncRequestQueue's error handler
  }
}



function validateResponseInJSONTOBeMyType(data: any): boolean {
  // For responseType interface specifically:
  try {
    if (typeof data !== 'object' || data === null) return false;
    // Check if the data has all required fields with correct types
    if (typeof data.order_id_for_recurring !== 'string') return false;
    console.log("---++----+++0000", data.plan_type);
    if (typeof data.message !== 'string') return false;
    console.log("-----message is string------", data.plan_type);
    if (typeof data.order_id_for_onetime !== 'string') return false;
    console.log("---++----+++0000", data.plan_type);

    if (typeof data.status_code !== 'number') return false;
    console.log("---++----+++0000----");
    console.log(`the data plan type ->${data.plan_type}`);


    // if(data.plan_type !==  "onetime" &&  data.plan_type !== "recurringpayment") return false
    console.log("\n\n we are returning true\n\n\n");

    return true;
  } catch (error) {
    console.log("error during checking json ->", error);
    return false
  }
}

function factoryPromiseArrayForMakeAPayment(userKey: string) {
  let reqBody: requestType = { user_key: userKey }
  let promiseArray = [
    fetch('/api/makeAPayment', {
      headers: { 'Content-Type': 'application/json' }
      , method: "POST",
      body: JSON.stringify(reqBody)
    }),
  ]
  return promiseArray
}

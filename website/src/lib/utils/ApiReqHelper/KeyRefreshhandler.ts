import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { getNewKeyAndSaveItToTheGlobalState } from "../ApiReqHelper/updateTheKeyOnItsDeperaction.svelte";
import { AsyncRequestQueue, type funcToProcessIndividualPromise } from "../newAsyncRequestQueue";

type ProcessPromiseFunction<T> = (promise: Promise<Response>) => Promise<T>;

interface RefreshKeyHandlerOptions<T> {
  keyStateObj: keyStateObject;
  promiseArray: Promise<Response>[];
  processPromise: ProcessPromiseFunction<T>;
  onSuccess?: (result: T) => void;
  validateResult?: (result: T) => boolean;
  maxRetries?: number;
}


/**
 * note: the AsyncRequestQueue class should be loaded like the promiseArray should be there in the class
  */
export async function executeWithKeyRefresh<Response, R>(
  keyStateObj: keyStateObject,
  asyncQueue: AsyncRequestQueue<Response, R>,
  funcToProcessIndividualPromise: funcToProcessIndividualPromise<Response, R>,
  promiseArray: Promise<Response>[],
  keyformTheChromeExtension?: string,
  promiseArray2?: Promise<Response>[]
): Promise<{ success: boolean, result: R | null, error: Error | null }> {
  console.log(`in the executeWithKeyRefresh `);

  let oldKey = keyformTheChromeExtension ? keyformTheChromeExtension : keyStateObj.key;
  if (oldKey === null || oldKey === "") {
    console.log(`the oldKey11 is not there and returning`);
    return { success: false, result: null, error: new Error("Key is empty or null") };
  }

  console.log(`+1=1=1=1=1=1=1=1=1=1=1=1=1=1=1`);

  try {
    console.log(`11111111111111111111`);
    let result = await asyncQueue.process(funcToProcessIndividualPromise, promiseArray);
    console.log(`22222222222222222`);
    console.log(`the result index length is ${result.length} and the index in promise queue is ${asyncQueue.promiseQueueSubmittedByUser.length} and in the promise array in the func is ${promiseArray.length}`);

    // Check if result exists and has the first element
    if (!result || result.length === 0) {
      return { success: false, result: null, error: new Error("No results returned from process") };
    }

    console.log(`the res[0] is ${JSON.stringify(result[0])}`);

    // Process the first result (assuming that's what you want)

    const value = result[0];
    console.log(`the shape of the result object -> ${JSON.stringify(value)}`);
    console.log(`is the error!== null ${value.error !== null}`);


    if (value.error !== null || value.result === null) {
      console.log(`\n\\n\\n\n\n\n\n\n\-------there is a error in the result ->${value.error} \n\n\n\n\n\\\n\n\n\n\n\n\n`);

      // Check if we have an original response to work with
      if (value.ifErrorThenOriginalPromiseResponse === null) {
        return {
          success: false,
          result: null,
          error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and we did not got the response form the api call`)
        };
      }

      // Process the original response to check for 426 error
      let responseWeGotInReq = await value.ifErrorThenOriginalPromiseResponse;
      console.log(`trying to access property status on responseWeGotInReq.status`);

      //@ts-ignore
      if (responseWeGotInReq.status !== 426) {
        return {
          success: false,
          result: null,
          error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and the response is neither 200 nor 426`)
        };
      } else {
        //@ts-ignore
        console.log(`we got the response form the server to be 426 and we are in process of updating the keys, status is ${responseWeGotInReq.status}`);

        // Try to get a new key
        let responseFromKeyUpdate = await getNewKeyAndSaveItToTheGlobalState(oldKey);
        if (responseFromKeyUpdate !== null) {
          return {
            success: false,
            result: null,
            error: responseFromKeyUpdate
          };
        }
        // now lets re run the same fetch again 
        // console.log(`>> fetching the old fetch again===`);
        // let promiseArray2ndTime = promiseArray2 ? promiseArray2 : promiseArray
        // console.log(`the length of promise array 2 is ${promiseArray2?.length}`);
        // promiseArray2ndTime = [...promiseArray]


        let result = await asyncQueue.process(funcToProcessIndividualPromise, promiseArray);

        const value = result[0];
        console.log(`the shape of the result object -> ${JSON.stringify(value)}`);
        console.log(`the result form the re ran fetch is ${JSON.stringify(result[0])}`);
        if (value.error !== null) {
          console.log(`there is a error in re ran fetch and it is ->${value.error}`);
          let resp = await value.ifErrorThenOriginalPromiseResponse
          if (value.ifErrorThenOriginalPromiseResponse) {
            // there is a error in fetching this time too so lets return
            console.log(`(value.ifErrorThenOriginalPromiseResponse) is there `);
            const response: Response = await value.ifErrorThenOriginalPromiseResponse;
            try {
              //@ts-ignore
              console.log(`Response status: ${response.status}`);
              //@ts-ignore
              console.log(`Response ok: ${response.ok}`);
              console.log(`Response is String is ${JSON.stringify(response)}`);
            } catch (error) {
              console.log(` error is in gettign the status of ifErrorThenOriginalPromiseResponse field ->`, error);
            }
            return {
              success: false,
              result: null,
              error: new Error(`there is a error in getting the request after key fetch 2nd time so we are returning`)
            };
          }
        }
        // Key update success, return success indicator
        return { success: true, result: null, error: null };
      }
    } else {
      console.log(`the error is probably not there `);
      return { success: true, result: value.result, error: null };
    }
  } catch (error) {
    console.log(` error is ->`, error);
    return {
      success: false,
      result: null,
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}







// ------ old one ------



// Static variables to handle multiple concurrent refresh attempts
// let refreshInProgress = false;
// let lastRefreshTime = 0;
// const REFRESH_COOLDOWN = 0; // 5 seconds cooldown between refresh attempts
//
// /**
//  * Generic handler that performs API calls with automatic key refresh capability
//  * when the server responds with a 426 status (key expired)
//  * @param options Configuration options for the handler
//  * @param validateResult is a function that is used to tell us when is the response ok, and we return
//  * @returns Object containing success flag, result, and error
//  */
// export async function executeWithKeyRefresh<T>({
//   keyStateObj,
//   promiseArray,
//   processPromise,
//   onSuccess,
//   validateResult = (result) => (result as any).status_code === 200,
//   maxRetries = 1
// }: RefreshKeyHandlerOptions<T>): Promise<{
//   success: boolean;
//   result: T | null;
//   error: Error | null;
// }> {
//   const maxRetry = 1
//   let oldKey = keyStateObj.key
//   try {
//     // implementing one myself
//     // 1st fetch the thing with the help of the promise array
//     if (oldKey === "" || oldKey === null) {
//       return { success: false, result: null, error: new Error("the key is not there") }
//     }
//     const asyncReqQueue = new AsyncRequestQueue<Response, T>(10);
//
//     for (let index = 0; index < maxRetries; index++) {
//       console.log(`on the loop iteration ${index} and max reties are ${maxRetries}`);
//
//       let result = await asyncReqQueue.process(processPromise, promiseArray)
//       result.forEach(async (value, index) => {
//         if (value.error !== null || value.result === null) {
//           console.log(`\n\\n\\n\n\n\n\n\n\-------there is a error in the result ->${value.error} \n\n\n\n\n\\\n\n\n\n\n\n\n`);
//           // return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null`) }
//           // now we need to check if the error is cause of 426 or a general error
//           if (value.ifErrorThenOriginalPromiseResponse === null) {
//             return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and we did not got the response form the api call`) }
//           }
//           let responseWeGotInReq = await value.ifErrorThenOriginalPromiseResponse
//           if (responseWeGotInReq.status !== 426) {
//             return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and the respons is neither 200 nor 426`) }
//           } else {
//             console.log(`we got the response form the servet to be 426 and we are in processs of updating the keys, status is ${responseWeGotInReq.status}`);
//             let responseFromKeyUpdate = await getNewKeyAndSaveItToTheGlobalState(oldKey)
//             if (responseFromKeyUpdate !== null) {
//               return { success: false, result: null, error: responseFromKeyUpdate }
//             }
//             // key update success refetch the old query
//             // or we can get form the chrome extension
//           }
//         }
//       })
//     }
//
//     // remove this
//     let error = new Error()
//     return {
//       success: false,
//       result: null,
//       error: error instanceof Error ? error : new Error(String(error))
//     }
//
//   } catch (error) {
//     console.log(`Unexpected error in executeWithKeyRefresh: ${error}`);
//     return {
//       success: false,
//       result: null,
//       error: error instanceof Error ? error : new Error(String(error))
//     };
//   }
// }





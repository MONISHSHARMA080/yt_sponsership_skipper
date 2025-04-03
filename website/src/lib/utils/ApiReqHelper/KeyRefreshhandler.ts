import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { getNewKeyAndSaveItToTheGlobalState } from "../ApiReqHelper/updateTheKeyOnItsDeperaction.svelte";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";

type ProcessPromiseFunction<T> = (promise: Promise<Response>) => Promise<T>;

interface RefreshKeyHandlerOptions<T> {
  keyStateObj: keyStateObject;
  promiseArray: Promise<Response>[];
  processPromise: ProcessPromiseFunction<T>;
  onSuccess?: (result: T) => void;
  validateResult?: (result: T) => boolean;
  maxRetries?: number;
}

// Static variables to handle multiple concurrent refresh attempts
let refreshInProgress = false;
let lastRefreshTime = 0;
const REFRESH_COOLDOWN = 0; // 5 seconds cooldown between refresh attempts

/**
 * Generic handler that performs API calls with automatic key refresh capability
 * when the server responds with a 426 status (key expired)
 * @param options Configuration options for the handler
 * @param validateResult is a function that is used to tell us when is the response ok, and we return
 * @returns Object containing success flag, result, and error
 */
export async function executeWithKeyRefresh<T>({
  keyStateObj,
  promiseArray,
  processPromise,
  onSuccess,
  validateResult = (result) => (result as any).status_code === 200,
  maxRetries = 1
}: RefreshKeyHandlerOptions<T>): Promise<{
  success: boolean;
  result: T | null;
  error: Error | null;
}> {
  const maxRetry = 1
  let oldKey = keyStateObj.key
  try {
    // implementing one myself
    // 1st fetch the thing with the help of the promise array
    if (oldKey === "" || oldKey === null) {
      return { success: false, result: null, error: new Error("the key is not there") }
    }
    const asyncReqQueue = new AsyncRequestQueue<Response, T>(10);

    for (let index = 0; index < maxRetries; index++) {
      console.log(`on the loop iteration ${index} and max reties are ${maxRetries}`);

      let result = await asyncReqQueue.process(promiseArray, processPromise)
      result.forEach(async (value, index) => {
        if (value.error !== null || value.result === null) {
          console.log(`\n\\n\\n\n\n\n\n\n\-------there is a error in the result ->${value.error} \n\n\n\n\n\\\n\n\n\n\n\n\n`);
          // return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null`) }
          // now we need to check if the error is cause of 426 or a general error
          if (value.ifErrorThenOriginalPromiseResponse === null) {
            return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and we did not got the response form the api call`) }
          }
          let responseWeGotInReq = await value.ifErrorThenOriginalPromiseResponse
          if (responseWeGotInReq.status !== 426) {
            return { success: false, result: null, error: value.error instanceof Error ? value.error : new Error(`the result is equal to null and the respons is neither 200 nor 426`) }
          } else {
            console.log(`we got the response form the servet to be 426 and we are in processs of updating the keys, status is ${responseWeGotInReq.status}`);
            let responseFromKeyUpdate = await getNewKeyAndSaveItToTheGlobalState(oldKey)
            if (responseFromKeyUpdate !== null) {
              return { success: false, result: null, error: responseFromKeyUpdate }
            }
            // key update success refetch the old query
            // or we can get form the chrome extension


          }

        }

      })

    }











    // remove this
    let error = new Error()
    return {
      success: false,
      result: null,
      error: error instanceof Error ? error : new Error(String(error))
    }

  } catch (error) {
    console.log(`Unexpected error in executeWithKeyRefresh: ${error}`);
    return {
      success: false,
      result: null,
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

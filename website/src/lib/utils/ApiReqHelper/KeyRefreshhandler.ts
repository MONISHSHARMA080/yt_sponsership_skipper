import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { getNewKeyAndSaveItToTheGlobalState } from "../ApiReqHelper/updateTheKeyOnItsDeperaction.svelte";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";

type ApiCallFunction<T> = (key: string) => Promise<T[]>;
type ProcessPromiseFunction<T> = (promise: Promise<Response>) => Promise<T>;

interface RefreshKeyHandlerOptions<T> {
  keyStateObj: keyStateObject;
  createPromiseArray: (key: string) => Promise<Response>[];
  processPromise: ProcessPromiseFunction<T>;
  onSuccess?: (result: T) => void;
  validateResult?: (result: T) => boolean;
}

/**
 * Generic handler that performs API calls with automatic key refresh capability
 * when the server responds with a 426 status (key expired)
 */
export async function executeWithKeyRefresh<T>({
  keyStateObj,
  createPromiseArray,
  processPromise,
  onSuccess,
  validateResult = (result) => (result as any).status_code === 200
}: RefreshKeyHandlerOptions<T>): Promise<{
  success: boolean;
  result: T | null;
  error: Error | null;
}> {
  try {
    // Validate key state
    if (!keyStateObj.isValidatedThroughBackend || !keyStateObj.key) {
      console.log(`The key is not validated or is empty: ${keyStateObj.key}`);
      return { success: false, result: null, error: new Error("Invalid key state") };
    }

    // Create async request queue
    const asyncReqQueue = new AsyncRequestQueue<Response, T>(10);

    // First attempt with current key
    let promiseArray = createPromiseArray(keyStateObj.key);
    let results = await asyncReqQueue.process(promiseArray, processPromise);
    let responseObj = results[0];

    // Check for errors - specifically handle 426 (expired key)
    if (responseObj.error !== null || responseObj.result === null) {
      console.log(`Error in result: ${responseObj.error}, server message: ${responseObj?.result}`);

      // Handle key expiration (426)
      if (responseObj.ifErrorThenOriginalPromiseResponse !== null) {
        const response = await responseObj.ifErrorThenOriginalPromiseResponse;
        console.log(`Response from awaiting the result in error: ${JSON.stringify(response)}`);

        if (response.status === 426) {
          // Key expired - get a new one
          const refreshError = await getNewKeyAndSaveItToTheGlobalState(keyStateObj.key);

          if (refreshError !== null) {
            console.log(`Error getting new key: ${refreshError}`);
            return { success: false, result: null, error: new Error(`Failed to refresh key: ${refreshError}`) };
          }

          // Retry with new key
          console.log('New key obtained, retrying API call');
          const newKey = keyFromChromeExtensionState.key;

          if (!newKey) {
            console.log('New key is null or empty');
            return { success: false, result: null, error: new Error("New key is empty") };
          }

          // Retry the operation with new key
          promiseArray = createPromiseArray(newKey);
          results = await asyncReqQueue.process(promiseArray, processPromise);
          responseObj = results[0];

          // Check if retry was successful
          if (responseObj.error !== null || responseObj.result === null) {
            console.log(`Error in result after key refresh: ${responseObj.error}`);
            return { success: false, result: null, error: responseObj.error };
          }
        }
      }

      // If we get here with error still present, return failure
      if (responseObj.error !== null || responseObj.result === null) {
        return { success: false, result: null, error: responseObj.error };
      }
    }

    // Validate result
    if (!validateResult(responseObj.result)) {
      console.log(`Response validation failed: ${JSON.stringify(responseObj.result)}`);
      return {
        success: false,
        result: responseObj.result,
        error: new Error(`Invalid response: ${(responseObj.result as any).status_code}`)
      };
    }

    // Process success callback if provided
    if (onSuccess) {
      onSuccess(responseObj.result);
    }

    return { success: true, result: responseObj.result, error: null };
  } catch (error) {
    console.log(`Unexpected error in executeWithKeyRefresh: ${error}`);
    return {
      success: false,
      result: null,
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

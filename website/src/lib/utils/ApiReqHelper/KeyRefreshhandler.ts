import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { getNewKeyAndSaveItToTheGlobalState } from "../ApiReqHelper/updateTheKeyOnItsDeperaction.svelte";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";

type ProcessPromiseFunction<T> = (promise: Promise<Response>) => Promise<T>;

interface RefreshKeyHandlerOptions<T> {
  keyStateObj: keyStateObject;
  createPromiseArray: (key: string) => Promise<Response>[];
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
 * @returns Object containing success flag, result, and error
 */
export async function executeWithKeyRefresh<T>({
  keyStateObj,
  createPromiseArray,
  processPromise,
  onSuccess,
  validateResult = (result) => (result as any).status_code === 200,
  maxRetries = 1
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

    // Track retry count
    let retryCount = 0;
    let currentKey = keyStateObj.key;
    let responseObj;

    // Loop with retry limit
    while (retryCount <= maxRetries) {
      console.log(`Attempt ${retryCount + 1} of ${maxRetries + 1} with key: ${currentKey.substring(0, 5)}...`);

      // Create promise array with current key
      const promiseArray = createPromiseArray(currentKey);
      const results = await asyncReqQueue.process(promiseArray, processPromise);
      responseObj = results[0];

      // Check for success
      if (responseObj.error === null && responseObj.result !== null) {
        // We have a result - validate it
        if (validateResult(responseObj.result)) {
          // Valid successful result
          if (onSuccess) {
            onSuccess(responseObj.result);
          }
          return { success: true, result: responseObj.result, error: null };
        } else {
          // Invalid result
          console.log(`Response validation failed: ${JSON.stringify(responseObj.result)}`);
          return {
            success: false,
            result: responseObj.result,
            error: new Error(`Invalid response: ${(responseObj.result as any).status_code}`)
          };
        }
      }

      // We have an error - check if it's a key expiration (426)
      if (responseObj.ifErrorThenOriginalPromiseResponse !== null && retryCount < maxRetries) {
        const response = await responseObj.ifErrorThenOriginalPromiseResponse;

        if (response.status === 426) {
          // Key expired - get a new one, but first check if another refresh is already in progress
          console.log('Key expired (426), attempting to refresh key');

          // Prevent multiple simultaneous refresh attempts and respect cooldown
          const now = Date.now();
          if (refreshInProgress) {
            console.log('Another key refresh is already in progress. Waiting...');
            // Wait a bit and then check if the global key has been updated
            await new Promise(resolve => setTimeout(resolve, 1000));

            // Check if the global key has been updated by another instance
            if (keyFromChromeExtensionState.key !== currentKey && keyFromChromeExtensionState.key !== null) {
              currentKey = keyFromChromeExtensionState.key;
              retryCount++;
              console.log(`Key was refreshed by another process, will retry with new key (${retryCount}/${maxRetries})`);
              continue;
            }
          } else if (now - lastRefreshTime < REFRESH_COOLDOWN) {
            console.log(`Refresh attempted too soon after previous refresh. Waiting for cooldown...`);
            await new Promise(resolve => setTimeout(resolve, 1000)); // Wait a bit

            // Check if the global key has changed during the wait
            if (keyFromChromeExtensionState.key !== currentKey && keyFromChromeExtensionState.key !== null) {
              currentKey = keyFromChromeExtensionState.key;
              retryCount++;
              console.log(`Key was refreshed while waiting, will retry with new key (${retryCount}/${maxRetries})`);
              continue;
            }
          }

          // Set the refreshInProgress flag
          refreshInProgress = true;

          try {
            const refreshError = await getNewKeyAndSaveItToTheGlobalState(currentKey);

            if (refreshError !== null) {
              console.log(`Error getting new key: ${refreshError}`);
              return {
                success: false,
                result: null,
                error: new Error(`Failed to refresh key: ${refreshError}`)
              };
            }

            // Get the new key
            const newKey = keyFromChromeExtensionState.key;

            // Update refresh time
            lastRefreshTime = Date.now();

            if (!newKey) {
              console.log('New key is null or empty');
              return { success: false, result: null, error: new Error("New key is empty") };
            }

            // Update current key for next attempt
            currentKey = newKey;
            retryCount++;
            console.log(`Key refreshed, will retry (${retryCount}/${maxRetries})`);
          } finally {
            // Reset the refreshInProgress flag
            refreshInProgress = false;
          }

          continue; // Try again with new key
        }
      }

      // If we get here, we either have a non-426 error or we've hit the retry limit
      break;
    }

    // If we get here, we've either exhausted retries or had a non-recoverable error
    return {
      success: false,
      result: null,
      error: responseObj?.error || new Error("Maximum retry attempts reached")
    };

  } catch (error) {
    console.log(`Unexpected error in executeWithKeyRefresh: ${error}`);
    return {
      success: false,
      result: null,
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

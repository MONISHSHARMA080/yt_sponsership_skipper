import { PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH } from "$env/static/public";
import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { executeWithKeyRefresh } from "./ApiReqHelper/KeyRefreshhandler";
import { AsyncRequestQueue } from "./newAsyncRequestQueue";
// import { AsyncRequestQueue } from "./asyncRequestQueue";
import { z } from "zod";


interface ResponseData {
  message: string;
  status_code: number;
  success: boolean;
  encrypted_key: string;
  email: string;
  name: string;
  is_user_on_paid_tier: boolean
}

const ResponseDataSchema = z.object({
  message: z.string(),
  status_code: z.number(),
  success: z.boolean(),
  encrypted_key: z.string(),
  email: z.string(),
  name: z.string(),
  is_user_on_paid_tier: z.boolean()
});

interface apiResponseAndError {
  result: ResponseData | null,
  error: Error | null
}

export class checkIfKeyIsValidAndUpdateTheState {

  async seeIfKeyIsValid(key: string): Promise<apiResponseAndError> {
    console.log(`the key in the seeIfKeyIsValid method is ->${key}`);

    try {
      let asyncRequestQueue = new AsyncRequestQueue<Response, ResponseData>(100)
      let promiseQueue = [
        fetch(`/api/checkIfKeyIsValid`, {
          method: "POST",
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH
          },
          body: JSON.stringify({ key: key })
        })]
      let promiseQueue1 = [
        fetch(`/api/checkIfKeyIsValid`, {
          method: "POST",
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH
          },
          body: JSON.stringify({ key: "thisKeyIsFake" })
        })]
      let res = await executeWithKeyRefresh(keyFromChromeExtensionState, asyncRequestQueue, this.processIndividualPromise, promiseQueue, key, promiseQueue1)
      console.log(`--+++==>>the result form checking the key is valid is ${JSON.stringify(res)}`);
      if (res.error !== null || res.success === false || res.result === null || res.result === undefined) {
        console.log(`the error we got is ${res.error}`);
        return { result: null, error: res.error instanceof Error ? res.error : Error("error executeWithKeyRefresh func  on the /api/checkIfKeyIsValid route and it is ->" + res.error) }
      } else {
        return { result: res.result, error: null, }

      }

      // const asyncRequestQueue = new AsyncRequestQueue<ResponseData>(10)
      // asyncRequestQueue.addToQueue([
      // () => 
      //     .then(async (resp) => {
      //       if (!resp.ok) {
      //         throw new Error(`HTTP error! Status: ${resp.status} and body ->${resp.body}`);
      //       }
      //       return await resp.json();
      //     })
      //     .then(data => data as ResponseData)
      // ]);
      // let result = await asyncRequestQueue.processQueue()
      // console.log("the result form the result is ->", result[0]);
      // // updating the globalState and then the local storage
      // // this one is deperecated as now writing to shared state will write to the storage
      // this.updateGlobalStateAndWriteToTheStorageIfKeysAreValid(result[0], key)
      // // this.updateTheLocalStorage(keyFromChromeExtensionState)
      // return { result: result[0].result, error: result[0].error }
      //
      let error = ""
      return { result: null, error: error instanceof Error ? error : Error("error in checking if key is valid and updating it func and it is ->" + error) }
    } catch (error) {
      console.log("error occurred in the seeIfKeyISValidFunc ->", error)
      return { result: null, error: error instanceof Error ? error : Error("error in checking if key is valid and updating it func and it is ->" + error) }
    }
  }

  private updateGlobalStateAndWriteToTheStorageIfKeysAreValid(res: apiResponseAndError, key: string) {
    if (res.error !== null && res.result === null) {
      console.log(" the error form the backend is(and not updating the state) ->", res.error)
      console.log(" the result is  ", res.result)
      return
    }
    if (res.result === null) {
      return
    }
    // I do not want to write to the state multiple times, so that's why we are cloning it and will also then assing it 
    let cloneObjOfSharedState = Object.assign({}, keyFromChromeExtensionState)
    cloneObjOfSharedState.isValidatedThroughBackend = true
    cloneObjOfSharedState.email = res.result.email
    cloneObjOfSharedState.key = key
    cloneObjOfSharedState.isPaidUser = res.result.is_user_on_paid_tier
    Object.assign(keyFromChromeExtensionState, cloneObjOfSharedState)
    // keyFromChromeExtensionState.isValidatedThroughBackend = true
    // keyFromChromeExtensionState.email = res.result.email
    // keyFromChromeExtensionState.name = res.result.name
    // keyFromChromeExtensionState.key = key
    console.log("the key form the backend is ->", res.result.encrypted_key === key);
  }
  private async processIndividualPromise<T>(resp1: Promise<Response>): Promise<T> {
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
      if (this.validateResponseInJSONTOBeMyType(responseInJson)) {
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

  private validateResponseInJSONTOBeMyType(response: any): boolean {
    return ResponseDataSchema.safeParse(response).success
  }



}

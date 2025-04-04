
// --------

import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { z } from "zod";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";


/**
 * this function is to run when we have gotten a 403 status code, get a new key for now and when the user gets a new key it will
 * be override form next time with chrome extension one; 
 */

interface responseTypeForNewKey {
  new_encrypted_key: string,
  message: string,
  status_code: number,
}
const responseTypeForNewKeySchema = z.object({
  new_encrypted_key: z.string(),
  message: z.string(),
  status_code: z.number()
});


// this func is not in the executeWithKeyRefresh as it is used there 
async function getNewKey(oldKey: string): Promise<[string, Error | null]> {
  let asyncReqQueue = new AsyncRequestQueue<Response, responseTypeForNewKey>(10)
  let promiseArray = [
    fetch('/api/getNewKey', {
      headers: { 'Content-Type': 'application/json' }
      , method: "POST",
      body: JSON.stringify({ user_key: oldKey })
    }),
  ]
  let result = await asyncReqQueue.process((promiseToProcess) => processIndividualPromise<responseTypeForNewKey>(promiseToProcess), promiseArray)
  let res = result[0]
  if (res.error !== null || res.result === null) {
    console.log(`there is a error in getting new key from the backend ->${res.error}`);
    return ["", res.error instanceof Error ? res.error : new Error(`there is a error in getting the response or it might be null `)]
  }
  if (res.result.status_code !== 200) {
    console.log(`successfully got the response for the new key but it is not 200  the response is ${JSON.stringify(res.result)} `);
    return ["", new Error(`the response form server is not 200  here is the response ->${JSON.stringify(res.result)}`)]
  }
  return [res.result.new_encrypted_key, null]
}


/**will get the key and set it on the keyStateObject and writitng it to storage will eventually happen*/
export async function getNewKeyAndSaveItToTheGlobalState(oldKey: string): Promise<Error | null> {
  try {
    let [newKey, err] = await getNewKey(oldKey)

    if (err !== null || newKey === "") {
      console.log(`there is a error in getting new key ${err}`);
      return err instanceof Error ? err : new Error(`there is a error in getting new key or the key is empty, error is  ${err} and the key is ${newKey}`)
    }

    keyFromChromeExtensionState.key = newKey
    return null
  } catch (error) {
    console.log(`there is a error in try catch of the getNewKeyAndSaveItToTheGlobalState() and it is ${error}`);
    return new Error(`there is a error in getting new key or the key is empty, `)
  }
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

function validateResponseInJSONTOBeMyType(responseInJson: any): boolean {
  // Attempt to parse the JSON using the schema
  try {
    const result = responseTypeForNewKeySchema.safeParse(responseInJson);
    return result.success;
  } catch (error) {
    console.error("Validation error:", error);
    return false;
  }
}

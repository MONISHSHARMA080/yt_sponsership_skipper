import { PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH } from "$env/static/public";
import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { AsyncRequestQueue } from "./asyncRequestQueue";


interface  ResponseData{ 
  message: string;
  status_code: number;
  success: boolean;
  encrypted_key: string;
  email: string;
  name: string;
}
interface apiResponseAndError{
    result:ResponseData|null,
    error:Error|null
}

export class checkIfKeyIsValidAndUpdateTheState{

    async seeIfKeyIsValid(key:string):Promise<apiResponseAndError>{
        try {
       const asyncRequestQueue = new AsyncRequestQueue<ResponseData>(10)
        asyncRequestQueue.addToQueue([
            () => fetch(`/api/checkIfKeyIsValid`, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Origin':PUBLIC_BACKEND_URL_WITHOUT_BACKSLASH
            },
            body: JSON.stringify({ key: key })
            })
            .then(async (resp) => {
            if (!resp.ok) {
                throw new Error(`HTTP error! Status: ${resp.status} and body ->${resp.body}`);
            }
                return await resp.json();
            })
                .then(data => data as ResponseData)
        ]);
            let result = await asyncRequestQueue.processQueue()
            console.log("the result form the result is ->", result[0]); 
            this.updateGlobalStateIfKeysAreValid(result[0])
            return{result:result[0].result, error:result[0].error}

        } catch (error) {
            console.log("error occurred in the seeIfKeyISValidFunc ->",error)
            return {result:null, error:error instanceof Error ? error:Error("error in checking if key is valid and updating it func and it is ->"+ error) }
           
        }
    }

    private updateGlobalStateIfKeysAreValid(res:apiResponseAndError) {


        if(res.error !== null && res.result === null){
            console.log(" the error form the backend is(and not updating the state) ->",res.error)
            console.log(" the result is  ", res.result)
            return
        }
        if (res.result === null) {
            return
        }

        keyFromChromeExtensionState.isValidatedThroughBackend = true
        keyFromChromeExtensionState.email = res.result.email
        keyFromChromeExtensionState.name = res.result.name
    }
}



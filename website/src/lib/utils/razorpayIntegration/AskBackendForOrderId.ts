import { razorpayOrderId } from "$lib/sharedState/razorPayKey.svelte";
import type { keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { AsyncRequestQueue } from "../newAsyncRequestQueue";

    

interface requestType {
    user_key:string
}
interface responseType{
    order_id_for_recurring:string,
    order_id_for_onetime:string,
    message:string,
    status_code:number,
    // plan_type:"onetime"|"recurringpayment"
}

/** this func is designed to run when the state changes(derived/effect), run the function and ask the bakend for the order Id and store it in the global state 
 * 
*/
export async function askBackendForOrderId(keyStateObj:keyStateObject, ){
  console.log(`1111`);
  
    try {
        if( !keyStateObj.isValidatedThroughBackend  || keyStateObj.key ===""  || keyStateObj.key === null){
          console.log(`the key is not validated or is empty or null we are returning ->`,keyStateObj.key);
          
            return 
        }
        // if we are validated 
        let reqBody:requestType = { user_key:keyStateObj.key} 
        let asyncReqQueue = new AsyncRequestQueue<Response,responseType>(10)
        let promiseArray = [
            fetch('/api/makeAPayment',{
                headers:{'Content-Type': 'application/json'}
                ,method:"POST",
                body: JSON.stringify(reqBody)
            }),
            // fetch('/api/makeAPayment',{
            //     headers:{'Content-Type': 'application/json'}
            //     ,method:"POST",
            //     body: JSON.stringify(reqBody)
            // })
        ]
  console.log(`22222`);
        let result = await asyncReqQueue.process(promiseArray,(promiseToProcess)=>processIndividualPromise<responseType>(promiseToProcess))
        
        // for (let index = 0; index < result.length; index++) {
            let res = result[0]
            if (res.error !== null) {
                console.log(`there is a error in the result array at ${0} and is ->`,res.error, "\n and the message form the server is ->",res.result?.message);
                return
            }
            console.log(`about to store the response in razorpay storage and it is -> ${res.result?.order_id_for_recurring} and the one time is ${res.result?.order_id_for_onetime}`);
            
          //  if (res.result?.order_id_for_onetime ){
          //       razorpayOrderId.orderIdForOnetime = res.result.order_id_for_onetime
          //  }else if (res.result?.order_id_for_recurring ){
          //       razorpayOrderId.orderIdForRecurring = res.result.order_id_for_recurring
          //  }
        if (res.result?.order_id_for_onetime) {
          razorpayOrderId.orderIdForOnetime = res.result.order_id_for_onetime;
        }
        if (res.result?.order_id_for_recurring) {
          razorpayOrderId.orderIdForRecurring = res.result.order_id_for_recurring;
        }
        // }
        console.log("\n\n\n\n\n result array is ->", result,"\n\n\n\n\n and the razor pay state is 1 time ", razorpayOrderId.orderIdForOnetime, " and recurr ", razorpayOrderId.orderIdForRecurring);
        
    } catch (error) {
        console.log(`there is a error in asking backend for the order Id ->`,error);
        return
    }
}

// async function processIndividualPromise<T>(resp1:Promise<Response>):Promise<T>{
//     let resp = await resp1
//     return new Promise((resolve, reject)=>{
//         if(!resp.ok){
//             console.error("the response is not ok form the backend->", resp);
//             throw "the response is not ok bruh"
//         }
//         resp.json().then((responseInJson)=>{
//             console.log(`the json response is ->`,responseInJson);
            
//             if (validateResponseInJSONTOBeMyType(responseInJson)) {
//                 // JSON is valid   
//                 let respReceived:T = responseInJson
//                 console.log("the valid JSON is ->", responseInJson);
//                 resolve(respReceived)
//             }else{
//                 throw " the json is not even valid or of the same type in the response "
//             }
//         })
//     })
// }

// throw insde the catch cause I want my error to be caught by the async req queue 
async function processIndividualPromise<T>(resp1: Promise<Response>): Promise<T> {
  try {
    let resp = await resp1;
    
    if (!resp.ok) {
      console.error("the response is not ok form the backend->", resp);
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



function validateResponseInJSONTOBeMyType(data:any):boolean {
        // For responseType interface specifically:
        try {
            if (typeof data !== 'object' || data === null) return false;
            // Check if the data has all required fields with correct types
            if (typeof data.order_id_for_recurring !== 'string') return false;
            console.log("---++----+++0000",data.plan_type);
            if (typeof data.message !== 'string') return false;
            console.log("-----message is string------",data.plan_type);
            if (typeof data.order_id_for_onetime !== 'string') return false;
            console.log("---++----+++0000",data.plan_type);
            
            if (typeof data.status_code !== 'number') return false;
            console.log("---++----+++0000----");
            console.log(`the data plan type ->${data.plan_type}`);
            
            
            // if(data.plan_type !==  "onetime" &&  data.plan_type !== "recurringpayment") return false
            console.log("\n\n we are returning true\n\n\n");
            
            return true;
        } catch (error) {
            console.log("error during checking json ->",error);
            return false   
        }
    }
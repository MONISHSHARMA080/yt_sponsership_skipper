import { shouldWeGetOrderIdRecursively } from "$lib/sharedState/getOrderIdRecursively.svelte";
import { razorpayOrderId } from "$lib/sharedState/razorPayKey.svelte";
import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { askBackendForOrderId } from "../razorpayIntegration/AskBackendForOrderId.svelte";


export default getOrderIdRecursively;

// -----------f it writing one myself
// instead of this make a call to chromeExtensions background and if not there display the message to the user, now 
// if there is a error we can just simply ask for the key and that way we wouldn't have diverging keys
async function getOrderIdRecursively(callerName: String) {

  console.log("((((((((((((((((----Are we Already in fetch cycle,", JSON.stringify(razorpayOrderId));
  try {
    let timeToWaitBeforeEachRequest = 2000;
    razorpayOrderId.fetchingStatus = "fetching"; // Set this ONCE at the beginning
    razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle = true

    let numberOfIter = 4
    for (let index = 0; index < numberOfIter; index++) {
      console.log(`in the iteration ${index} of getOrderIdRecursively and the caller is ${callerName}`);


      let res = await askBackendForOrderId(keyFromChromeExtensionState);
      console.log(`!!!the success is ${res} !and we are on the loop iter ${index}`);

      if (res) {
        console.log(`we got the order id Successfully and will quit`);
        razorpayOrderId.fetchingStatus = "success";
        razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle = false
        shouldWeGetOrderIdRecursively.shouldWeDoIt = false
        return;
      } else {
        console.log(`order id is not received form the loop`);

      }

      razorpayOrderId.fetchingStatus = "fetching"; // Set back to fetching for next attempt
      // Only set to error if we're on the final attempt OR before waiting
      if (index === numberOfIter - 1) {
        razorpayOrderId.fetchingStatus = "error";
        razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle = false
        return; // Exit after final attempt
      } else {
        razorpayOrderId.fetchingStatus = "error"; // Show error before waiting
        console.log("Waiting for 4 seconds before retrying...");
        let a = await new Promise(resolve => setTimeout(() => {
          console.log(` !!++!! waiting for the promsise to finish`);
          resolve(999)
        }, timeToWaitBeforeEachRequest));
        console.log(`!!++!!the waiting finished`);

      }
    }
    console.log(`the fetching completed and the fetching status is ->${razorpayOrderId.fetchingStatus}`);
    razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle = false
  } catch (error) {
    razorpayOrderId.areWeInAMiddleOfMultipleFetchCycle = false
    console.log(` error is in getOrderIdRecursively's try catch block ->`, error);
  }
}

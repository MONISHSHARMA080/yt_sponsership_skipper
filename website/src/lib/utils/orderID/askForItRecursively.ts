import { razorpayOrderId } from "$lib/sharedState/razorPayKey.svelte";
import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { askBackendForOrderId } from "../razorpayIntegration/AskBackendForOrderId.svelte";

// function getOrderIdRecursively(maxRetries = 6) {
//   let retriesLeft = maxRetries;
//   let succeeded = false;
//   let requestInProgress = false; // Track if a request is currently in progress
//
//   // Function to execute the attempt
//   function attemptFetch() {
//     // Don't start a new request if one is already in progress or we've succeeded
//     if (requestInProgress || succeeded || retriesLeft <= 0) {
//       return;
//     }
//
//     try {
//       requestInProgress = true; // Mark that a request is in progress
//
//       askBackendForOrderId(keyFromChromeExtensionState)
//         .then((val) => {
//           requestInProgress = false; // Request completed
//
//           console.log(`The order ID from the backend's value is -> ${val}`);
//           if (val) {
//             // Success! We got an order ID
//             succeeded = true;
//             console.log("Successfully received order ID");
//           } else {
//             // No value, retry if we have attempts left
//             retriesLeft--;
//             console.log(`No order ID received. ${retriesLeft} attempts left.`);
//             if (retriesLeft > 0) {
//               setTimeout(attemptFetch, 2000); // Try again in 2 seconds
//             }
//           }
//         })
//         .catch((error) => {
//           requestInProgress = false; // Request completed (with error)
//
//           // Error occurred, retry if we have attempts left
//           retriesLeft--;
//           console.error("Error fetching order ID:", error);
//           console.log(`${retriesLeft} attempts left.`);
//           if (retriesLeft > 0) {
//             setTimeout(attemptFetch, 3000); // Try again in 3 seconds
//           }
//         });
//     } catch (e) {
//       requestInProgress = false; // Request failed to start
//
//       // Handle any synchronous errors
//       retriesLeft--;
//       console.error("Unexpected error:", e);
//       if (retriesLeft > 0) {
//         setTimeout(attemptFetch, 3000);
//       }
//     }
//   }
//
//   // Start the first attempt
//   attemptFetch();
//
//   // Return a function to check the status and retry if needed
//   return function checkAndRetry() {
//     if (!succeeded && !requestInProgress && retriesLeft > 0) {
//       attemptFetch(); // Try again if not succeeded and no request in progress
//     }
//     return succeeded;
//   };
// }

export default getOrderIdRecursively;




// -----------f it writing one myself
async function getOrderIdRecursively() {
  let timeToWaitBeforeEachRequest = 2000;
  razorpayOrderId.fetchingStatus = "fetching"; // Set this ONCE at the beginning

  let numberOfIter = 1.2
  for (let index = 0; index < numberOfIter; index++) {
    console.log(`in the iteration ${index} of getOrderIdRecursively`);

    let res = await askBackendForOrderId(keyFromChromeExtensionState);
    console.log(`!!!the success is ${res} !and we are on the loop iter ${index}`);

    if (res) {
      console.log(`we got the order id Successfully and will quit`);
      razorpayOrderId.fetchingStatus = "success";
      return;
    }

    // Only set to error if we're on the final attempt OR before waiting
    if (index === numberOfIter - 1) {
      razorpayOrderId.fetchingStatus = "error";
      return; // Exit after final attempt
    } else {
      razorpayOrderId.fetchingStatus = "error"; // Show error before waiting
      console.log("Waiting for 4 seconds before retrying...");
      await new Promise(resolve => setTimeout(resolve, timeToWaitBeforeEachRequest));
      razorpayOrderId.fetchingStatus = "fetching"; // Set back to fetching for next attempt
    }
  }
  console.log(`the fetching completed and the fetching status is ->${razorpayOrderId.fetchingStatus}`);

}

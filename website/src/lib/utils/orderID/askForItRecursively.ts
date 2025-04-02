import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";
import { askBackendForOrderId } from "../razorpayIntegration/AskBackendForOrderId";

// let getOrderIdRecursively = () => {
//   return askBackendForOrderId(keyFromChromeExtensionState)
//     .then((val) => {
//       console.log(`the order id from the backend's value is ->${val}`);
//
//       // If no value received, wait and try again
//       if (!val) {
//         console.log("No order ID received, retrying in 2 seconds...");
//         return new Promise((resolve) => {
//           setTimeout(() => {
//             resolve(getOrderIdRecursively());
//           }, 2000); // Retry after 2 seconds
//         });
//       }
//
//       // If we have a value, return it
//       return val;
//     })
//     .catch((error) => {
//       console.error("Error fetching order ID:", error);
//       console.log("Retrying in 3 seconds due to error...");
//
//       // On error, wait and try again
//       return new Promise((resolve) => {
//         setTimeout(() => {
//           resolve(getOrderIdRecursively());
//         }, 3000); // Retry after 3 seconds on error
//       });
//     });
// };
//
 function getOrderIdRecursively(maxRetries = 6) {
  let retriesLeft = maxRetries;
  let succeeded = false;

  // Function to execute the attempt
  function attemptFetch() {
    if (succeeded || retriesLeft <= 0) {
      return; // Stop if we've succeeded or run out of retries
    }

    try {
      askBackendForOrderId(keyFromChromeExtensionState)
        .then((val) => {
          console.log(`The order ID from the backend's value is -> ${val}`);

          if (val) {
            // Success! We got an order ID
            succeeded = true;
            console.log("Successfully received order ID");
          } else {
            // No value, retry if we have attempts left
            retriesLeft--;
            console.log(`No order ID received. ${retriesLeft} attempts left.`);

            if (retriesLeft > 0) {
              setTimeout(attemptFetch, 2000); // Try again in 2 seconds
            }
          }
        })
        .catch((error) => {
          // Error occurred, retry if we have attempts left
          retriesLeft--;
          console.error("Error fetching order ID:", error);
          console.log(`${retriesLeft} attempts left.`);

          if (retriesLeft > 0) {
            setTimeout(attemptFetch, 3000); // Try again in 3 seconds
          }
        });
    } catch (e) {
      // Handle any synchronous errors
      retriesLeft--;
      console.error("Unexpected error:", e);

      if (retriesLeft > 0) {
        setTimeout(attemptFetch, 3000);
      }
    }
  }

  // Start the first attempt
  attemptFetch();

  // Return a function to check the status
  return function checkSucceeded() {
    return succeeded;
  };
}
//

export default getOrderIdRecursively
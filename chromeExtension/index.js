// @ts-check
/// <reference types="chrome" />
console.log("hi form the index.js");

// try {
//   let element = document.getElementById("status");
//   if (element !== null) {
//     element.addEventListener("click", async function a() {
//       let [key, error] = await getKeyFromStorageOrBackend();
//       if (error) {
//         console.log(
//           "error is there in getting the key and it is -->",
//           error,
//           "\n\n the key is",
//           key,
//         );
//         return;
//       }
//       console.log("the key is -->", key);
//     });
//   }
// } catch (e) {
//   console.log("error in the index.js ", e);
// }

try {
  let element = document.getElementById("websiteButton");
  if (element !== null) {
    element.addEventListener("click", (event) => {
      let response = chrome.runtime.sendMessage({ action: "navigateToWebsite" })
    })
  }

} catch (error) {
  console.log("there is a error in the index.js in the website button ->", error);
}

// /**
//  * Authenticates user and retrieves a key; will return empty string on the error
//  * @typedef {Object} AuthResponse
//  * @property {string} message - Response message
//  * @property {number} status_code - HTTP status code
//  * @property {boolean} success - Success indicator
//  * @property {string} encrypted_key - The encrypted key
//  *
//  * @typedef {Object} UserDetail
//  * @property {string} account_id - The account ID of the user that is unique for every account.
//  * @property {string} user_token - The authentication token for the user.
//  *
//  * @throws {Error} If there's an error getting the token
//  *
//  * @returns {Promise<[string, Error|null]>} key - The authentication key
//  */

// async function userAuthAndGetTheKey() {
//   /** @type {UserDetail} */
//   const UserDetail = {
//     account_id: "",
//     user_token: "",
//   };

//   try {
//     // Get user info
//     const userInfo = await new Promise((resolve, reject) => {
//       chrome.identity.getProfileUserInfo({}, (userInfoFromChrome) => {
//         if (chrome.runtime.lastError) {
//           reject(chrome.runtime.lastError);
//         } else {
//           resolve(userInfoFromChrome);
//         }
//       });
//     });

//     // Set user info
//     UserDetail.account_id = userInfo.id;
//     console.log("email:", userInfo.email);
//     console.log("user info type -->", typeof userInfo.id);
//     console.log(
//       "id:",
//       UserDetail.account_id,
//       "  type of id is ",
//       typeof UserDetail.account_id,
//     );

//     // Get auth token
//     const authToken = await new Promise((resolve, reject) => {
//       chrome.identity.getAuthToken({ interactive: true }, (token) => {
//         if (chrome.runtime.lastError) {
//           console.log(
//             "error happened during the authToken part, error is -->",
//             chrome.runtime.lastError.message,
//           );
//           reject(chrome.runtime.lastError);
//         } else {
//           resolve(token);
//         }
//       });
//     });

//     UserDetail.user_token = authToken;
//     console.log("auth token:", authToken);

//     try {
//       console.log("config is -->", config.BACKEND_URL);
//     } catch (error) {
//       console.log("error in printing the config-->", error);
//       return ["", error];
//     }

//     const response = await fetch(config.BACKEND_URL + "/signup", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify(UserDetail),
//     });

//     if (!response.ok && response.status !== 200) {
//       return ["", new Error(`HTTP error! Status: ${response.status}`)];
//     }

//     /** @type {AuthResponse} */
//     const data = await response.json();
//     console.log("Success:", data);
//     return [data.encrypted_key, null];
//   } catch (error) {
//     console.error("Error:", error);
//     return ["", error];
//   }
// }
// // break down the function in 3 parts, 1) that get the key form the backend, 2) that sets it in the localstorage,
// /**
//  * takes in a key,value and will store it in the storage
//  * @param {string} key
//  * @param {string} value
//  * @param {Function} functionToRun
//  * @returns {Error|null} error
//  */
// function saveValueToTheStorage(key, value, functionToRun) {
//   if (value === "") {
//     return new Error("value can't be empty");
//   }
//   try {
//     chrome.storage.sync.set({ [key]: value }, () => {
//       functionToRun();
//     });
//     return null;
//   } catch (error) {
//     return error;
//   }
// }

// /**
//  * Gets a value from Chrome storage by key and processes it with the provided function
//  * @param {string} key
//  * @param {Function} functionToRun -
//  * @returns {Promise<[string|null, Error|null]>} A tuple containing [value, error]
//  */
// function getValueFromTheStorage(key, functionToRun) {
//   return new Promise((resolve) => {
//     try {
//       chrome.storage.sync.get([key], (item) => {
//         // Handle runtime errors
//         if (chrome.runtime.lastError) {
//           console.error(
//             "Error retrieving data from storage:",
//             chrome.runtime.lastError,
//           );
//           resolve([null, new Error(chrome.runtime.lastError.message)]);
//           return;
//         }
//         const value = item[key];
//         try {
//           functionToRun();
//         } catch (callbackError) {
//           resolve([null, callbackError]);
//           return;
//         }
//         resolve([value ?? null, null]);
//       });
//     } catch (error) {
//       resolve([null, error]);
//     }
//   });
// }

// /**
//  * returns the key, will try the storage, if not there then will fetch it, it can also return string and error if the key is fetched and not able to stored in the storage
//  * @returns {Promise<[string,Error|null]>}
//  */
// async function getKeyFromStorageOrBackend() {
//   try {
//     let [valueFromStorage, error] = await getValueFromTheStorage(
//       "key",
//       () => {},
//     );
//     if (error) {
//       return ["", error];
//     }
//     if (valueFromStorage != "") {
//       return [valueFromStorage, null];
//     }
//     let [encryptedKey, errorFromKey] = await userAuthAndGetTheKey();
//     if (errorFromKey || encryptedKey === "") {
//       return ["", errorFromKey];
//     }
//     error = saveValueToTheStorage("key", encryptedKey, () => {
//       console.log("storing the key +");
//     });
//     if (error) {
//       return [encryptedKey, error];
//     } else {
//       return [encryptedKey, null];
//     }
//   } catch (error) {
//     return ["", error];
//   }
// }

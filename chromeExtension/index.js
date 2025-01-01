// @ts-check
/// <reference types="chrome" />
console.log("hi form the index.js");

try {
  let element = document.getElementById("status");
  if (element !== null) {
    element.addEventListener("click", async function a() {
      console.log("hi from the index.js button-- ");
      const encryptedKey = await userAuthAndGetTheKey();
      console.log("encrypted key is ->", encryptedKey);
    });
  }
} catch (e) {
  console.log("error ++--++", e);
}

/**
 * Authenticates user and retrieves a key; will return empty string
 * @typedef {Object} AuthResponse
 * @property {string} message - Response message
 * @property {number} status_code - HTTP status code
 * @property {boolean} success - Success indicator
 * @property {string} encrypted_key - The encrypted key
 *
 * @typedef {Object} UserDetail
 * @property {string} account_id - The account ID of the user that is unique for every account.
 * @property {string} user_token - The authentication token for the user.
 *
 * @throws {Error} If there's an error getting the token
 *
 * @returns {Promise<string>} key - The authentication key
 */

async function userAuthAndGetTheKey() {
  /** @type {UserDetail}  */
  const UserDetail = {
    account_id: "",
    user_token: "",
  };
  try {
    let userInfo;
    chrome.identity.getProfileUserInfo({}, (userInfoFromChrome) => {
      userInfo = userInfoFromChrome;
      UserDetail.account_id = userInfo.id; // This ensures we send a number, not a string
      console.log("email:", userInfo.email);
      console.log("user info type -->", typeof userInfo.id);
      console.log(
        "id:",
        UserDetail.account_id,
        "  type of id is ",
        typeof UserDetail.account_id,
      );
    });
    /**
     * Gets the authentication token for the current user
     * @returns {Promise<string>} A promise that resolves with the auth token
     * @throws {Error} If there's an error getting the token
     */
    let authToken;
    chrome.identity.getAuthToken({ interactive: true }, (token) => {
      if (chrome.runtime.lastError) {
        console.log(
          "error happened during the authToken part, error is -->",
          chrome.runtime.lastError.message,
        );
        return "";
      }
      authToken = token;
    });
    UserDetail.user_token = authToken;
    console.log("auth token:", authToken);

    const response = await fetch(config.BACKEND_URL + "/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(UserDetail),
    });
    try {
      console.log("config is -->", config.BACKEND_URL);
    } catch (error) {
      console.log("error in printing the config-->", error);
      return "";
    }

    if (!response.ok && response.status != 200) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    /** @type {AuthResponse} */
    const data = await response.json();
    console.log("Success:", data);
    return data.encrypted_key;
  } catch (error) {
    console.error("Error:", error);
    return "";
  }
}
// break down the function in 3 parts, 1) that get the key form the backend, 2) that sets it in the localstorage,
/**
 * takes in a key,value and will store it in the storage
 * @param {string} key
 * @param {string} value
 * @param {Function} functionToRun
 * @returns {Error|null} error
 */
function saveValueToTheStorage(key, value, functionToRun) {
  try {
    chrome.storage.sync.set({ [key]: value }, () => {
      functionToRun();
    });
    return null;
  } catch (error) {
    return error;
  }
}

/**
 * Gets a value from Chrome storage by key and processes it with the provided function
 * @param {string} key
 * @param {Function} functionToRun -
 * @returns {Promise<[any|null, Error|null]>} A tuple containing [value, error]
 */
function getValueFromTheStorage(key, functionToRun) {
  return new Promise((resolve) => {
    try {
      chrome.storage.sync.get([key], (item) => {
        // Handle runtime errors
        if (chrome.runtime.lastError) {
          console.error(
            "Error retrieving data from storage:",
            chrome.runtime.lastError,
          );
          resolve([null, new Error(chrome.runtime.lastError.message)]);
          return;
        }
        const value = item[key];
        try {
          functionToRun();
        } catch (callbackError) {
          resolve([null, callbackError]);
          return;
        }
        resolve([value ?? null, null]);
      });
    } catch (error) {
      resolve([null, error]);
    }
  });
}

// @ts-check
/// <reference types="chrome-types" />
const config = {
  BACKEND_URL: "http://localhost:8080",
};
/**
 * Authenticates user and retrieves a key; will return empty string on the error
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
 * @param {Config} config - Configuration object containing backend URL
 *
 * @throws {Error} If there's an error getting the token
 *
 * @returns {Promise<[string, Error|null]>} key - The authentication key
 */

async function userAuthAndGetTheKey(config) {
  /** @type {UserDetail} */
  const UserDetail = {
    account_id: "",
    user_token: "",
  };

  try {
    // Get user info
    const userInfo = await new Promise((resolve, reject) => {
      chrome.identity.getProfileUserInfo(
        { accountStatus: "SYNC" },
        (userInfoFromChrome) => {
          if (chrome.runtime.lastError) {
            reject(chrome.runtime.lastError);
          } else {
            resolve(userInfoFromChrome);
          }
        },
      );
    });

    // Set user info
    UserDetail.account_id = userInfo.id;
    console.log("email:", userInfo.email);
    console.log("user info type -->", typeof userInfo.id);
    console.log(
      "id:",
      UserDetail.account_id,
      "  type of id is ",
      typeof UserDetail.account_id,
    );

    // Get auth token
    const authToken = await new Promise((resolve, reject) => {
      chrome.identity.getAuthToken({ interactive: true }, (token) => {
        if (chrome.runtime.lastError) {
          console.log(
            "error happened during the authToken part, error is -->",
            chrome.runtime.lastError.message,
          );
          reject(chrome.runtime.lastError);
        } else {
          resolve(token);
        }
      });
    });

    UserDetail.user_token = authToken;
    console.log("auth token:", authToken);

    try {
      console.log("config is -->", config.BACKEND_URL);
    } catch (error) {
      console.log("error in printing the config-->", error);
      if (error instanceof Error) {
        return ["", error];
      }
      return ["", new Error(String(error))];
    }

    const response = await fetch(config.BACKEND_URL + "/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(UserDetail),
    });

    if (!response.ok && response.status !== 200) {
      return ["", new Error(`HTTP error! Status: ${response.status}`)];
    }

    /** @type {AuthResponse} */
    const data = await response.json();
    console.log("Success:", data);
    return [data.encrypted_key, null];
  } catch (error) {
    console.error("Error:", error);
    if (error instanceof Error) {
      return ["", error];
    }
    return ["", new Error(String(error))];
  }
}
// break down the function in 3 parts, 1) that get the key form the backend, 2) that sets it in the localstorage,
/**
 * @template T the type parameter name
 * takes in a key,value and will store it in the storage
 * @param {string} key
 * @param {T} value
 * @param {Function} [functionToRun]
 * @returns {Error|null} error
 */
export function saveValueToTheStorage(key, value, functionToRun) {
  if (value === "") {
    return new Error("value can't be empty");
  }
  try {
    chrome.storage.local.set({ [key]: value }, () => {
      functionToRun ? functionToRun() : null;
    });
    return null;
  } catch (error) {
    if (error instanceof Error) {
      return error;
    }
    // If it's not an Error object, create a new Error
    return new Error(String(error));
  }
}

/**
 * Gets a value from Chrome storage by key and processes it with the provided function
 * @param {string} key
 * @param {Function} functionToRun -
 * @returns {Promise<[string|null, Error|null]>} A tuple containing [value, error]
 */
export function getValueFromTheStorage(key, functionToRun) {
  return new Promise((resolve) => {
    try {
      chrome.storage.local.get([key], (item) => {
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
          if (callbackError instanceof Error) {
            resolve(["", new Error(String(callbackError))]);
            return;
          }
          // If it's not an Error object, create a new Error
          resolve(["", new Error(String(callbackError))]);
          return;
        }
        resolve([value ?? null, null]);
      });
    } catch (e) {
      if (e instanceof Error) {
        return ["", e];
      }
      // If it's not an Error object, create a new Error
      return ["", new Error(String(e))];
    }
  });
}

/**
 * returns the key, will try the storage, if not there then will fetch it, it can also return string and error if the key is fetched and not able to stored in the storage
 *
 * @typedef {Object} Config
 * @property {string} BACKEND_URL - The base URL for the backend API server
 *
 * @param {Config} config - Configuration object for the application
 *
 * @returns {Promise<[string,Error|null]>}
 */
export async function getKeyFromStorageOrBackend(config) {
  console.log("in getKeyFromStorageOrBackend");
  try {
    let [valueFromStorage, error] = await getValueFromTheStorage(
      "key",
      () => { },
    );
    console.log(" after the value");
    if (error !== null) {
      console.log(
        "there is  a error and that is -->",
        error,
        "and the value is ->",
        valueFromStorage,
      );
      return ["", error];
    }
    if (valueFromStorage !== null && valueFromStorage !== "") {
      console.log("about to return the value and that is ->", valueFromStorage);
      return [valueFromStorage, null];
    }
    console.log("didn't find error in the ");
    // the value is blank so lets fetch
    // if (valueFromStorage !== "") {
    //   return [valueFromStorage, null];
    // }
    console.log("fetching from storage");
    let [encryptedKey, errorFromKey] = await userAuthAndGetTheKey(config);
    if (errorFromKey || encryptedKey === "") {
      return ["", errorFromKey];
    }
    error = saveValueToTheStorage("key", encryptedKey, () => {
      console.log("storing the key +");
    });
    if (error) {
      return [encryptedKey, error];
    } else {
      return [encryptedKey, null];
    }
  } catch (e) {
    if (e instanceof Error) {
      return ["", e];
    }
    // If it's not an Error object, create a new Error
    return ["", new Error(String(e))];
  }
}

export function sayHi() {
  console.log("hi form the helper file --++:)");
}

/**
 * @param {string} pathWithoutBackSlash - url path of the function
 * @param {"POST"|"GET"} method - request method
 * @param {Object} [header={ "Content-Type": "application/json" } ] - request header
 * @param {Object} bodyOBJ - request body
 *
 * @returns {() => Promise<Response>}
 */

export function fetchFunctionBuilder(
  pathWithoutBackSlash,
  method,
  bodyOBJ,
  header = { "Content-Type": "application/json" },
) {
  return () => {
    return fetch(config.BACKEND_URL + "/" + pathWithoutBackSlash, {
      method: method,
      // @ts-ignore
      headers: header,
      body: JSON.stringify(bodyOBJ),
    });
  };
}

/**
 * @typedef {Object} ResponseObject
 * @property {number} status - The status code of the response
 * @property {string} message - A message providing additional information
 * @property {number} startTime - The start time in milliseconds
 * @property {number} endTime - The end time in milliseconds
 * @property {boolean} containSponserSubtitle - Whether the video has sponsorship subtitle
 * @property {string} [error] - Optional error message
 * @export
 */

/**
 * @typedef {Object} bodyOfTheRequest
 * @property {string} youtube_Video_Id -- video ID
 * @property {string} encrypted_string -- video ID
 *
 *
 * @param {string} key -- key form the backend
 * @param {string} videoID -- key form the backend
 *
 * @returns {Promise<[ResponseObject|null, Error|null]>}
 */
export async function getWhereToSkipInYtVideo(key, videoID) {
  /** @type bodyOfTheRequest */
  let requestBody = { youtube_Video_Id: videoID, encrypted_string: key };
  let fetchRequestToBackend = fetchFunctionBuilder(
    "youtubeVideo",
    "POST",
    { "Content-Type": "application/json" },
    requestBody,
  );
  console.log(` sending a fetch request to the backend`)
  try {
    let response = await fetchRequestToBackend();
    console.log("the response form the yt video api is -->", response);
    response.status
    /** @type ResponseObject */
    let responseOBJ = await response.json();
    if (responseOBJ.status !== 200) {
      console.log("there is a error in the yt api -->", responseOBJ);
      return [null, new Error(response.statusText)];
    }
    return [responseOBJ, null];
  } catch (e) {
    console.log("error in getWhereToSkipInYtVideo -->", e);
    if (e instanceof Error) {
      return [null, e];
    }
    // If it's not an Error object, create a new Error
    return [null, new Error(String(e))];
  }
}

// alwaysSkipTheSponsorAndDoNotShowTheModal :false

export async function getDefaultValueOfToSkipTheSponsorAndShowTheModal() /** @returns {Promise<[Boolean,Error|null]>} */ {
  let key = "alwaysSkipTheSponsorAndDoNotShowTheModal";
  let [valueFromStorage, error] = await getValueFromTheStorage(key, () => { });
  if (
    error !== null ||
    valueFromStorage === null ||
    valueFromStorage === "" ||
    typeof valueFromStorage !== "boolean"
  ) {
    // return [false, error]
    // going to set the value too
    let error = saveValueToTheStorage(key, false, () => {
      console.log(
        " the value for the default ",
        key,
        " is not found so saving the default to be false",
      );
    });
    if (error) {
      console.log(
        "there is a error in the default value for the default  in the ",
        key,
        " -->",
        error,
      );
      return [false, new Error(String(error))];
    }
    return [false, null];
  }
  return [Boolean(valueFromStorage), null];
}


/**
 * run this function after the response is received if the response.status is 426(upgrade required) we will fetch the new key and store it in the storage 
 * if it is alright we will return the response as it is 
 *
 * NOTE:  any other reponse than 200 or 426, will be treated same as 200 by us and the same repsonse will be returned, you have to gandle that case
 *
 * @argument {Response} response 
 * @returns {Response}
 *
 */
function updateTheKeyToNewValueIfUpgradeIsRequeired(response) {
  if (response.status !== 426) {
    return response
  }

  fetchFunctionBuilder("", "POST", {})
  return response
}

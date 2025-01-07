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
      chrome.identity.getProfileUserInfo({accountStatus:"SYNC"}, (userInfoFromChrome) => {
        if (chrome.runtime.lastError) {
          reject(chrome.runtime.lastError);
        } else {
          resolve(userInfoFromChrome);
        }
      });
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
      return ["", error];
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
    return ["", error];
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
  if (value === "") {
    return new Error("value can't be empty");
  }
  try {
    chrome.storage.local.set({ [key]: value }, () => {
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
 * @returns {Promise<[string|null, Error|null]>} A tuple containing [value, error]
 */
function getValueFromTheStorage(key, functionToRun) {
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
      () => {},
    );
    console.log(" after the value")
    if (error !== null) {
      console.log("there is  a error and that is -->" , error, "and the value is ->",valueFromStorage);
      return ["", error];
    }
    if (valueFromStorage !== null && valueFromStorage !== "") {
      console.log("about to return the value and that is ->",valueFromStorage)
      return [valueFromStorage, null]
    }
    console.log("didn't find error in the ")
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
  } catch (error) {
    return ["", error];
  }
}

export function sayHi(){
  console.log("hi form the helper file --++:)")
}

/**
 * @param {string} pathWithoutBackSlash - url path of the function
 * @param {"POST"|"GET"} method - request method
 * @param {Object} header - request header
 * @param {Object} bodyOBJ - request body
 *
 * @returns {() => Promise<Response>}
 */

export   function fetchFunctionBuilder(pathWithoutBackSlash, method, header, bodyOBJ) {
  return ()=>{
    return  fetch(config.BACKEND_URL+"/"+pathWithoutBackSlash, {
      method: method,
      headers: header,
      body: JSON.stringify(bodyOBJ),
    })
  }
}

/**
 * @typedef {Object} bodyOfTheRequest
 * @property {string} youtube_Video_Id -- video ID
 * @property {string} encrypted_string -- video ID
 *
 * @typedef {Object} responseObject
 * @property {number} status - The status code of the response.
 * @property {string} message - A message providing additional information about the response.
 * @property {number} startTime - The start time in milliseconds for where the skip should begin.
 * @property {number} endTime - The end time in milliseconds for where the skip should end.
 * @property {boolean} containSponserSubtitle- does the video has sponsorship subtitle
 * @property {string} [error] - An optional error message if something goes wrong.
 *
 * @param {string} key -- key form the backend
 * @param {string} videoID -- key form the backend
 *
 * @returns {Promise<[responseObject|null, Error|null]>}
 */
export async  function getWhereToSkipInYtVideo(key, videoID) {
  /** @type bodyOfTheRequest */
  let  requestBody ={youtube_Video_Id:videoID, encrypted_string:key};
 let fetchRequestToBackend = fetchFunctionBuilder("youtubeVideo", "POST", {'Content-Type': 'application/json'}, requestBody )
  try {
    let response = await fetchRequestToBackend()
    console.log("the response form the yt video api is -->", response)
    /** @type responseObject */
    let responseOBJ = await response.json();
    if (responseOBJ.status !== 200) {
      console.log("there is a error in the yt api -->", responseOBJ);
      return [null, new Error(response.statusText)];
    }
    return [responseOBJ, null];
  }catch (e) {
    console.log("error in getWhereToSkipInYtVideo -->",e);
    return [null, e];
  }
}

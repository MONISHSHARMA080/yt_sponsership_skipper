// @ts-check
/// <reference types="chrome-types" />
// const config = {
//   BACKEND_URL: "http://localhost:8080",
// };
import { config } from "./config.js";
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
        { accountStatus: "ANY" },
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
 * makes a fetch request to get the transcript form the url in it
 * @param {captionTracks} captionTrack  
 *  @returns {Promise<[string|null, Error|null]>} - the function will fetch the transcript and will return it 
 */
async function fetchTheTranscript(captionTrack) {
  try {

    // first we are going to get the transcript clicking the caption button and 

    // let response = await chrome.runtime.sendMessage({ action: "clickCCButtonAndGetCaptions" })

    // console.log(`the response from send message is ${response}\n`)

    //


    let result = await fetch(captionTrack.baseUrl)
    console.log(`got the result in fetching the caption's tracks and it is Status:${result.status}, status code : ${result.statusText} `)
    if (!result.ok) {
      console.log(` there is a error in fetching the transcript and it is not 200`)
      return [null, new Error(`there is a error in fetching the transcript and the response is not ok , it is ${result.status}`)]
    }
    let res = await result.text()
    if (res === "") {
      return [null, new Error('the fetched transcript is empty\n')]
    }
    console.log(`@ the result is ${res}`)
    return [res, null]
  } catch (err) {
    return [null, err instanceof Error ? err : new Error(`there is a error in fetching the transcript and it is :->${err}`)]
  }

}

/**
 *
 * gets the preferable captions first if not there then gets the auto gen and if that is also not there then gets the one in the
 * first index(0) , or error if any one of those is not there
 *
 * @param {string} jsonStringifiedCaptions  
 *  @returns {Promise<[string|null, Error|null]>} - the function will fetch the transcript and will return it 
 */
async function GetTheTranscriptFromTheCaptions(jsonStringifiedCaptions) {
  console.log(`the captions tracks we got in the helper file is ->${jsonStringifiedCaptions}`)
  try {
    /** @type Captions */
    let captions
    captions = JSON.parse(jsonStringifiedCaptions)
    // console.log(`the lenght of the captionTracks is ${captions}`)
    // /** @type {captionTracks|null} */
    // let firstChoiceEnSub = null
    // /** @type {captionTracks|null} */
    // let secondChoiceEnAutoGenSub = null
    // /** @type {captionTracks|null} */
    // let thirdChoiceChooisingTheFirstOne = null
    //
    // captions.captionTracks.forEach((value, index) => {
    //   if (index === 0) {
    //     console.log(`got the captions at 0th index and the url here is ${value.baseUrl} \n`)
    //     thirdChoiceChooisingTheFirstOne = value
    //   }
    //   if (!value.name.simpleText.includes("English (auto-generated)") && value.name.simpleText.includes("English")) {
    //     console.log(`got the captions at English(this  one includes English in the value.name.simpleText and does not contain English (auto-generated) it is  ${value.name.simpleText}) one and the url is ${value.baseUrl} \n`)
    //     firstChoiceEnSub = value
    //   } else if (value.name.simpleText === "English (auto-generated)") {
    //     console.log(`got the captions at English (auto-generated) one  ${value.baseUrl} \n`)
    //     secondChoiceEnAutoGenSub = value
    //   }
    // })
    console.log(`about to fetch`)
    /** @type {[string|null, Error|null]} */
    let resultFromFetching

    if (firstChoiceEnSub !== null || secondChoiceEnAutoGenSub !== null && thirdChoiceChooisingTheFirstOne !== null) {
      if (firstChoiceEnSub !== null) {
        // call the func and then return
        // if there is a error then we are going to the other one too
        console.log(`fetching the English `)
        let res = await fetchTheTranscript(firstChoiceEnSub)
        resultFromFetching = res
      } else if (secondChoiceEnAutoGenSub !== null) {
        console.log(`fetching the English (auto-generated)`)
        // call the func and then return
        // if there is a error then we are going to the other one too
        let res = await fetchTheTranscript(secondChoiceEnAutoGenSub)
        resultFromFetching = res
      } else if (thirdChoiceChooisingTheFirstOne !== null) {
        // call the func and then return
        // if there is a error then we are going to the other one too
        console.log(`fetching the 3rd choice as other onces are null`)
        let res = await fetchTheTranscript(thirdChoiceChooisingTheFirstOne)
        resultFromFetching = res
      } else {
        console.error(`there is a error in fetching in the transcript as all the choices are null`)
        return [null, new Error(`there is a error in fetching in the transcript as all the choices are null`)];
      }
      console.log(`the result from the fetch transcript functin is ${JSON.stringify(resultFromFetching)}`)
      // handle the error
      if (resultFromFetching[1] !== null || resultFromFetching[0] === null) {
        let e = resultFromFetching[1]
        return [null, e instanceof Error ? e : new Error(`there is a error in fetching in the transcript and it is:=> ${e}`)];
      } else {
        console.log(`the response form fetching the transcript is ok +`)
        return resultFromFetching;
      }
    } else {
      return [null, new Error(`there is a error in getting the captions tracks from the parsed captions object `)];
    }



  } catch (error) {
    console.log(`there is a error in pasing the CaptionTracks form the string and it is -> ${error}`)
    return [null, new Error(`there is a error in pasing the CaptionTracks form the string and it is -> ${error}`)];
  }
}



/**
 * @typedef {object} Captions
 * @property {captionTracks[]} captionTracks 
 * @property {object[]} audioTracks
 * @property {number[]} audioTracks.captionTrackIndices
 * @property {object[]} translationLanguages
 * @property {string} translationLanguages.languageCode
 * @property {object} translationLanguages.languageName
 * @property {string} translationLanguages.languageName.simpleText
 * @property {number} defaultAudioTrackIndex
 *
 *
 * @typedef {object} captionTracks 
 * @property {string} baseUrl
 * @property {object} name
 * @property {string} name.simpleText
 * @property {string} vssId
 * @property {string} languageCode
 * @property {string} kind
 * @property {boolean} isTranslatable
 * @property {string} trackName
 *
 */

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
 * @property {string} transcript -- video ID
 *
 *
 * @param {string} key -- key form the backend
 * @param {string} videoID -- key form the backend
 * @param {string} transcript -- fetched transcrip 
 *
 * @returns {Promise<[ResponseObject|null, Error|null]>}
 */

export async function getWhereToSkipInYtVideo(key, videoID, transcript) {
  console.log(`the captions tracks we got in the helper file is ->${transcript}`)
  // let res = await GetTheTranscriptFromTheCaptions(transcript)
  // console.log(`the res form getting the transcript form the captions is ${res}---- ${JSON.stringify(res)}`)
  // console.log("------\n\n\n\n\n\n\n")
  // if (res[0] === null || res[0] === "" || res[1] !== null) {
  //   console.log(`there is a error in getting the transcript form the captions object and it is -> ${res[1]}`)
  //   return [null, res[1] instanceof Error ? res[1] : new Error(` there is a error in getting the result from the transcript and it is ->${res[1]} `)]
  // }


  /** @type bodyOfTheRequest */
  let requestBody = { youtube_Video_Id: videoID, encrypted_string: key, transcript: transcript };
  let fetchRequestToBackend = fetchFunctionBuilder(
    "youtubeVideo",
    "POST",
    requestBody,
    { "Content-Type": "application/json" },
  );
  console.log(` sending a fetch request to the backend`)
  try {
    let response = await fetchRequestToBackend();
    console.log("the response form the yt video api is -->", response, `the origianl request(before the uodate the key func's)req body is ${JSON.stringify(requestBody)} `);
    response.status
    let newRespOrError = await updateTheKeyToNewValueIfUpgradeIsRequeired(response, key, fetchRequestToBackend)
    if (newRespOrError instanceof Error) {
      console.error(`there is a error in getting the new key token so we are just returning ; the error is ${newRespOrError}`)
      return [null, newRespOrError]
    } else {
      console.log(`the new reponse is a instanceof of Response ${newRespOrError instanceof Response}`)
      response = newRespOrError
    }
    /** @type ResponseObject */
    let responseOBJ = await response.json();
    if (responseOBJ.status !== 200) {
      console.log("there is a error in the yt api(response is not 200) -->", JSON.stringify(responseOBJ));
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
 * Error case: in case we get 403 in the (forbidden) we are reuring a error as we can't update the key and we got told to update the key(catch 22)
 * so error or a error in try catch
 *
 * NOTE:  any other reponse than 200 or 426, will be treated same as 200 by us and the same repsonse will be returned, you have to gandle that case
 *
 *
 *
 * @typedef {Object} ApiResponse
 * @property {string} message
 * @property {string} new_encrypted_key 
 * @property {number} status_code 
 *
 *
 * @argument {Response} response 
 * @argument {string} oldKey 
 * @argument {()=>Promise<Response>} previousFetchFunc 
 *
 * @returns {Promise<Response|Error>}
 *
 */
async function updateTheKeyToNewValueIfUpgradeIsRequeired(response, oldKey, previousFetchFunc) {
  console.log(`in the update the key func `)
  if (response.status !== 426) {
    console.log(`the staus code is not eqaul to 426 and returning the response`)
    return response
  }

  let fetchFunc = fetchFunctionBuilder("getNewKey", "POST", { "user_key": oldKey })
  let responseFromNewKey = await fetchFunc()

  // the repsonse form this can be 400, 500, 403, 200 
  // in the case of 403 I am updatign the key when It is valid, in that case log out a error as this condition should not happen and not we can't update the 
  // key and we got a error respnse now, so return err in this case as we can't return the response too

  /** @type ApiResponse*/
  let apiRespBody = await responseFromNewKey.json()

  console.log(` \n\n the api respnse status is ${apiRespBody.status_code} -- message is ${apiRespBody.message} -- new_encrypted_key is ${apiRespBody.new_encrypted_key} \n\n `)

  if (responseFromNewKey.status === 403) {
    console.error(`\n\n we got the error as the status code is 403 and the response body is ${JSON.stringify(responseFromNewKey)} \n\n `)
    return new Error(" we got the error as the status code is 403, we can't update the key as they are still valid")
  } else if (responseFromNewKey.status === 500) {
    console.error(`\n\n we got the error as the status code is 500 and the response body is ${JSON.stringify(responseFromNewKey)} \n\n `)
    return new Error(" we got the error as the status code is 500 ")
  } else if (responseFromNewKey.status === 400) {
    console.error(`\n\n we got the error as the status code is 400 and the response body is ${JSON.stringify(responseFromNewKey)} \n\n `)
    return new Error(" we got the error as the status code is 400 ")
  }

  // the response is probally 200
  if (responseFromNewKey.status !== 200 || apiRespBody.new_encrypted_key === "") {
    console.error(`\n\n we got the error as the status code is ${responseFromNewKey.status} and the response body is ${JSON.stringify(responseFromNewKey)} \n\n `)
    return new Error(" we got the error as the status code is not 200(it's not 400,403,500 either) or the key form the backend is empty ")
  }

  // the key is there and the response is 200
  // saving it to the local storage
  let err = saveValueToTheStorage("key", apiRespBody.new_encrypted_key, () => { })
  if (err !== null) {
    console.error(`there is a error in saving the new key in the localStorage ->${err}`);
    return err
  }
  let newResponseWithoutNewKeyError = await previousFetchFunc()
  console.log(`the response form re running the previousFetchFunc is ${newResponseWithoutNewKeyError.status}-- ok -> ${newResponseWithoutNewKeyError.ok}  `)
  return newResponseWithoutNewKeyError
}


// --------


/**
 * this function is to run when we have gotten a 403 status code, get a new key for now and when the user gets a new key it will
 * be override form next time with chrome extension one; 
 */








// -----

let config = {}

/**
 * @argument {string} pathWithoutBackSlash - url path of the function
 * @param {"POST"|"GET"} method - request method
 * @param {Object} [header={ "Content-Type": "application/json" } ] - request header
 * @param {Object} bodyOBJ - request body
 *
 * @returns {() => Promise<Response>}
 */
export function fetchFunctionBuilder(
  pathWithoutBackSlash: string,
  method: "POST" | "GET",
  bodyOBJ: Object,
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

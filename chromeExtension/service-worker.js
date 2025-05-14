// @ts-check

import {
  getDefaultValueOfToSkipTheSponsorAndShowTheModal,
  getKeyFromStorageOrBackend,
  getValueFromTheStorage,
  getWhereToSkipInYtVideo,
  saveValueToTheStorage
} from './helper.js';
import { config } from './config.js';


console.log("hi from the service worker and will run say hi() now");
/**
 * Configuration object for authentication
 * @typedef {Object} Config
 * @property {string} BACKEND_URL - Backend service URL
 */
// const config = {
//     BACKEND_URL: "http://localhost:8080",
// };

// Initialize the config
/**
 * @typedef {Object} GetKeyMessage
 * @property {'getKeyFromStorageOrBackend'} type - Message type identifier
 */

/**
 * Message handler for Chrome extension background script
 * @param {GetKeyMessage} request - The message request object
 * @param {chrome.runtime.MessageSender} sender - Message sender information
 * @param {(response?: any) => void} sendResponse - Callback to send response
 * @returns {boolean} - Return true to indicate async response
 */
chrome.runtime.onMessage.addListener((
  /**
   * Message handler for Chrome extension background script
   * @type {GetKeyMessage} request - The message request object
   * @param {chrome.runtime.MessageSender} sender - Message sender information
   * @param {(response?: any) => void} sendResponse - Callback to send response
   * @returns {boolean} - Return true to indicate async response
   */
  // @ts-ignore
  request, sender, sendResponse) => {

  if (request.type === "getKeyFromStorageOrBackend") {
    console.log("Received message in background script:", request);
    // Execute the key fetch function and handle the response
    getKeyFromStorageOrBackend(config)
      .then(([key, error]) => {
        console.log("Key fetch completed", { key: key?.substring(0, 10), error });
        sendResponse([key, error]);
      })
      .catch(error => {
        console.error("Error in background script:", error);
        sendResponse(["", error]);
      });

    // Return true to indicate we will send response asynchronously
    return true;
  }
});

/**
 * @typedef {Object} GetKeyMessage2
 * @property {'getKeyFromStorage'} type - Message type identifier
 */

/**
 * Message handler for Chrome extension background script
 * @param {GetKeyMessage} request - The message request object
 * @param {chrome.runtime.MessageSender} sender - Message sender information
 * @param {(response?: any) => void} sendResponse - Callback to send response
 * @returns {boolean} - Return true to indicate async response
 */
chrome.runtime.onMessage.addListener((
  /**
   * Message handler for Chrome extension background script
   * @type {GetKeyMessage2} request - The message request object,
   * @typedef {chrome.runtime.MessageSender} sender - Message sender information
   * @typedef {(response?: any) => void} sendResponse - Callback to send response
   * @returns {boolean} - Return true to indicate async response
   */
  // @ts-ignore
  request, sender, sendResponse) => {

  if (request.type === "getKeyFromStorage") {
    console.log("Received message in background script:", request);
    // Execute the key fetch function and handle the response
    getValueFromTheStorage("key", () => { })
      .then(([key, error]) => {
        console.log("Key fetch completed for the event ->", { key: key?.substring(0, 10), error });
        sendResponse([key, error]);
      })
      .catch(error => {
        console.error("Error in background script:", error);
        sendResponse(["", error]);
      });

    // Return true to indicate we will send response asynchronously
    return true;
  }
});


/**
 * @typedef {Object} ResponseObject
 * @property {number[]} [skipPoints] - Array of timestamps where to skip
 * @property {string} [error] - Error message if something went wrong
 */



/**
 * @typedef {Object} MessageRequest
 * @property {string} type - The type of message being sent
 * @property {string} encKey - The encryption key
 * @property {string} videoID - The ID of the YouTube video
 * @property {string} jsonStringifiedCaptions - JSON Stringifiewd captions 
 */

/**
 * @typedef {function([ResponseObject|null, Error|null]): void} sendResponse
 */

/**
 * @callback MessageCallback
 * @param {MessageRequest} request - The message request object
 * @param {chrome.runtime.MessageSender} sender - The message sender
 * @param {function([ResponseObject|null, Error|null]): void} sendResponse - Callback to send response
 * @returns {boolean} - Return true to indicate async response
 */

/**
 * Handles incoming messages in the Chrome extension background script
 * @type {MessageCallback}
 */

chrome.runtime.onMessage.addListener((
  /**@type {MessageRequest} request */
  request,
  /** @type {chrome.runtime.MessageSender} sender */
  sender,
  /**@type {sendResponse} sendResponse*/
  sendResponse) => {

  if (request.type === "getWhereToSkipInYtVideo") {
    // Execute the key fetch function and handle the response
    // before we call this we would need to get the youtbe captions
    (async () => {
      console.log(`getting the active tab`)
      let activeTab = await chrome.tabs.query({ active: true, currentWindow: true })
      console.log(`the active tab is ${JSON.stringify(activeTab)} and it's length is ${activeTab.length}`)
      if (activeTab.length < 1) {
        sendResponse([null, new Error(`we were not able to find the active Tab as the value is <1`)]);
        return
      }
      let doesActiveTabHasVideoId
      /** @type {number|undefined}*/
      let activeTabId
      activeTab.find((tab) => {
        if (tab.url?.includes(request.videoID)) {
          console.log(`the active tab has the video id and it is ->${JSON.stringify(tab)}`)
          doesActiveTabHasVideoId = true
          activeTabId = tab.id
        }
      })
      if (!doesActiveTabHasVideoId && activeTabId === undefined) {
        console.log(`the active tab is is undefined or not there`)
        sendResponse([null, new Error(`the active tab is is undefined or not there`)]);
        return
      }

      const res = await chrome.scripting.executeScript({
        target: { tabId: activeTabId, allFrames: false },
        world: "MAIN",
        func: () => {
          console.log(`------Hi form the service worker in the youtube ----`)
          // @ts-ignore
          const tracklist = window.ytInitialPlayerResponse.captions.playerCaptionsTracklistRenderer;
          return JSON.stringify(tracklist) || null;
        }

      });
      if (res.length === 0) {
        console.log(`there is a error in getting the result form the youtube script `)
        sendResponse([null, new Error(`there is a error in getting the result form the youtube script `)]);
        return
      }
      let result = res[0]
      if (result.result === null || result.result === undefined) {
        sendResponse([null, new Error(`there is a error in getting the result form the youtube script (it returned null)`)]);
        return
      }
      console.log(`the script executioon is completed and it is ${res}`)
      console.log(`the script executioon is completed and it is (in json string) ->\n\n  ${JSON.stringify(res)} \n\n\n`)

      console.log("Received message in background script:", request);
      getWhereToSkipInYtVideo(request.encKey, request.videoID, result.result)
        .then(([responseObject, error]) => {
          console.log("Key fetch completed for where to skip in the video", { key: responseObject, error });
          sendResponse([responseObject, error]);
        })
        .catch(error => {
          console.error("Error in background script:", error);
          sendResponse([null, error]);
        });

    })();
    // Return true to indicate we will send response asynchronously
    return true;
  }
});
/**
 * @typedef {Object} MessageRequest2
 * @property {string} type - The type of message being sent
 *
 */
chrome.runtime.onMessage.addListener((
    /** @type {{ type: string; }} */ request, /** @type {any} */ sender, /** @type {( response:[boolean, Error|null] ) => void } */ sendResponse) => {
  if (request.type === "alwaysSkipTheSponsorAndDoNotShowTheModal") {
    console.log(`got the request in the  alwaysSkipTheSponsorAndDoNotShowTheModal`, request);

    getDefaultValueOfToSkipTheSponsorAndShowTheModal().then(([value, error]) => {
      console.log("Error in background script while getting the default value of to skip modal or not:->", error, "  and the value is -->", value);
      if (error !== null && error instanceof Error) {
        return sendResponse([Boolean(value), error])
      }
      return sendResponse([Boolean(value), null])
    })
    return true;
  }
  return false;

})

/**
 * @typedef {Object} saveValueInStorage
 * @property {'saveValueInStorage'} type - The type identifier for saving value to storage
 * @property {string} key - The storage key to save the value under
 * @property {any} value - The value to be saved in storage
 * @export
 */

/**
 * Example usage in the message listener:
 */
chrome.runtime.onMessage.addListener((
    /** @type { saveValueInStorage} */ request,
    /** @type {any} */ sender,
    /** @type {(response: Error|null) => void} */ sendResponse
) => {
  if (request.type === "saveValueInStorage") {
    try {
      const error = saveValueToTheStorage(request.key, request.value);
      console.log("error in storing the value in the db is -->", error, "\n and the key ->", request.key, " and the value was ->", request.value);
      return sendResponse(error);
    } catch (error) {
      console.log("error in the try catch -->", error);
      return error
    }
  }
});

/**
 * @typedef {Object} InstallDetails
 * @property {'install' | 'update' | 'chrome_update' | 'shared_module_update'} reason - The reason for the installation event
 * @property {string} [previousVersion] - Previous version of the extension, if this is an update
 * @property {string} [id] - ID of the shared module that was updated
 */
/**
 * Handles the extension's installation or update events.
 * Initializes storage with default values and shows welcome page on first install.
 * @type {InstallDetails} detail - Details about the installation event
 */

chrome.runtime.onInstalled.addListener((
  /** @type {InstallDetails} detail - Details about the installation event */
  detail) => {
  if (detail.reason === 'install') {
    getKeyFromStorageOrBackend(config)
      .then(([key, error]) => {
        console.log("Key fetch completed on the first run ", { key: key?.substring(0, 10), error });
      })
      .catch(error => {
        console.error("Error in background script while getting the key on first run:", error);
      });

  }
})


/**
 * @typedef {Object} MessageRequest3
 * @property {string} action - The action to perform (e.g., 'getKey')
 */

/**
 * @typedef {Object} MessageResponse
 * @property {string} data - The sensitive key data
 */

/**
 * @typedef {Object} StorageData
 * @property {string} key - The sensitive key stored in extension storage
 */

/**
 * Handles external messages from whitelisted domains
 * @param {MessageRequest3} request - The request message from the website
 * @param {chrome.runtime.MessageSender} sender - Information about the message sender
 * @param {function(MessageResponse?): void} sendResponse - Callback to send response
 * @returns {boolean} - Must return true if response is async
 */
chrome.runtime.onMessageExternal.addListener(
  function (
    /**@type {MessageRequest3} */ request,
    /**@type {chrome.runtime.MessageSender} */ sender,
    /**@type  {function(MessageResponse?): void}  */ sendResponse) {
    // Validate sender origin
    if (sender.origin !== config.websiteURL) {
      console.error(`Unauthorized access attempt from ${sender.origin}`);
      sendResponse(null);
      return false;
    }

    // Validate request action
    if (request.action !== "getKey") {
      console.error(`Invalid action requested: ${request.action}`);
      sendResponse(null);
      return false;
    }

    /**
     * Retrieve key from storage and send response
     * @param {StorageData} result - The data retrieved from storage
     */
    function handleStorageData(result) {
      if (chrome.runtime.lastError) {
        console.error('Storage error:', chrome.runtime.lastError);
        sendResponse(null);
        return;
      }

      sendResponse({ data: result.key });
    }

    // Get data from storage
    chrome.storage.local.get(['key'], handleStorageData);
    return true; // Required for async response
  }
);

/**
 * @typedef {Object} RequestObjForChangingKeyOnPayment
 * @property {'paymentReceivedChangeTheKey'} type - The type of the request
 * @property {string} key - The new key to be stored
 */

/**
 * @typedef {Object} responseFromChangingKeyOnPayment
 * @property {boolean} success - Indicates whether the key was successfully changed
 * @export 
 */

/**
 * @typedef {Object} KeyUpdateStatusMessage
 * @property {'keyUpdateStatus'} type - The type of the status update message
 * @property {boolean} success - Whether the key update was successful
 */



chrome.runtime.onMessage.addListener(
  /**
* Message listener for handling key change requests from content script
* @param {RequestObjForChangingKeyOnPayment} request - The request object from content script
* @param {chrome.runtime.MessageSender} sender - Information about the sender
* @param {function(responseFromChangingKeyOnPayment): void} sendResponse - Function to call with the response
//   @ returns {boolean} - Return true to indicate async response
*/
  (request, sender, sendResponse) => {
    if (request.type === "paymentReceivedChangeTheKey") {
      // Save the new key to storage
      console.log("in the paymentReceivedChangeTheKey", request);
      if (request.key === "" || request.key === null) {
        return sendResponse({ success: false });
      }
      chrome.storage.local.set({ key: request.key }, () => {
        if (chrome.runtime.lastError) {
          console.error('Error setting key:', chrome.runtime.lastError);
          return sendResponse({ success: false });
        }
        // Send response to content script
        sendResponse({ success: true });
      });
    }
    return true; // for the async support 
  }
);



//@ts-ignore
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "navigateToWebsite") {
    try {
      console.log(`the user send us a action to go to our website`)
      // Create a new tab with the specified URL
      chrome.tabs.create({ url: config.websiteURL });
      sendResponse({ success: true });
    } catch (error) {
      console.log(`the error in sending user to the new page is ${error}`)
      sendResponse({ success: false, error: error })
    }
  }
  // Return true to indicate you'll send a response asynchronously
  return true;
});


// get the 
//@ts-ignore
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === "ee") {
    try {
      console.log(`00-0------00 injectig the script`)
      console.log(`the user send us a action to go to our website`)
      // Create a new tab with the specified URL
      chrome.tabs.create({ url: config.websiteURL });
      sendResponse({ success: true });
    } catch (error) {
      console.log(`the error in sending user to the new page is ${error}`)
      sendResponse({ success: false, error: error })
    }
  }
  // Return true to indicate you'll send a response asynchronously
  return true;
});





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

// NOTE:
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
      // let activeTab = await chrome.tabs.query({ active: true, currentWindow: true })
      let activeTab = await chrome.tabs.query({})
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

      // 2. --- SETUP WEB REQUEST INTERCEPTOR (The Core Logic) ---
      let isIntercepted = false;
      let storedCaptionsData = null;
      /** @type {Promise<string>} */
      let interceptPromise = new Promise((resolve) => {

        /** @type {function(chrome.webRequest.WebRequestBodyEvent): Promise<void>} */
        const listener = async (details) => {
          // Check if it's the specific URL we need
          if (details.url.startsWith("https://www.youtube.com/api/timedtext")) {
            if (isIntercepted) {
              return;
            }
            chrome.webRequest.onBeforeSendHeaders.removeListener(listener);
            console.log("Interceptor Fired: Captions Request Found!", details.requestHeaders);

            const reqHeaders = details.requestHeaders;
            // 1. Make the new fetch request to your service using the original URL
            const fetchUrl = details.url;
            isIntercepted = true;
            try {
              const fetchResponse = await fetch(fetchUrl);
              // 2. Store the response body in a variable
              // This is the response body of the original YouTube Captions request!
              const captionsData = await fetchResponse.json();
              // const a = await fetchResponse.json()

              // Store the data in the global variable and resolve the promise
              try {
                storedCaptionsData = captionsData.events;
              } catch (e) {
                storedCaptionsData = "";
              }
              resolve(storedCaptionsData);

            } catch (e) {
              console.error("Fetch failed during interception:", e);
              resolve(""); // Resolve with empty string or handle error
            }

            // 3. IMPORTANT: Remove the listener immediately after the first intercept
            chrome.webRequest.onBeforeSendHeaders.removeListener(listener);
          }
        };

        // Add the listener. Using onBeforeSendHeaders gives you request headers.
        // The URL is sufficient to make a new request, but we capture the headers anyway.
        chrome.webRequest.onBeforeSendHeaders.addListener(
          //@ts-ignore
          listener,
          { urls: ["*://www.youtube.com/api/timedtext?*"] },
          ["requestHeaders"] // Option needed to see request headers
        );
      });

      console.log("Web Request Interceptor Registered.");

      // 3. --- TRIGGER THE REQUEST VIA SCRIPTING ---

      await chrome.scripting.executeScript({
        target: { tabId: activeTabId, allFrames: false },
        world: "MAIN",
        func: () => {
          console.log("Attempting to click the CC button...");
          /** @type {HTMLElement|null} */
          const ccButton = document.querySelector("#movie_player > div.ytp-chrome-bottom > div.ytp-chrome-controls > div.ytp-right-controls > div.ytp-right-controls-left > button.ytp-subtitles-button.ytp-button");
          if (ccButton === null) {
            console.error("CC button not found");
            return false;
          }
          // This click triggers the network request you want to intercept
          ccButton.click();
          return true;
        }
      });
      // 4. --- WAIT FOR THE INTERCEPTION AND FETCH TO COMPLETE ---

      // Wait for the Promise to resolve, meaning the fetch and storage is done.
      const req = await interceptPromise;
      if (req === "") {
        sendResponse([null, new Error(`getting caotions or fetching it failled(resoonse text is "")`)]);
      }

      console.log("Intercepted Request Data is now stored:", storedCaptionsData ? "Yes" : "No", "and it is  ", req, `\n and it's type is ${typeof req} `, "\n and req is ", request)
      // sendResponse([req, null]);
      getWhereToSkipInYtVideo(request.encKey, request.videoID, req)
        .then(([responseObject, error]) => {
          console.log("Key fetch completed for where to skip in the video", { key: responseObject, error });
          sendResponse([responseObject, error]);
        })
        .catch(error => {
          console.error("Error in background script:", error);
          sendResponse([null, error]);
        });
      return true;
    })();
    // Return true to indicate we will send response asynchronously
    return false;
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



// claude one 
/**
 * @typedef {Object} ClickCCButtonRequest
 * @property {'clickCCButtonAndGetCaptions'} type - Message type identifier
 * @property {string} videoID - The ID of the YouTube video
 */
/**
 * @typedef {Object} ClickCCButtonResponse
 * @property {boolean} success - Indicates if the script execution was attempted.
 * @property {string} [error] - Error message if something went wrong.
 * @property {string} [captionData] - The intercepted caption data from timedtext API.
 */
/**
 * Message handler to click the CC button in the active tab and intercept captions.
 */
chrome.runtime.onMessage.addListener(
  /**
   * @param {ClickCCButtonRequest} request - Information about the sender
   * @param {chrome.runtime.MessageSender} sender - Information about the sender
   * @param {function(ClickCCButtonResponse): void} sendResponse - Function to call with the response
   */
  (request, sender, sendResponse) => {
    if (request.type === "clickCCButtonAndGetCaptions") {
      (async () => {
        console.log("Received request to click CC button:", request);

        // 1. Find the tab containing the videoID
        let activeTab = await chrome.tabs.query({});
        /** @type {number|undefined} */
        let activeTabId;

        // Find the tab with the correct video ID
        activeTab.find((tab) => {
          if (tab.url?.includes(request.videoID)) {
            activeTabId = tab.id;
            return true;
          }
          return false;
        });

        if (activeTabId === undefined) {
          console.error(`Tab with video ID ${request.videoID} not found.`);
          return sendResponse({ success: false, error: "Tab not found" });
        }

        try {
          // 2. Attach debugger to intercept network requests
          await chrome.debugger.attach({ tabId: activeTabId }, "1.3");
          await chrome.debugger.sendCommand({ tabId: activeTabId }, "Network.enable");

          // 3. Set up listener for network responses
          let captionData = null;
          let capturedRequestId = null;

          const networkListener = (source, method, params) => {
            if (source.tabId !== activeTabId) return;

            // Capture the request ID when timedtext request is made
            if (method === "Network.requestWillBeSent" &&
              params.request.url.includes("timedtext") &&
              params.request.url.includes("video.google.com")) {
              console.log("Detected timedtext request:", params.request.url);
              capturedRequestId = params.requestId;
            }

            // Capture the response body when it arrives
            if (method === "Network.loadingFinished" && params.requestId === capturedRequestId) {
              chrome.debugger.sendCommand(
                { tabId: activeTabId },
                "Network.getResponseBody",
                { requestId: params.requestId }
              ).then(response => {
                captionData = response.body;
                console.log("Successfully captured caption data");
              }).catch(err => {
                console.error("Error getting response body:", err);
              });
            }
          };

          chrome.debugger.onEvent.addListener(networkListener);

          // 4. Click the CC button
          await chrome.scripting.executeScript({
            target: { tabId: activeTabId, allFrames: false },
            world: "MAIN",
            func: () => {
              console.log("Attempting to click the CC button...");
              /** @type {HTMLElement|null} */
              const ccButton = document.querySelector("#movie_player > div.ytp-chrome-bottom > div.ytp-chrome-controls > div.ytp-right-controls > div.ytp-right-controls-left > button.ytp-subtitles-button.ytp-button");

              if (ccButton === null) {
                console.error("CC button not found");
                return false;
              }

              ccButton.click();
              return true;
            }
          });

          console.log("CC button clicked, waiting for timedtext response...");

          // 5. Wait for caption data to be captured (with timeout)
          await new Promise(resolve => setTimeout(resolve, 3000));

          // 6. Clean up
          chrome.debugger.onEvent.removeListener(networkListener);
          await chrome.debugger.detach({ tabId: activeTabId });

          if (captionData) {
            sendResponse({ success: true, captionData: captionData });
          } else {
            //@ts-ignore
            sendResponse({ success: true, captionData: null });
          }

        } catch (error) {
          console.error("Error in CC button click and caption capture:", error);
          // Try to detach debugger if still attached
          try {
            await chrome.debugger.detach({ tabId: activeTabId });
          } catch (e) { }
          //@ts-ignore
          sendResponse({ success: false, error: error.toString() });
        }
      })();

      // Return true to indicate we will send response asynchronously
      return true;
    }
  }
);



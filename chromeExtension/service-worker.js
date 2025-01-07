import {getKeyFromStorageOrBackend, getWhereToSkipInYtVideo} from './helper.js';

console.log("hi from the service worker and will run say hi() now");
/**
 * Configuration object for authentication
 * @typedef {Object} Config
 * @property {string} BACKEND_URL - Backend service URL
 */
const config = {
    BACKEND_URL: "http://localhost:8080",
};

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
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log("Received message in background script:", request);

    if (request.type === "getKeyFromStorageOrBackend") {
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
 * @typedef {Object} MessageRequest
 * @property {string} type - The type of message being sent
 * @property {string} encKey - The encryptionkey
 * @property {string} videoID - The ID of the ytVideo
 *
 * @callback MessageCallback
 * @param {MessageRequest} request - The message request object
 * @param {chrome.runtime.MessageSender} sender - The message sender
 * @param {(response: [responseObject|null, Error|null]) => void} sendResponse - Callback to send response
 * @returns {boolean} - Return true to indicate async response
 */

/** @type {MessageCallback} */
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log("Received message in background script:", request);

    if (request.type === "getWhereToSkipInYtVideo") {
        // Execute the key fetch function and handle the response
        getWhereToSkipInYtVideo(request.encKey, request.videoID)
            .then(([responseObject, error]) => {
                console.log("Key fetch completed", { key: responseObject, error });
                sendResponse([responseObject, error]);
            })
            .catch(error => {
                console.error("Error in background script:", error);
                sendResponse([null, error]);
            });

        // Return true to indicate we will send response asynchronously
        return true;
    }
});

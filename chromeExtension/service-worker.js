// @ts-check

import {
    getDefaultValueOfToSkipTheSponsorAndShowTheModal,
    getKeyFromStorageOrBackend,
    getWhereToSkipInYtVideo,
saveValueToTheStorage
} from './helper.js';


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
 * @typedef {Object} ResponseObject
 * @property {number[]} [skipPoints] - Array of timestamps where to skip
 * @property {string} [error] - Error message if something went wrong
 */

/**
 * @typedef {Object} MessageRequest
 * @property {string} type - The type of message being sent
 * @property {string} encKey - The encryption key
 * @property {string} videoID - The ID of the YouTube video
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
    console.log("Received message in background script:", request);
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
/**
 * @typedef {Object} MessageRequest2
 * @property {string} type - The type of message being sent
 *
 */
chrome.runtime.onMessage.addListener((
    /** @type {{ type: string; }} */ request, /** @type {any} */ sender, /** @type {( response:[boolean, Error|null] ) => void } */ sendResponse) => {
    if (request.type === "alwaysSkipTheSponsorAndDoNotShowTheModal") {
        console.log(`got the request in the  alwaysSkipTheSponsorAndDoNotShowTheModal`, request);
        
        getDefaultValueOfToSkipTheSponsorAndShowTheModal().then(([value, error] )=>{
        console.log("Error in background script while getting the default value of to skip modal or not:->", error,"  and the value is -->",value);
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
        console.log("error in storing the value in the db is -->", error, "\n and the key ->", request.key, " and the value was ->",request.value);
        return sendResponse(error);
        } catch (error) {
            console.log("error in the try catch -->",error);
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
    detail)=>{
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
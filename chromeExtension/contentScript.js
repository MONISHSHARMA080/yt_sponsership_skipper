// @ts-check

// the problem seems to be that the value of videoSkipper just goes to be true when we execute the func that is executed 10 sec before
// and when we click the do modal button the modal does not dissappears SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal === false
// --> well the problem is that the eventlistenr is calling the func many times, so maybe add another state or jut make a way to create the modal once


async function main() {
  console.log("in the main");
  let [key, error] = await chrome.runtime.sendMessage({
    type: "getKeyFromStorageOrBackend",
  });
  console.log("the modal is inserted")
  console.log("lets see what we got ");
  if (error) {
    console.log("error is there in getting the key and it is -->", error, "\n\n the key is", key,);
    return;
  }
  console.log("the key is  -->", key, "from the  content script, error is -->", error);
  let [videoID, errorFormGettingVideoID] = getVideoID();
  if (errorFormGettingVideoID) {
    console.log("videoID is  -->", videoID, "and the error getting it is -->", errorFormGettingVideoID);
    return;
  }
  let videoPlayer = getVideoPlayer();
  if (videoPlayer === null) {
    console.log("video player is not there");
    return;
  }
  let [responseObject, errorFromYTApi] = await chrome.runtime.sendMessage({ type: "getWhereToSkipInYtVideo", encKey: key, videoID: videoID, });
  if (errorFromYTApi || responseObject === null) {
    console.log("there is a error in the yt api -->", errorFromYTApi);
    return;
  }
  console.log("the response object is -->", responseObject);
  /** @type {import("./helper").ResponseObject} responseObject */
  let responseOBjectFromYt = responseObject;
  if (responseOBjectFromYt.containSponserSubtitle === false) {
    console.log("the video does not have a sponsorship subtitle");
    return;
  }
  console.log("the func will be called on the time-> ", responseOBjectFromYt.startTime - 10,);
  let [defalutValueToSkipTheModal, errorFormSkippingTheModal] = await chrome.runtime.sendMessage({ type: "alwaysSkipTheSponsorAndDoNotShowTheModal" });
  if (errorFormSkippingTheModal) {
    console.log("there is an error in getting the vlaue form the value of wether we should skip the modal or not -->", errorFormSkippingTheModal, "value of the modal is ");
  }
  console.log(`the value of the alwaysSkipTheSponsorAndDoNotShowTheModal form the storage is ${defalutValueToSkipTheModal} `);
  if (typeof defalutValueToSkipTheModal !== "boolean") {
    defalutValueToSkipTheModal = false
    console.log(`value of the alwaysSkipTheSponsorAndDoNotShowTheModal was not bool so we made it ${defalutValueToSkipTheModal}`);
    chrome.runtime.sendMessage({ type: "saveValueInStorage", key: "alwaysSkipTheSponsorAndDoNotShowTheModal", value: false, reason: "didn't find the value to be bool so have to set it " })
  }

  /** @type {skippedTheSponser} */
  const SkippedVideoSponsorOBJ = { videoSponsorSkipper: false, callBackBeforeSomeTimeOfSponsor: false, alwaysSkipTheSponsorAndDoNotShowTheModal: defalutValueToSkipTheModal, userOptedForSkippingTheModal: false, nowShowingTheModal: false };
  // videoSponsorSkipper is to know if we have skipped the sponsor

  videoPlayer.addEventListener("timeupdate", (event) => {
    console.log("\n Current time :-->", videoPlayer.currentTime); // working
    // the end time is not <= and is < because it will  not move forward if we did not do that, and just to be sure lets make a var too
    skipTheVideo(responseOBjectFromYt, videoPlayer, SkippedVideoSponsorOBJ);
    beforeSomeTimeExecuteSomething(responseOBjectFromYt.startTime - 10, videoPlayer, () => {
      // make the if statement check
      console.log("condition of the modal -->", SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal === false && SkippedVideoSponsorOBJ.userOptedForSkippingTheModal === false);

      if (SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal === false && SkippedVideoSponsorOBJ.userOptedForSkippingTheModal === false) {
        // to do: update the modal state form here on and when the modal is created then update the if conditon as I do not want to createmultiple of it or maybe do it in the modal func
        embeddTheModalInThePage(SkippedVideoSponsorOBJ)
      }
      if (SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal === true && SkippedVideoSponsorOBJ.userOptedForSkippingTheModal === true) {
        // should make it more explicit (use it)
        console.log("in the else/other block");
        SkippedVideoSponsorOBJ.videoSponsorSkipper = false // meaning: bro just skip the video 
      }
    },
      SkippedVideoSponsorOBJ
    );
  });
  // now make 2 function to abstract the logic here, one is for skipping and one takes in  a callback that will execute, if the video is certain sec
  // before , eg 10 sec before certain time do something
}
try {
  let isEventListenerAdded = { isEventListenerAdded: false }
  listenAndReplyToTheSvelteMessage(isEventListenerAdded).then((a) => { console.log("the listen to svelte function is over ->", a) })
  // addEventListenerForClosingAllEventListener(isEventListenerAdded)
  addEventListenerForChangingTheKeyFromSvelte()
} catch (error) {
  console.log("there is a error in listen and replyu to svelte function ->", error);
}
try {
  Promise.all([
    main(),

  ]).then(([mainResult]) => {
    console.log("main finished and the returned value is -->", mainResult, "\n\n and the ");
  })
} catch (e) {
  console.log("error in main script:", e);
}

/**
 * @return {[string,Error|null]}
 */
function getVideoID() {
  let url = window.location.href;
  console.log("the url in getVideoID():", url);
  if (url === null || url === "" || url === undefined) {
    return ["", new Error("the url is not there")];
  }
  return [url, null];
}


function getVideoPlayer() /** @return {HTMLVideoElement|null}*/ {
  try {
    // const video = document.querySelector("video");
    // video.addEventListener("timeupdate", (event) => {
    //   console.log("The currentTime attribute has been updated. Again.-->", event);
    // });
    /** @type {HTMLVideoElement |null} */
    const videoElement = document.querySelector("video.html5-main-video");
    if (videoElement === null || videoElement === undefined) {
      return null;
    }
    //  videoElement.addEventListener("timeupdate", (event) => {
    //   console.log("\n Current time :-->", videoElement.currentTime); // working
    // });
    return videoElement;
  } catch (e) {
    console.log(
      "error occurred in getting the current time form the player->",
      e,
    );
    return null;
  }
}

/**
 *@typedef {Object} skippedTheSponser - obj to see (pass by ref) if the video is skipped
 * @property {Boolean}videoSponsorSkipper
 * @property {Boolean} callBackBeforeSomeTimeOfSponsor
 * @property {Boolean} alwaysSkipTheSponsorAndDoNotShowTheModal
 * @property {Boolean} userOptedForSkippingTheModal
 * @property {Boolean} nowShowingTheModal 
 *
 * @argument {skippedTheSponser} SkippedVideoSponsorOBJ
 * @argument {import("./helper").ResponseObject}responseObjectFromYt
 * @argument {HTMLVideoElement} videoPlayer
 */
function skipTheVideo(responseObjectFromYt, videoPlayer, SkippedVideoSponsorOBJ) {
  console.log("in the Skip the video func\n\n")
  console.log(`the cureent time is ${videoPlayer.currentTime} and end time is ${responseObjectFromYt.endTime} and start Time is ${responseObjectFromYt.startTime} and video skipper is ${SkippedVideoSponsorOBJ.videoSponsorSkipper} `)
  if (videoPlayer.currentTime >= responseObjectFromYt.startTime && videoPlayer.currentTime < responseObjectFromYt.endTime && SkippedVideoSponsorOBJ.videoSponsorSkipper === false) {
    console.log(`the video player time is greater that or = the start time of from the backend -->time in the
       video player ${videoPlayer.currentTime} ----start time is  ${responseObjectFromYt.startTime} --  `);
    console.log("\n now going to skip in the video to ", responseObjectFromYt.endTime,);
    videoPlayer.currentTime = responseObjectFromYt.endTime;
    SkippedVideoSponsorOBJ.videoSponsorSkipper = true;
  }
}

/**
 * @argument {Number} timeToCallTheFunc - call the func at sponsorShipStart - 10
 * @argument {HTMLVideoElement} videoPlayer
 * @param {Function} callbackFunction
 * @param {skippedTheSponser} skippedTheSponsorOBJ
 */
// * // @argument {Number} callTheFuncBeforeThisTime - eg call the function before the
// * @param {skippedTheSponsor} SkippedVideoSponsorOBJ - pass by reference boolean value in a Object
function beforeSomeTimeExecuteSomething(timeToCallTheFunc, videoPlayer, callbackFunction, skippedTheSponsorOBJ) {
  // if the time of video is > than the time to call  and is lower than the time to to call before (eg before 10 sec) execute the function
  console.log("in the video player func ->");

  if (videoPlayer.currentTime >= timeToCallTheFunc && skippedTheSponsorOBJ.videoSponsorSkipper === false) {
    // skippedTheSponsorOBJ.videoSponsorSkipper = false; // first as if the func throw we would not be able to update the state
    try {
      callbackFunction()
    } catch (e) {
      console.log(`error in the callback function:-> ${e} `)
    }
  }
}

/**
 * I should probably store some var
 * @param {Function} onCloseFunction - func that will execute when the user closes or declines the modal
 * @param {Function} onUserClickDoNotSkipTheSponsorShipFunc - func that will execute when the user wants to not skip the sponsor on the video
 * @param {Function} doNotShowThisModalAgainAlwaysSkipFunction - func that will execute when the user wants does not want to see this again (skip the sponsor everytime)
 * @returns {HTMLDivElement}
 */
function createSponsorShipModalToTellUserWeAreAboutToSkip(onCloseFunction, onUserClickDoNotSkipTheSponsorShipFunc, doNotShowThisModalAgainAlwaysSkipFunction) {
  const modalContainer = document.createElement('div');
  modalContainer.style.position = 'fixed';
  modalContainer.style.top = '18px';
  modalContainer.style.right = '24px';
  modalContainer.style.width = '256px';
  modalContainer.style.backgroundColor = '#01044a';
  modalContainer.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
  modalContainer.style.borderRadius = '28px';
  modalContainer.style.overflow = 'hidden';
  modalContainer.style.zIndex = '90000';
  // modalContainer.style.display = 'none';
  // id
  modalContainer.id = "sponserShipModal"
  // Add animation properties
  modalContainer.style.transition = 'opacity 0.3s, transform 0.3s';
  modalContainer.style.opacity = '0';
  modalContainer.style.transform = 'translateY(-20px)';

  // Create the content wrapper
  const contentWrapper = document.createElement('div');
  contentWrapper.style.padding = '16px';
  modalContainer.appendChild(contentWrapper);

  // Create the header section
  const header = document.createElement('div');
  header.style.display = 'flex';
  header.style.justifyContent = 'space-between';
  header.style.alignItems = 'center';
  header.style.marginBottom = '8px';
  contentWrapper.appendChild(header);

  // Add the title
  const title = document.createElement('h');
  title.style.fontSize = '20px';
  title.style.fontWeight = '730';
  title.style.color = '#d9d9d9';
  title.textContent = 'Sponsorship Alert';
  header.appendChild(title);

  // Add the close button
  const closeButton = document.createElement('button');
  closeButton.style.position = 'relative';
  closeButton.style.color = '#d9d9d9';
  closeButton.style.cursor = 'pointer';
  closeButton.style.border = 'none';
  closeButton.style.background = 'transparent';
  closeButton.style.padding = '8px';
  closeButton.style.borderRadius = '50%';
  closeButton.style.transition = 'background-color 0.3s, transform 0.2s';
  closeButton.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" style="width: 16px; height: 16px;"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>';

  closeButton.addEventListener('mouseover', () => {
    closeButton.style.backgroundColor = 'rgba(255, 255, 255, 0.2)';
    closeButton.style.transform = 'scale(1.1)';
  });

  closeButton.addEventListener('mouseout', () => {
    closeButton.style.backgroundColor = 'transparent';
    closeButton.style.transform = 'scale(1)';
  });

  // Add click event to hide the modal with animation
  closeButton.addEventListener('click', () => {
    modalContainer.style.opacity = '0';
    modalContainer.style.transform = 'translateY(-20px)';

    onCloseFunction()
    setTimeout(() => {
      modalContainer.style.display = 'none';
    }, 300); // Wait for animation to complete
  });
  header.appendChild(closeButton);

  // Add the message
  const message = document.createElement('p');
  message.style.fontSize = '16px';
  message.style.color = '#d9d9d9';
  message.style.marginBottom = '12px';
  message.textContent = "A sponsorship is coming up. We're about to skip it for you.";
  contentWrapper.appendChild(message);

  // Add the "Don't Skip" button
  const dontSkipButton = document.createElement('button');
  dontSkipButton.style.display = 'block';
  dontSkipButton.style.width = '100%';
  dontSkipButton.style.padding = '8px';
  dontSkipButton.style.fontSize = '12px';
  dontSkipButton.style.fontWeight = '500';
  dontSkipButton.style.color = 'white';
  dontSkipButton.style.backgroundColor = '#3904c9';
  dontSkipButton.style.border = 'none';
  dontSkipButton.style.borderRadius = '24px';
  dontSkipButton.style.cursor = 'pointer';
  dontSkipButton.style.transition = 'background-color 0.3s';

  dontSkipButton.addEventListener('mouseover', () => {
    dontSkipButton.style.backgroundColor = '#1e40af';
  });

  dontSkipButton.addEventListener('mouseout', () => {
    dontSkipButton.style.backgroundColor = '#1d4ed8';
  });

  dontSkipButton.textContent = "Don't skip this sponsorship";
  dontSkipButton.addEventListener('click', () => {
    onUserClickDoNotSkipTheSponsorShipFunc ? onUserClickDoNotSkipTheSponsorShipFunc() : null; // if I did not pass it then do not play it
    console.log("User chose not to skip the sponsorship");
    modalContainer.style.opacity = '0';
    modalContainer.style.transform = 'translateY(-20px)';
    setTimeout(() => {
      modalContainer.style.display = 'none';
    }, 300);
  });
  contentWrapper.appendChild(dontSkipButton);
  // Add the "Don't Skip" button
  const alwaysSkipTheSponsorButton = document.createElement('button');
  alwaysSkipTheSponsorButton.style.display = 'block';
  alwaysSkipTheSponsorButton.style.width = '100%';
  alwaysSkipTheSponsorButton.style.padding = '8px';
  alwaysSkipTheSponsorButton.style.fontSize = '12px';
  alwaysSkipTheSponsorButton.style.fontWeight = '500';
  alwaysSkipTheSponsorButton.style.marginTop = '8px';
  alwaysSkipTheSponsorButton.style.color = 'white';
  alwaysSkipTheSponsorButton.style.backgroundColor = '#3904c9';
  alwaysSkipTheSponsorButton.style.border = 'none';
  alwaysSkipTheSponsorButton.style.borderRadius = '24px';
  alwaysSkipTheSponsorButton.style.cursor = 'pointer';
  alwaysSkipTheSponsorButton.style.transition = 'background-color 0.3s';

  alwaysSkipTheSponsorButton.addEventListener('mouseover', () => {
    alwaysSkipTheSponsorButton.style.backgroundColor = '#1e40af';
  });

  alwaysSkipTheSponsorButton.addEventListener('mouseout', () => {
    alwaysSkipTheSponsorButton.style.backgroundColor = '#1d4ed8';
  });

  alwaysSkipTheSponsorButton.textContent = "Always skip the sponsorship";
  alwaysSkipTheSponsorButton.addEventListener('click', () => {
    doNotShowThisModalAgainAlwaysSkipFunction(); // if I did not pass it then do not play it
    console.log("in the always skip the sponsor function");
    modalContainer.style.opacity = '0';
    modalContainer.style.transform = 'translateY(-20px)';
    setTimeout(() => {
      modalContainer.style.display = 'none';
    }, 300);
  });
  contentWrapper.appendChild(alwaysSkipTheSponsorButton);
  // Trigger entrance animation after a brief delay
  setTimeout(() => {
    modalContainer.style.opacity = '1';
    modalContainer.style.transform = 'translateY(0)';
  }, 100);
  // close the modal after 10 sec
  setTimeout(() => {
    modalContainer.style.display = 'none';
  }, 10000); // Wait for animation to complete
  return modalContainer;
}

/** @param {skippedTheSponser} SkippedVideoSponsorOBJ */
function embeddTheModalInThePage(SkippedVideoSponsorOBJ) {
  let moadlElement = document.getElementById("sponserShipModal")
  if (moadlElement !== null) {
    console.log("the modal is there and not goung to create a new one, returning");

    return
  }
  document.body.appendChild(
    createSponsorShipModalToTellUserWeAreAboutToSkip(
      () => {
        console.log("on close func, doing nothing as we want to skip the sponsor")
      }, () => {
        console.log("in the do not skip this sponsor function");
        SkippedVideoSponsorOBJ.videoSponsorSkipper = true // meaning do not skip the sponsor
        SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal = false// doing this to make the modal dissappear after clicking on it 
      }, () => {
        // just raw dawg the storage update or do it in the backgroundScript for the clean code -->  key=alwaysSkipTheSponsorAndDoNotShowTheModal
        console.log("\n\n ++++++++++++");

        console.log("in the do not show this modal ever again function");
        chrome.runtime.sendMessage({ type: "saveValueInStorage", key: "alwaysSkipTheSponsorAndDoNotShowTheModal", value: true })
        console.log("the video skipper is-->", SkippedVideoSponsorOBJ.videoSponsorSkipper);

        SkippedVideoSponsorOBJ.videoSponsorSkipper = false // meaning: bro just skip the video 
        // turning it to be true for the modal
        SkippedVideoSponsorOBJ.userOptedForSkippingTheModal = true// doing this to make the modal dissappear after clicking on it 
        SkippedVideoSponsorOBJ.alwaysSkipTheSponsorAndDoNotShowTheModal = true // we can do it cause user just said so
        console.log("the video skipper is-->", SkippedVideoSponsorOBJ.videoSponsorSkipper);
        console.log("++++++++++++\n\n");
      }
    ))
}

/** 
 * @typedef {Object} ObjToSeeIfEventListenerIsAdded 
 * @property {boolean} isEventListenerAdded
 * 
 * @argument {ObjToSeeIfEventListenerIsAdded} ObjToSeeIfEventListenerIsAdded
*/
async function listenAndReplyToTheSvelteMessage(ObjToSeeIfEventListenerIsAdded) {
  try {
    const response = await chrome.runtime.sendMessage({ type: "getKeyFromStorage" });
    const [key, error] = response;

    console.log("the response is ->", key, " error is -> ", error);

    if (error !== null || key === "") {
      console.log("Error getting key:", error);
      return;
    }

    if (ObjToSeeIfEventListenerIsAdded.isEventListenerAdded === true) return;

    // Define message handler
    // @ts-ignore
    const messageHandler = (event) => {
      // Verify origin
      if (event.origin !== window.location.origin) {
        console.log("the event is from a different origin and the origin is ->", event.origin);
        return;
      }
      console.log("event is ->", event, "\n\n\n and the event event.data is ->", event.data, " -- and event.data.type is ->", event.data.type);

      // Handle GET_KEY message
      if (event.data.type === "GET_KEY") {
        window.postMessage(
          { type: "GET_KEY", key: key },
          window.location.origin
        );
        // Remove the listener after sending the key
        window.removeEventListener('message', messageHandler);
        ObjToSeeIfEventListenerIsAdded.isEventListenerAdded = false;
      }
    };

    // Add the event listener
    window.addEventListener('message', messageHandler);
    ObjToSeeIfEventListenerIsAdded.isEventListenerAdded = true;
    return true;
  } catch (error) {
    console.error("Error in listenAndReplyToTheSvelteMessage:", error);
    return false;
  }

}


/**
 * Removes the event listener
 * @param {ObjToSeeIfEventListenerIsAdded} ObjToSeeIfEventListenerIsAdded
 * @argument {any}  messagehandler
 */
function removeMessageListener(ObjToSeeIfEventListenerIsAdded, messagehandler) {
  try {
    if (!ObjToSeeIfEventListenerIsAdded.isEventListenerAdded) return;
    window.removeEventListener('message', messagehandler);
    ObjToSeeIfEventListenerIsAdded.isEventListenerAdded = false;
    console.log("removing all the eventlistener from the contentScript");
  } catch (error) {
    console.error("Error removing event listener:", error);
  }

}

/**
 * Removes the event listener
 * @argument {ObjToSeeIfEventListenerIsAdded} ObjToSeeIfEventListenerIsAdded
 * @argument {(Event)} messageHandler
 */
function addEventListenerForClosingAllEventListener(ObjToSeeIfEventListenerIsAdded, messageHandler) {
  try {
    // @ts-ignore
    const cleanupHandler = (event) => {
      if (event.data.type === 'removeAllEventListener') {
        console.log("the message received for removing all the eventlistener in the contentscript");
        removeMessageListener(ObjToSeeIfEventListenerIsAdded, messageHandler);
        // Remove this cleanup handler as well
        window.removeEventListener('message', cleanupHandler);
      }
    };

    window.addEventListener('message', cleanupHandler);
  } catch (error) {
    console.log("error in func that listens to the eventlistener and will remove all the eventlistener ->", error);
  }
}


/**
 * function add event listener such that svelte one will send us a event to change the key when the user updates the teir (makes a payment) and will send us 
 * a key for that
 */
function addEventListenerForChangingTheKeyFromSvelte() {
  // how do I make sure that only extension receives this message and not any other extension
  try {
    // @ts-ignore
    const messageHandler = (event) => {
      // Verify origin
      if (event.origin !== window.location.origin) {
        console.log("the event is from a different origin and it is ", event.origin);
        return;
      }
      // Handle GET_KEY message
      if (event.data.type === "paymentReceivedChangeTheKey") {
        // @ts-ignore
        console.log("event is ->", event, "\n\n\n and the event event.data is ->", event.data);
        const newKey = event.data.key;
        if (newKey === "" || newKey === null || newKey === undefined) {
          console.log("the new key is empty");
          return;
        }
        console.log("the new key is ->", newKey);
        // Save the new key to storage
        chrome.runtime.sendMessage({ type: "paymentReceivedChangeTheKey", key: newKey }).then(
          /**@param {import("./service-worker").responseFromChangingKeyOnPayment} response */
          (response) => {
            // toDo: add some logic in the service worker to update the key in the storage
            // close the event listener here 
            // send a message to the svelte that the key has been updated or not
            console.log("response from the service worker on changing the key on payment received ->", response);
            console.log("about to send this this the svelte -> ", { type: "keyChangedOnPaymentReceived", key: newKey, success: response.success },);

            // event that sends a message to the svelte that the key has been updated or not, there make/register the event listener after sending the message
            window.postMessage(
              { type: "keyChangedOnPaymentReceived", key: newKey, success: response.success },
              window.location.origin
            );
            window.removeEventListener('message', messageHandler);

          })
      };
    }
    // Add the event listener
    window.addEventListener('message', messageHandler);
    return true;     // async support  
  } catch (error) {

  }
}

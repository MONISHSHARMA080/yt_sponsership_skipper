// @ts-check
async function main() {
  console.log("in the main");
  let [key, error] = await chrome.runtime.sendMessage({
    type: "getKeyFromStorageOrBackend",
  });
  console.log("the modal is inserted")
  console.log("lets see what we got ");
  if (error) {
    console.log("error is there in getting the key and it is -->", error, "\n\n the key is", key,);return;
  }
  console.log("the key is  -->", key, "from the  content script, error is -->", error);
  let [videoID, errorFormGettingVideoID] = getVideoID();
  if (errorFormGettingVideoID) {
    console.log( "videoID is  -->",  videoID, "and the error getting it is -->", errorFormGettingVideoID);
    return;
  }
  let videoPlayer = getVideoPlayer();
  if (videoPlayer === null) {
    console.log("video player is not there");
    return;
  }
  let [responseObject, errorFromYTApi] = await chrome.runtime.sendMessage({type: "getWhereToSkipInYtVideo", encKey: key, videoID: videoID,});
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

  //   {
  //     "containSponserSubtitle": 200,
  //     "containSponserSubtitle": "subtitles found",
  //     "containSponserSubtitle": 447,
  //     "containSponserSubtitle": 554,
  //     "containSponserSubtitle": true
  // }
  console.log("the func will be called on the time-> ", responseOBjectFromYt.startTime - 10,);

  /** @type {skippedTheSponser} */
  const SkippedVideoSponsorOBJ = { videoSponsorSkipper:false, callBackBeforeSomeTimeOfSponsor:false, alwaysSkipTheSponsorAndDoNotShowTheModal:false};
  videoPlayer.addEventListener("timeupdate", (event) => {
    console.log("\n Current time :-->", videoPlayer.currentTime); // working
    // the end time is not <= and is < because it will  not move forward if we did not do that, and just to be sure lets make a var too
    skipTheVideo(responseOBjectFromYt, videoPlayer, SkippedVideoSponsorOBJ);
    beforeSomeTimeExecuteSomething(responseOBjectFromYt.startTime - 10, videoPlayer, () => {console.log(" hi this func will execute before certain time (10 sec)  ",);
          document.body.appendChild(createSponsorShipModalToTellUserWeAreAboutToSkip(()=>{
              console.log("on close func, doing nothing as we want to skip the sponsor")
              },()=>{
                  SkippedVideoSponsorOBJ.videoSponsorSkipper = true
              }, ()=>{}
              ) )
        // change  the cost obj and probably write it in the storage, so that it can persist
      },
        SkippedVideoSponsorOBJ
    );
  });
  // now make 2 function to abstract the logic here, one is for skipping and one takes in  a callback that will execute, if the video is certain sec
  // before , eg 10 sec before certain time do something
}
try {
  main().then((result) => {
    console.log("main finished and the returned value is -->", result);
  });
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
 *
 * @argument {skippedTheSponser} SkippedVideoSponsorOBJ
 * @argument {import("./helper").ResponseObject}responseObjectFromYt
 * @argument {HTMLVideoElement} videoPlayer
 */
function skipTheVideo(responseObjectFromYt, videoPlayer, SkippedVideoSponsorOBJ) {
  console.log("in the Skip the video func\n\n")
  console.log(`the cureent time is ${videoPlayer.currentTime} and end time is ${responseObjectFromYt.endTime} and start Time is ${responseObjectFromYt.startTime} and video skipper is ${SkippedVideoSponsorOBJ.videoSponsorSkipper} `)
  if (
    videoPlayer.currentTime >= responseObjectFromYt.startTime &&
    videoPlayer.currentTime < responseObjectFromYt.endTime &&
    SkippedVideoSponsorOBJ.videoSponsorSkipper === false
  ) {
    console.log(`the video player time is greater that or = the start time of from the backend -->time in the
       video player ${videoPlayer.currentTime} ----start time is  ${responseObjectFromYt.startTime} --  `);
    console.log(
      "\n now going to skip in the video to ",
      responseObjectFromYt.endTime,
    );
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
  if (videoPlayer.currentTime >= timeToCallTheFunc && skippedTheSponsorOBJ.videoSponsorSkipper === false ) {
    skippedTheSponsorOBJ.videoSponsorSkipper = true; // first as if the func throw we would not be able to update the state
    try{
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
 * */
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

    onCloseFunction?onCloseFunction():null;
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
        onUserClickDoNotSkipTheSponsorShipFunc?onUserClickDoNotSkipTheSponsorShipFunc():null; // if I did not pass it then do not play it
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
        doNotShowThisModalAgainAlwaysSkipFunction?doNotShowThisModalAgainAlwaysSkipFunction():null; // if I did not pass it then do not play it
        console.log("User chose always to skip the sponsorship");
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

    return modalContainer;
}

// /**
//  * @typedef {Object} ResponseObject
//  * @property {number} status - The status code of the response
//  * @property {string} message - A message providing additional information
//  * @property {number} startTime - The start time in milliseconds
//  * @property {number} endTime - The end time in milliseconds
//  * @property {boolean} containSponserSubtitle - Whether the video has sponsorship subtitle
//  * @property {string} [error] - Optional error message
//  */
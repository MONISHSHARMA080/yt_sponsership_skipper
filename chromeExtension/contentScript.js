// @ts-check
console.log("hi form the contentScript :)");

async function main() {
  console.log("in the main");
  let [key, error] = await chrome.runtime.sendMessage({
    type: "getKeyFromStorageOrBackend",
  });
  // let [key, error] = await getKeyFromStorageOrBackend();
  console.log("lets see what we got ");
  if (error) {
    console.log(
      "error is there in getting the key and it is -->",
      error,
      "\n\n the key is",
      key,
    );
    return;
  }

  console.log(
    "the key is  -->",
    key,
    "from the  content script, error is -->",
    error,
  );
  let [videoID, errorFormGettingVideoID] = getVideoID();
  if (errorFormGettingVideoID) {
    console.log(
      "videoID is  -->",
      videoID,
      "and the error getting it is -->",
      errorFormGettingVideoID,
    );
    return;
  }
  let videoPlayer = getVideoPlayer();
  if (videoPlayer === null) {
    console.log("video player is not there");
    return;
  }
  let [responseObject, errorFromYTApi] = await chrome.runtime.sendMessage({
    type: "getWhereToSkipInYtVideo",
    encKey: key,
    videoID: videoID,
  });
  if (errorFromYTApi || responseObject === null) {
    console.log("there is a error in the yt api -->", errorFromYTApi);
    return;
  }
  console.log("the response object is -->", responseObject);
  /** @type {ResponseObject} responseObject */
  let responseOBjectFromYt = responseObject;
  if (responseOBjectFromYt.containSponserSubtitle === false) {
    console.log("the video does not have a sponsership subtitle");
    return;
  }

  //   {
  //     "status": 200,
  //     "message": "subtitles found",
  //     "startTime": 447,
  //     "endTime": 554,
  //     "containSponserSubtitle": true
  // }
  console.log(
    "the func will be called on the time-> ",
    responseOBjectFromYt.startTime - 10,
  );
  /** @type {skippedTheSponser} */
  const SkippedVideoSponsorOBJ = { videoSponsorSkipper: false, callBackBeforeSomeTimeOfSponsor : false };
  videoPlayer.addEventListener("timeupdate", (event) => {
    console.log("\n Current time :-->", videoPlayer.currentTime); // working
    // the end time is not <= and is < because it will  not move forward if we did not do that, and just to be sure lets make a var too
    skipTheVideo(responseOBjectFromYt, videoPlayer, SkippedVideoSponsorOBJ);
    beforeSomeTimeExecuteSomething(
      responseOBjectFromYt.startTime - 10,
      videoPlayer,
      () => {console.log(" hi this func will execute before certain time (10 sec)  ",);
      },
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

function getVideoPlayer() /** @return {Element|null}*/ {
  try {
    // const video = document.querySelector("video");
    // video.addEventListener("timeupdate", (event) => {
    //   console.log("The currentTime attribute has been updated. Again.-->", event);
    // });
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
 *
 * @argument {skippedTheSponser} SkippedVideoSponsorOBJ
 * @argument {ResponseObject}responseObjectFromYt
 * @argument {Element} videoPlayer
 */
function skipTheVideo(
  responseObjectFromYt,
  videoPlayer,
  SkippedVideoSponsorOBJ,
) {
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
 * @argument {Number} timeToCallTheFunc - call the func at sponserShipStart - 10
 * @argument {Element} videoPlayer
 * @param {Function} callbackFunction
 * @param {skippedTheSponser} skippedTheSponsorOBJ
 */
// * // @argument {Number} callTheFuncBeforeThisTime - eg call the function before the
// * @param {skippedTheSponser} SkippedVideoSponsorOBJ - pass by reference boolean value in a Object
function beforeSomeTimeExecuteSomething(
  timeToCallTheFunc,
  videoPlayer,
  callbackFunction,
  skippedTheSponsorOBJ
) {
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

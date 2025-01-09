
console.log("hi form the contentScript :)");

async function main() {
  console.log("in the main");
  let [key, error] = await chrome.runtime.sendMessage({type: "getKeyFromStorageOrBackend"});
  // let [key, error] = await getKeyFromStorageOrBackend();
  console.log("lets see what we got ")
  if (error) {
    console.log(
      "error is there in getting the key and it is -->", error, "\n\n the key is", key,);
    return;
  }

  console.log("the key is  -->", key, "from the  content script, error is -->", error);
  let [videoID, errorFormGettingVideoID] =getVideoID()
  if (errorFormGettingVideoID) {
    console.log("videoID is  -->", videoID, "and the error getting it is -->", errorFormGettingVideoID);
    return
  }

  let videoPlayer =getVideoPlayer()
  if (videoPlayer === null) {
    console.log("video player is not there")
    return;
  }
  let [responseObject, errorFromYTApi]= await chrome.runtime.sendMessage({type: "getWhereToSkipInYtVideo", encKey:key, videoID:videoID});
  if (errorFromYTApi || responseObject === null ){
    console.log("there is a error in the yt api -->",errorFromYTApi)
    return
  }
  console.log("the response object is -->", responseObject)
  /** @type {ResponseObject} responseObject */
  let responseOBjectFromYt = responseObject
  if(responseOBjectFromYt.containSponserSubtitle === false){
    console.log("the video does not have a sponsership subtitle")
    return
  }

//   {
//     "status": 200,
//     "message": "subtitles found",
//     "startTime": 447,
//     "endTime": 554,
//     "containSponserSubtitle": true
// }
let videoSponserSkipped = false
  videoPlayer.addEventListener("timeupdate", (event) => {
      console.log("\n Current time :-->", videoPlayer.currentTime); // working
    // the end time is not <= and is < cause it will  not move forward if we did not do that, and jsut to be sure lets make a var too
    if(videoPlayer.currentTime >= responseOBjectFromYt.startTime && videoPlayer.currentTime < responseOBjectFromYt.endTime && videoSponserSkipped === false ){
      console.log(`the video player time is greater that or = the start time of from the backend -->time in the videoplayer ${videoPlayer.currentTime} ----start time is  ${responseOBjectFromYt.startTime} --  `);
      console.log("\n now going to skip in the video to ", responseOBjectFromYt.endTime)
      videoPlayer.currentTime = responseOBjectFromYt.endTime;
      videoSponserSkipper = true;
    }
    });

  // now control the YT player
}
try {
  main().then((result) => {
    console.log("main finished -->",result);});
}catch (e) {
  console.log("error in main script:", e);
}





















/**
 * @return {[string,Error|null]}
 */
 function getVideoID(){
  let url = window.location.href
  console.log("the url in getVideoID():", url);
  if (url === null || url === ""||url === undefined) {
    return ["", new Error("the url is not there")]
  }
  return [url, null];
}








function getVideoPlayer() /** @return {Element|null}*/ {
  try{
    // const video = document.querySelector("video");
// video.addEventListener("timeupdate", (event) => {
//   console.log("The currentTime attribute has been updated. Again.-->", event);
// });
    const videoElement = document.querySelector('video.html5-main-video');
    if (videoElement === null || videoElement === undefined) {
      return null
    }
    //  videoElement.addEventListener("timeupdate", (event) => {
    //   console.log("\n Current time :-->", videoElement.currentTime); // working
    // });
    return videoElement;
  }catch (e) {
    console.log("error occurred in getting the current time form the player->", e);
    return null;
  }
}
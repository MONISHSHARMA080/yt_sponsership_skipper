
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
  let [responseObject, errorFromYTApi]= await chrome.runtime.sendMessage({type: "getWhereToSkipInYtVideo", encKey:key, videoID:videoID});
  if (errorFromYTApi || responseObject === null ){
    console.log("there is a error in the yt api -->",errorFromYTApi)
    return
  }
  console.log("the response object is -->", responseObject)
  if(responseObject){}
}
try {
  main();
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

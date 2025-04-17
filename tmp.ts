async function getTheVideoPlaybackTimeForItsLife() {
  try {
    const videoPlayer = document.querySelector('video');
    if (!videoPlayer) {
      return { result: null, error: "The video player is not found" };
    }

    let resultArray = [];
    let i = 0;

    videoPlayer.addEventListener("timeupdate", () => {
      resultArray.push(videoPlayer.currentTime);
      console.log("Adding currentTime to array, index:" + i);
      i++
    });

    return await new Promise((resolve) => {
      videoPlayer.addEventListener("ended", () => {
        console.log("Video has ended, returning resultArray");
        resolve({ result: resultArray, error: "" });
      });
    });
  } catch (error) {
    console.log("Error:", error);
    return { result: null, error: "There was an error: " + error };
  }
}


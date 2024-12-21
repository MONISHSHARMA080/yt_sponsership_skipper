console.log("hi form the background.js");
chrome.action.onClicked.addListener(function () {
  chrome.tabs.create({ url: "index1.html" });
});

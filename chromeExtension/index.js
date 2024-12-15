console.log("hi form the index.js");

try {
  let element = document.getElementById("status");
  if (element !== null) {
    element.addEventListener("click", function a() {
      console.log("hi from the index.js button-- ");
      getUserDetailsForSignUp();
    });
  }
} catch (e) {
  console.log(e);
}

function getUserDetailsForSignUp() {
  chrome.identity.getProfileUserInfo({}, (userInfo) => {
    console.log("email I got is -->", userInfo.email);
    console.log("id I got is -->", userInfo.id);
    console.log("whole obj-->", userInfo);
  });
  chrome.identity.getProfileUserInfo({}, (ProfileUserInfo) => {
    console.log("profile user info is ", ProfileUserInfo);
  });
  chrom.identity.getAuthToken({}, (GetAuthTokenResult) => {});
}

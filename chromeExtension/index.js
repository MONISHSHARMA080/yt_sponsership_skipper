// @ts-check
/// <reference types="chrome" />

console.log("hi form the index.js");

try {
  let element = document.getElementById("status");
  if (element !== null) {
    element.addEventListener("click", async function a() {
      console.log("hi from the index.js button-- ");
      await getUserDetailsForSignUp();
    });
  }
} catch (e) {
  console.log(e);
}

async function getUserDetailsForSignUp() {
  /**
   * @typedef {Object} UserDetail
   * @property {string} account_id - The account ID of the user that is unique for every account.
   * @property {string} user_token - The authentication token for the user.
   */
  /** @type {UserDetail}  */
  const UserDetail = {
    account_id: "",
    user_token: "",
  };
  try {
    const userInfo = await new Promise((resolve) => {
      chrome.identity.getProfileUserInfo({}, resolve);
    });

    UserDetail.account_id = userInfo.id;
    console.log("email:", userInfo.email);
    console.log("id:", userInfo.id);

    // Convert getAuthToken to Promise
    const authToken = await new Promise((resolve, reject) => {
      chrome.identity.getAuthToken({ interactive: true }, (token) => {
        if (chrome.runtime.lastError) {
          reject(chrome.runtime.lastError);
        } else {
          resolve(token);
        }
      });
    });

    UserDetail.user_token = authToken;
    console.log("auth token:", authToken);

    const response = await fetch("http://localhost:8080/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(UserDetail),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    console.log("Success:", data);
  } catch (error) {
    console.error("Error:", error);
  }
}

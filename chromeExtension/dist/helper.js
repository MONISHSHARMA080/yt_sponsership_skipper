const l = {
  BACKEND_URL: "http://localhost:8080"
};
async function g() {
  try {
    let [e, r] = await c("key");
    if (r || e == null || e == "")
      return console.log("there is an error in getting the key form the storage, and that is -->", r), [null, r];
    let [o, t] = await u();
    return t || o === null || o === "" ? (console.log("there is an error in getting the key (or it is null or '') ->", t), [null, t]) : (console.log("the key form the backend is -->", o), [o, null]);
  } catch (e) {
    console.log("error in the getKeyFromStorageOrBackend() and that is ->", e);
    const r = e instanceof Error ? e.message : typeof e == "string" ? e : "An unknown error occurred";
    return [null, new Error(r)];
  }
}
async function c(e) {
  try {
    const o = (await new Promise((t, s) => {
      chrome.storage.local.get([e], (n) => {
        chrome.runtime.lastError ? s(chrome.runtime.lastError) : t(n);
      });
    }))[e];
    return console.log("Value from storage:", o), [o || null, null];
  } catch (r) {
    return console.log("Error getting key from storage:", r), [null, r instanceof Error ? r : new Error(String(r))];
  }
}
async function u() {
  try {
    const e = await new Promise((n) => {
      chrome.identity.getProfileUserInfo(n);
    }), r = await new Promise((n, i) => {
      chrome.identity.getAuthToken({ interactive: !0 }, (a) => {
        a ? n(a) : i(new Error("Token is undefined"));
      });
    }), o = {
      account_id: e.id,
      user_token: r
    }, t = await fetch(`${l.BACKEND_URL}/signup`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(o)
    });
    if (!t.ok)
      throw new Error(`HTTP error! Status: ${t.status}`);
    const s = await t.json();
    return console.log("Success:", s), [s.encrypted_key, null];
  } catch (e) {
    return console.log("Error getting key from backend:", e), [null, e instanceof Error ? e : new Error(String(e))];
  }
}
export {
  l as config,
  g as getKeyFromStorageOrBackend
};

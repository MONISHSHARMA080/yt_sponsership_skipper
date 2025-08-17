
## Launching soon

[Try it (it's live)](https://chromewebstore.google.com/detail/youtube-sponsorship-skipp/dpkehfkmkhmbmbofjaikccklcdpdmhpl?authuser=2&hl=en)

# -- captions solution --
if the captions are not in the youtube script then it is definately getting fetched via the api response , see the network tab as it is being fetched via the timed text api
--> implement this and it will be working again, intercept it and then copy the response

## -- update svelte communication --
Both your extension and your Svelte site talk to the same API backend. So for eg if the user goes to the website then make a fetch req to the backend and if the user's token is on let's say version 2 and the server version is 15 then just give them the latest version and save the same in the storage, or you can just get the version from the server in both chrome extension and site, and if the version is same in both then we can keep it  

Extension logs in → updates state on server → website fetches from server.

Keeps state authoritative and avoids race conditions.

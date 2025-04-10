import test from "@playwright/test";
import { initializeBrowesrForChromeExtension } from "../helperFiles/initializeBrowser";


test("set key in chrome extension", async () => {
  let contextOfBrowser = await initializeBrowesrForChromeExtension()

  // approach 1

  let extensionId: string = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"
  const page = await contextOfBrowser.newPage();

  try {

    const client = await page.context().newCDPSession(page);
    const storage = await client.send('Storage.getStorageKeyForFrame',);
    console.log(`the storage key is ${storage.storageKey} and the storage object is \n\n\n ${JSON.stringify(storage)}  \n\n\n`);
  } catch (error) {
    console.log(`first approach's  error is ->`, error);
  }


  // approach 2
  await page.goto(`chrome-extension://${extensionId}/index.html`);

  // Execute JavaScript in the context of the extension page to access storage
  const result = await page.evaluate(() => {
    return new Promise((resolve) => {
      //@ts-ignore
      chrome.storage.local.get(null, (data) => {
        resolve(data);
      });
    });
  });
  console.log('Extension storage data:', result);

})

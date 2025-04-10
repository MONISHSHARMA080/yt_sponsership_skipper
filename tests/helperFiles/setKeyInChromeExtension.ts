import { BrowserContext, Page, test } from "@playwright/test";


export async function setKeyInChromeExtensionsLocalStorage(key: string = "key", value: string, page: Page, extensionID: string = "dpkehfkmkhmbmbofjaikccklcdpdmhpl", browserContext: BrowserContext) {
  try {
    let serviceWorkers = browserContext.serviceWorkers()
    serviceWorkers.forEach((value, index) => {
      console.log(`the service worker at ${index} is ${JSON.stringify(value)} -- and non json stringify is ${value}`);

    })

    let bgPage = browserContext.backgroundPages().find((page, index) => {
      page.url().startsWith('chrome-extension://')
    })
    if (!bgPage) {
      return new Error(`didn't find the bgPage in the browser context`)
    }
    let storageState = await bgPage.context().storageState()
    console.log(`the storage state before the setting key for chrome extension is ${JSON.stringify(storageState)} and we are have key:-> ${key} --- and value:-> ${value} `);
    await bgPage.evaluate(() => {
      // localStorage.setItem(key, value)
      //@ts-ignore
      chrome.storage.local.set({ key, value })
    })
    console.log(`the storage state after setting key for chrome extension is ${await bgPage.context().storageState()} and we are have key:-> ${key} --- and value:-> ${value} `);
    await page.close()
  } catch (error) {
    console.log(` error in storing the key in chrome-extension is ->`, error);
    await page.close()
  }

}
class InteractWithStorageHelper {

  constructor(parameters) {

  }
}

import { Page, test } from "@playwright/test";


export async function setKeyInChromeExtensionsLocalStorage(key: string = "key", value: string, page: Page) {
  let storageState = await page.context().storageState()
  console.log(`the storage state before the setting key for chrome extension is ${JSON.stringify(storageState)} and we are have key:-> ${key} --- and value:-> ${value} `);
  await page.evaluate(() => {
    localStorage.setItem(key, value)
  })

}

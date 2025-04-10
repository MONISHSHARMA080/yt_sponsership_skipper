import { BrowserContext, chromium } from "@playwright/test";
import path from "path";

/**initialize browser context for chrome Extension*/
export async function initializeBrowesrForChromeExtension(extensionPath: string = '../../chromeExtension/', extensionId: string = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"): Promise<BrowserContext> {
  const context = await chromium.launchPersistentContext('', {
    headless: false,
    args: [
      `--disable-extension-except=${path.join(__dirname, '../../chromeExtension')}`,
      `--load-extension=${extensionId}`
    ]
  })
  return context
}

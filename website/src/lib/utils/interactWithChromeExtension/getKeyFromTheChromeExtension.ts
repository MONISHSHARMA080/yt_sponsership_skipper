import { browser } from "$app/environment";
import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { accessLocalStorage } from "../localStorage";
import { checkIfKeyIsValidAndUpdateTheState } from "../seeIfTheKeyIsValidByBackend";

type funcToRunWhenWeGetTheKey = (key: string) => void

/** the clas is there so that you can call start() and get the new key form the chrome extensio, note that the function in 
 * start() will not run when key is received by the local storage
 * 
 * the class will also on starting extract the user obj form the localstorage  in the start func 
 *   
 * this class is meant to be used in a way where you execute it and see the change (or $derived/effect) in the  keyFromChromeExtensionState
 */
export class interactWithTheChromeExtensionAndStoreIt {
  private callBackAfterKeyIsReceived: null | funcToRunWhenWeGetTheKey = null
  private checkIfKeyIsValid = new checkIfKeyIsValidAndUpdateTheState()
  private localStorageHelper = new accessLocalStorage()
  private keyStateFromStorage: keyStateObject | null = null

  constructor() {
    this.messageHandler = this.messageHandler.bind(this);
  }

  private messageHandler(event: MessageEvent) {
    console.log(`in the message handler and the event is ->`, event);

    if (event.origin !== window.location.origin) {
      return
    }
    if (event.data.type === "GET_KEY") {
      console.log(event.data.type, "---the key is  ->", event.data.key);
      if (event.data.key === undefined || event.data.key === null || event.data.key === "") {
        console.log("the key is not there (or undefined or null) and the event.data is ->", event.data);
        // possibly throw an error form here
        return
      }
      // doing both as I want to keep the fufute extensibility too and why not it is the same thing 
      // keyFromChromeExtensionState.key = event.data.key 

      // --- new code ----
      console.log(`the key form the chrome extension is ->`, event.data.key);

      if (this.keyStateFromStorage !== null) {
        console.log(`the key form the local storage is 234->`, this.keyStateFromStorage);


        if (this.keyStateFromStorage.key === event.data.key) {
          // the key is same and we don't need to do anything
          console.log("the key is same and we don't need to do anything, the key is ->", event.data.key, "and the key state from storage is ->", this.keyStateFromStorage.key);
          // update the global state 
          this.saveKeyObjFromLocalStorageToGlobalState(this.keyStateFromStorage)
          this.removeAndCloseEventListeners()
          return
        }
      }
      // the key is null as it might be the first  time or the key is not same 
      // the key is diff. and we need to do the whole thing again
      // this method will also update the key in local storage so don't worry about it
      console.log("the key is diff and we are checking if it is valid, ");
      console.log(`+++++++++++++++++++++++++++++++++++++++++++++++++++++2222222222222222222222222`);

      this.checkIfKeyIsValid.seeIfKeyIsValid(event.data.key)


      if (this.callBackAfterKeyIsReceived !== null) {
        try {
          this.callBackAfterKeyIsReceived(event.data.key);
        } catch (error) {
          console.log("the error is in the callback after the key fucntion ->", error);
        }
      } else {
        console.log(" the key uodate function is not there");
      }
      if (!event.data.key) return;
      // Send message to remove listeners and remove our own listener
      console.log("closing all the event listeners as the key is received");
      // window.postMessage({
      // 	type: "removeAllEventListener"
      // }, window.location.origin);
      // // Remove this event listener since we got the key
      // if (typeof window !== 'undefined') {
      // 	window.removeEventListener('message', this.messageHandler);
      // }
      this.removeAndCloseEventListeners()
    }
  }

  public cleanup() {
    console.log("cleaning the internact with chrome extension class");
    window.removeEventListener('message', this.messageHandler);
  }

  /** @param {(key:string)=>void} param - the function as a param will not run when the key is received by the local storage 
   * 
   * the function will also try to get the user obj form the local storage and put it in the global
   */
  public start(funcToRunWhenWeGetTheKey: funcToRunWhenWeGetTheKey): Error | null {
    try {
      if (!browser) {
        return new Error("the app is not running in the browser(ssr) and we can't start listenig to the chrome extension ")
      }
      let [keyObj, error] = this.getKeyObjFromLocalStorage()
      if (error !== null || keyObj === null) {
        console.log("the key object is not there in the local storage or the error is ->", error);
        // if we don't have the key then we should clear the local storage
        // this.localStorageHelper.removeFromLocalStorage("KEY")
      }
      console.log("the key object is  in the start()->", keyObj);
      console.log(`key form the storage is ->`, this.keyStateFromStorage);
      // if the key is present 
      if (keyObj?.key !== null && keyObj?.key !== undefined && keyObj?.key !== "" && keyObj?.key.length > 4){
        if (keyObj.isValidatedThroughBackend){
          // then go ahead and set the global state to it 
          console.log(`the key object is already in the storage and is validated through backend so we are assigning `);
           Object.assign(keyFromChromeExtensionState, keyObj)
        }
      }

      // this.keyStateFromStorage = keyObj
      // this.saveKeyObjFromLocalStorageToGlobalState(keyObj)
      // if the keyObj is not null then set it as the global export state
      // now when  we get the keys form the chrome extension check the key is same or not, if not then 
      // this.saveKeyObjFromLocalStorageToState(keyObj)

      // wait and get the key form the chrome extension if it is same then quit and if diff. then do the whole thing and deleate the key

      this.callBackAfterKeyIsReceived = funcToRunWhenWeGetTheKey;
      window.addEventListener('message', this.messageHandler);
      window.postMessage({ type: 'GET_KEY' }, window.location.origin);
      // maybe do a settimeout where we clean after 3 min of sleep 

      // commentinng this as this works in the bg 
      // setTimeout(()=>{this.saveKeyObjFromLocalStorageToGlobalState(keyObj);this.cleanup(); console.log("time passes and we are closing the evente listender with set timeout") ; return null }, 8000)
      return null
    } catch (error) {
      if (error instanceof Error) {
        return error
      } else {
        return new Error("error occurred -> " + error)
      }
    }
  }

  private getKeyObjFromLocalStorage(): [keyStateObject | null, Error | null] {
    // if the key returned is not there or null then it is the first time 
    try {
      let [keyObjInStr, error] = this.localStorageHelper.getFromLocalStorage("KEY")
      if (error !== null || keyObjInStr === "" || keyObjInStr === null) {
        return [null, error instanceof Error ? error : new Error("either the key not there or it is the forst time , error I got is ->" + error)]
      }
      let KeyState: keyStateObject
      try {
        KeyState = JSON.parse(keyObjInStr)
      } catch (error) {
        // clear the local storage as will change it further and the 
        this.localStorageHelper.removeFromLocalStorage("KEY")
        return [null, new Error("error in parsing the key object from the local storage ->" + error)]
      }
      return [KeyState, null]
    } catch (error) {
      return [null, error instanceof Error ? error : new Error("error in getting the key object from the local storage ->" + error)]
    }
  }

  private saveKeyObjFromLocalStorageToGlobalState(keyObj: keyStateObject | null) {
    if (keyObj !== null && keyObj.key !== "" && keyObj.key !== null) {
      console.log(`changine the keyFromChromeExtensionState and the key Obj to replace it is ${JSON.stringify(keyObj)} ---- and the new one is ${JSON.stringify(keyFromChromeExtensionState)}`);
      Object.assign(keyFromChromeExtensionState, keyObj)
   }
  }


  private removeAndCloseEventListeners() {
    // also send the message to the chrome extension 
    window.postMessage({
      type: "removeAllEventListener"
    }, window.location.origin);
    // Remove this event listener since we got the key
    if (typeof window !== 'undefined') {
      window.removeEventListener('message', this.messageHandler);
    }
  }
}

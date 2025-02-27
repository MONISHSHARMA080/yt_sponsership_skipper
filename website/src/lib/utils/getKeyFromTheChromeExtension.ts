import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { accessLocalStorage } from "./localStorage";
import { checkIfKeyIsValidAndUpdateTheState } from "./seeIfTheKeyIsValidByBackend";

type funcToRunWhenWeGetTheKey = (key:string)=>void 

export class interactWithTheChromeExtensionAndStoreIt{
   private callBackAfterKeyIsReceived :null|funcToRunWhenWeGetTheKey = null
   private checkIfKeyIsValid = new checkIfKeyIsValidAndUpdateTheState()
   private localStorageHelper = new accessLocalStorage()
   private keyStateFromStorage :keyStateObject|null = null

   constructor() {
      this.messageHandler = this.messageHandler.bind(this);
   }

   private  messageHandler(event: MessageEvent) {
		if (event.origin !== window.location.origin){
         return
      }
		if (event.data.type === "GET_KEY") {
         console.log(event.data.type,"---the key is  ->", event.data.key);
         if (event.data.key === undefined || event.data.key === null || event.data.key === "") {
            console.log("the key is not there (or undefined or null) and the event.data is ->",event.data);
            // possibly throw an error form here
            return
         }
         // doing both as I want to keep the fufute extensibility too and why not it is the same thing 
         // keyFromChromeExtensionState.key = event.data.key 


         // --- new code ----
         if (this.keyStateFromStorage !== null) {
           if (this.keyStateFromStorage.key === event.data.key) {
               // the key is same and we don't need to do anything
               console.log("the key is same and we don't need to do anything, the key is ->", event.data.key, "and the key state from storage is ->", this.keyStateFromStorage.key);
               return
           }
         }else{
               // the key is null as it might be the first  time or the key is not same 
               // the key is diff. and we need to do the whole thing again
               // this method will also update the key in local storage so don't worry about it
               console.log("the key is diff and we are checking if it is valid, ");
               
               this.checkIfKeyIsValid.seeIfKeyIsValid(event.data.key)
           } 






         if (this.callBackAfterKeyIsReceived !== null) {
            try {
               this.callBackAfterKeyIsReceived(event.data.key);
            } catch (error) {
               console.log("the error is in the callback after the key fucntion ->", error);   
            }
         }else{
            console.log(" the key uodate function is not there");
         }
			if (!event.data.key) return;
			// Send message to remove listeners and remove our own listener
			console.log("closing all the event listeners as the key is received");
			window.postMessage({
				type: "removeAllEventListener"
			}, window.location.origin);
			// Remove this event listener since we got the key
			if (typeof window !== 'undefined') {
				window.removeEventListener('message', this.messageHandler);
			}
		}
	}

   public  cleanup(){
      console.log("cleaning the internact with chrome extension class");
      window.removeEventListener('message', this.messageHandler);
   }

   public start(funcToRunWhenWeGetTheKey:funcToRunWhenWeGetTheKey) : Error|null {
      try {
         let [keyObj, error] = this.getKeyObjFromLocalStorage()
         if( error !== null || keyObj === null){
            console.log("the key object is not there in the local storage or the error is ->", error);
            // if we don't have the key then we should clear the local storage
            this.localStorageHelper.removeFromLocalStorage("KEY")
         }
         console.log("the key object is ->", keyObj);
         this.keyStateFromStorage = keyObj
         // if the keyObj is not null then set it as the global export state
         // now when  we get the keys form the chrome extension check the key is same or not, if not then 
         // this.saveKeyObjFromLocalStorageToState(keyObj)


         // wait and get the key form the chrome extension if it is same then quit and if diff. then do the whole thing and deleate the key








         this.callBackAfterKeyIsReceived = funcToRunWhenWeGetTheKey;
         window.addEventListener('message', this.messageHandler);
         window.postMessage({ type: 'GET_KEY' }, window.location.origin);
         // maybe do a settimeout where we clean after 3 min of sleep 
         return null
      } catch (error) {
         if( error instanceof Error){
            return error
         }else{
            return new Error("error occurred -> "+ error)
         }
      }
   }

   private getKeyObjFromLocalStorage():[keyStateObject|null,Error|null] {
      // if the key returned is not there or null then it is the first time 
     try {
      let [keyObjInStr, error] = this.localStorageHelper.getFromLocalStorage("KEY")
      if( error !== null || keyObjInStr === "" || keyObjInStr === null){
         return [null,error instanceof Error ? error: new Error("either the key not there or it is the forst time , error I got is ->"+error)]
      }
      let KeyState :keyStateObject
      try {
         KeyState =  JSON.parse(keyObjInStr)
      } catch (error) {
         // clear the local storage as will change it further and the 
         this.localStorageHelper.removeFromLocalStorage("KEY")
         return [null, new Error("error in parsing the key object from the local storage ->"+error)]
      }
      return [KeyState, null]
     } catch (error) {
         return [null, error instanceof Error ? error : new Error("error in getting the key object from the local storage ->"+ error)]
     }
   }

   private saveKeyObjFromLocalStorageToState(keyObj:keyStateObject|null){
      if( keyObj !== null){
         keyFromChromeExtensionState.key = keyObj.key
         keyFromChromeExtensionState.isValidatedThroughBackend = keyObj.isValidatedThroughBackend
         keyFromChromeExtensionState.name = keyObj.name
         keyFromChromeExtensionState.email = keyObj.email
      }
   }
}
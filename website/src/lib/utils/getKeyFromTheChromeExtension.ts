import { keyFromChromeExtensionState } from "$lib/sharedState/sharedKeyState.svelte";

type funcToRunWhenWeGetTheKey = (key:string)=>void 

export class interactWithTheChromeExtensionAndStoreItInTheState{
   private callBackAfterKeyIsReceived :null|funcToRunWhenWeGetTheKey = null

   constructor() {
      this.messageHandler = this.messageHandler.bind(this);
   }

   private  messageHandler(event: MessageEvent) {
		if (event.origin !== window.location.origin){
         return
      }
		if (event.data.type === "GET_KEY") {
         console.log(event.data,"---the key is  ->", event.data.key);
         if (event.data.key === undefined || event.data.key === null || event.data.key === "") {
            console.log("the key is not there (or undefined or null) and the event.data is ->",event.data);
            // possibly throw an error form here
            return
         }
         // doing both as I want to keep the fufute extensibility too and why not it is the same thing 
         keyFromChromeExtensionState.key = event.data.key 

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
}
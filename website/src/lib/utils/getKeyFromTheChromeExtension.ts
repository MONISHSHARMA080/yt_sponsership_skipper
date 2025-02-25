export class interactWithTheChromeExtension{
   private  messageHandler(event: MessageEvent) {
		if (event.origin !== window.location.origin) return;
		
		if (event.data.type === "GET_KEY") {
			console.log("the key is  ->", event.data.key);
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

    cleanup(){
      window.removeEventListener('message', this.messageHandler);
   }

    start():Error|null{
      try {
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
export class sendChromeExtensionNewKey{

    private newKey: string;
    private resolveFunction: ((value: [Error | null, boolean]) => void) | null = null;

    constructor(newKey: string){
        this.newKey = newKey;
        this.handleEvent = this.handleEvent.bind(this);
        console.log("clearing the console just for the testing purpose");
        setTimeout(()=>{}, 5000) // just for testing purpose
        console.clear() // just for testing purpose
        console.log("cleared the console just for the testing purpose");
    }
    private handleEvent = (event:MessageEvent)=>{
        if (this.resolveFunction === null){ 
            throw Error("resolve function is not set");
            
        }
        try{
            if (event.origin !== window.location.origin){
                this.clearEventListener()
                console.log("the origing is not same");
                this.resolveFunction([Error("origin is not same as the window location origin"), false]);
            }
            if (event.data.type === "keyChangedOnPaymentReceived") {
                console.log("the key is received by the chrome extension ->", event.data.key, " and the success is ->", event.data.success);
                this.clearEventListener()   
                const newKeyReceived = event.data.key;
                const succesReceived =  event.data.success;

                if (!this.checkIfDataReceivedIsCorrect(newKeyReceived, succesReceived)){
                    this.resolveFunction([Error("there is a probelm with key received is not the same or the success is false"), false]);
                }
                this.resolveFunction([null, succesReceived]);
            }
            console.log("idk why");
            
        }catch(error){
            this.clearEventListener()
            return error instanceof Error ? error : Error("error in handling the event ->"+ error);

        }
    }

    /** fuc does not has the reject func run so get the failure or succes via the value of the .then or resolve func only */
    public async sendKey(): Promise<[Error | null, boolean]> {
       return new Promise<[Error|null, boolean]>((resolve, reject) => {
        try {
            this.resolveFunction = resolve;
            window.addEventListener("message", this.handleEvent )
            window.postMessage({type: "paymentReceivedChangeTheKey", key: this.newKey}, window.location.origin);
            return [null, false];
        } catch (error) {
           resolve ([error instanceof Error ? error : Error("error in sending the key ->"+ error), false]); 
        }
       })
    }
    
    public clearEventListener(): void {
        window.removeEventListener("message", this.handleEvent);
        console.log("Event listener cleared");
    }

    public checkIfDataReceivedIsCorrect(newKeyReceived: string, succesReceived: boolean): boolean { 
        if (newKeyReceived === this.newKey && succesReceived === true){
            return true;
        }
        console.log("data received is not correct as are the keys same ->", this.newKey === newKeyReceived, " and the success is ->", succesReceived)
        
        return false;
    }

}
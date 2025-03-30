import { keyFromChromeExtensionState, type keyStateObject } from "$lib/sharedState/sharedKeyState.svelte";
import { z } from "zod";
import { accessLocalStorage } from "./localStorage";

/**@ */
interface resultType <T>{
    error:Error|null,
    result:T|null
}
const keyStateSchema = z.object({
    key: z.string().nullable(),
    isValidatedThroughBackend: z.boolean(),
    name: z.string().nullable(),
    email: z.string().nullable(),
    isPaidUser: z.boolean()
});

export class KeyUpdate{
    /**
     * the goal of this class is to update the key in the local storage and also in the shared state, we are not running the is valid through backend function on it 
     * as we are assuming that the key is valid , and you would update it only if it comes form the backend
     * 
     * if the key on the object form the storage is "" we will just put it there 
     * 
     * also update the shared state outside of it 
     */
    public UpdateKey(newKey:string, isUserOnPaidTier:boolean = false):Error|null{
        if (!this.isKeyValid(newKey)){
            return new Error("The key is not valid")
        }
        // now setting the new local storage key
        let storage = new accessLocalStorage
        let [keyFormStorage, error] = storage.getFromLocalStorage("KEY")
        if (error !== null || keyFormStorage ==="" || keyFormStorage=== null ){
            console.log("the key is not what we expect it to be so returning");
            return error instanceof Error ? error : new Error("there is a error in getting the key form the storage, either the key is not there or empty")
        }

        let result =this.parseTheKeyFromStorageIntoJSONObj(keyFormStorage)

        if (result.error !== null || result.result === null){
            console.log(`can't parse the key form the storage into the type JSON schema ->${result.error},==++ result is ${result.result} `);
            return result.error instanceof Error ? error : new Error("either we can't parse the key into the json schema or the keyStateObj is empty/null ")
        }
        result.result.key= newKey
        result.result.isPaidUser= isUserOnPaidTier

        let [success, error2] = storage.setInLocalStorage("KEY", JSON.stringify(result.result))
        if (error2 !== null ){
            return error2
        }
        console.log("successfully changed the key in the storage, the new key will be--=--",newKey);
        // updating the state
        // keyFromChromeExtensionState.key = result.result.key
        // keyFromChromeExtensionState.isPaidUser = result.result.isPaidUser
        Object.assign(keyFromChromeExtensionState, result.result);
        return null

    }
    private isKeyValid(newKey:string):boolean{
        if (newKey === null || newKey === undefined || newKey === "" ){
            return false;
        }else{
            return true
        }
    }
    private parseTheKeyFromStorageIntoJSONObj(key: string): resultType<keyStateObject> {
        try {
            // First try to parse the string to JSON
            const parsedData = JSON.parse(key);
            
            // Then validate the parsed data against our Zod schema
            const validationResult = keyStateSchema.safeParse(parsedData);
            
            if (!validationResult.success) {
                // If validation fails, return an error with details
                return {
                    result: null, 
                    error: new Error(`Invalid key structure: ${validationResult.error.message}`)
                };
            }
            
            // Return the validated data
            return {
                result: validationResult.data,
                error: null
            };
        } catch (error) {
            return {
                result: null, 
                error: error instanceof Error 
                    ? error 
                    : new Error("Error parsing key from storage into JSON: " + error)
            };
        }
    }

}
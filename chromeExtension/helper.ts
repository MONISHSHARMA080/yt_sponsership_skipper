import { AuthResponse, UserDetail } from "./types";
// import { config } from "./config";

export const config = {
  BACKEND_URL: "http://localhost:8080",
};



export async function getKeyFromStorageOrBackend():Promise<[String|null, Error|null ]>{
   try {
    let [valueOfTheKey, errorinGettingKeyFromTYheStorage] = await getKeyFromTheStorage("key")
    if (errorinGettingKeyFromTYheStorage || valueOfTheKey == null || valueOfTheKey =="" ){
        console.log("there is an error in getting the key form the storage, and that is -->",errorinGettingKeyFromTYheStorage);
        return [null, errorinGettingKeyFromTYheStorage]
    }
    let [keyFromTheBackend, error] = await getKeyFromTheBackend()
    if (error || keyFromTheBackend ===null || keyFromTheBackend ==="" ){
        console.log("there is an error in getting the key (or it is null or '') ->",error);
       return [null, error] 
    }
    console.log("the key form the backend is -->",keyFromTheBackend);
    return [keyFromTheBackend, null]
   } catch (error) {
   console.log("error in the getKeyFromStorageOrBackend() and that is ->", error);
   const errorMessage = error instanceof Error 
       ? error.message 
       : typeof error === 'string' 
           ? error 
           : "An unknown error occurred";

   return [null, new Error(errorMessage)];
   } 
}


async function getKeyFromTheStorage(key: string): Promise<[string | null, Error | null]> {
    try {
        // Fixed Promise type declaration
        const result = await new Promise<{ [key: string]: string }>((resolve, reject) => {
            chrome.storage.local.get([key], (items) => {
                if (chrome.runtime.lastError) {
                    reject(chrome.runtime.lastError);
                } else {
                    resolve(items);
                }
            });
        });
        
        const value = result[key];
        console.log("Value from storage:", value);
        
        return [value || null, null];
    } catch (error) {
        console.log("Error getting key from storage:", error);
        return [null, error instanceof Error ? error : new Error(String(error))];
    }
}




async function getKeyFromTheBackend(): Promise<[string | null, Error | null]> {
    try {
        // Convert chrome.identity.getProfileUserInfo to Promise
        const userInfo = await new Promise<chrome.identity.UserInfo>((resolve) => {
            chrome.identity.getProfileUserInfo(resolve);
        });
        
        // Convert chrome.identity.getAuthToken to Promise
        const token = await new Promise<string>((resolve, reject) => {
            chrome.identity.getAuthToken({ interactive: true }, (token) => {
                if (!token) {
                    reject(new Error("Token is undefined"));
                } else {
                    resolve(token);
                }
            });
        });
        
        const userDetail: UserDetail = {
            account_id: userInfo.id,
            user_token: token
        };
        
        const response = await fetch(`${config.BACKEND_URL}/signup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(userDetail),
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const data: AuthResponse = await response.json();
        console.log("Success:", data);
        
        return [data.encrypted_key, null];
    } catch (error) {
        console.log("Error getting key from backend:", error);
        return [null, error instanceof Error ? error : new Error(String(error))];
    }
}

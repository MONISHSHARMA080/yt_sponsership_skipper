
export class accessLocalStorage{

    private isThisServerEnv():Boolean{
        return typeof window === 'undefined' || typeof window.localStorage === 'undefined';
    }

    getFromLocalStorage(key:string) :[string|null, Error|null]{
       try {
            if (this.isThisServerEnv()) {
                return ["",new Error("window is undefined")]
            } 
            let value = window.localStorage.getItem(key)
            return [value, null]
       } catch (error) {
            return [null, error instanceof Error ? error : new Error("error in getting the value from local storage ->"+ error)]
       } 
    }
    setInLocalStorage(key:string, value:string):[boolean, Error|null]{
        try {
            if (this.isThisServerEnv()) {
                return [false,new Error("window is undefined")]
            } 
            window.localStorage.setItem(key, value)
            return [true, null]
        } catch (error) {
            return [false, error instanceof Error ? error : new Error("error in setting the value in local storage ->"+ error)]
        }
    }
    removeFromLocalStorage(key:string):[boolean, Error|null]{
        try {
            if (this.isThisServerEnv()) {
                return [false,new Error("window is undefined")]
            }
            window.localStorage.removeItem(key)
            return [true, null]
            
        } catch (error) {
            return [false, error instanceof Error ? error : new Error("error in removing the value from local storage ->"+ error)]
        }
    }
}
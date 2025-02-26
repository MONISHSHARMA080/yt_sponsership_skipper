type keyStateObject = {key:string|null, isValidatedThroughBackend:boolean}

export const keyFromChromeExtensionState = $state<keyStateObject>({key:null, isValidatedThroughBackend:false})

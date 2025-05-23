export type keyStateObject = { key: string | null, isValidatedThroughBackend: boolean, name: string | null, email: string | null, isPaidUser: boolean }

export const keyFromChromeExtensionState = $state<keyStateObject>(
  { key: null, isValidatedThroughBackend: false, name: null, email: null, isPaidUser: false }
)

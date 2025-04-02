
// this class will take a shared stateObject and will watch it usign the $derived rune and if the value change we will write it to 
// the local storage that way the state will remain in sync

import { accessLocalStorage } from "../localStorage";
import { z } from "zod";

export class WriteSharedStateToStorageWhenItChanges<T extends object> {
  private sharedStateToWriteToTheStorage: T;
  private keyToWriteInTheStorage: string;
  private interactWithStorageClass = new accessLocalStorage;
  private schema: z.ZodType<T>;

  constructor(sharedStateToWriteToTheStorageOnChange: T, objectKeyInStorage: string) {
    this.sharedStateToWriteToTheStorage = sharedStateToWriteToTheStorageOnChange;
    this.keyToWriteInTheStorage = objectKeyInStorage;
    // Generate a Zod schema from the initial state object
    this.schema = this.generateSchemaFromObject(sharedStateToWriteToTheStorageOnChange);
  }

  /**
   * Automatically generates a Zod schema based on the structure of the provided object
   */
  private generateSchemaFromObject(obj: unknown): z.ZodType {
    if (obj === null) {
      return z.null();
    }

    if (typeof obj !== 'object') {
      return z.any();
    }

    if (Array.isArray(obj)) {
      return z.array(z.any());
    }

    const schemaShape: Record<string, z.ZodTypeAny> = {};

    for (const [key, value] of Object.entries(obj)) {
      if (value === null) {
        // Handle null values (could be null or the original type)
        schemaShape[key] = z.any().nullable();
      } else if (typeof value === 'string') {
        schemaShape[key] = z.string();
      } else if (typeof value === 'number') {
        schemaShape[key] = z.number();
      } else if (typeof value === 'boolean') {
        schemaShape[key] = z.boolean();
      } else if (Array.isArray(value)) {
        schemaShape[key] = z.array(z.any()); // Simple array validation
      } else if (typeof value === 'object') {
        // Recursively handle nested objects
        schemaShape[key] = this.generateSchemaFromObject(value);
      } else {
        schemaShape[key] = z.any();
      }
    }
    return z.object(schemaShape);
  }

  /**
   * This method watches the value and will write it to storage when it changes
   * @argument functoTellWhenShouldYouUpdateIt - when this func determines(/returns true) we write state to the storage
   */
  public wathcAndSaveOnChange(functoTellWhenShouldYouUpdateIt: (sharedStateToGetAsACallback: T) => boolean) {
    console.log(`in the wathcAndSaveOnChange`);

    // Try to get from storage but don't fail if it's not there
    let keyFromStorage = this.getKeyObjFromTheStorage();
    if (!(keyFromStorage instanceof Error) && this.validation(keyFromStorage)) {
      // If we have valid storage data, you could potentially merge it with state here
      console.log("Found valid data in storage");
    }

    // Set up the effect regardless of storage state
    let changedValue = $derived(this.sharedStateToWriteToTheStorage);
    $effect(() => {
      if (functoTellWhenShouldYouUpdateIt(this.sharedStateToWriteToTheStorage)) {
        let newValueOfTheChangedObjInStr = JSON.stringify(changedValue);
        console.log(`Saving to storage: ${this.keyToWriteInTheStorage} = ${newValueOfTheChangedObjInStr}`);
        let [success, err] = this.interactWithStorageClass.setInLocalStorage(
          this.keyToWriteInTheStorage,
          newValueOfTheChangedObjInStr
        );
        if (err !== null || !success) {
          console.error(`Error saving to storage: ${err}`);
        } else {
          console.log(`Successfully saved to storage`);
        }
      }
    });
  }
  /**
   * Validates that the stored object has the same structure as the shared state object using Zod
   * @param storedObj The object retrieved from local storage
   * @returns true if the structure matches, false otherwise
   */
  private validation(storedObj: unknown): boolean {
    try {
      // Use the schema to validate the stored object
      const result = this.schema.safeParse(storedObj);
      return result.success;
    } catch (error) {
      console.warn("Validation error:", error);
      return false;
    }
  }

  private getKeyObjFromTheStorage(): Object | Error {
    let [keyInstr, err] = this.interactWithStorageClass.getFromLocalStorage(this.keyToWriteInTheStorage);
    if (err != null || keyInstr === "" || keyInstr === null) {
      console.log(`didn't get the value of what you said in the js`);
      return err instanceof Error ? err : new Error("error in getting the value form the storage on the given key ->" + this.keyToWriteInTheStorage);
    }

    try {
      const parsedObj = JSON.parse(keyInstr);

      // Validate the parsed object against our schema
      if (this.validation(parsedObj)) {
        return parsedObj;
      } else {
        return new Error("The stored object doesn't match the expected structure");
      }
    } catch (err) {
      console.log(`the error in parsing the key in str -->${keyInstr} --- into json ${err}`);
      return err instanceof Error ? err : new Error(`there is a error in parsing the key into json ->${err}`);
    }
  }
}

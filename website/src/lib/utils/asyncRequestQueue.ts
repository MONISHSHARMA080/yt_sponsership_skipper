// export class AsyncRequestQueue<T> {
//   private queue: (() => Promise<T>)[] = [];
//   private activeCount = 0;
//   private results: ({result:T|null,error:Error|null})[];
//   private readonly concurrencyLimit: number;

//   constructor(concurrencyLimit: number) {
//     this.concurrencyLimit = concurrencyLimit;
//     this.results = []
//   }

//   addToQueue(PromiseFunc: (() => Promise<T>)[]) {
//     this.queue = [...this.queue, ...PromiseFunc];
//   }

//   private async processItem(currentIndex: number) {
//     if (this.queue.length === 0) return;

//     // Get and remove the next promise function from the queue
//     const promiseFunc = this.queue.shift();
//     if (!promiseFunc) return;

//     this.activeCount++;
//     this.results[currentIndex] = {result:null, error:null};


//     try {
//       // Execute the promise and store its result
//       const result = await promiseFunc();
//       console.log("the result from the promiseFunc() is ->", result);
      
//       this.results[currentIndex].result = result;
//       console.log("--> in the index ", currentIndex, " in the queue <--");
//     } catch (error) {
//       // Store any errors that occur
//       console.log("the error we got is ->", error);
//       console.log("the results array at "+currentIndex+" is ->", this.results[currentIndex], this.results);
      
      
//       if (error instanceof Error) {
//         this.results[currentIndex].error = error ;
//       }else{
//         this.results[currentIndex].error = Error("error is ->"+error)
//       }
//     }

//     this.activeCount--;

//     // If there are more items in the queue, process the next one
//     if (this.queue.length > 0) {
//       await this.processItem(currentIndex + 1);
//     }
//   }

//   async processQueue(): Promise<({result:T|null,error:Error|null})[]> {
//     let totalPromises = this.queue.length;
//     this.results = new Array(totalPromises);
//     const initialBatch = Math.min(totalPromises, this.concurrencyLimit);

//     const promiseStarter = Array(initialBatch)
//       .fill(0)
//       .map((func, index) => {
//         console.log("the index in the promiseStarter is ->" + index);
//         this.processItem(index);
//       });
//     let res = await Promise.all(promiseStarter);
//     console.log( "result of the promise starter is ->",res," and the promise starter is ", promiseStarter,);
//     return this.results;
//   }
// }


// async function main() {
//   const queue = new AsyncRequestQueue<string>(1);

//   const delay = (ms: number) =>
//     new Promise((resolve) => setTimeout(resolve, ms));

//   const promiseFunctions = [
//     async () => {
//       await delay(700);
//       console.log("\n\n  primise 1 completed \n\n");
//       return "A";
//     }, // Takes 1 second
//     async () => {
//       await delay(600);
//       console.log("\n\n  primise 2 completed \n\n");
//       return "B";
//     }, // Takes 1 second
//     async () => {
//       await delay(400);
//       console.log("\n\n  primise 3 completed \n\n");
//       return "C";
//     }, // Takes 1 second
//     async () => {
//       await delay(320);
//       console.log("\n\n  primise 4 completed \n\n");
//       return "D";
//     }, // Takes 1 second
//     async () => {
//       await delay(500);
//       console.log("\n\n  primise 5 completed \n\n");
//       return "E";
//     }, // Takes 1 second
//     async () => {
//       await delay(100);
//       console.log("\n\n  primise 6 completed \n\n");
//       return "F";
//     }, // Takes 1 second
//     async () => {
//       await delay(100);
//       console.log("\n\n  primise 7 completed \n\n");
//       return "G";
//     }, // Takes 1 second
//   ];

//   queue.addToQueue(promiseFunctions);
//   console.log("about to start processing promises");
//   const result = await queue.processQueue();
//   console.log("the result array is ->", result);
// }

// main();



export class AsyncRequestQueue<T> {
  private queue: (() => Promise<T>)[] = [];
  private activeCount = 0;
  private results: ({ result: T | null, error: Error | null })[] = [];
  private readonly concurrencyLimit: number;
  
  constructor(concurrencyLimit: number) {
    this.concurrencyLimit = concurrencyLimit;
  }
  
  addToQueue(promiseFuncs: (() => Promise<T>)[]) {
    this.queue = [...this.queue, ...promiseFuncs];
  }
  
  private async processItem(index: number): Promise<void> {
    if (this.queue.length === 0) return;
    
    // Get and remove the next promise function from the queue
    const promiseFunc = this.queue.shift();
    if (!promiseFunc) return;
    
    this.activeCount++;
    this.results[index] = { result: null, error: null };
    
    try {
      // Execute the promise and store its result
      const result = await promiseFunc();
      console.log(`Result for index ${index}:`, result);
      this.results[index].result = result;
    } catch (error) {
      console.log(`Error for index ${index}:`, error);
      
      if (error instanceof Error) {
        this.results[index].error = error;
      } else {
        this.results[index].error = new Error(`Unknown error: ${error}`);
      }
    } finally {
      this.activeCount--;
      
      // Process next item in queue if there are any
      if (this.queue.length > 0) {
        await this.processItem(this.results.length);
      }
    }
  }
  
  async processQueue(): Promise<({ result: T | null, error: Error | null })[]> {
    const totalPromises = this.queue.length;
    if (totalPromises === 0) return [];
    
    this.results = new Array(totalPromises).fill(null).map(() => ({ result: null, error: null }));
    const initialBatch = Math.min(totalPromises, this.concurrencyLimit);
    
    // Start processing initial batch of promises
    const processingPromises = Array(initialBatch)
      .fill(0)
      .map((_, index) => this.processItem(index));
    
    // Wait for all initial processors to complete
    await Promise.all(processingPromises);
    
    return this.results;
  }
}
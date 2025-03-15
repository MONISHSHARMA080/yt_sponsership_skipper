type processResult<T, R> = {
  result: R | null;
  error: Error | null;
  ifErrorThenOriginalPromise: Promise<T> | null;
};

type promiseArray<T> = Promise<T>[];
type processingQueue<T> = {
  promiseRunning: Promise<T>;
  indexOfThePromise: number;
}[];
type resultArray<T, R> = processResult<T, R>[];
type funcToProcessIndividualPromise<T, R> = (value: Promise<T>) => Promise<R>;

/** there are 2 generics cause, take for eg fetch if I have it then I have a promise that returns response but let's say I want
 * it to return a thing in form of json  form the body, then I can do it
 *
 * @param {R} R - it is the the return value/type when the function to run on individual promise is provided with the promise<T> and will be stored in the return array
 * @param {T} T - it is the type of what an individual promise returns , it will be given to yout function
 * */
export class AsyncRequestQueue<T, R> {
  private concurrencyLimit: number;
  private resultArray: resultArray<T, R>;
  private processingQueue: processingQueue<T>;
  private promiseQueue: promiseArray<T>;
  private promiseQueueSubmittedByUser: promiseArray<T>;

  constructor(concurrencyLimit: number = 5) {
    if (concurrencyLimit <= 0) {
      throw "concurrency can't be <= 0";
    }
    this.concurrencyLimit = concurrencyLimit;
    this.resultArray = [];
    this.processingQueue = [];
    this.promiseQueue= []
    this.promiseQueueSubmittedByUser = []
  }

  public process(
    promiseArray: promiseArray<T>,
    funcToProcessIndividualPromise: funcToProcessIndividualPromise<T, R>,
  ): Promise<resultArray<T, R>> {
    return new Promise((resolve, reject) => {
      this.promiseQueue = [...promiseArray]; // Create a copy
      this.promiseQueueSubmittedByUser = promiseArray;
      this.resultArray = new Array(promiseArray.length);
      this.processAll(resolve, funcToProcessIndividualPromise);
    });
  }

  private async processAll(
    resolveFunc: (value: resultArray<T, R>) => void,
    funcToProcessIndividualPromise: funcToProcessIndividualPromise<T, R>,
  ) {
    // if we have reached the end then
    if (this.promiseQueue.length === 0 && this.processingQueue.length === 0) {
      console.log(
        `length of the processqueue ->${this.promiseQueue.length} as we have reached the end`,
      );
      resolveFunc(this.resultArray);
    }

    // if we have reached the processing/concurrecy limit  then return or if we are at the last promise
    if (
      this.processingQueue.length >= this.concurrencyLimit ||
      this.promiseQueue.length === 0
    ) {
      console.log(`-- reached the concurrencyLimit -- `);
      console.log(
        `number of running promises are -> ${this.processingQueue.length} and the concurrency limit is ${this.concurrencyLimit} `,
      );

      return;
    }

    const promiseFromTheQueue = this.promiseQueue.shift();
    if (promiseFromTheQueue === undefined || promiseFromTheQueue === null) {
      return;
    }
    //pushing it in to the processing queue

    const indexOfThePromise = this.promiseQueueSubmittedByUser.indexOf(promiseFromTheQueue);

    this.processingQueue.push({promiseRunning:promiseFromTheQueue, indexOfThePromise:indexOfThePromise});

    this.processIndividualPromiseAndRemoveItFormTheProcessingQueue(
      promiseFromTheQueue,
      funcToProcessIndividualPromise,
      indexOfThePromise,
    ).finally(() => {
      console.log( `promise number ${indexOfThePromise} was completed  and now recursing `);
      this.processAll(resolveFunc, funcToProcessIndividualPromise);
    });
    this.processAll(resolveFunc, funcToProcessIndividualPromise);
  }

  /** process indvidual promises and then remove it form the  */
  private async processIndividualPromiseAndRemoveItFormTheProcessingQueue(
    promiseToProcess: Promise<T>,
    funcToProcessIndividualPromise: funcToProcessIndividualPromise<T, R>,
    indexOfPromise: number,
  ) {
    console.log(
      `processsig the promiose no. ${indexOfPromise}  and number of running promises are -> ${this.processingQueue.length} `,
    );

    try {
      let resultFormFunc =
        await funcToProcessIndividualPromise(promiseToProcess);
      console.log(" \n result for the func is ->", resultFormFunc);

      this.resultArray[indexOfPromise] = {
        result: resultFormFunc,
        error: null,
        ifErrorThenOriginalPromise: null,
      };
      this.removeFromTheProcessingQueue(promiseToProcess, indexOfPromise);
    } catch (error) {
      console.log(
        ` error in individual item at index ${indexOfPromise} is  ->`,
        error,
      );
      let errorInExecution =
        error instanceof Error
          ? error
          : new Error("there is a error executing the function->"+ error);
      this.resultArray[indexOfPromise] = {
        result: null,
        error: errorInExecution,
        ifErrorThenOriginalPromise: promiseToProcess,
      };
      this.removeFromTheProcessingQueue(promiseToProcess, indexOfPromise);
    }

    console.log(
      ` what do we have at this destination -> ${this.resultArray[indexOfPromise].result} `,
    );
  }

  private removeFromTheProcessingQueue(promiseToProcess: Promise<T>,indexOfPromise: number) {
    let index = this.processingQueue.findIndex(  item => item.indexOfThePromise === indexOfPromise);

    if (index < 0) {
      console.error("\nthe index of promise in the processing queue is <0 (the promise was not there), the index of promise in the promsie queue was --> ", indexOfPromise,"\n");
      return;
    }
    this.processingQueue.splice(index,1);
  
  }
}

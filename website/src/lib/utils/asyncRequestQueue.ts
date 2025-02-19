class AsyncRequestQueue<T> {
  private queue: (() => Promise<T>)[] = [];
  private activeCount = 0;
  private results: (T | Error)[];
  private readonly concurrencyLimit: number;

  constructor(concurrencyLimit: number) {
    this.concurrencyLimit = concurrencyLimit;
  }

  addToQueue(PromiseFunc: (() => Promise<T>)[]) {
    this.queue = [...this.queue, ...PromiseFunc];
  }

  private async processItem(currentIndex: number) {
    if (this.queue.length === 0) return;

    // Get and remove the next promise function from the queue
    const promiseFunc = this.queue.shift();
    if (!promiseFunc) return;

    this.activeCount++;

    try {
      // Execute the promise and store its result
      const result = await promiseFunc();
      this.results[currentIndex] = result;
      console.log("--> in the index ", currentIndex, " in the queue <--");
    } catch (error) {
      // Store any errors that occur
      this.results[currentIndex] = error as Error;
    }

    this.activeCount--;

    // If there are more items in the queue, process the next one
    if (this.queue.length > 0) {
      await this.processItem(currentIndex + 1);
    }
  }

  async processQueue(): Promise<(T | Error)[]> {
    let totalPromises = this.queue.length;
    this.results = new Array(totalPromises);
    const initialBatch = Math.min(totalPromises, this.concurrencyLimit);

    const promiseStarter = Array(initialBatch)
      .fill(0)
      .map((func, index) => {
        console.log("the index in the promiseStarter is ->" + index);
        this.processItem(index);
      });
    let res = await Promise.all(promiseStarter);
    console.log(
      "result of the promise starter is ->",
      res,
      " and the promise starter is ",
      promiseStarter,
    );

    return this.results;
  }
}

async function main() {
  const queue = new AsyncRequestQueue<string>(1);

  const delay = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  const promiseFunctions = [
    async () => {
      await delay(700);
      console.log("\n\n  primise 1 completed \n\n");
      return "A";
    }, // Takes 1 second
    async () => {
      await delay(600);
      console.log("\n\n  primise 2 completed \n\n");
      return "B";
    }, // Takes 1 second
    async () => {
      await delay(400);
      console.log("\n\n  primise 3 completed \n\n");
      return "C";
    }, // Takes 1 second
    async () => {
      await delay(320);
      console.log("\n\n  primise 4 completed \n\n");
      return "D";
    }, // Takes 1 second
    async () => {
      await delay(500);
      console.log("\n\n  primise 5 completed \n\n");
      return "E";
    }, // Takes 1 second
    async () => {
      await delay(100);
      console.log("\n\n  primise 6 completed \n\n");
      return "F";
    }, // Takes 1 second
    async () => {
      await delay(100);
      console.log("\n\n  primise 7 completed \n\n");
      return "G";
    }, // Takes 1 second
  ];

  queue.addToQueue(promiseFunctions);
  console.log("about to start processing promises");
  const result = await queue.processQueue();
  console.log("the result array is ->", result);
}

main();

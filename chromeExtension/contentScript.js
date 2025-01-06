console.log("hi form the contentScript :)");


async function main() {
  console.log("in the main");
  let [key, error] = await getKeyFromStorageOrBackend();
  if (error) {
    console.log(
      "error is there in getting the key and it is -->",
      error,
      "\n\n the key is",
      key,
    );
    return;
  }
  console.log("the key is  -->", key, "from the  content script, error is -->", error);
}
main();

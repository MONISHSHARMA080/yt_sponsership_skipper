import { sayHi, getKeyFromStorageOrBackend } from './helper.js';

console.log("hi from the service worker and will run say hi() now");
sayHi();
async  function sayHi1() {

    await getKeyFromStorageOrBackend()
}
sayHi1().then(r => console.log("value is r in " ,r)).catch(e => console.log("error is -->",e));
